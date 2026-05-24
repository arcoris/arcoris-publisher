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

	"arcoris.dev/arcoris-publisher/internal/app"
	"arcoris.dev/arcoris-publisher/internal/cli"
	"arcoris.dev/arcoris-publisher/internal/report"
	"arcoris.dev/arcoris-publisher/internal/workflow"
)

func TestDefaultOptions(t *testing.T) {
	t.Parallel()

	opts := DefaultOptions()
	if opts.CLI.ManifestPath == "" || opts.CLI.TargetRootDir == "" {
		t.Fatalf("DefaultOptions().CLI is not initialized: %+v", opts.CLI)
	}
}

func TestNormalizeOptionsRestoresCLIDefaults(t *testing.T) {
	t.Parallel()

	opts := normalizeOptions(Options{})
	if opts.CLI.ManifestPath == "" || opts.CLI.SourceRepositoryDir == "" || opts.CLI.TargetRootDir == "" {
		t.Fatalf("normalizeOptions() did not restore CLI defaults: %+v", opts.CLI)
	}
}

func TestNormalizeOptionsCopiesEnvironment(t *testing.T) {
	t.Parallel()

	env := []string{"A=B"}
	opts := normalizeOptions(Options{Env: env})
	env[0] = "A=mutated"
	if opts.Env[0] != "A=B" {
		t.Fatalf("normalizeOptions().Env = %v", opts.Env)
	}
}

func TestNormalizeOptionsPreservesExplicitCLI(t *testing.T) {
	t.Parallel()

	opts := normalizeOptions(Options{CLI: cli.Options{ManifestPath: "custom.yaml"}})
	if opts.CLI.ManifestPath != "custom.yaml" {
		t.Fatalf("normalizeOptions() ManifestPath = %q", opts.CLI.ManifestPath)
	}
}

func TestNormalizeOptionsPreservesReportPathPolicy(t *testing.T) {
	t.Parallel()

	opts := normalizeOptions(Options{
		CLI: cli.Options{
			Report: report.Options{IncludeLocalPaths: true},
		},
	})
	if !opts.CLI.Report.IncludeLocalPaths {
		t.Fatal("normalizeOptions() dropped IncludeLocalPaths")
	}
	if opts.CLI.Report.Format == "" {
		t.Fatal("normalizeOptions() did not restore report format")
	}
}

func TestCLIOptionsApplyRuntimeAppOptions(t *testing.T) {
	t.Parallel()

	rt := New(Options{App: app.Options{Workflow: workflow.Options{DryRun: true}}})
	if !rt.CLIOptions().App.Workflow.DryRun {
		t.Fatalf("CLIOptions() did not apply runtime app options")
	}
}
