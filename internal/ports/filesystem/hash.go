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

// TreeHash is a deterministic content hash of a filesystem tree.
type TreeHash string

// String returns the stable hash text exactly as supplied by an adapter.
func (h TreeHash) String() string {
	return string(h)
}

// Hasher describes deterministic tree hashing operations.
type Hasher interface {
	TreeHash(ctx context.Context, root string, opts TreeHashOptions) (TreeHash, error)
}

// TreeHashOptions configures deterministic tree hashing.
type TreeHashOptions struct {
	// Include limits hashing to matching relative paths when non-empty.
	Include []string
	// Exclude removes matching relative paths after include filtering.
	Exclude []string
	// IncludeFileMode includes executable and permission bits in the digest.
	IncludeFileMode bool
	// IncludeSymlinks includes symbolic-link metadata in the digest.
	IncludeSymlinks bool
	// SymlinkPolicy controls whether links are rejected, recorded, or followed.
	SymlinkPolicy SymlinkPolicy
}
