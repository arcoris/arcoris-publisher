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

import "strconv"

// String returns the canonical version string.
func (v Version) String() string { return string(v) }

// Major returns the semantic major component.
func (v Version) Major() int { return parseSemanticPart(v, majorSubmatchIndex) }

// Minor returns the semantic minor component.
func (v Version) Minor() int { return parseSemanticPart(v, minorSubmatchIndex) }

// Patch returns the semantic patch component.
func (v Version) Patch() int { return parseSemanticPart(v, patchSubmatchIndex) }

// parseSemanticPart extracts one numeric semantic version component.
func parseSemanticPart(version Version, index int) int {
	match := semanticVersionRE.FindStringSubmatch(string(version))
	if len(match) == 0 || index < majorSubmatchIndex || index > patchSubmatchIndex {
		return 0
	}
	value, _ := strconv.Atoi(match[index])
	return value
}
