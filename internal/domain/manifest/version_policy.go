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

// VersionPolicy describes how publish versions are assigned to modules.
type VersionPolicy string

const (
	// VersionPolicyReleaseTrain assigns one release version to all published modules.
	VersionPolicyReleaseTrain VersionPolicy = "release-train"
	// VersionPolicySnapshot assigns development snapshot versions.
	VersionPolicySnapshot VersionPolicy = "snapshot"
)

// ParseVersionPolicy validates a version policy.
func ParseVersionPolicy(value string) (VersionPolicy, error) {
	switch VersionPolicy(value) {
	case VersionPolicyReleaseTrain, VersionPolicySnapshot:
		return VersionPolicy(value), nil
	default:
		return "", fmt.Errorf("unsupported version policy %q", value)
	}
}
