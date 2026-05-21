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

package source

import (
	"testing"

	portgit "arcoris.dev/arcoris-publisher/internal/ports/git"
)

func TestRepositorySnapshotAccessors(t *testing.T) {
	snap := RepositorySnapshot{
		repositoryDir: "/repo",
		stagingDir:    "/repo/staging",
		head:          portgit.CommitHash("abcdef"),
		branch:        portgit.BranchName("main"),
		status:        portgit.Status{Clean: true},
	}

	if snap.RepositoryDir() != "/repo" {
		t.Fatalf("RepositoryDir() = %q", snap.RepositoryDir())
	}
	if snap.StagingDir() != "/repo/staging" {
		t.Fatalf("StagingDir() = %q", snap.StagingDir())
	}
	if snap.Head() != "abcdef" {
		t.Fatalf("Head() = %q", snap.Head())
	}
	if snap.Branch() != "main" {
		t.Fatalf("Branch() = %q", snap.Branch())
	}
	if snap.Dirty() {
		t.Fatal("Dirty() = true")
	}
}

func TestRepositorySnapshotStatusIsDetached(t *testing.T) {
	snap := RepositorySnapshot{status: portgit.Status{
		Clean: true,
		Entries: []portgit.StatusEntry{{
			Path: "go.mod",
			Code: " M",
		}},
	}}

	status := snap.Status()
	status.Entries[0].Path = "go.sum"

	if snap.Status().Entries[0].Path != "go.mod" {
		t.Fatal("Status() returned attached entries")
	}
}

func TestRepositorySnapshotDirty(t *testing.T) {
	notClean := RepositorySnapshot{status: portgit.Status{Clean: false}}
	if !notClean.Dirty() {
		t.Fatal("Dirty() = false for not-clean status")
	}

	withEntries := RepositorySnapshot{status: portgit.Status{
		Clean: true,
		Entries: []portgit.StatusEntry{{
			Path: "go.mod",
			Code: " M",
		}},
	}}
	if !withEntries.Dirty() {
		t.Fatal("Dirty() = false for status entries")
	}
}
