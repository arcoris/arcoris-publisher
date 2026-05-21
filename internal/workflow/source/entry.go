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

// EntrySnapshot captures the source-side state of one explicit publish entry.
type EntrySnapshot struct {
	// entry is the original explicit publish entry from the plan.
	entry manifest.PublishEntry

	// sourcePath is the absolute inspected source path.
	sourcePath string

	// targetPath is the target repository relative path.
	targetPath manifest.RelativePath

	// present reports whether sourcePath existed during inspection.
	present bool

	// hash is the content hash for present entries when hashing is enabled.
	hash Hash
}

// Entry returns the original explicit publish entry from the plan.
func (s EntrySnapshot) Entry() manifest.PublishEntry { return s.entry }

// Kind returns the explicit publish entry kind.
func (s EntrySnapshot) Kind() manifest.PublishEntryKind { return s.entry.Kind() }

// SourcePath returns the absolute inspected source path.
func (s EntrySnapshot) SourcePath() string { return s.sourcePath }

// TargetPath returns the target repository relative path declared by the entry.
func (s EntrySnapshot) TargetPath() manifest.RelativePath { return s.targetPath }

// Optional reports whether the entry may be absent.
func (s EntrySnapshot) Optional() bool { return s.entry.Optional() }

// Present reports whether the entry source path exists.
func (s EntrySnapshot) Present() bool { return s.present }

// Hash returns the entry content hash when the entry is present and hash
// computation is enabled.
func (s EntrySnapshot) Hash() Hash { return s.hash }
