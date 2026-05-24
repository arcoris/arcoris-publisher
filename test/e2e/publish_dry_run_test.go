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

package e2e_test

import "testing"

func TestPublishDryRunMinimalFixture(t *testing.T) {
	root := copyFixture(t, "minimal")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/foundation", "arcoris/control")

	result := runArcpub(t,
		"publish",
		"--dry-run",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
		"--output", "json",
	)
	assertExitCode(t, result, 0)
	decoded := assertJSON(t, result.Stdout)
	if decoded["status"] != "verified" {
		t.Fatalf("dry-run workflow status = %#v, want verified\n%s", decoded["status"], result.Stdout)
	}
}
