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

package manifest_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestParseRelativePathNormalizesSafePaths(t *testing.T) {
	got, err := manifest.ParseRelativePath("path", "module//go.mod", false)
	if err != nil {
		t.Fatalf("ParseRelativePath returned error: %v", err)
	}
	if got.String() != "module/go.mod" {
		t.Fatalf("unexpected normalized path: %q", got.String())
	}
}

func TestParseRelativePathHandlesDotPolicy(t *testing.T) {
	if _, err := manifest.ParseRelativePath("path", ".", false); err == nil {
		t.Fatalf("expected dot rejection when allowDot is false")
	}
	got, err := manifest.ParseRelativePath("path", ".", true)
	if err != nil {
		t.Fatalf("ParseRelativePath returned error for dot: %v", err)
	}
	if got.String() != "." {
		t.Fatalf("unexpected dot path: %q", got.String())
	}
}

func TestParseRelativePathRejectsUnsafePaths(t *testing.T) {
	for _, value := range []string{
		"",
		"../secret",
		"safe/../secret",
		"/absolute",
		"C:/absolute",
		"dir\\file",
		" dir/file",
	} {
		if _, err := manifest.ParseRelativePath("path", value, false); err == nil {
			t.Fatalf("ParseRelativePath(%q) returned nil error", value)
		}
	}
}
