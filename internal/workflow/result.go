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
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

// Result contains completed stage results.
type Result struct {
	// source is present after source inspection succeeds.
	source source.Snapshot

	// target is present after target preparation succeeds.
	target target.WorkspaceSet

	// construct is present after explicit projection construction succeeds.
	construct construct.Result

	// moduleFile is present after go.mod rewriting succeeds.
	moduleFile modulefile.Result

	// verify is present after verification executes.
	verify verify.Result

	// publish is present only when publication was requested and verification
	// passed.
	publish publish.Result
}

// Source returns the source inspection snapshot.
func (r Result) Source() source.Snapshot { return r.source }

// Target returns prepared target workspaces.
func (r Result) Target() target.WorkspaceSet { return r.target }

// Construct returns explicit projection construction results.
func (r Result) Construct() construct.Result { return r.construct }

// ModuleFile returns go.mod rewrite results.
func (r Result) ModuleFile() modulefile.Result { return r.moduleFile }

// Verify returns verification results.
func (r Result) Verify() verify.Result { return r.verify }

// Publish returns publication results.
func (r Result) Publish() publish.Result { return r.publish }
