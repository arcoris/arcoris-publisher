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

// Package git defines the infrastructure port for Git repository operations.
//
// This package describes Git capabilities only. It MUST NOT contain ARCORIS
// publishing operations such as module publication, staging synchronization, or
// Go module rewriting.
//
// Adapters may call the Git CLI, libgit2, go-git, or another implementation, but
// they should expose the same semantics through these contracts: context-aware
// operations, stable typed refs, redaction of sensitive values, and structured
// errors using this package's error codes.
package git

// WorktreeClient groups Git capabilities that operate on repository worktrees
// and remotes without tag publication.
type WorktreeClient interface {
	RepositoryReader
	RepositoryWriter
	RemoteClient
}

// Client groups all Git capabilities required by publisher workflows.
type Client interface {
	WorktreeClient
	TagClient
}
