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

import (
	"strings"
	"testing"
)

func TestRenderTreeHashSortsEntriesAndIncludesModeOptionally(t *testing.T) {
	entries := []hashEntry{
		{rel: "b.txt", kind: "file", mode: 0o600, content: "b"},
		{rel: "a.txt", kind: "file", mode: 0o644, content: "a"},
	}
	sortedEntries := []hashEntry{
		{rel: "a.txt", kind: "file", mode: 0o644, content: "a"},
		{rel: "b.txt", kind: "file", mode: 0o600, content: "b"},
	}

	if got, want := renderTreeHash(entries, false), renderTreeHash(sortedEntries, false); got != want {
		t.Fatalf("renderTreeHash() should sort entries: %q != %q", got, want)
	}
	if renderTreeHash(entries, false) == renderTreeHash(entries, true) {
		t.Fatalf("renderTreeHash() should change when mode is included")
	}
	if got := renderTreeHash(entries, false).String(); !strings.HasPrefix(got, "sha256:") {
		t.Fatalf("renderTreeHash() = %q, want sha256 prefix", got)
	}
}
