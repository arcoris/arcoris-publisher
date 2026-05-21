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

func TestParseBranchNameAcceptsSafeBranch(t *testing.T) {
	got, err := manifest.ParseBranchName("release/v1")
	if err != nil {
		t.Fatalf("ParseBranchName returned error: %v", err)
	}
	if got.String() != "release/v1" {
		t.Fatalf("unexpected branch: %q", got.String())
	}
}

func TestParseBranchNameRejectsUnsafeRefs(t *testing.T) {
	for _, value := range []string{
		"",
		"bad branch",
		"-branch",
		"feature..main",
		"feature//main",
		"/main",
		"main/",
		"main.lock",
		"main@{upstream}",
	} {
		if _, err := manifest.ParseBranchName(value); err == nil {
			t.Fatalf("ParseBranchName(%q) returned nil error", value)
		}
	}
}
