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

	"arcoris.dev/arcoris-publisher/internal/workflow/preflight"
)

// Preflight builds a plan and checks whether publish can safely start without
// constructing target trees or writing transaction state.
func (a App) Preflight(ctx context.Context, req Request) (Result, error) {
	plan, err := a.BuildPlan(ctx, req.ManifestPath, req.Version)
	if err != nil {
		return Result{}, err
	}

	result, err := preflight.New(a.workflowDeps.Preflight, a.workflowOptions.Preflight).Check(ctx, preflight.Request{
		Plan:                plan,
		SourceRepositoryDir: req.SourceRepositoryDir,
		StagingDir:          req.StagingDir,
		TargetRootDir:       req.TargetRootDir,
	})
	if err != nil {
		return Result{plan: plan, preflight: result}, fmt.Errorf("preflight: %w", err)
	}

	return Result{plan: plan, preflight: result}, nil
}
