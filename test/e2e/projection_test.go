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

import (
	"path/filepath"
	"testing"
)

func TestVerifyCopiesOnlyExplicitEntries(t *testing.T) {
	root := copyFixture(t, "projection")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/projection")

	result := runArcpub(t,
		"verify",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
	)
	assertExitCode(t, result, 0)

	worktree := targetWorktreePath(targetRoot, "arcoris/projection")
	assertFileExists(t, filepath.Join(worktree, "go.mod"))
	assertFileExists(t, filepath.Join(worktree, "public", "doc.go"))
	assertPathMissing(t, filepath.Join(worktree, "secret.txt"))
	assertPathMissing(t, filepath.Join(worktree, "private", "secret.go"))
	assertPathMissing(t, filepath.Join(worktree, "private"))
}
