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

package app

import (
	"context"
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

// PrepareTargets builds a plan and prepares local target Git worktrees. It may
// clone, fetch, add remotes, and checkout branches; it never publishes refs or
// writes transaction state.
func (a App) PrepareTargets(ctx context.Context, req Request) (Result, error) {
	plan, err := a.BuildPlan(ctx, req.ManifestPath, req.Version)
	if err != nil {
		return Result{}, err
	}

	template, hasTemplate, err := targetRemoteTemplate(plan.TargetPolicy(), req.TargetRemoteTemplate)
	if err != nil {
		return Result{plan: plan}, err
	}

	opts := prepareTargetOptions(a.workflowOptions.Target)
	result, err := target.New(a.workflowDeps.Target, opts).PrepareTargets(ctx, target.PrepareRequest{
		Plan:              plan,
		RootDir:           req.TargetRootDir,
		RemoteTemplate:    template,
		HasRemoteTemplate: hasTemplate,
	})
	if err != nil {
		return Result{plan: plan, targetPrepare: result}, fmt.Errorf("target prepare: %w", err)
	}

	return Result{plan: plan, targetPrepare: result}, nil
}

func targetRemoteTemplate(policy manifest.TargetPolicy, override string) (manifest.RemoteTemplate, bool, error) {
	if override != "" {
		template, err := manifest.ParseRemoteTemplate(override)
		return template, err == nil, err
	}
	template, ok := policy.RemoteTemplate()
	return template, ok, nil
}

func prepareTargetOptions(opts target.Options) target.Options {
	defaults := target.DefaultOptions()
	if opts.RemoteName == "" {
		opts.RemoteName = defaults.RemoteName
	}
	opts.CheckoutBranch = true
	opts.CreateMissing = false
	opts.Fetch = true
	opts.FetchRequired = true
	opts.RequireClean = true
	return opts
}
