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

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
)

// Service publishes verified target repositories.
type Service struct {
	// deps contains infrastructure ports used by publication.
	deps Dependencies

	// opts contains normalized publication options.
	opts Options
}

// New returns a publication service.
func New(deps Dependencies, opts Options) Service {
	defaults := DefaultOptions()
	if opts.RemoteName == "" {
		opts.RemoteName = defaults.RemoteName
	}
	if opts.RollbackMode == "" {
		opts.RollbackMode = defaults.RollbackMode
	}
	return Service{deps: deps, opts: opts}
}

// Publish commits, tags, and pushes every changed verified module.
func (s Service) Publish(ctx context.Context, req Request) (Result, error) {
	if err := s.validateRequest(req); err != nil {
		return Result{}, err
	}

	preflight, err := s.preflightModules(ctx, req)
	if err != nil {
		return Result{}, err
	}

	if s.opts.DryRun {
		results := make([]ModuleResult, 0, len(preflight))
		for _, item := range preflight {
			results = append(results, ModuleResult{module: item.mod.Name(), skipped: item.skip})
		}
		return Result{modules: results}, nil
	}

	return s.publishTransaction(ctx, req, preflight)
}

// validateRequest rejects invalid publication inputs before Git is mutated.
func (s Service) validateRequest(req Request) error {
	if req.Plan.Empty() {
		return &Error{Code: CodeInvalidRequest, Message: "plan is empty"}
	}
	if req.Verify.Failed() {
		return &Error{
			Code:    CodeVerificationFailed,
			Message: "verification result contains failed checks",
		}
	}
	if s.deps.Git == nil {
		return &Error{Code: CodeInvalidRequest, Message: "git dependency is required"}
	}

	return nil
}

// preflightModules validates all modules before any Git mutation is attempted.
func (s Service) preflightModules(ctx context.Context, req Request) ([]modulePreflight, error) {
	preflight := make([]modulePreflight, 0, req.Plan.Len())
	for _, mod := range req.Plan.Modules() {
		item, err := s.preflightModule(ctx, req, mod)
		if err != nil {
			return nil, err
		}
		preflight = append(preflight, item)
	}

	return preflight, nil
}

func (s Service) publishTransaction(
	ctx context.Context,
	req Request,
	preflight []modulePreflight,
) (Result, error) {
	stateDir := deriveStateDir(s.opts.StateDir, preflight)
	if stateDir == "" {
		return Result{}, &Error{Code: CodeJournalFailed, Message: "transaction state dir is unavailable"}
	}
	store := NewFileJournalStore(stateDir)
	if pending, ok, err := store.HasPending(ctx); err != nil {
		return Result{}, &Error{Code: CodeJournalFailed, Message: "pending transaction lookup failed", Cause: err}
	} else if ok {
		return Result{}, &Error{
			Code:    CodePendingTransaction,
			Message: fmt.Sprintf("pending publish transaction %s has status %s; run arcpub transactions show %s or arcpub rollback --transaction %s", pending.ID, pending.Status, pending.ID, pending.ID),
		}
	}

	now := time.Now().UTC()
	idFunc := s.opts.TransactionIDFunc
	if idFunc == nil {
		idFunc = defaultTransactionID
	}
	journal := newTransactionJournal(idFunc(TransactionInput{
		Version: transactionVersion(req),
	}), req, s.opts.RemoteName, preflight, now)

	lock, err := acquireTransactionLock(ctx, stateDir, journal.ID, now)
	if err != nil {
		return Result{}, &Error{Code: CodeLockFailed, Message: "publish transaction lock failed", Cause: err}
	}
	defer func() { _ = lock.Release() }()

	if err := store.Create(ctx, journal); err != nil {
		return Result{}, &Error{Code: CodeJournalFailed, Message: "create transaction journal failed", Cause: err}
	}

	tx := transactionRunner{service: s, request: req, store: store, journal: journal}
	return tx.run(ctx, preflight)
}

type modulePreflight struct {
	mod          plan.ModulePlan
	worktree     string
	sourceModule source.ModuleSnapshot
	skip         bool
}

// preflightModule validates one module without mutating Git.
func (s Service) preflightModule(
	ctx context.Context,
	req Request,
	mod plan.ModulePlan,
) (modulePreflight, error) {
	name := mod.Name()
	worktree, ok := targetWorktree(req, name)
	if !ok {
		return modulePreflight{}, &Error{
			Code:    CodeInvalidRequest,
			Message: fmt.Sprintf("target workspace for %s is missing", name),
		}
	}
	if len(mod.Branches()) != 1 {
		return modulePreflight{}, &Error{
			Code:    CodePreflightFailed,
			Message: fmt.Sprintf("module %s has %d branch mappings; multi-branch publish is not supported yet", name, len(mod.Branches())),
		}
	}

	sourceModule, ok := sourceModuleForPublish(req.Source, name)
	if !ok {
		return modulePreflight{}, &Error{
			Code:    CodeMissingSourceSnapshot,
			Message: fmt.Sprintf("source snapshot for %s is missing", name),
		}
	}

	skip, err := s.shouldSkipModule(ctx, req, name, worktree)
	if err != nil {
		return modulePreflight{}, err
	}

	if err := s.preflightTag(ctx, req, mod, worktree); err != nil {
		return modulePreflight{}, err
	}

	return modulePreflight{
		mod:          mod,
		worktree:     worktree,
		sourceModule: sourceModule,
		skip:         skip,
	}, nil
}

// targetWorktree returns the prepared worktree path for name.
func targetWorktree(req Request, name manifest.ModuleName) (string, bool) {
	ws, ok := req.Targets.WorkspaceByModule(name)
	if !ok {
		return "", false
	}

	return ws.WorktreeDir(), true
}

// shouldSkipModule reports whether the final target worktree has no publishable
// changes. Git status is the source of truth when available because construct
// and modulefile results can include no-op writes or later cleanup.
func (s Service) shouldSkipModule(
	ctx context.Context,
	req Request,
	name manifest.ModuleName,
	worktree string,
) (bool, error) {
	status, err := s.deps.Git.Status(ctx, worktree)
	if err == nil {
		return !status.IsDirty(), nil
	}
	if !s.opts.AllowStatusFallback {
		return false, &Error{
			Code:    CodePreflightFailed,
			Message: fmt.Sprintf("module %s target status failed", name),
			Cause:   err,
		}
	}

	return !stageResultsChanged(req, name), nil
}

// stageResultsChanged is the documented fallback for tests and degraded Git
// ports that cannot return status after construction.
func stageResultsChanged(req Request, name manifest.ModuleName) bool {
	if result, ok := req.Construct.ModuleByName(name); ok && result.Changed() {
		return true
	}
	if result, ok := req.ModuleFile.ModuleByName(name); ok && result.Changed() {
		return true
	}

	return false
}

// commitWorktree stages and commits target worktree changes.
func (s Service) commitWorktree(
	ctx context.Context,
	worktree string,
	mod plan.ModulePlan,
	sourceModule source.ModuleSnapshot,
	req Request,
) (git.CommitHash, error) {
	if err := s.deps.Git.AddAll(ctx, worktree); err != nil {
		return "", &Error{Code: CodePublishFailed, Message: "git add failed", Cause: err}
	}

	commit, err := s.deps.Git.Commit(
		ctx,
		worktree,
		commitMessage(mod, sourceModule, req),
		git.CommitOptions{AllowEmpty: s.opts.AllowEmptyCommits},
	)
	if err != nil {
		return "", &Error{Code: CodePublishFailed, Message: "git commit failed", Cause: err}
	}

	return commit, nil
}

func sourceModuleForPublish(
	snapshot source.Snapshot,
	name manifest.ModuleName,
) (source.ModuleSnapshot, bool) {
	for _, sourceModule := range snapshot.Modules() {
		if sourceModule.Name() == name {
			return sourceModule, true
		}
	}

	return source.ModuleSnapshot{}, false
}

// preflightTag rejects local or remote release tag collisions before branch
// refs can be mutated.
func (s Service) preflightTag(
	ctx context.Context,
	req Request,
	mod plan.ModulePlan,
	worktree string,
) error {
	if !req.Plan.PublishPolicy().Tags().Enabled() {
		return nil
	}

	tag := git.TagName(mod.Version().String())
	localExists, err := s.deps.Git.TagExists(ctx, worktree, tag)
	if err != nil {
		return &Error{
			Code:    CodePreflightFailed,
			Message: fmt.Sprintf("module %s local tag lookup failed for %s", mod.Name(), tag),
			Cause:   err,
		}
	}
	if localExists {
		return &Error{
			Code:    CodePreflightFailed,
			Message: fmt.Sprintf("module %s local tag %s already exists", mod.Name(), tag),
		}
	}

	remoteRef := "refs/tags/" + tag.String()
	remoteExists, err := s.deps.Git.RemoteRefExists(ctx, worktree, s.opts.RemoteName, remoteRef)
	if err != nil {
		return &Error{
			Code:    CodePreflightFailed,
			Message: fmt.Sprintf("module %s remote tag lookup failed for %s", mod.Name(), tag),
			Cause:   err,
		}
	}
	if remoteExists {
		return &Error{
			Code:    CodePreflightFailed,
			Message: fmt.Sprintf("module %s remote tag %s already exists", mod.Name(), tag),
		}
	}

	return nil
}

// pushOptions maps the resolved publication push policy to Git push flags.
func pushOptions(req Request) git.PushOptions {
	opts := git.PushOptions{}
	if req.Plan.PublishPolicy().PushPolicy() == manifest.PushPolicyForceWithLease {
		opts.ForceWithLease = true
	}
	return opts
}
