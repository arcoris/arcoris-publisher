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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/config"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
)

func TestWorkflowDependenciesUseProvidedPorts(t *testing.T) {
	t.Parallel()

	fs := porttest.NewFileSystem()
	git := porttest.NewGit()
	goToolchain := &porttest.GoToolchain{}
	rt := NewWithDependencies(Dependencies{
		FileSystem: fs,
		Git:        git,
		Go:         goToolchain,
	}, Options{})

	deps := rt.WorkflowDependencies()
	if deps.Source.Git != git || deps.Target.Git != git || deps.Verify.Git != git || deps.Publish.Git != git {
		t.Fatalf("workflow git ports were not shared: %+v", deps)
	}
	if deps.Source.FS != fs || deps.Target.FS != fs || deps.Construct.FS != fs || deps.ModuleFile.FS != fs || deps.Verify.FS != fs {
		t.Fatalf("workflow filesystem ports were not shared: %+v", deps)
	}
	if deps.Verify.Go != goToolchain {
		t.Fatalf("workflow go toolchain was not wired")
	}
}

func TestAppDependenciesUseProvidedLoader(t *testing.T) {
	t.Parallel()

	loader := config.NewLoader(config.LoaderOptions{})
	rt := NewWithDependencies(Dependencies{Loader: &loader}, Options{})
	deps := rt.AppDependencies()
	if deps.Loader != &loader {
		t.Fatal("AppDependencies replaced provided loader")
	}
	if deps.Workflow.Source.Git == nil || deps.Workflow.Publish.Git == nil {
		t.Fatalf("app workflow dependencies are incomplete: %+v", deps.Workflow)
	}
}

func TestAppFactoryBuildsApplicationsWithCommandOptions(t *testing.T) {
	t.Parallel()

	rt := NewWithDependencies(Dependencies{}, Options{})
	app, err := rt.AppFactory()(rt.CLIOptions().App)
	if err != nil {
		t.Fatalf("AppFactory() error = %v", err)
	}
	if app == nil {
		t.Fatal("AppFactory() returned nil application")
	}
}
