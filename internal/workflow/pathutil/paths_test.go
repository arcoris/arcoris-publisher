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

package pathutil

import (
	"path/filepath"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestCleanAbsRejectsBlankPath(t *testing.T) {
	if _, err := CleanAbs(" "); err == nil {
		t.Fatal("CleanAbs(blank) error = nil")
	}
}

func TestEnsureInside(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "repo")
	child := filepath.Join(root, "module")
	outside := filepath.Join(string(filepath.Separator), "other")
	if err := EnsureInside(root, child); err != nil {
		t.Fatalf("EnsureInside(child) error = %v", err)
	}
	if err := EnsureInside(root, outside); err == nil {
		t.Fatal("EnsureInside(escaped) error = nil")
	}
}

func TestJoinRelative(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "repo")
	if got := JoinRelative(root, manifest.RelativePath(".")); got != filepath.Clean(root) {
		t.Fatalf("JoinRelative(.) = %q", got)
	}
	want := filepath.Join(root, "pkg", "api")
	if got := JoinRelative(root, manifest.RelativePath("pkg/api")); got != want {
		t.Fatalf("JoinRelative(pkg/api) = %q, want %q", got, want)
	}
}
