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

// PublishEntryKind identifies the kind of explicit publication entry.
type PublishEntryKind string

const (
	// PublishEntryFile publishes exactly one file.
	PublishEntryFile PublishEntryKind = "file"
	// PublishEntryDirectory publishes one explicitly declared directory.
	PublishEntryDirectory PublishEntryKind = "directory"
)

// ParsePublishEntryKind validates an explicit publication entry kind.
func ParsePublishEntryKind(value string) (PublishEntryKind, error) {
	switch PublishEntryKind(value) {
	case PublishEntryFile, PublishEntryDirectory:
		return PublishEntryKind(value), nil
	default:
		return "", fmt.Errorf("unsupported publish entry type %q", value)
	}
}
