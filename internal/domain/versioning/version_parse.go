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

// ParseVersion validates a Go module version string.
//
// The parser keeps policy local to the domain layer: it rejects missing values,
// surrounding whitespace, non-canonical numeric components, malformed
// prerelease identifiers, and malformed build metadata. It does not inspect
// repository tags or module import paths.
func ParseVersion(value string) (Version, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("version is required")
	}
	if value != strings.TrimSpace(value) {
		return "", fmt.Errorf("version must not contain surrounding whitespace")
	}

	match := semanticVersionRE.FindStringSubmatch(value)
	if len(match) == 0 {
		return "", fmt.Errorf("version %q must use vMAJOR.MINOR.PATCH form", value)
	}
	if err := validatePrerelease(match[prereleaseSubmatchIndex]); err != nil {
		return "", err
	}
	if err := validateBuildMetadata(match[buildMetadataSubmatchIndex]); err != nil {
		return "", err
	}
	return Version(value), nil
}

// MustVersion parses value and panics if it is invalid.
//
// MustVersion is intended for tests and static wiring. Runtime code should call
// ParseVersion and return diagnostics to the caller.
func MustVersion(value string) Version {
	version, err := ParseVersion(value)
	if err != nil {
		panic(err)
	}
	return version
}
