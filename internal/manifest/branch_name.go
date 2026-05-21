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

import (
	"fmt"
	"strings"
)

// BranchName is a lightweight-validated Git branch name used by manifest policy.
type BranchName string

// ParseBranchName validates a branch name used in publication policy.
func ParseBranchName(value string) (BranchName, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("branch is required")
	}
	if strings.TrimSpace(value) != value || strings.ContainsAny(value, " \t\n\r") {
		return "", fmt.Errorf("branch must not contain whitespace")
	}
	if strings.HasPrefix(value, "-") {
		return "", fmt.Errorf("branch must not start with '-'")
	}
	if hasRejectedBranchShape(value) {
		return "", fmt.Errorf("branch contains characters rejected by policy")
	}
	return BranchName(value), nil
}

// String returns the branch name string.
func (b BranchName) String() string { return string(b) }
