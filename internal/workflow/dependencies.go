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
	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/preflight"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

// Dependencies contains ports for every workflow stage.
type Dependencies struct {
	// Source wires source inspection.
	Source source.Dependencies

	// Target wires target preparation.
	Target target.Dependencies

	// Construct wires explicit projection construction.
	Construct construct.Dependencies

	// ModuleFile wires go.mod rewriting.
	ModuleFile modulefile.Dependencies

	// Verify wires target verification.
	Verify verify.Dependencies

	// Preflight wires read-only publish readiness checks.
	Preflight preflight.Dependencies

	// Publish wires optional Git publication.
	Publish publish.Dependencies
}
