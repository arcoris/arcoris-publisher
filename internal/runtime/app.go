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

package runtime

import (
	"arcoris.dev/arcoris-publisher/internal/app"
	"arcoris.dev/arcoris-publisher/internal/workflow"
	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/preflight"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

// WorkflowDependencies builds the workflow port graph from runtime adapters.
func (r Runtime) WorkflowDependencies() workflow.Dependencies {
	return workflow.Dependencies{
		Source: source.Dependencies{
			Git: r.deps.Git,
			FS:  r.deps.FileSystem,
		},
		Target: target.Dependencies{
			Git: r.deps.Git,
			FS:  r.deps.FileSystem,
		},
		Construct: construct.Dependencies{
			FS: r.deps.FileSystem,
		},
		ModuleFile: modulefile.Dependencies{
			FS: r.deps.FileSystem,
		},
		Verify: verify.Dependencies{
			FS:  r.deps.FileSystem,
			Git: r.deps.Git,
			Go:  r.deps.Go,
		},
		Preflight: preflight.Dependencies{
			Source: source.Dependencies{
				Git: r.deps.Git,
				FS:  r.deps.FileSystem,
			},
			Git: r.deps.Git,
			FS:  r.deps.FileSystem,
		},
		Publish: publish.Dependencies{
			Git: r.deps.Git,
		},
	}
}

// AppDependencies builds the application dependency graph.
func (r Runtime) AppDependencies() app.Dependencies {
	return app.Dependencies{
		Loader:   r.deps.Loader,
		Workflow: r.WorkflowDependencies(),
	}
}

// App creates the high-level application service.
func (r Runtime) App() app.App {
	return app.New(r.AppDependencies(), r.opts.App)
}
