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

// Options configures each workflow stage without hiding stage-specific knobs.
type Options struct {
	// Source configures source inspection.
	Source source.Options

	// Target configures target preparation.
	Target target.Options

	// Construct configures explicit target construction.
	Construct construct.Options

	// ModuleFile configures go.mod rewriting.
	ModuleFile modulefile.Options

	// Verify configures target verification.
	Verify verify.Options

	// Publish configures optional publication.
	Publish publish.Options

	// DryRun disables publication mutations through the publish stage options.
	DryRun bool
}
