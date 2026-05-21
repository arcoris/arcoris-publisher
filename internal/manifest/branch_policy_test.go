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

import "testing"

func TestHasRejectedBranchShapeRejectsGitRefSyntaxHazards(t *testing.T) {
	for _, value := range []string{"main?bad", "feature..main", "feature//main", "main.lock", "main@{upstream}", "/main", "main/", "main.", "."} {
		if !hasRejectedBranchShape(value) {
			t.Fatalf("hasRejectedBranchShape(%q) = false", value)
		}
	}
}

func TestHasRejectedBranchShapeAcceptsSimpleBranchPaths(t *testing.T) {
	if hasRejectedBranchShape("release/v1") {
		t.Fatalf("safe branch path was rejected")
	}
}
