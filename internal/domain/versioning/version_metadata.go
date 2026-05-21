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

package versioning

import "strings"

// Prerelease returns the optional prerelease suffix without the leading hyphen.
func (v Version) Prerelease() string {
	match := semanticVersionRE.FindStringSubmatch(string(v))
	if len(match) == 0 {
		return ""
	}
	return match[prereleaseSubmatchIndex]
}

// BuildMetadata returns the optional build metadata suffix without the leading plus.
func (v Version) BuildMetadata() string {
	match := semanticVersionRE.FindStringSubmatch(string(v))
	if len(match) == 0 {
		return ""
	}
	return match[buildMetadataSubmatchIndex]
}

// WithoutBuildMetadata returns the same version without build metadata.
//
// Build metadata is useful for humans, but Go pseudo-version bases and module
// requirements must compare without build metadata influencing the generated
// version string.
func (v Version) WithoutBuildMetadata() Version {
	value := string(v)
	if idx := strings.IndexByte(value, '+'); idx >= 0 {
		return Version(value[:idx])
	}
	return v
}
