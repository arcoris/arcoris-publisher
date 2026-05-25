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
	"strings"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

func (r *transactionRunner) rollback(ctx context.Context) error {
	r.journal.Rollback = RollbackStatusPending
	_ = r.setStatus(ctx, TransactionStatusRollingBack)

	for i := len(r.journal.Modules) - 1; i >= 0; i-- {
		r.rollbackModule(ctx, &r.journal.Modules[i])
		_ = r.update(ctx)
	}

	if len(r.journal.ManualActions) > 0 {
		r.journal.Rollback = RollbackStatusFailed
		_ = r.setStatus(ctx, TransactionStatusRollbackFailed)
		return &Error{
			Code:    CodeRollbackFailed,
			Message: fmt.Sprintf("transaction %s rollback requires manual recovery", r.journal.ID),
		}
	}

	r.journal.Rollback = RollbackStatusSucceeded
	if err := r.setStatus(ctx, TransactionStatusRolledBack); err != nil {
		return err
	}
	return nil
}

// rollbackModule reverses side effects from most-public to least-public state:
// tags first, then final branches, then candidate refs, then local worktree
// state. This order avoids leaving a release tag pointing at a branch that has
// already been restored.
func (r *transactionRunner) rollbackModule(ctx context.Context, mod *ModuleTransactionState) {
	if mod.Skipped {
		return
	}
	r.rollbackRemoteTag(ctx, mod)
	r.rollbackFinalBranch(ctx, mod)
	r.rollbackCandidate(ctx, mod)
	r.rollbackLocalTag(ctx, mod)
	r.rollbackLocalWorktree(ctx, mod)
}

// rollbackRemoteTag deletes only a tag that still appears to be the tag created
// by this transaction. If the remote tag moved, the function records a manual
// recovery action instead of deleting an operator-owned ref.
func (r *transactionRunner) rollbackRemoteTag(ctx context.Context, mod *ModuleTransactionState) {
	if !mod.RemoteTagPushed || mod.FinalTagRef == "" {
		return
	}
	if mod.RemoteTagHash != "" {
		current, ok, err := r.service.deps.Git.RemoteRefHash(ctx, mod.WorktreeDir, r.service.opts.RemoteName, mod.FinalTagRef)
		if err != nil {
			r.recordRollbackFailure(mod, "remote tag lookup failed", err)
			return
		}
		if ok && current != mod.RemoteTagHash {
			r.recordManualAction(*mod, mod.FinalTagRef, mod.RemoteTagHash, "", "remote tag changed after transaction; refusing to delete")
			return
		}
	}
	if err := r.service.deps.Git.DeleteRemoteRef(ctx, mod.WorktreeDir, r.service.opts.RemoteName, mod.FinalTagRef, git.PushOptions{}); err != nil {
		r.recordRollbackFailure(mod, "remote tag delete failed", err)
		return
	}
	mod.Rollback.RemoteTagDeleted = true
	mod.RemoteTagPushed = false
}

// rollbackFinalBranch restores or deletes a final branch only when the current
// remote ref still equals the transaction-created commit. If somebody moved the
// branch after our push, rollback refuses to overwrite that newer state.
func (r *transactionRunner) rollbackFinalBranch(ctx context.Context, mod *ModuleTransactionState) {
	if !mod.FinalBranchPromoted {
		return
	}
	current, ok, err := r.service.deps.Git.RemoteRefHash(ctx, mod.WorktreeDir, r.service.opts.RemoteName, mod.FinalBranchRef)
	if err != nil {
		r.recordRollbackFailure(mod, "remote branch lookup failed", err)
		return
	}
	if !ok || current != mod.CreatedCommit {
		r.recordManualAction(*mod, mod.FinalBranchRef, mod.CreatedCommit, mod.RemoteBaseCommit, "remote branch moved after transaction; refusing to overwrite")
		return
	}

	if mod.RemoteBaseExists {
		refspec := git.RefSpec(mod.RemoteBaseCommit.String() + ":" + mod.FinalBranchRef)
		if err := r.service.deps.Git.Push(ctx, mod.WorktreeDir, r.service.opts.RemoteName, refspec, git.PushOptions{ForceWithLease: true}); err != nil {
			r.recordRollbackFailure(mod, "remote branch restore failed", err)
			return
		}
	} else if err := r.service.deps.Git.DeleteRemoteRef(ctx, mod.WorktreeDir, r.service.opts.RemoteName, mod.FinalBranchRef, git.PushOptions{}); err != nil {
		r.recordRollbackFailure(mod, "remote branch delete failed", err)
		return
	}

	mod.Rollback.FinalBranchRestored = true
	mod.FinalBranchPromoted = false
}

func (r *transactionRunner) rollbackCandidate(ctx context.Context, mod *ModuleTransactionState) {
	if !mod.CandidatePushed {
		return
	}
	if err := r.service.deps.Git.DeleteRemoteRef(ctx, mod.WorktreeDir, r.service.opts.RemoteName, mod.CandidateBranchRef, git.PushOptions{}); err != nil {
		r.recordRollbackFailure(mod, "candidate ref delete failed", err)
		return
	}
	mod.Rollback.CandidateDeleted = true
	mod.CandidatePushed = false
}

func (r *transactionRunner) rollbackLocalTag(ctx context.Context, mod *ModuleTransactionState) {
	if !mod.LocalTagCreated || mod.FinalTagRef == "" {
		return
	}
	if err := r.service.deps.Git.DeleteTag(ctx, mod.WorktreeDir, git.TagName(strings.TrimPrefix(mod.FinalTagRef, "refs/tags/"))); err != nil {
		r.recordRollbackFailure(mod, "local tag delete failed", err)
		return
	}
	mod.Rollback.LocalTagDeleted = true
	mod.LocalTagCreated = false
}

func (r *transactionRunner) rollbackLocalWorktree(ctx context.Context, mod *ModuleTransactionState) {
	if mod.LocalBaseHead == "" {
		return
	}
	if mod.LocalBaseBranch != "" {
		if err := r.service.deps.Git.Checkout(ctx, mod.WorktreeDir, mod.LocalBaseBranch.String(), git.CheckoutOptions{Force: true}); err != nil {
			r.recordRollbackFailure(mod, "local branch checkout failed", err)
			return
		}
	}
	if err := r.service.deps.Git.ResetHard(ctx, mod.WorktreeDir, mod.LocalBaseHead.String()); err != nil {
		r.recordRollbackFailure(mod, "local reset failed", err)
		return
	}
	if err := r.service.deps.Git.Clean(ctx, mod.WorktreeDir, git.CleanOptions{
		RemoveUntracked: true,
		RemoveIgnored:   true,
		Directories:     true,
		Force:           true,
	}); err != nil {
		r.recordRollbackFailure(mod, "local clean failed", err)
		return
	}
	mod.Rollback.LocalWorktreeReset = true
}

func (r *transactionRunner) recordRollbackFailure(mod *ModuleTransactionState, action string, err error) {
	message := action + ": " + err.Error()
	mod.Rollback.FailedActions = append(mod.Rollback.FailedActions, message)
	r.journal.ManualActions = append(r.journal.ManualActions, ManualRecoveryAction{
		Module:     mod.Module,
		Repository: mod.Repository,
		Ref:        mod.FinalBranchRef,
		Message:    message,
	})
}

func (r *transactionRunner) recordManualAction(
	mod ModuleTransactionState,
	ref string,
	expected git.CommitHash,
	desired git.CommitHash,
	message string,
) {
	command := ""
	switch {
	case desired != "":
		command = fmt.Sprintf("git push %s %s:%s", r.service.opts.RemoteName, desired, ref)
	default:
		command = fmt.Sprintf("git push %s :%s", r.service.opts.RemoteName, ref)
	}
	r.journal.ManualActions = append(r.journal.ManualActions, ManualRecoveryAction{
		Module:       mod.Module,
		Repository:   mod.Repository,
		Ref:          ref,
		ExpectedHash: expected,
		DesiredHash:  desired,
		Message:      message,
		Command:      command,
	})
}
