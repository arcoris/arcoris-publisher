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

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
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
	if opts.RemoteName == "" {
		opts.RemoteName = DefaultOptions().RemoteName
	}
	return Service{deps: deps, opts: opts}
}

// Publish commits, tags, and pushes every changed verified module.
func (s Service) Publish(ctx context.Context, req Request) (Result, error) {
	if req.Plan.Empty() {
		return Result{}, &Error{Code: CodeInvalidRequest, Message: "plan is empty"}
	}
	if req.Verify.Failed() {
		return Result{}, &Error{
			Code:    CodeVerificationFailed,
			Message: "verification result contains failed checks",
		}
	}
	if s.deps.Git == nil {
		return Result{}, &Error{Code: CodeInvalidRequest, Message: "git dependency is required"}
	}
	results := make([]ModuleResult, 0, req.Plan.Len())
	for _, mod := range req.Plan.Modules() {
		mr, err := s.publishModule(ctx, req, mod.Name())
		if err != nil {
			return Result{}, err
		}
		results = append(results, mr)
	}
	return Result{modules: results}, nil
}

// publishModule publishes one changed module target worktree.
func (s Service) publishModule(
	ctx context.Context,
	req Request,
	name manifest.ModuleName,
) (ModuleResult, error) {
	mod, _ := req.Plan.ModuleByName(name)
	ws, ok := req.Targets.WorkspaceByModule(name)
	if !ok {
		return ModuleResult{}, &Error{
			Code:    CodeInvalidRequest,
			Message: fmt.Sprintf("target workspace for %s is missing", name),
		}
	}
	constructResult, changed := req.Construct.ModuleByName(name)
	if !changed || !constructResult.Changed() {
		return ModuleResult{module: name, skipped: true}, nil
	}
	if s.opts.DryRun {
		return ModuleResult{module: name, skipped: false}, nil
	}
	if err := s.deps.Git.AddAll(ctx, ws.WorktreeDir()); err != nil {
		return ModuleResult{}, &Error{Code: CodePublishFailed, Message: "git add failed", Cause: err}
	}
	commit, err := s.deps.Git.Commit(
		ctx,
		ws.WorktreeDir(),
		commitMessage(mod, req.Source),
		git.CommitOptions{AllowEmpty: s.opts.AllowEmptyCommits},
	)
	if err != nil {
		return ModuleResult{}, &Error{Code: CodePublishFailed, Message: "git commit failed", Cause: err}
	}
	tags := []git.TagName{}
	if req.Plan.PublishPolicy().Tags().Enabled() {
		tag := git.TagName(mod.Version().String())
		if err := s.deps.Git.CreateTag(ctx, ws.WorktreeDir(), tag, commit, git.TagOptions{
			Annotated: true,
			Message:   "release " + mod.Version().String(),
		}); err != nil {
			return ModuleResult{}, &Error{Code: CodePublishFailed, Message: "git tag failed", Cause: err}
		}
		if err := s.deps.Git.PushTag(
			ctx,
			ws.WorktreeDir(),
			s.opts.RemoteName,
			tag,
			git.PushOptions{},
		); err != nil {
			return ModuleResult{}, &Error{
				Code:    CodePublishFailed,
				Message: "git push tag failed",
				Cause:   err,
			}
		}
		tags = append(tags, tag)
	}
	for _, branch := range mod.Branches() {
		refspec := git.RefSpec(
			"refs/heads/" + branch.Target().String() + ":refs/heads/" + branch.Target().String(),
		)
		opts := git.PushOptions{}
		if req.Plan.PublishPolicy().PushPolicy() == manifest.PushPolicyForceWithLease {
			opts.ForceWithLease = true
		}
		if err := s.deps.Git.Push(ctx, ws.WorktreeDir(), s.opts.RemoteName, refspec, opts); err != nil {
			return ModuleResult{}, &Error{
				Code:    CodePublishFailed,
				Message: "git push branch failed",
				Cause:   err,
			}
		}
	}
	return ModuleResult{module: name, commit: commit, tags: tags, pushed: true}, nil
}
