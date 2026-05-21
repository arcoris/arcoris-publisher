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

package manifest

import "fmt"

// PublishEntrySpec is the raw explicit file or directory publication entry.
type PublishEntrySpec struct {
	Type      string `json:"type" yaml:"type"`
	From      string `json:"from" yaml:"from"`
	To        string `json:"to" yaml:"to"`
	Optional  *bool  `json:"optional,omitempty" yaml:"optional,omitempty"`
	Recursive *bool  `json:"recursive,omitempty" yaml:"recursive,omitempty"`
}

// PublishEntry is a validated explicit file or directory publication entry.
type PublishEntry struct {
	kind      PublishEntryKind
	from      RelativePath
	to        RelativePath
	optional  bool
	recursive bool
}

// NewPublishEntry validates a raw explicit publication entry.
func NewPublishEntry(spec PublishEntrySpec) (PublishEntry, error) {
	kind, err := ParsePublishEntryKind(spec.Type)
	if err != nil {
		return PublishEntry{}, err
	}
	from, err := ParseRelativePath("entry.from", spec.From, false)
	if err != nil {
		return PublishEntry{}, err
	}
	to, err := ParseRelativePath("entry.to", spec.To, true)
	if err != nil {
		return PublishEntry{}, err
	}
	recursiveDefault := kind == PublishEntryDirectory
	entry := PublishEntry{kind: kind, from: from, to: to, optional: boolValue(spec.Optional, false), recursive: boolValue(spec.Recursive, recursiveDefault)}
	if kind == PublishEntryFile && entry.recursive {
		return PublishEntry{}, fmt.Errorf("file entry must not be recursive")
	}
	return entry, nil
}

// Kind returns the entry kind.
func (e PublishEntry) Kind() PublishEntryKind { return e.kind }

// From returns the source path relative to the module root.
func (e PublishEntry) From() RelativePath { return e.from }

// To returns the target path relative to the target repository root.
func (e PublishEntry) To() RelativePath { return e.to }

// Optional reports whether a missing source entry should be tolerated.
func (e PublishEntry) Optional() bool { return e.optional }

// Recursive reports whether a directory entry should be copied recursively.
func (e PublishEntry) Recursive() bool { return e.recursive }

// Spec returns a serializable publish entry representation.
func (e PublishEntry) Spec() PublishEntrySpec {
	optional := e.optional
	recursive := e.recursive
	return PublishEntrySpec{Type: string(e.kind), From: string(e.from), To: string(e.to), Optional: &optional, Recursive: &recursive}
}
