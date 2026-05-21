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

import "testing"

func TestVersionMetadataReturnsEmptyForInvalidVersion(t *testing.T) {
	version := Version("not-a-version")

	if version.Prerelease() != "" {
		t.Fatalf("Prerelease() = %q, want empty", version.Prerelease())
	}
	if version.BuildMetadata() != "" {
		t.Fatalf("BuildMetadata() = %q, want empty", version.BuildMetadata())
	}
	if version.Major() != 0 || version.Minor() != 0 || version.Patch() != 0 {
		t.Fatalf("invalid semantic parts = %d.%d.%d, want 0.0.0", version.Major(), version.Minor(), version.Patch())
	}
	if got := parseSemanticPart(MustVersion("v1.2.3"), 99); got != 0 {
		t.Fatalf("parseSemanticPart(invalid index) = %d, want 0", got)
	}
}

func TestParseVersionRejectsMalformedMetadataIdentifiers(t *testing.T) {
	for _, value := range []string{" v1.2.3", "v1.2.3-alpha..1", "v1.2.3+build..1"} {
		if _, err := ParseVersion(value); err == nil {
			t.Fatalf("ParseVersion(%q) succeeded, expected error", value)
		}
	}
}

func TestMustVersionReturnsValidVersion(t *testing.T) {
	if got := MustVersion("v1.2.3").String(); got != "v1.2.3" {
		t.Fatalf("MustVersion() = %q", got)
	}
}

func TestIsNumericRejectsEmptyString(t *testing.T) {
	if isNumeric("") {
		t.Fatalf("isNumeric(empty) = true, want false")
	}
}
