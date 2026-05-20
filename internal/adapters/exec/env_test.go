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

package exec

import "testing"

func TestMergeEnvOverlaysAssignments(t *testing.T) {
	base := []string{"A=1", "B=2"}
	override := []string{"B=3", "C=4"}

	got := MergeEnv(base, override)
	want := []string{"A=1", "B=3", "C=4"}

	assertStringSlice(t, got, want)
	base[1] = "B=mutated"
	if got[1] != "B=3" {
		t.Fatalf("MergeEnv should detach base slice, got %#v", got)
	}
}

func TestMergeEnvKeepsMalformedEntries(t *testing.T) {
	got := MergeEnv([]string{"A=1"}, []string{"NO_EQUALS", "=empty"})
	want := []string{"A=1", "NO_EQUALS", "=empty"}

	assertStringSlice(t, got, want)
}

func assertStringSlice(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("item %d = %q, want %q in %#v", i, got[i], want[i], got)
		}
	}
}
