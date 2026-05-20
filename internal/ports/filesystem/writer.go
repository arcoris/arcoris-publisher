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

package filesystem

import (
	"context"
	"io/fs"
)

// Writer describes mutating filesystem operations.
type Writer interface {
	// WriteFile creates or replaces a file according to opts.
	//
	// Implementations should reject overwrites unless opts.Overwrite is true and
	// should create parent directories only when opts.CreateDirs is true.
	WriteFile(ctx context.Context, path string, data []byte, opts WriteFileOptions) error
	// MkdirAll creates path and any missing parents.
	//
	// Calling MkdirAll for an existing directory should be a successful no-op.
	// Calling it for an existing non-directory should return a structured error.
	MkdirAll(ctx context.Context, path string, opts MkdirOptions) error
	// RemoveAll recursively removes path according to opts.
	//
	// Adapters must enforce SafetyRoot before deleting anything. Missing paths
	// are successful only when opts.AllowMissing is true.
	RemoveAll(ctx context.Context, path string, opts RemoveOptions) error
}

// WriteFileOptions configures file creation or replacement.
type WriteFileOptions struct {
	// Perm is the file mode to use when creating the file.
	//
	// Zero lets the adapter choose a safe default such as 0o644.
	Perm fs.FileMode
	// Overwrite permits replacing an existing file.
	Overwrite bool
	// CreateDirs creates missing parent directories before writing.
	CreateDirs bool
}

// MkdirOptions configures recursive directory creation.
type MkdirOptions struct {
	// Perm is the mode to use for created directories.
	//
	// Zero lets the adapter choose a safe default such as 0o755.
	Perm fs.FileMode
}

// RemoveOptions configures recursive removal.
type RemoveOptions struct {
	// SafetyRoot confines deletion to a known parent directory.
	SafetyRoot string
	// AllowMissing treats a missing path as a successful no-op.
	AllowMissing bool
}
