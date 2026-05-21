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

// Options configures source inspection behavior.
type Options struct {
	// AllowDetachedHEAD allows CurrentBranch to be empty. Detached source
	// checkouts have weaker provenance because the source branch is unknown, so
	// the safe default is to reject them.
	AllowDetachedHEAD bool

	// DisableHashes disables deterministic entry and module hash computation.
	// Hashes are enabled by default because they are useful for provenance,
	// no-op detection, and diagnostics.
	DisableHashes bool
}

// DefaultOptions returns conservative source inspection defaults.
func DefaultOptions() Options {
	return Options{}
}
