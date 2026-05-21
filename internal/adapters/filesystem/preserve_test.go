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

package filesystem

import "testing"

func TestNormalizePreserve(t *testing.T) {
	got := normalizePreserve([]string{"./b/c", "", ".", "a"})
	want := []string{"a", "b/c"}
	if len(got) != len(want) {
		t.Fatalf("normalizePreserve() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("normalizePreserve()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestPreservePredicates(t *testing.T) {
	preserve := []string{"keep", "nested/file.txt"}

	if !shouldPreserve("keep/child.txt", preserve) {
		t.Fatalf("shouldPreserve() should accept descendants")
	}
	if shouldPreserve("other", preserve) {
		t.Fatalf("shouldPreserve() should reject unrelated paths")
	}
	if !hasPreservedChild("nested", preserve) {
		t.Fatalf("hasPreservedChild() should detect descendants")
	}
	if hasPreservedChild("nested/file.txt", preserve) {
		t.Fatalf("hasPreservedChild() should reject exact file path")
	}
}
