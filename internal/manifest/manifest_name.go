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

// ManifestName is a stable manifest resource name.
type ManifestName string

// ParseManifestName validates a manifest resource name.
func ParseManifestName(value string) (ManifestName, error) {
	if err := validateSimpleName("metadata.name", value); err != nil {
		return "", err
	}
	return ManifestName(value), nil
}
