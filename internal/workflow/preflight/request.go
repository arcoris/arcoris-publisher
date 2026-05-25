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

import "arcoris.dev/arcoris-publisher/internal/plan"

// Request describes the publish attempt whose start safety is checked.
type Request struct {
	// Plan is the executable publication plan.
	Plan plan.Plan

	// SourceRepositoryDir is the local source Git checkout root.
	SourceRepositoryDir string

	// StagingDir is the staging root containing module sources.
	StagingDir string

	// TargetRootDir contains existing target worktrees.
	TargetRootDir string
}
