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

func TestParseValidVersions(t *testing.T) {
	inputs := []string{
		"v0.1.0",
		"v1.2.3",
		"v1.2.3-rc.1",
		"v0.0.0-20260521143000-abcdefabcdef",
	}
	for _, input := range inputs {
		version, err := Parse(input)
		if err != nil {
			t.Fatalf("Parse(%q) error = %v", input, err)
		}
		if version.String() != input {
			t.Fatalf("String() = %q, want %q", version.String(), input)
		}
	}
}

func TestParseRejectsInvalidVersions(t *testing.T) {
	inputs := []string{
		"",
		"0.1.0",
		"v1",
		"v1.2",
		"v1.2.x",
		"v1.2.3+build",
		" latest ",
		"main",
	}
	for _, input := range inputs {
		if _, err := Parse(input); err == nil {
			t.Fatalf("Parse(%q) succeeded, want error", input)
		}
	}
}

func TestVersionKind(t *testing.T) {
	release := Must("v1.2.3-rc.1")
	if !release.IsRelease() || release.IsPseudo() {
		t.Fatalf("release kind mismatch")
	}
	pseudo := Must("v0.0.0-20260521143000-abcdefabcdef")
	if !pseudo.IsPseudo() || pseudo.IsRelease() {
		t.Fatalf("pseudo kind mismatch")
	}
}
