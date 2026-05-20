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
	WriteFile(ctx context.Context, path string, data []byte, opts WriteFileOptions) error
	MkdirAll(ctx context.Context, path string, opts MkdirOptions) error
	RemoveAll(ctx context.Context, path string, opts RemoveOptions) error
}

// WriteFileOptions configures file creation or replacement.
type WriteFileOptions struct {
	// Perm is the file mode to use when creating the file.
	Perm fs.FileMode
	// Overwrite permits replacing an existing file.
	Overwrite bool
	// CreateDirs creates missing parent directories before writing.
	CreateDirs bool
}

// MkdirOptions configures recursive directory creation.
type MkdirOptions struct {
	// Perm is the mode to use for created directories.
	Perm fs.FileMode
}

// RemoveOptions configures recursive removal.
type RemoveOptions struct {
	// SafetyRoot confines deletion to a known parent directory.
	SafetyRoot string
	// AllowMissing treats a missing path as a successful no-op.
	AllowMissing bool
}
