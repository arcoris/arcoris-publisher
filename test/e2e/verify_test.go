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
	"strings"
	"testing"
)

func TestVerifyMinimalFixture(t *testing.T) {
	requireGitAndGo(t)
	root := copyFixture(t, "minimal")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/foundation", "arcoris/control")

	result, decoded := runArcpubJSON(t,
		0,
		"verify",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
		"--output", "json",
	)
	if decoded["status"] != "verified" {
		t.Fatalf("workflow status = %#v, want verified\n%s", decoded["status"], result.Stdout)
	}
	assertVerifiedWorkflowReport(t, decoded, 2)
	assertFileExists(t, filepath.Join(targetWorktreePath(targetRoot, "arcoris/foundation"), "go.mod"))
	assertFileExists(t, filepath.Join(targetWorktreePath(targetRoot, "arcoris/control"), "go.mod"))
}

func TestVerifyBadFixtureReturnsVerificationFailed(t *testing.T) {
	requireGitAndGo(t)
	root := copyFixture(t, "bad-verification")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/broken")

	result, decoded := runArcpubJSON(t,
		2,
		"verify",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
		"--output", "json",
	)
	if decoded["status"] != "verification_failed" {
		t.Fatalf("workflow status = %#v, want verification_failed\n%s", decoded["status"], result.Stdout)
	}
	assertNotContains(t, result.Stderr, "panic")
}

func TestVerifyDoesNotLeakLocalPathsByDefault(t *testing.T) {
	requireGitAndGo(t)
	root := copyFixture(t, "minimal")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/foundation", "arcoris/control")

	result, _ := runArcpubJSON(t,
		0,
		"verify",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
		"--output", "json",
	)
	assertNoLocalPathLeak(t, result.Stdout, root, targetRoot)
}

func TestVerifyIncludeLocalPaths(t *testing.T) {
	requireGitAndGo(t)
	root := copyFixture(t, "minimal")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/foundation", "arcoris/control")

	result := runArcpub(t,
		"verify",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
		"--output", "json",
		"--include-local-paths",
	)
	assertExitCode(t, result, 0)
	if !strings.Contains(result.Stdout, root) && !strings.Contains(result.Stdout, targetRoot) {
		t.Fatalf("include-local-paths output did not include source or target paths:\n%s", result.Stdout)
	}
}

func TestVerifyMissingTargetWorktreeReturnsError(t *testing.T) {
	requireGitAndGo(t)
	root := copyFixture(t, "minimal")
	initGitRepo(t, root)
	targetRoot := t.TempDir()

	result := runArcpub(t,
		"verify",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
	)
	assertExitCode(t, result, 1)
	assertContains(t, result.Stderr, "target")
	assertContains(t, result.Stderr, "worktree")
}

func assertVerifiedWorkflowReport(t *testing.T, report map[string]any, moduleCount float64) {
	t.Helper()
	if stringField(t, report, "status") != "verified" {
		t.Fatalf("workflow status = %#v, want verified", report["status"])
	}
	source := objectField(t, report, "source")
	target := objectField(t, report, "target")
	construct := objectField(t, report, "construct")
	moduleFile := objectField(t, report, "moduleFile")
	verify := objectField(t, report, "verify")
	publish := objectField(t, report, "publish")

	if stringField(t, source, "status") != "present" || floatField(t, source, "moduleCount") != moduleCount {
		t.Fatalf("unexpected source report: %#v", source)
	}
	if stringField(t, target, "status") != "present" || floatField(t, target, "workspaceCount") != moduleCount {
		t.Fatalf("unexpected target report: %#v", target)
	}
	if stringField(t, construct, "status") != "present" || floatField(t, construct, "moduleCount") != moduleCount {
		t.Fatalf("unexpected construct report: %#v", construct)
	}
	if stringField(t, moduleFile, "status") != "present" || floatField(t, moduleFile, "moduleCount") != moduleCount {
		t.Fatalf("unexpected moduleFile report: %#v", moduleFile)
	}
	if stringField(t, verify, "status") != "passed" || floatField(t, verify, "moduleCount") != moduleCount {
		t.Fatalf("unexpected verify report: %#v", verify)
	}
	if stringField(t, publish, "status") != "empty" {
		t.Fatalf("unexpected publish report: %#v", publish)
	}
}
