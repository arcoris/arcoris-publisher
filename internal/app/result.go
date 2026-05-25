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
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/workflow"
	"arcoris.dev/arcoris-publisher/internal/workflow/preflight"
)

// Result describes a high-level workflow use case result.
type Result struct {
	// plan is the executable plan used for the workflow run.
	plan plan.Plan

	// workflow contains completed stage results.
	workflow workflow.Result

	// preflight contains read-only publish readiness checks.
	preflight preflight.Result
}

// Plan returns the executable plan used for the workflow run.
func (r Result) Plan() plan.Plan { return r.plan }

// Workflow returns completed workflow stage results.
func (r Result) Workflow() workflow.Result { return r.workflow }

// Preflight returns read-only publish readiness checks.
func (r Result) Preflight() preflight.Result { return r.preflight }
