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

package module

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// PublishSpec is the raw explicit module publication declaration.
type PublishSpec struct {
	Entries []manifest.PublishEntrySpec `json:"entries" yaml:"entries"`
}

// Publish is the validated explicit module publication declaration.
type Publish struct {
	entries []manifest.PublishEntry
}

// NewPublish validates spec and returns Publish.
func NewPublish(spec PublishSpec) (Publish, error) {
	var collector manifest.IssueCollector
	if len(spec.Entries) == 0 {
		collector.Add(manifest.IssueMissingField, "entries", "at least one explicit publish entry is required")
	}
	entries := make([]manifest.PublishEntry, 0, len(spec.Entries))
	targets := make(map[manifest.RelativePath]int, len(spec.Entries))
	for i, entrySpec := range spec.Entries {
		entry, err := manifest.NewPublishEntry(entrySpec)
		if err != nil {
			collector.AddError(fmt.Sprintf("entries[%d]", i), err)
			continue
		}
		if prev, exists := targets[entry.To()]; exists {
			collector.Add(manifest.IssueDuplicateValue, fmt.Sprintf("entries[%d].to", i), "duplicate target path %q previously declared at entries[%d]", entry.To(), prev)
			continue
		}
		targets[entry.To()] = i
		entries = append(entries, entry)
	}
	if err := collector.Err(); err != nil {
		return Publish{}, err
	}
	return Publish{entries: entries}, nil
}

// Entries returns detached explicit publication entries.
func (p Publish) Entries() []manifest.PublishEntry { return manifest.ClonePublishEntries(p.entries) }
