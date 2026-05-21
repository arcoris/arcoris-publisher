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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// ModuleSnapshot captures source-side state for one planned module.
type ModuleSnapshot struct {
	// name is the planned module name.
	name manifest.ModuleName

	// sourceDir is the absolute source directory for the module.
	sourceDir string

	// moduleRootDir is the absolute root used to resolve publish entry sources.
	moduleRootDir string

	// entries contains inspected explicit publish entries.
	entries []EntrySnapshot

	// hash is the combined hash of present entry hashes.
	hash Hash
}

// Name returns the module name.
func (s ModuleSnapshot) Name() manifest.ModuleName { return s.name }

// SourceDir returns the absolute source directory for the module.
func (s ModuleSnapshot) SourceDir() string { return s.sourceDir }

// ModuleRootDir returns the absolute module root used to resolve publish entry
// source paths.
func (s ModuleSnapshot) ModuleRootDir() string { return s.moduleRootDir }

// Entries returns detached explicit publish entry snapshots.
func (s ModuleSnapshot) Entries() []EntrySnapshot { return cloneEntrySnapshots(s.entries) }

// Hash returns the combined hash of present explicit publish entries when hash
// computation is enabled.
func (s ModuleSnapshot) Hash() Hash { return s.hash }

// cloneEntrySnapshots detaches entry snapshot slices before returning them.
func cloneEntrySnapshots(in []EntrySnapshot) []EntrySnapshot {
	out := make([]EntrySnapshot, len(in))
	copy(out, in)
	return out
}
