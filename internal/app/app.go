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
	"arcoris.dev/arcoris-publisher/internal/config"
	"arcoris.dev/arcoris-publisher/internal/workflow"
)

// App coordinates publisher use cases.
type App struct {
	// loader loads and resolves manifest configuration.
	loader config.Loader

	// workflowDeps are ports passed into workflow stages.
	workflowDeps workflow.Dependencies

	// workflowOptions configure workflow stages.
	workflowOptions workflow.Options
}

// New creates an application service.
func New(deps Dependencies, opts Options) App {
	loader := config.NewLoader(config.LoaderOptions{})
	if deps.Loader != nil {
		loader = *deps.Loader
	}

	return App{
		loader:          loader,
		workflowDeps:    deps.Workflow,
		workflowOptions: opts.Workflow,
	}
}
