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

// Package filesystem implements the filesystem port using the local
// operating-system filesystem.
//
// The adapter focuses on deterministic tree operations and defensive path
// handling. It intentionally keeps publisher policy out of filesystem code:
// callers choose include/exclude patterns, safety roots, symlink policy, and
// preservation rules through the port option structs.
package filesystem

import fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"

// FileSystem implements safe local filesystem operations.
//
// FileSystem is stateless and safe to reuse. All operation-specific behavior is
// supplied through method arguments and option structs.
type FileSystem struct{}

// New creates a local filesystem adapter.
func New() *FileSystem { return &FileSystem{} }

var _ fsport.FileSystem = (*FileSystem)(nil)
