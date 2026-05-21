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
	"regexp"
	"strings"
)

var simpleNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// validateSimpleName validates lower-kebab identifiers used by manifest names.
//
// Manifest names and module names intentionally share one policy: a lowercase
// DNS-like token that is stable in JSON/YAML, filenames, logs, and branch
// labels. Keeping the rule centralized avoids two public parsers silently
// accepting different identifier shapes.
func validateSimpleName(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	if strings.TrimSpace(value) != value {
		return fmt.Errorf("%s must not have surrounding whitespace", field)
	}
	if !simpleNamePattern.MatchString(value) {
		return fmt.Errorf("%s must match %s", field, simpleNamePattern.String())
	}
	return nil
}
