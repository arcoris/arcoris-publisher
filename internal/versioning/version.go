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
	"regexp"
	"strings"
)

var (
	// semverPattern accepts canonical Go module SemVer values without build
	// metadata.
	semverPattern = regexp.MustCompile(
		`^v(0|[1-9][0-9]*)\.` +
			`(0|[1-9][0-9]*)\.` +
			`(0|[1-9][0-9]*)` +
			`(?:-[0-9A-Za-z][0-9A-Za-z.-]*)?$`,
	)

	// pseudoVersionPattern accepts Go pseudo-versions after the base SemVer
	// shape has already been validated.
	pseudoVersionPattern = regexp.MustCompile(
		`^v(0|[1-9][0-9]*)\.` +
			`(0|[1-9][0-9]*)\.` +
			`(0|[1-9][0-9]*)` +
			`(?:-[0-9A-Za-z][0-9A-Za-z.-]*)?` +
			`(?:\.0)?-[0-9]{14}-[0-9a-fA-F]{12,40}$`,
	)
)

// Version is a validated Go module version used by ARCORIS Publisher.
//
// Version intentionally models the string form used in go.mod requirements and
// Git tags. Build metadata is rejected because Go module versions and tags used
// by this publisher should be stable SemVer or pseudo-version strings without a
// non-canonical +metadata suffix.
type Version string

// Parse validates value as a Go module publication version.
func Parse(value string) (Version, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("version is required")
	}

	if trimmed != value || strings.ContainsAny(value, " \t\n\r") {
		return "", fmt.Errorf("version must not contain whitespace")
	}

	if strings.Contains(value, "+") {
		return "", fmt.Errorf("version build metadata is not supported")
	}

	if !semverPattern.MatchString(value) {
		return "", fmt.Errorf(
			"version %q must be a canonical vMAJOR.MINOR.PATCH SemVer value",
			value,
		)
	}

	return Version(value), nil
}

// Must parses value and panics when it is invalid.
func Must(value string) Version {
	version, err := Parse(value)
	if err != nil {
		panic(err)
	}
	return version
}

// String returns the version string.
func (v Version) String() string { return string(v) }

// IsZero reports whether v is the empty version value.
func (v Version) IsZero() bool { return v == "" }

// IsPseudo reports whether v looks like a Go pseudo-version.
func (v Version) IsPseudo() bool { return pseudoVersionPattern.MatchString(string(v)) }

// IsRelease reports whether v is a non-pseudo SemVer publication version.
func (v Version) IsRelease() bool {
	return !v.IsZero() && !v.IsPseudo() && semverPattern.MatchString(string(v))
}
