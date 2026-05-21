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

import (
	"fmt"
	"strings"
)

// validatePrerelease checks SemVer prerelease identifier rules not captured by the regex.
func validatePrerelease(value string) error {
	if value == "" {
		return nil
	}
	for _, segment := range strings.Split(value, ".") {
		if segment == "" {
			return fmt.Errorf("prerelease contains an empty identifier")
		}
		if isNumeric(segment) && len(segment) > 1 && segment[0] == '0' {
			return fmt.Errorf("numeric prerelease identifier %q must not have leading zeroes", segment)
		}
	}
	return nil
}

// validateBuildMetadata checks build metadata identifier rules not captured by the regex.
func validateBuildMetadata(value string) error {
	if value == "" {
		return nil
	}
	for _, segment := range strings.Split(value, ".") {
		if segment == "" {
			return fmt.Errorf("build metadata contains an empty identifier")
		}
	}
	return nil
}

// isNumeric reports whether value contains only ASCII decimal digits.
func isNumeric(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
