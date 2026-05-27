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

package preflight

import (
	"context"
	"fmt"
	"path/filepath"

	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

// Service checks whether publish can safely start without mutating target
// files, local refs, remote refs, journals, or locks.
type Service struct {
	deps Dependencies
	opts Options
}

// New returns a preflight service.
func New(deps Dependencies, opts Options) Service {
	defaults := DefaultOptions()
	if opts.RemoteName == "" {
		opts.RemoteName = defaults.RemoteName
	}
	if deps.Source.FS == nil {
		deps.Source.FS = deps.FS
	}
	if deps.Source.Git == nil {
		deps.Source.Git = deps.Git
	}
	return Service{deps: deps, opts: opts}
}

// Check validates current publish readiness. It may query remotes for refs, but
// it never constructs files, checks out branches, fetches, pushes, or writes
// transaction state.
func (s Service) Check(ctx context.Context, req Request) (Result, error) {
	if err := s.validateDependencies(); err != nil {
		return Result{}, err
	}

	builder := resultBuilder{result: Result{version: versionOf(req.Plan)}}
	if req.Plan.Empty() {
		builder.addGlobal(failed("plan", "empty_plan", "publication plan is empty"))
		return builder.build(), nil
	}
	builder.addGlobal(passed("plan", "publication plan built"))

	sourceSnapshot, ok := s.checkSource(ctx, req, &builder)
	if ok {
		builder.addGlobal(passed("source-repository", "source repository inspected"))
		for _, issue := range sourceSnapshot.Warnings() {
			builder.addGlobal(warning("source-repository", string(issue.Code), issue.Message))
		}
	}

	root, ok := s.checkTargetRoot(ctx, req, &builder)
	stateDir := s.stateDir(root)
	s.checkTransactions(ctx, stateDir, &builder)

	if ok {
		for _, mod := range req.Plan.Modules() {
			builder.addModule(s.checkModule(ctx, root, req.Plan, mod))
		}
	}

	return builder.build(), nil
}

func (s Service) validateDependencies() error {
	if s.deps.FS == nil {
		return fmt.Errorf("preflight filesystem dependency is required")
	}
	if s.deps.Git == nil {
		return fmt.Errorf("preflight git dependency is required")
	}
	return nil
}

func (s Service) checkSource(ctx context.Context, req Request, builder *resultBuilder) (source.Snapshot, bool) {
	snapshot, err := source.New(s.deps.Source, s.opts.Source).Inspect(ctx, source.Request{
		Plan:          req.Plan,
		RepositoryDir: req.SourceRepositoryDir,
		StagingDir:    req.StagingDir,
	})
	if err != nil {
		builder.addGlobal(failed("source-repository", "source_inspection_failed", "source repository inspection failed"))
		return source.Snapshot{}, false
	}
	return snapshot, true
}

func (s Service) checkTargetRoot(ctx context.Context, req Request, builder *resultBuilder) (string, bool) {
	root, err := pathutil.CleanAbs(req.TargetRootDir)
	if err != nil {
		builder.addGlobal(failed("target-root", "invalid_target_root", "target root path is invalid"))
		return "", false
	}
	exists, err := s.deps.FS.Exists(ctx, root)
	if err != nil {
		builder.addGlobal(failed("target-root", "target_root_lookup_failed", "target root lookup failed"))
		return root, false
	}
	if !exists {
		builder.addGlobal(failed("target-root", "target_root_missing", "target root does not exist"))
		return root, false
	}
	isDir, err := s.deps.FS.IsDir(ctx, root)
	if err != nil {
		builder.addGlobal(failed("target-root", "target_root_lookup_failed", "target root lookup failed"))
		return root, false
	}
	if !isDir {
		builder.addGlobal(failed("target-root", "target_root_not_directory", "target root is not a directory"))
		return root, false
	}
	builder.addGlobal(pathCheck(passed("target-root", "target root exists"), root))
	return root, true
}

func (s Service) checkTransactions(ctx context.Context, stateDir string, builder *resultBuilder) {
	if stateDir == "" {
		builder.addGlobal(failed("pending-transactions", "state_dir_missing", "transaction state directory is unavailable"))
		return
	}

	diagnostics, err := publish.InspectTransactionState(ctx, stateDir)
	s.checkPendingTransactions(diagnostics, err, builder)
	s.checkPublishLock(diagnostics, err, builder)
}

func (s Service) checkPendingTransactions(diagnostics publish.TransactionStateDiagnostics, err error, builder *resultBuilder) {
	if blocker, ok := firstTransactionJournalBlocker(diagnostics.Blockers); ok {
		builder.addGlobal(pendingTransactionCheck(blocker))
		return
	}
	if err != nil && diagnostics.Lock.Status != publish.LockShowStatusFailed {
		builder.addGlobal(failed("pending-transactions", "journal_lookup_failed", "transaction journal lookup failed"))
		return
	}
	builder.addGlobal(passed("pending-transactions", "no pending transactions"))
}

func firstTransactionJournalBlocker(blockers []publish.TransactionStateBlocker) (publish.TransactionStateBlocker, bool) {
	for _, blocker := range blockers {
		switch blocker.Kind {
		case publish.TransactionBlockerActiveJournal,
			publish.TransactionBlockerFailedJournal,
			publish.TransactionBlockerRollbackFailed,
			publish.TransactionBlockerCorruptJournal,
			publish.TransactionBlockerJournalFileReadFailed,
			publish.TransactionBlockerJournalDirectoryReadFailed:
			return blocker, true
		}
	}
	return publish.TransactionStateBlocker{}, false
}

func pendingTransactionCheck(blocker publish.TransactionStateBlocker) CheckResult {
	switch blocker.Kind {
	case publish.TransactionBlockerCorruptJournal:
		return failed("pending-transactions", "transaction_journal_corrupt", "transaction journal is corrupt")
	case publish.TransactionBlockerJournalDirectoryReadFailed:
		return failed("pending-transactions", "transaction_journal_directory_read_failed", "transaction journal directory lookup failed")
	case publish.TransactionBlockerJournalFileReadFailed:
		return failed("pending-transactions", "transaction_journal_file_read_failed", "transaction journal file lookup failed")
	case publish.TransactionBlockerFailedJournal, publish.TransactionBlockerRollbackFailed:
		return failed(
			"pending-transactions",
			"transaction_recovery_required",
			fmt.Sprintf("publish transaction %s has recovery status %s", blocker.TransactionID, blocker.Status),
		)
	default:
		return failed(
			"pending-transactions",
			"pending_transaction",
			fmt.Sprintf("pending publish transaction %s has status %s", blocker.TransactionID, blocker.Status),
		)
	}
}

func (s Service) checkPublishLock(diagnostics publish.TransactionStateDiagnostics, err error, builder *resultBuilder) {
	if blocker, ok := firstLockBlocker(diagnostics); ok {
		builder.addGlobal(lockBlockerCheck(blocker))
		return
	}
	if diagnostics.Lock.Status == publish.LockShowStatusAbsent {
		builder.addGlobal(passed("publish-lock", "publish lock absent"))
		return
	}
	if err != nil {
		builder.addGlobal(failed("publish-lock", "lock_lookup_failed", "publish lock lookup failed"))
		return
	}
	builder.addGlobal(failed("publish-lock", "lock_lookup_failed", "publish lock state is unavailable"))
}

func firstLockBlocker(diagnostics publish.TransactionStateDiagnostics) (publish.TransactionStateBlocker, bool) {
	lockID := diagnostics.Lock.Lock.ID
	for _, blocker := range diagnostics.Blockers {
		switch blocker.Kind {
		case publish.TransactionBlockerPublishLock,
			publish.TransactionBlockerMissingLockJournal,
			publish.TransactionBlockerCorruptLock,
			publish.TransactionBlockerLockReadFailed:
			return blocker, true
		case publish.TransactionBlockerCorruptJournal,
			publish.TransactionBlockerJournalFileReadFailed:
			if lockID != "" && blocker.TransactionID == lockID {
				return blocker, true
			}
		}
	}
	return publish.TransactionStateBlocker{}, false
}

func lockBlockerCheck(blocker publish.TransactionStateBlocker) CheckResult {
	switch blocker.Kind {
	case publish.TransactionBlockerPublishLock:
		switch blocker.Reason {
		case publish.TransactionBlockerReasonStaleTerminalLock:
			return failed("publish-lock", "stale_publish_lock_terminal_transaction", fmt.Sprintf("publish lock references terminal transaction %s with status %s", blocker.TransactionID, blocker.Status))
		case publish.TransactionBlockerReasonRecoveryLock:
			return failed("publish-lock", "publish_lock_recovery_required", fmt.Sprintf("publish lock references recovery transaction %s with status %s", blocker.TransactionID, blocker.Status))
		default:
			return failed("publish-lock", "publish_lock_exists", fmt.Sprintf("publish lock exists for active transaction %s with status %s", blocker.TransactionID, blocker.Status))
		}
	case publish.TransactionBlockerMissingLockJournal:
		return failed("publish-lock", "stale_publish_lock_journal_missing", fmt.Sprintf("publish lock references missing transaction journal %s", blocker.TransactionID))
	case publish.TransactionBlockerCorruptLock:
		return failed("publish-lock", "publish_lock_corrupt", "publish lock is not parseable")
	case publish.TransactionBlockerLockReadFailed:
		return failed("publish-lock", "publish_lock_read_failed", "publish lock lookup failed")
	case publish.TransactionBlockerCorruptJournal:
		return failed("publish-lock", "publish_lock_journal_corrupt", fmt.Sprintf("publish lock references corrupt transaction journal %s", blocker.TransactionID))
	case publish.TransactionBlockerJournalFileReadFailed:
		return failed("publish-lock", "publish_lock_journal_read_failed", fmt.Sprintf("publish lock references unreadable transaction journal %s", blocker.TransactionID))
	default:
		return failed("publish-lock", "lock_lookup_failed", "publish lock state is unavailable")
	}
}

func (s Service) checkModule(ctx context.Context, root string, p plan.Plan, mod plan.ModulePlan) ModuleResult {
	worktree := target.RepositoryWorktree(root, mod.Repository())
	result := ModuleResult{
		name:        mod.Name(),
		repository:  mod.Repository(),
		worktreeDir: worktree,
	}

	result.checks = append(result.checks, s.checkBranches(mod)...)
	result.checks = append(result.checks, s.checkWorktree(ctx, mod, worktree)...)
	result.checks = append(result.checks, skipped("target-fetch", "fetch skipped by read-only preflight"))
	result.checks = append(result.checks, s.checkTags(ctx, p, mod, worktree)...)
	result.checks = append(result.checks, s.checkCandidateRef(mod))
	result.checks = append(result.checks, passed("publish-entries", "explicit publish entries are planned"))
	if p.PublishPolicy().Provenance().FileEnabled() {
		result.checks = append(result.checks, passed("provenance-path", "provenance path is valid in the plan"))
	}

	return result
}

func (s Service) checkBranches(mod plan.ModulePlan) []CheckResult {
	branches := mod.Branches()
	if len(branches) != 1 {
		return []CheckResult{failed(
			"multi-branch",
			"unsupported_multi_branch",
			fmt.Sprintf("module has %d branch mappings; multi-branch publish is not supported", len(branches)),
		)}
	}
	ref := publish.BranchRef(branches[0].Target())
	if err := publish.ValidateGitRef(ref); err != nil {
		return []CheckResult{failed("target-branch", "unsafe_branch_ref", "target branch is not a safe Git ref")}
	}
	return []CheckResult{passed("multi-branch", "single target branch"), passed("target-branch", "target branch ref is safe")}
}

func (s Service) checkWorktree(ctx context.Context, mod plan.ModulePlan, worktree string) []CheckResult {
	checks := []CheckResult{}
	exists, err := s.deps.FS.Exists(ctx, worktree)
	if err != nil {
		return []CheckResult{pathCheck(failed("target-worktree", "worktree_lookup_failed", "target worktree lookup failed"), worktree)}
	}
	if !exists {
		return []CheckResult{pathCheck(failed("target-worktree", "worktree_missing", "target worktree is missing"), worktree)}
	}
	isDir, err := s.deps.FS.IsDir(ctx, worktree)
	if err != nil {
		return []CheckResult{pathCheck(failed("target-worktree", "worktree_lookup_failed", "target worktree lookup failed"), worktree)}
	}
	if !isDir {
		return []CheckResult{pathCheck(failed("target-worktree", "worktree_not_directory", "target worktree is not a directory"), worktree)}
	}
	checks = append(checks, pathCheck(passed("target-worktree", "target worktree exists"), worktree))

	status, err := s.deps.Git.Status(ctx, worktree)
	if err != nil {
		return append(checks, pathCheck(failed("target-status", "target_status_failed", "target Git status failed"), worktree))
	}
	if status.IsDirty() {
		return append(checks, pathCheck(failed("target-status", "target_worktree_dirty", "target worktree is dirty"), worktree))
	}
	checks = append(checks, pathCheck(passed("target-status", "target worktree is clean"), worktree))
	checks = append(checks, s.checkCommitIdentity(ctx, worktree))

	if len(mod.Branches()) == 1 {
		ref := publish.BranchRef(mod.Branches()[0].Target())
		exists, err := s.deps.Git.RefExists(ctx, worktree, ref)
		if err != nil {
			checks = append(checks, failed("target-branch", "target_branch_lookup_failed", "target branch lookup failed"))
		} else if !exists {
			checks = append(checks, failed("target-branch", "target_branch_missing", "target branch is missing locally"))
		} else {
			checks = append(checks, passed("target-branch", "target branch exists locally"))
		}
		if _, _, err := s.deps.Git.RemoteRefHash(ctx, worktree, s.opts.RemoteName, ref); err != nil {
			checks = append(checks, failed("remote-branch", "remote_branch_lookup_failed", "remote branch lookup failed"))
		} else {
			checks = append(checks, passed("remote-branch", "remote branch lookup succeeded"))
		}
	}

	return checks
}

func (s Service) checkCommitIdentity(ctx context.Context, worktree string) CheckResult {
	check := target.CheckCommitIdentity(ctx, s.deps.Git, worktree)
	if check.Passed() {
		return passed("commit-identity", "Git commit identity is configured")
	}
	return failed("commit-identity", check.Code(), check.Message())
}

func (s Service) checkTags(ctx context.Context, p plan.Plan, mod plan.ModulePlan, worktree string) []CheckResult {
	if !p.PublishPolicy().Tags().Enabled() {
		return []CheckResult{skipped("local-tag", "tag publishing disabled"), skipped("remote-tag", "tag publishing disabled")}
	}

	tag := git.TagName(mod.Version().String())
	checks := []CheckResult{}
	localExists, err := s.deps.Git.TagExists(ctx, worktree, tag)
	if err != nil {
		checks = append(checks, failed("local-tag", "local_tag_lookup_failed", "local tag lookup failed"))
	} else if localExists {
		checks = append(checks, failed("local-tag", "local_tag_exists", "local release tag already exists"))
	} else {
		checks = append(checks, passed("local-tag", "local release tag absent"))
	}

	remoteRef := "refs/tags/" + tag.String()
	remoteExists, err := s.deps.Git.RemoteRefExists(ctx, worktree, s.opts.RemoteName, remoteRef)
	if err != nil {
		checks = append(checks, failed("remote-tag", "remote_tag_lookup_failed", "remote tag lookup failed"))
	} else if remoteExists {
		checks = append(checks, failed("remote-tag", "remote_tag_exists", "remote release tag already exists"))
	} else {
		checks = append(checks, passed("remote-tag", "remote release tag absent"))
	}
	return checks
}

func (s Service) checkCandidateRef(mod plan.ModulePlan) CheckResult {
	ref := publish.CandidateRef("tx-preflight", mod.Name())
	if err := publish.ValidateGitRef(ref); err != nil {
		return failed("candidate-ref", "unsafe_candidate_ref", "candidate ref would be unsafe")
	}
	return passed("candidate-ref", "candidate ref namespace is safe")
}

func (s Service) stateDir(targetRoot string) string {
	if s.opts.StateDir != "" {
		return s.opts.StateDir
	}
	if targetRoot == "" {
		return ""
	}
	return filepath.Join(targetRoot, ".arcpub", "state")
}

func versionOf(p plan.Plan) string {
	for _, mod := range p.Modules() {
		return mod.Version().String()
	}
	return ""
}
