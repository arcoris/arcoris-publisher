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

// SourceDir is a normalized relative path to a staged module root.
//
// It deliberately wraps RelativePath because source directories have a tighter
// semantic contract than generic manifest paths: they point at module roots and
// may not be ".". The separate type makes that rule visible at call sites.
type SourceDir string

// ParseSourceDir validates a staged module source directory path.
func ParseSourceDir(value string) (SourceDir, error) {
	p, err := ParseRelativePath("sourceDir", value, false)
	if err != nil {
		return "", err
	}
	return SourceDir(p), nil
}

// String returns the normalized source directory path string.
func (d SourceDir) String() string { return string(d) }

// ToRelativePath converts the source directory into a generic relative path.
func (d SourceDir) ToRelativePath() RelativePath { return RelativePath(d) }
