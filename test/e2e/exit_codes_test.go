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

func TestUnknownCommandExitCode(t *testing.T) {
	result := runArcpub(t, "unknown")
	assertExitCode(t, result, 64)
	assertContains(t, result.Stderr, "unknown command")
}

func TestPlanMissingVersionExitCode(t *testing.T) {
	root := copyFixture(t, "minimal")
	result := runArcpub(t, "plan", "--manifest", e2eManifest(root))
	assertExitCode(t, result, 64)
}

func TestPlanInvalidVersionExitCode(t *testing.T) {
	root := copyFixture(t, "minimal")
	result := runArcpub(t, "plan", "--manifest", e2eManifest(root), "--version", "not-a-version")
	assertExitCode(t, result, 64)
}

func TestPlanMissingManifestExitCode(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing", "arcpub.yaml")
	result := runArcpub(t, "plan", "--manifest", missing, "--version", "v0.1.0")
	assertExitCode(t, result, 1)
}

func TestVerifyFailedExitCode(t *testing.T) {
	root := copyFixture(t, "bad-verification")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/broken")

	result := runArcpub(t,
		"verify",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
	)
	assertExitCode(t, result, 2)
	assertContains(t, result.Stdout, "Workflow")
	assertContains(t, result.Stdout, "Verification: failed")
}

func TestPublishDryRunFailedVerificationExitCode(t *testing.T) {
	root := copyFixture(t, "bad-verification")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/broken")

	result := runArcpub(t,
		"publish",
		"--dry-run",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
	)
	assertExitCode(t, result, 2)
	assertContains(t, result.Stdout, "Workflow")
	assertContains(t, result.Stdout, "Verification: failed")
}
