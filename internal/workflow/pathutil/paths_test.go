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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestCleanAbsRejectsBlankPath(t *testing.T) {
	if _, err := CleanAbs(" "); err == nil {
		t.Fatal("CleanAbs(blank) error = nil")
	}
}

func TestEnsureInside(t *testing.T) {
	if err := EnsureInside("/repo", "/repo/module"); err != nil {
		t.Fatalf("EnsureInside(child) error = %v", err)
	}
	if err := EnsureInside("/repo", "/other"); err == nil {
		t.Fatal("EnsureInside(escaped) error = nil")
	}
}

func TestJoinRelative(t *testing.T) {
	if got := JoinRelative("/repo", manifest.RelativePath(".")); got != "/repo" {
		t.Fatalf("JoinRelative(.) = %q", got)
	}
	if got := JoinRelative("/repo", manifest.RelativePath("pkg/api")); got != "/repo/pkg/api" {
		t.Fatalf("JoinRelative(pkg/api) = %q", got)
	}
}
