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

package source

import "arcoris.dev/arcoris-publisher/internal/plan"

// Request describes the source checkout to inspect for a publication plan.
type Request struct {
	// Plan is the immutable executable publication plan produced by the plan
	// package. Source inspection consumes only plan-level effective values.
	Plan plan.Plan

	// RepositoryDir is the absolute or process-relative path to the source Git
	// checkout root.
	RepositoryDir string

	// StagingDir is the absolute or process-relative path to the staging root
	// that contains module SourceDir paths from the plan.
	StagingDir string
}
