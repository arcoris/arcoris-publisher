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

import "context"

// Reader describes read-only filesystem operations.
type Reader interface {
	// Exists reports whether path is present.
	//
	// A missing path should return (false, nil). Other stat failures, including
	// permission errors or context cancellation, should be returned as errors.
	Exists(ctx context.Context, path string) (bool, error)
	// IsDir reports whether path exists and is a directory.
	//
	// A missing path should return (false, nil). Adapters should not follow
	// symlinks unless their documented filesystem behavior explicitly does so.
	IsDir(ctx context.Context, path string) (bool, error)
	// ReadFile returns the full contents of a regular file.
	//
	// Implementations should reject directories and should return a detached byte
	// slice that callers can retain after the operation finishes.
	ReadFile(ctx context.Context, path string) ([]byte, error)
}
