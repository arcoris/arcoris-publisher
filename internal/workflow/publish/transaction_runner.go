// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package publish

import (
	"context"
	"fmt"
	"time"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

type transactionRunner struct {
	service Service
	request Request
	store   JournalStore
	journal TransactionJournal
}

// run executes the publish saga.
//
// The ordering is deliberate. Candidate refs are pushed to every repository
// before any final branch is promoted, and release tags are pushed only after
// every final branch points at the transaction commit. That does not make Git
// publication ACID across repositories, but it gives rollback a durable map of
// which compensating actions are safe to attempt.
func (r *transactionRunner) run(ctx context.Context, preflight []modulePreflight) (Result, error) {
	if err := r.setStatus(ctx, TransactionStatusPreflighted); err != nil {
		return Result{}, err
	}
	if err := r.snapshot(ctx); err != nil {
		return r.fail(ctx, err)
	}
	if err := r.commitLocal(ctx, preflight); err != nil {
		return r.fail(ctx, err)
	}
	if err := r.pushCandidates(ctx); err != nil {
		return r.fail(ctx, err)
	}
	if err := r.promoteBranches(ctx); err != nil {
		return r.fail(ctx, err)
	}
	if err := r.publishTags(ctx); err != nil {
		return r.fail(ctx, err)
	}
	if err := r.cleanupCandidates(ctx); err != nil {
		r.journal.Warnings = append(r.journal.Warnings, err.Error())
	}
	if err := r.setStatus(ctx, TransactionStatusCommitted); err != nil {
		return Result{modules: r.results(), transaction: r.journal}, err
	}

	return Result{modules: r.results(), transaction: r.journal}, nil
}

func (r *transactionRunner) snapshot(ctx context.Context) error {
	for i := range r.journal.Modules {
		mod := &r.journal.Modules[i]
		if mod.Skipped {
			continue
		}
		head, err := r.service.deps.Git.Head(ctx, mod.WorktreeDir)
		if err != nil {
			return transactionError(CodeTransactionFailed, mod.Module, "local HEAD snapshot failed", err)
		}
		branch, err := r.service.deps.Git.CurrentBranch(ctx, mod.WorktreeDir)
		if err != nil {
			return transactionError(CodeTransactionFailed, mod.Module, "local branch snapshot failed", err)
		}
		remoteHead, exists, err := r.service.deps.Git.RemoteRefHash(ctx, mod.WorktreeDir, r.service.opts.RemoteName, mod.FinalBranchRef)
		if err != nil {
			return transactionError(CodeTransactionFailed, mod.Module, "remote branch snapshot failed", err)
		}
		mod.LocalBaseHead = head
		mod.LocalBaseBranch = branch
		mod.RemoteBaseCommit = remoteHead
		mod.RemoteBaseExists = exists
	}
	return r.setStatus(ctx, TransactionStatusSnapshotted)
}

// commitLocal creates local commits but does not expose them through final
// remote refs. If a later candidate push fails, rollback can reset the worktree
// to LocalBaseHead without touching public branches.
func (r *transactionRunner) commitLocal(ctx context.Context, preflight []modulePreflight) error {
	byModule := make(map[string]modulePreflight, len(preflight))
	for _, item := range preflight {
		byModule[item.mod.Name().String()] = item
	}
	for i := range r.journal.Modules {
		mod := &r.journal.Modules[i]
		if mod.Skipped {
			continue
		}
		item := byModule[mod.Module.String()]
		commit, err := r.service.commitWorktree(ctx, mod.WorktreeDir, item.mod, item.sourceModule, r.request)
		if err != nil {
			return transactionError(CodePublishFailed, mod.Module, "local commit failed", err)
		}
		mod.CreatedCommit = commit
		if err := r.update(ctx); err != nil {
			return err
		}
	}
	return r.setStatus(ctx, TransactionStatusCommittedLocally)
}

// pushCandidates publishes transaction-private refs for every changed module.
// A failure here is still cheap to roll back because no final branch or tag has
// been updated.
func (r *transactionRunner) pushCandidates(ctx context.Context) error {
	for i := range r.journal.Modules {
		mod := &r.journal.Modules[i]
		if mod.Skipped {
			continue
		}
		refspec := git.RefSpec(mod.CreatedCommit.String() + ":" + mod.CandidateBranchRef)
		if err := r.service.deps.Git.Push(ctx, mod.WorktreeDir, r.service.opts.RemoteName, refspec, git.PushOptions{}); err != nil {
			return transactionError(CodeCandidatePushFailed, mod.Module, "candidate ref push failed", err)
		}
		mod.CandidatePushed = true
		if err := r.update(ctx); err != nil {
			return err
		}
	}
	return r.setStatus(ctx, TransactionStatusCandidatesPushed)
}

// promoteBranches updates final branches only after all candidate refs exist.
// Each promotion rechecks the captured remote base so an external push during
// the transaction fails before we overwrite it.
func (r *transactionRunner) promoteBranches(ctx context.Context) error {
	if err := r.setStatus(ctx, TransactionStatusPromoting); err != nil {
		return err
	}
	for i := range r.journal.Modules {
		mod := &r.journal.Modules[i]
		if mod.Skipped {
			continue
		}
		if err := r.ensureRemoteBaseUnchanged(ctx, mod); err != nil {
			return err
		}
		refspec := git.RefSpec(mod.CreatedCommit.String() + ":" + mod.FinalBranchRef)
		opts := pushOptions(r.request)
		if opts.ForceWithLease && mod.RemoteBaseExists {
			opts.ForceWithLeaseRef = mod.FinalBranchRef
			opts.ForceWithLeaseExpect = mod.RemoteBaseCommit
		}
		if err := r.service.deps.Git.Push(ctx, mod.WorktreeDir, r.service.opts.RemoteName, refspec, opts); err != nil {
			return transactionError(CodePromotionFailed, mod.Module, "final branch promotion failed", err)
		}
		mod.FinalBranchPromoted = true
		if err := r.update(ctx); err != nil {
			return err
		}
	}
	return r.setStatus(ctx, TransactionStatusBranchesPromoted)
}

// publishTags is intentionally last. Tags are harder for users to reason about
// once observed, so branches must already be consistently promoted before tags
// are created and pushed.
func (r *transactionRunner) publishTags(ctx context.Context) error {
	if !r.request.Plan.PublishPolicy().Tags().Enabled() {
		return nil
	}
	if err := r.setStatus(ctx, TransactionStatusTagging); err != nil {
		return err
	}
	for i := range r.journal.Modules {
		mod := &r.journal.Modules[i]
		if mod.Skipped {
			continue
		}
		tag := git.TagName(r.tagName(mod).String())
		if err := r.service.deps.Git.CreateTag(ctx, mod.WorktreeDir, tag, mod.CreatedCommit, git.TagOptions{
			Annotated: true,
			Message:   "release " + tag.String(),
		}); err != nil {
			return transactionError(CodeTagPublishFailed, mod.Module, "local tag creation failed", err)
		}
		mod.LocalTagCreated = true
		if err := r.update(ctx); err != nil {
			return err
		}
		if err := r.service.deps.Git.PushTag(ctx, mod.WorktreeDir, r.service.opts.RemoteName, tag, git.PushOptions{}); err != nil {
			return transactionError(CodeTagPublishFailed, mod.Module, "remote tag push failed", err)
		}
		mod.RemoteTagPushed = true
		hash, ok, err := r.service.deps.Git.RemoteRefHash(ctx, mod.WorktreeDir, r.service.opts.RemoteName, mod.FinalTagRef)
		if err != nil {
			return transactionError(CodeTagPublishFailed, mod.Module, "remote tag snapshot failed after push", err)
		}
		if !ok {
			return transactionError(CodeTagPublishFailed, mod.Module, "remote tag snapshot missing after push", fmt.Errorf("ref %s is absent", mod.FinalTagRef))
		}
		mod.RemoteTagHash = hash
		if err := r.update(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *transactionRunner) cleanupCandidates(ctx context.Context) error {
	for i := range r.journal.Modules {
		mod := &r.journal.Modules[i]
		if !mod.CandidatePushed {
			continue
		}
		if err := r.service.deps.Git.DeleteRemoteRef(ctx, mod.WorktreeDir, r.service.opts.RemoteName, mod.CandidateBranchRef, git.PushOptions{}); err != nil {
			return transactionError(CodePublishFailed, mod.Module, "candidate ref cleanup failed", err)
		}
		mod.CandidatePushed = false
	}
	return r.update(ctx)
}

func (r *transactionRunner) ensureRemoteBaseUnchanged(ctx context.Context, mod *ModuleTransactionState) error {
	current, exists, err := r.service.deps.Git.RemoteRefHash(ctx, mod.WorktreeDir, r.service.opts.RemoteName, mod.FinalBranchRef)
	if err != nil {
		return transactionError(CodePromotionFailed, mod.Module, "remote branch lease lookup failed", err)
	}
	if exists != mod.RemoteBaseExists || current != mod.RemoteBaseCommit {
		return &Error{
			Code: CodePromotionFailed,
			Message: fmt.Sprintf(
				"module %s remote branch %s changed during transaction",
				mod.Module,
				mod.FinalBranchRef,
			),
		}
	}
	return nil
}

func (r *transactionRunner) fail(ctx context.Context, err error) (Result, error) {
	r.journal.Failure = err.Error()
	_ = r.setStatus(ctx, TransactionStatusFailed)
	if r.service.opts.RollbackMode == RollbackAutomatic {
		rollbackErr := r.rollback(ctx)
		if rollbackErr != nil {
			return Result{modules: r.results(), transaction: r.journal}, rollbackErr
		}
		return Result{modules: r.results(), transaction: r.journal}, err
	}
	if r.service.opts.RollbackMode == RollbackManual {
		r.journal.Rollback = RollbackStatusPending
		_ = r.update(ctx)
	}
	return Result{modules: r.results(), transaction: r.journal}, err
}

func (r *transactionRunner) setStatus(ctx context.Context, status TransactionStatus) error {
	r.journal.Status = status
	r.journal.UpdatedAt = time.Now().UTC()
	if err := r.store.Update(ctx, r.journal); err != nil {
		return &Error{Code: CodeJournalFailed, Message: "update transaction journal failed", Cause: err}
	}
	return nil
}

func (r *transactionRunner) update(ctx context.Context) error {
	r.journal.UpdatedAt = time.Now().UTC()
	if err := r.store.Update(ctx, r.journal); err != nil {
		return &Error{Code: CodeJournalFailed, Message: "update transaction journal failed", Cause: err}
	}
	return nil
}

func (r *transactionRunner) results() []ModuleResult {
	out := make([]ModuleResult, 0, len(r.journal.Modules))
	for _, mod := range r.journal.Modules {
		result := ModuleResult{
			module:  mod.Module,
			commit:  mod.CreatedCommit,
			pushed:  r.journal.Status == TransactionStatusCommitted && mod.FinalBranchPromoted,
			skipped: mod.Skipped,
		}
		if result.pushed && mod.FinalTagRef != "" {
			result.tags = []git.TagName{git.TagName(r.tagName(&mod))}
		}
		out = append(out, result)
	}
	return out
}

func (r *transactionRunner) tagName(mod *ModuleTransactionState) git.TagName {
	if mod.FinalTagRef == "" {
		return ""
	}
	return git.TagName(stringsTrimPrefix(mod.FinalTagRef, "refs/tags/"))
}

func transactionError(code ErrorCode, module fmt.Stringer, message string, cause error) error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf("module %s %s", module.String(), message),
		Cause:   cause,
	}
}

func stringsTrimPrefix(value string, prefix string) string {
	if len(value) >= len(prefix) && value[:len(prefix)] == prefix {
		return value[len(prefix):]
	}
	return value
}
