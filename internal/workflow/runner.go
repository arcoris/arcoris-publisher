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

package workflow

import (
	"context"
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

// Runner executes workflow stages in dependency order.
type Runner struct {
	// deps contains stage infrastructure ports.
	deps Dependencies

	// opts contains stage-specific options.
	opts Options
}

// New returns a workflow runner.
func New(deps Dependencies, opts Options) Runner {
	if opts.DryRun {
		opts.Publish.DryRun = true
	}
	return Runner{deps: deps, opts: opts}
}

// Run executes source, target, construct, modulefile, verify, and optional
// publish stages.
func (r Runner) Run(ctx context.Context, req Request) (Result, error) {
	var result Result

	snapshot, err := r.inspectSource(ctx, req)
	if err != nil {
		return result, err
	}
	result.source = snapshot

	workspaces, err := r.prepareTargets(ctx, req)
	if err != nil {
		return result, err
	}
	result.target = workspaces

	constructed, err := r.constructTargets(ctx, req, snapshot, workspaces)
	if err != nil {
		return result, err
	}
	result.construct = constructed

	rewritten, err := r.rewriteModuleFiles(ctx, req, workspaces)
	if err != nil {
		return result, err
	}
	result.moduleFile = rewritten

	verified, err := r.verifyTargets(ctx, req, workspaces, constructed, rewritten)
	if err != nil {
		return result, err
	}
	result.verify = verified
	if verified.Failed() || !req.Publish || r.opts.DryRun {
		return result, nil
	}

	published, err := r.publishTargets(ctx, req, snapshot, workspaces, constructed, rewritten, verified)
	if err != nil {
		return result, err
	}
	result.publish = published

	return result, nil
}

func (r Runner) inspectSource(ctx context.Context, req Request) (source.Snapshot, error) {
	result, err := source.New(r.deps.Source, r.opts.Source).Inspect(ctx, source.Request{
		Plan:          req.Plan,
		RepositoryDir: req.SourceRepositoryDir,
		StagingDir:    req.StagingDir,
	})
	if err != nil {
		return source.Snapshot{}, fmt.Errorf("source: %w", err)
	}

	return result, nil
}

func (r Runner) prepareTargets(ctx context.Context, req Request) (target.WorkspaceSet, error) {
	result, err := target.New(r.deps.Target, r.opts.Target).Prepare(ctx, target.Request{
		Plan:    req.Plan,
		RootDir: req.TargetRootDir,
	})
	if err != nil {
		return target.WorkspaceSet{}, fmt.Errorf("target: %w", err)
	}

	return result, nil
}

func (r Runner) constructTargets(
	ctx context.Context,
	req Request,
	snapshot source.Snapshot,
	workspaces target.WorkspaceSet,
) (construct.Result, error) {
	result, err := construct.New(r.deps.Construct, r.opts.Construct).Construct(ctx, construct.Request{
		Plan:    req.Plan,
		Source:  snapshot,
		Targets: workspaces,
	})
	if err != nil {
		return construct.Result{}, fmt.Errorf("construct: %w", err)
	}

	return result, nil
}

func (r Runner) rewriteModuleFiles(
	ctx context.Context,
	req Request,
	workspaces target.WorkspaceSet,
) (modulefile.Result, error) {
	result, err := modulefile.New(r.deps.ModuleFile, r.opts.ModuleFile).Rewrite(ctx, modulefile.Request{
		Plan:    req.Plan,
		Targets: workspaces,
	})
	if err != nil {
		return modulefile.Result{}, fmt.Errorf("modulefile: %w", err)
	}

	return result, nil
}

func (r Runner) verifyTargets(
	ctx context.Context,
	req Request,
	workspaces target.WorkspaceSet,
	constructed construct.Result,
	rewritten modulefile.Result,
) (verify.Result, error) {
	result, err := verify.New(r.deps.Verify, r.opts.Verify).Verify(ctx, verify.Request{
		Plan:       req.Plan,
		Targets:    workspaces,
		Construct:  constructed,
		ModuleFile: rewritten,
	})
	if err != nil {
		return verify.Result{}, fmt.Errorf("verify: %w", err)
	}

	return result, nil
}

func (r Runner) publishTargets(
	ctx context.Context,
	req Request,
	snapshot source.Snapshot,
	workspaces target.WorkspaceSet,
	constructed construct.Result,
	rewritten modulefile.Result,
	verified verify.Result,
) (publish.Result, error) {
	result, err := publish.New(r.deps.Publish, r.opts.Publish).Publish(ctx, publish.Request{
		Plan:       req.Plan,
		Source:     snapshot,
		Targets:    workspaces,
		Construct:  constructed,
		ModuleFile: rewritten,
		Verify:     verified,
	})
	if err != nil {
		return publish.Result{}, fmt.Errorf("publish: %w", err)
	}

	return result, nil
}
