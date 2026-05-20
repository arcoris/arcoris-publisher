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

// Version is the manifest schema version.
type Version string

const (
	// VersionV1 is the first supported manifest schema version.
	VersionV1 Version = "v1"
)

// ParseVersion validates the manifest schema version.
func ParseVersion(value string) (Version, error) {
	if value == "" {
		return "", fmt.Errorf("version is required")
	}
	switch Version(value) {
	case VersionV1:
		return VersionV1, nil
	default:
		return "", fmt.Errorf("unsupported version %q", value)
	}
}
