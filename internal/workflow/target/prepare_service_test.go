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

package target

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
)

func TestPrepareTargetsClonesMissingWorktree(t *testing.T) {
	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	fs := porttest.NewFileSystem()
	git := porttest.NewGit()
	worktree := "/target/arcoris__foundation"
	git.RemoteRefs[porttest.RemoteRefKeyForRepo(worktree, "origin", "refs/heads/main")] = true
	template, err := manifest.ParseRemoteTemplate("file:///remotes/{name}.git")
	if err != nil {
		t.Fatal(err)
	}

	result, err := New(Dependencies{FS: fs, Git: git}, DefaultOptions()).PrepareTargets(
		context.Background(),
		PrepareRequest{Plan: p, RootDir: "/target", RemoteTemplate: template, HasRemoteTemplate: true},
	)
	if err != nil {
		t.Fatalf("PrepareTargets() error = %v", err)
	}
	if result.Failed() {
		t.Fatalf("PrepareTargets() failed: %#v", result)
	}
	if len(git.Calls) == 0 || git.Calls[0].Op != "clone" {
		t.Fatalf("git calls = %#v, want first clone", git.Calls)
	}
}

func TestPrepareTargetsRejectsMissingTemplate(t *testing.T) {
	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	result, err := New(Dependencies{FS: porttest.NewFileSystem(), Git: porttest.NewGit()}, DefaultOptions()).
		PrepareTargets(context.Background(), PrepareRequest{Plan: p, RootDir: "/target"})
	if err != nil {
		t.Fatalf("PrepareTargets() error = %v", err)
	}
	assertPrepareAction(t, result, "clone", "missing_remote_template")
}

func TestPrepareTargetsRejectsRemoteMismatch(t *testing.T) {
	p, fs, git, worktree := targetFixture(t)
	git.Refs[worktree+"\x00refs/heads/main"] = true
	git.RemoteURLs[worktree+"\x00origin"] = "file:///wrong.git"
	template, err := manifest.ParseRemoteTemplate("file:///remotes/{name}.git")
	if err != nil {
		t.Fatal(err)
	}

	result, err := New(Dependencies{FS: fs, Git: git}, DefaultOptions()).PrepareTargets(
		context.Background(),
		PrepareRequest{Plan: p, RootDir: "/target", RemoteTemplate: template, HasRemoteTemplate: true},
	)
	if err != nil {
		t.Fatalf("PrepareTargets() error = %v", err)
	}
	assertPrepareAction(t, result, "remote", "remote_mismatch")
}

func TestPrepareTargetsCreatesTrackingBranchWhenOnlyRemoteExists(t *testing.T) {
	p, fs, git, worktree := targetFixture(t)
	git.RemoteURLs[worktree+"\x00origin"] = "file:///foundation.git"
	git.RemoteRefs[porttest.RemoteRefKeyForRepo(worktree, "origin", "refs/heads/main")] = true

	result, err := New(Dependencies{FS: fs, Git: git}, DefaultOptions()).
		PrepareTargets(context.Background(), PrepareRequest{Plan: p, RootDir: "/target"})
	if err != nil {
		t.Fatalf("PrepareTargets() error = %v", err)
	}
	if result.Failed() {
		t.Fatalf("PrepareTargets() failed: %#v", result)
	}
	if len(git.Calls) < 2 || git.Calls[len(git.Calls)-2].Op != "create-branch" {
		t.Fatalf("git calls = %#v, want create-branch before checkout", git.Calls)
	}
}

func assertPrepareAction(t *testing.T, result PrepareResult, name string, code string) {
	t.Helper()
	for _, mod := range result.Modules() {
		for _, action := range mod.Actions() {
			if action.Name() == name && action.Code() == code {
				return
			}
		}
	}
	t.Fatalf("missing action %s code %s in %#v", name, code, result)
}
