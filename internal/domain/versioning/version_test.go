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

func TestParseVersionAcceptsCanonicalVersions(t *testing.T) {
	version, err := ParseVersion("v1.2.3-rc.1+build.5")
	if err != nil {
		t.Fatalf("ParseVersion returned error: %v", err)
	}
	if version.Major() != 1 || version.Minor() != 2 || version.Patch() != 3 {
		t.Fatalf("unexpected semantic parts: %d.%d.%d", version.Major(), version.Minor(), version.Patch())
	}
	if version.Prerelease() != "rc.1" {
		t.Fatalf("unexpected prerelease: %q", version.Prerelease())
	}
	if version.BuildMetadata() != "build.5" {
		t.Fatalf("unexpected build metadata: %q", version.BuildMetadata())
	}
	if version.WithoutBuildMetadata() != "v1.2.3-rc.1" {
		t.Fatalf("unexpected version without metadata: %q", version.WithoutBuildMetadata())
	}
}

func TestParseVersionRejectsInvalidVersions(t *testing.T) {
	cases := []string{"", "1.2.3", "v1.2", "v01.2.3", "v1.2.3-01", "v1.2.3+"}
	for _, tc := range cases {
		if _, err := ParseVersion(tc); err == nil {
			t.Fatalf("ParseVersion(%q) succeeded, expected error", tc)
		}
	}
}

func TestMustVersionPanicsOnInvalidVersion(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("MustVersion did not panic")
		}
	}()
	_ = MustVersion("bad")
}
