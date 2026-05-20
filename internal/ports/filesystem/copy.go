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

// Copier describes filesystem tree copy operations.
type Copier interface {
	// CopyTree copies entries from src to dst according to opts.
	//
	// The operation should be deterministic for a stable source tree and option
	// set. Adapters must enforce SafetyRoot for destination writes and apply
	// SymlinkPolicy consistently to every link encountered in the source tree.
	CopyTree(ctx context.Context, src string, dst string, opts CopyTreeOptions) (CopyTreeResult, error)
}

// CopyTreeOptions configures deterministic tree copying.
type CopyTreeOptions struct {
	// Include limits the copy to matching relative paths when non-empty.
	//
	// Include is evaluated before Exclude. Empty Include means all source entries
	// are candidates for copying.
	Include []string
	// Exclude skips matching relative paths after include filtering.
	//
	// Excluded files should be counted in CopyTreeResult.FilesSkipped when the
	// adapter can determine that they otherwise would have been copied.
	Exclude []string
	// PreserveMode keeps source file modes instead of adapter defaults.
	PreserveMode bool
	// PreserveTimes keeps source timestamps when the adapter supports it.
	PreserveTimes bool
	// SymlinkPolicy controls whether links are rejected, copied, or followed.
	SymlinkPolicy SymlinkPolicy
	// SafetyRoot confines writes to a known parent directory.
	//
	// Adapters must reject dst paths outside this root before creating files or
	// directories. Empty SafetyRoot delegates to the adapter's default policy.
	SafetyRoot string
}
