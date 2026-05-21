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

package construct

// Options configures target tree construction.
type Options struct {
	// PreserveGitDir keeps .git metadata when cleaning a target worktree.
	PreserveGitDir bool

	// GenerateProvenanceFile writes configured provenance metadata into target trees.
	GenerateProvenanceFile bool
}

// DefaultOptions returns conservative construction defaults.
func DefaultOptions() Options {
	return Options{PreserveGitDir: true}
}
