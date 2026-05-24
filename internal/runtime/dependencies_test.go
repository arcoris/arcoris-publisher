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
	"context"
	"testing"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

func TestNewDependenciesBuildsProductionPorts(t *testing.T) {
	t.Parallel()

	deps := NewDependencies(Options{})
	if deps.Process == nil {
		t.Fatal("Process dependency is nil")
	}
	if deps.FileSystem == nil {
		t.Fatal("FileSystem dependency is nil")
	}
	if deps.Git == nil {
		t.Fatal("Git dependency is nil")
	}
	if deps.Go == nil {
		t.Fatal("Go dependency is nil")
	}
	if deps.Loader == nil {
		t.Fatal("Loader dependency is nil")
	}
	if deps.BuildInfo == nil {
		t.Fatal("BuildInfo dependency is nil")
	}
}

func TestNormalizeDependenciesUsesInjectedProcessForDefaultAdapters(t *testing.T) {
	t.Parallel()

	runner := &recordingRunner{}
	rt := NewWithDependencies(
		Dependencies{Process: runner},
		Options{
			Env:       []string{"ARCORIS_TEST=1"},
			GitBinary: "test-git",
			GoBinary:  "test-go",
		},
	)

	if _, err := rt.Dependencies().Git.Head(context.Background(), "/repo"); err != nil {
		t.Fatalf("Git.Head() error = %v", err)
	}
	if _, err := rt.Dependencies().Go.Env(context.Background(), goport.EnvOptions{}); err != nil {
		t.Fatalf("Go.Env() error = %v", err)
	}

	if len(runner.specs) != 2 {
		t.Fatalf("recorded process specs = %d", len(runner.specs))
	}
	if runner.specs[0].Name != "test-git" {
		t.Fatalf("git spec Name = %q", runner.specs[0].Name)
	}
	if runner.specs[1].Name != "test-go" {
		t.Fatalf("go spec Name = %q", runner.specs[1].Name)
	}
	if !containsString(runner.specs[0].Env, "ARCORIS_TEST=1") || !containsString(runner.specs[1].Env, "ARCORIS_TEST=1") {
		t.Fatalf("process env was not propagated: %+v", runner.specs)
	}
}
