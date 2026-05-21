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

// DirtyPolicy controls whether a dirty source repository may be published.
type DirtyPolicy string

const (
	// DirtyPolicyFail rejects publication from a dirty source repository.
	DirtyPolicyFail DirtyPolicy = "fail"
	// DirtyPolicyWarn records a warning but allows publication from a dirty source repository.
	DirtyPolicyWarn DirtyPolicy = "warn"
	// DirtyPolicyAllow allows publication from a dirty source repository.
	DirtyPolicyAllow DirtyPolicy = "allow"
)

// ParseDirtyPolicy validates a source dirty policy.
func ParseDirtyPolicy(value string) (DirtyPolicy, error) {
	switch DirtyPolicy(value) {
	case DirtyPolicyFail, DirtyPolicyWarn, DirtyPolicyAllow:
		return DirtyPolicy(value), nil
	default:
		return "", fmt.Errorf("unsupported dirtyPolicy %q", value)
	}
}
