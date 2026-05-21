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

import portgit "arcoris.dev/arcoris-publisher/internal/ports/git"

// RepositorySnapshot captures source Git repository state used for provenance
// and dirty-check policy enforcement.
type RepositorySnapshot struct {
	// repositoryDir is the absolute inspected source repository root.
	repositoryDir string

	// stagingDir is the absolute inspected staging root.
	stagingDir string

	// head is the source repository HEAD commit.
	head portgit.CommitHash

	// branch is the current source branch, empty only for allowed detached HEAD.
	branch portgit.BranchName

	// status is the detached Git status snapshot.
	status portgit.Status
}

// RepositoryDir returns the inspected source repository root path.
func (s RepositorySnapshot) RepositoryDir() string { return s.repositoryDir }

// StagingDir returns the inspected staging root path.
func (s RepositorySnapshot) StagingDir() string { return s.stagingDir }

// Head returns the source HEAD commit.
func (s RepositorySnapshot) Head() portgit.CommitHash { return s.head }

// Branch returns the current source branch. It may be empty only when detached
// HEAD was explicitly allowed.
func (s RepositorySnapshot) Branch() portgit.BranchName { return s.branch }

// Status returns the raw Git status snapshot.
func (s RepositorySnapshot) Status() portgit.Status { return cloneStatus(s.status) }

// Dirty reports whether the source checkout had Git status entries or was not
// marked clean by the Git port.
func (s RepositorySnapshot) Dirty() bool { return !s.status.Clean || s.status.HasEntries() }

// cloneStatus detaches Git status entries before storing or returning them.
func cloneStatus(in portgit.Status) portgit.Status {
	out := portgit.Status{Clean: in.Clean}
	out.Entries = append([]portgit.StatusEntry(nil), in.Entries...)
	return out
}
