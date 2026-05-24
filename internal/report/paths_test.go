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

package report

import "testing"

func TestIncludePathHidesLocalAbsolutePathsByDefault(t *testing.T) {
	t.Parallel()

	for _, path := range []string{
		"/repo",
		"/target/module/go.mod",
		`C:\repo\target`,
		`C:/repo/target`,
		`\\server\share\repo`,
	} {
		if got := includePath(path, Options{}); got != "" {
			t.Fatalf("includePath(%q) = %q", path, got)
		}
	}
}

func TestIncludePathAllowsRelativeAndExplicitLocalPaths(t *testing.T) {
	t.Parallel()

	if got := includePath("src/arcoris.dev/foundation", Options{}); got != "src/arcoris.dev/foundation" {
		t.Fatalf("relative includePath() = %q", got)
	}
	if got := includePath("/repo", Options{IncludeLocalPaths: true}); got != "/repo" {
		t.Fatalf("explicit local includePath() = %q", got)
	}
}
