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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

type commandResult struct {
	Code   int
	Stdout string
	Stderr string
}

var builtArcpub struct {
	once sync.Once
	dir  string
	path string
	err  error
}

func TestMain(m *testing.M) {
	code := m.Run()
	if builtArcpub.dir != "" {
		_ = os.RemoveAll(builtArcpub.dir)
	}
	os.Exit(code)
}

func repoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	for {
		goMod := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(goMod)
		if err == nil && strings.Contains(string(data), "module arcoris.dev/arcoris-publisher") {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("repository root not found from %s", dir)
		}
		dir = parent
	}
}

func arcpubBinary(t *testing.T) string {
	t.Helper()
	requireExecutable(t, "go")

	root := repoRoot(t)
	builtArcpub.once.Do(func() {
		dir, err := os.MkdirTemp("", "arcpub-e2e-bin-*")
		if err != nil {
			builtArcpub.err = err
			return
		}
		builtArcpub.dir = dir

		name := "arcpub"
		if goruntime.GOOS == "windows" {
			name += ".exe"
		}
		builtArcpub.path = filepath.Join(dir, name)

		result := runCommand(t, root, nil, "go", "build", "-o", builtArcpub.path, "./cmd/arcpub")
		if result.Code != 0 {
			version := runCommand(t, root, nil, "go", "env", "GOVERSION")
			builtArcpub.err = fmt.Errorf(
				"go build -o %s ./cmd/arcpub failed in %s with code %d\ngo env GOVERSION: %s\nstdout:\n%s\nstderr:\n%s",
				builtArcpub.path,
				root,
				result.Code,
				strings.TrimSpace(version.Stdout+version.Stderr),
				result.Stdout,
				result.Stderr,
			)
		}
	})
	if builtArcpub.err != nil {
		t.Fatalf("build arcpub binary: %v", builtArcpub.err)
	}

	return builtArcpub.path
}

func runArcpub(t *testing.T, args ...string) commandResult {
	t.Helper()
	return runArcpubInDir(t, repoRoot(t), args...)
}

func runArcpubInDir(t *testing.T, dir string, args ...string) commandResult {
	t.Helper()
	return runCommand(t, dir, nil, arcpubBinary(t), args...)
}

func runArcpubJSON(t *testing.T, wantCode int, args ...string) (commandResult, map[string]any) {
	t.Helper()
	result := runArcpub(t, args...)
	assertExitCode(t, result, wantCode)
	return result, assertJSON(t, result.Stdout)
}

func runCommand(
	t *testing.T,
	dir string,
	env []string,
	name string,
	args ...string,
) commandResult {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = testEnv(env)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := commandResult{
		Code:   0,
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if err == nil {
		return result
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		result.Code = exitErr.ExitCode()
		return result
	}
	if ctx.Err() != nil {
		t.Fatalf("%s %v timed out: %v", name, args, ctx.Err())
	}
	t.Fatalf("%s %v failed to start: %v", name, args, err)
	return result
}

func testEnv(extra []string) []string {
	env := append([]string{}, os.Environ()...)
	env = append(env,
		"GIT_TERMINAL_PROMPT=0",
		"GIT_CONFIG_NOSYSTEM=1",
		"GOTOOLCHAIN=local",
	)
	return append(env, extra...)
}

func copyFixture(t *testing.T, name string) string {
	t.Helper()

	source := filepath.Join(repoRoot(t), "internal", "testdata", "e2e", name)
	target := t.TempDir()
	err := filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("fixture %s contains a symlink; e2e fixtures must be regular files and directories", path)
		}

		dst := filepath.Join(target, rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(dst, info.Mode())
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("fixture %s has unsupported file type %s", path, info.Mode().Type())
		}

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return copyFile(path, dst, info.Mode())
	})
	if err != nil {
		t.Fatalf("copy fixture %q: %v", name, err)
	}
	return target
}

func copyFile(source string, target string, mode os.FileMode) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	requireExecutable(t, "git")

	result := runCommand(t, dir, nil, "git", "init", "-b", "main")
	if result.Code != 0 {
		mustRun(t, dir, "git", "init")
		mustRun(t, dir, "git", "checkout", "-B", "main")
	}
	mustRun(t, dir, "git", "config", "user.name", "ARCORIS Test")
	mustRun(t, dir, "git", "config", "user.email", "arcoris-test@example.invalid")
	mustRun(t, dir, "git", "config", "core.autocrlf", "false")
	mustRun(t, dir, "git", "config", "commit.gpgsign", "false")
	mustRun(t, dir, "git", "config", "init.defaultBranch", "main")
	mustRun(t, dir, "git", "add", ".")
	mustRun(t, dir, "git", "commit", "-m", "test: seed fixture")
}

func initBareGitRepo(t *testing.T, dir string) {
	t.Helper()
	requireExecutable(t, "git")

	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		t.Fatalf("create bare repo parent: %v", err)
	}
	mustRun(t, repoRoot(t), "git", "init", "--bare", dir)
	mustRun(t, repoRoot(t), "git", "--git-dir", dir, "config", "receive.denyDeleteCurrent", "ignore")
}

func installBareHook(t *testing.T, bareRepo string, name string, script string) {
	t.Helper()
	if goruntime.GOOS == "windows" {
		t.Skip("bare repository hook e2e tests require POSIX shell hooks")
	}
	path := filepath.Join(bareRepo, "hooks", name)
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write hook %s: %v", path, err)
	}
}

func initTargetGitWorktree(t *testing.T, dir string, remote string) {
	t.Helper()
	requireExecutable(t, "git")

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create target worktree %s: %v", dir, err)
	}
	result := runCommand(t, dir, nil, "git", "init", "-b", "main")
	if result.Code != 0 {
		mustRun(t, dir, "git", "init")
		mustRun(t, dir, "git", "checkout", "-B", "main")
	}
	mustRun(t, dir, "git", "config", "user.name", "ARCORIS Test")
	mustRun(t, dir, "git", "config", "user.email", "arcoris-test@example.invalid")
	mustRun(t, dir, "git", "config", "core.autocrlf", "false")
	mustRun(t, dir, "git", "config", "commit.gpgsign", "false")
	if err := os.WriteFile(filepath.Join(dir, ".seed"), []byte("seed\n"), 0o644); err != nil {
		t.Fatalf("write target seed file: %v", err)
	}
	mustRun(t, dir, "git", "add", ".seed")
	mustRun(t, dir, "git", "commit", "-m", "test: seed target")
	if remote != "" {
		mustRun(t, dir, "git", "remote", "add", "origin", remote)
	}
}

func mustRun(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	result := runCommand(t, dir, nil, name, args...)
	if result.Code != 0 {
		t.Fatalf("%s %v failed with code %d\nstdout:\n%s\nstderr:\n%s", name, args, result.Code, result.Stdout, result.Stderr)
	}
}

func requireExecutable(t *testing.T, name string) {
	t.Helper()
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%s executable is required for executable e2e tests: %v", name, err)
	}
}

func requireGitAndGo(t *testing.T) {
	t.Helper()
	requireExecutable(t, "git")
	requireExecutable(t, "go")
}

func prepareTargetWorktrees(t *testing.T, targetRoot string, repositories ...string) {
	t.Helper()
	for _, repository := range repositories {
		initTargetGitWorktree(t, targetWorktreePath(targetRoot, repository), "")
	}
}

func targetWorktreePath(targetRoot string, repository string) string {
	return filepath.Join(targetRoot, strings.ReplaceAll(repository, "/", "__"))
}

func assertGitRefExists(t *testing.T, gitDir string, ref string) {
	t.Helper()
	result := runCommand(t, repoRoot(t), nil, "git", "--git-dir", gitDir, "show-ref", "--verify", "--quiet", ref)
	if result.Code != 0 {
		t.Fatalf("git ref %s missing in %s\nstdout:\n%s\nstderr:\n%s", ref, gitDir, result.Stdout, result.Stderr)
	}
}

func assertGitRefMissing(t *testing.T, gitDir string, ref string) {
	t.Helper()
	result := runCommand(t, repoRoot(t), nil, "git", "--git-dir", gitDir, "show-ref", "--verify", "--quiet", ref)
	if result.Code == 0 {
		t.Fatalf("git ref %s unexpectedly exists in %s", ref, gitDir)
	}
}

func assertWorktreeClean(t *testing.T, worktree string) {
	t.Helper()
	result := runCommand(t, worktree, nil, "git", "status", "--porcelain")
	if result.Code != 0 {
		t.Fatalf("git status failed in %s\nstdout:\n%s\nstderr:\n%s", worktree, result.Stdout, result.Stderr)
	}
	if strings.TrimSpace(result.Stdout) != "" {
		t.Fatalf("worktree %s is dirty:\n%s", worktree, result.Stdout)
	}
}

func gitLogMessage(t *testing.T, gitDir string, ref string) string {
	t.Helper()
	result := runCommand(t, repoRoot(t), nil, "git", "--git-dir", gitDir, "log", "-1", "--format=%B", ref)
	if result.Code != 0 {
		t.Fatalf("git log %s in %s failed\nstdout:\n%s\nstderr:\n%s", ref, gitDir, result.Stdout, result.Stderr)
	}
	return result.Stdout
}

func gitTreeContains(t *testing.T, gitDir string, ref string, path string) {
	t.Helper()
	result := runCommand(t, repoRoot(t), nil, "git", "--git-dir", gitDir, "cat-file", "-e", ref+":"+path)
	if result.Code != 0 {
		t.Fatalf("git tree %s:%s missing in %s\nstdout:\n%s\nstderr:\n%s", ref, path, gitDir, result.Stdout, result.Stderr)
	}
}

func gitTreeMissing(t *testing.T, gitDir string, ref string, path string) {
	t.Helper()
	result := runCommand(t, repoRoot(t), nil, "git", "--git-dir", gitDir, "cat-file", "-e", ref+":"+path)
	if result.Code == 0 {
		t.Fatalf("git tree %s:%s unexpectedly exists in %s", ref, path, gitDir)
	}
}

func assertExitCode(t *testing.T, result commandResult, want int) {
	t.Helper()
	if result.Code != want {
		t.Fatalf("exit code = %d, want %d\nstdout:\n%s\nstderr:\n%s", result.Code, want, result.Stdout, result.Stderr)
	}
}

func assertContains(t *testing.T, value string, want string) {
	t.Helper()
	if !strings.Contains(value, want) {
		t.Fatalf("output does not contain %q:\n%s", want, value)
	}
}

func assertNotContains(t *testing.T, value string, forbidden string) {
	t.Helper()
	if strings.Contains(value, forbidden) {
		t.Fatalf("output unexpectedly contains %q:\n%s", forbidden, value)
	}
}

func assertNoLocalPathLeak(t *testing.T, output string, paths ...string) {
	t.Helper()
	for _, path := range paths {
		for _, variant := range localPathVariants(path) {
			assertNotContains(t, output, variant)
		}
	}
}

func localPathVariants(path string) []string {
	if path == "" {
		return nil
	}
	clean := filepath.Clean(path)
	variants := []string{path, clean, filepath.ToSlash(clean)}
	if strings.Contains(clean, `\`) {
		variants = append(variants, strings.ReplaceAll(clean, `\`, `\\`))
	}
	if volume := filepath.VolumeName(clean); volume != "" {
		withoutVolume := strings.TrimPrefix(clean, volume)
		variants = append(variants, filepath.ToSlash(withoutVolume))
	}

	out := make([]string, 0, len(variants))
	seen := map[string]struct{}{}
	for _, variant := range variants {
		if variant == "." || variant == "" {
			continue
		}
		if _, ok := seen[variant]; ok {
			continue
		}
		seen[variant] = struct{}{}
		out = append(out, variant)
	}
	return out
}

func assertJSON(t *testing.T, output string) map[string]any {
	t.Helper()
	var decoded map[string]any
	if err := json.Unmarshal([]byte(output), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, output)
	}
	return decoded
}

func objectField(t *testing.T, object map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := object[key].(map[string]any)
	if !ok {
		t.Fatalf("field %q = %#v, want object", key, object[key])
	}
	return value
}

func arrayField(t *testing.T, object map[string]any, key string) []any {
	t.Helper()
	value, ok := object[key].([]any)
	if !ok {
		t.Fatalf("field %q = %#v, want array", key, object[key])
	}
	return value
}

func stringField(t *testing.T, object map[string]any, key string) string {
	t.Helper()
	value, ok := object[key].(string)
	if !ok {
		t.Fatalf("field %q = %#v, want string", key, object[key])
	}
	return value
}

func floatField(t *testing.T, object map[string]any, key string) float64 {
	t.Helper()
	value, ok := object[key].(float64)
	if !ok {
		t.Fatalf("field %q = %#v, want number", key, object[key])
	}
	return value
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected file %s: %v", path, err)
	}
	if info.IsDir() {
		t.Fatalf("expected file %s, got directory", path)
	}
}

func assertPathMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected path to be absent: %s", path)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", path, err)
	}
}

func e2eManifest(root string) string {
	return filepath.Join(root, "arcpub.yaml")
}
