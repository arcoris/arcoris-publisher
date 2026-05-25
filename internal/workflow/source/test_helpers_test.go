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

package source

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
	"arcoris.dev/arcoris-publisher/internal/plan"
	portfs "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	portgit "arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

// testModule describes one realistic module fixture passed through the real
// manifest, resolver, versioning, and plan constructors.
type testModule struct {
	// name is the resolved module name used by staging and module manifests.
	name string

	// dependencies are direct internal module names declared by the module
	// manifest.
	dependencies []string

	// entries overrides the default explicit go.mod and contracts publication
	// entries.
	entries []manifest.PublishEntrySpec
}

// defaultEntries returns the standard fixture publication content.
func defaultEntries() []manifest.PublishEntrySpec {
	return []manifest.PublishEntrySpec{
		{Type: string(manifest.PublishEntryFile), From: "go.mod", To: "go.mod"},
		{Type: string(manifest.PublishEntryDirectory), From: "contracts", To: "contracts"},
	}
}

// mustPlan builds a complete publication plan using the real upstream packages.
func mustPlan(t *testing.T, dirtyPolicy string, modules ...testModule) plan.Plan {
	t.Helper()
	stagingModules := make([]staging.ModuleSpec, 0, len(modules))
	moduleManifests := make([]modulemanifest.Manifest, 0, len(modules))
	for _, mod := range modules {
		stagingModules = append(stagingModules, staging.ModuleSpec{
			Name:       mod.name,
			SourceDir:  "src/arcoris.dev/" + mod.name,
			Repository: "arcoris/" + mod.name,
		})
		entries := mod.entries
		if entries == nil {
			entries = defaultEntries()
		}
		moduleManifest, err := modulemanifest.New(modulemanifest.Spec{
			APIVersion:   string(manifest.APIVersionV1Alpha1),
			Kind:         string(manifest.KindModuleManifest),
			Metadata:     manifest.MetadataSpec{Name: mod.name},
			Module:       manifest.ModuleIdentitySpec{Path: "arcoris.dev/" + mod.name},
			Dependencies: modulemanifest.DependenciesSpec{Internal: mod.dependencies},
			Publish:      modulemanifest.PublishSpec{Entries: entries},
		})
		if err != nil {
			t.Fatalf("module.New(%s) error = %v", mod.name, err)
		}
		moduleManifests = append(moduleManifests, moduleManifest)
	}
	var dirtyPolicyPtr *string
	if dirtyPolicy != "" {
		dirtyPolicyPtr = &dirtyPolicy
	}
	stagingManifest, err := staging.New(staging.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindStagingManifest),
		Metadata:   manifest.MetadataSpec{Name: "arcoris"},
		Source: manifest.SourceSpec{
			Repository:    "arcoris/arcoris",
			DefaultBranch: "main",
			DirtyPolicy:   dirtyPolicyPtr,
		},
		Modules: stagingModules,
	})
	if err != nil {
		t.Fatalf("staging.New() error = %v", err)
	}
	set, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest,
		Modules: moduleManifests,
	})
	if err != nil {
		t.Fatalf("resolved.Resolve() error = %v", err)
	}
	version := versioning.Must("v0.3.0")
	p, err := plan.FromPublicationSet(set, version)
	if err != nil {
		t.Fatalf("plan.FromPublicationSet() error = %v", err)
	}
	return p
}

// standardPlan returns a two-module plan with a public dependency edge.
func standardPlan(t *testing.T) plan.Plan {
	t.Helper()
	return mustPlan(
		t,
		"",
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
}

// standardRequest returns a source inspection request for the standard fixture
// repository.
func standardRequest(t *testing.T) Request {
	t.Helper()

	return Request{
		Plan:          standardPlan(t),
		RepositoryDir: "/repo",
		StagingDir:    "/repo/staging",
	}
}

// requestForPlan attaches a custom plan to the standard repository fixture.
func requestForPlan(p plan.Plan) Request {
	return Request{
		Plan:          p,
		RepositoryDir: "/repo",
		StagingDir:    "/repo/staging",
	}
}

// fileEntrySpec returns a file publish entry whose source and target match.
func fileEntrySpec(path string) manifest.PublishEntrySpec {
	return manifest.PublishEntrySpec{
		Type: string(manifest.PublishEntryFile),
		From: path,
		To:   path,
	}
}

// directoryEntrySpec returns a directory publish entry whose source and target
// match.
func directoryEntrySpec(path string) manifest.PublishEntrySpec {
	return manifest.PublishEntrySpec{
		Type: string(manifest.PublishEntryDirectory),
		From: path,
		To:   path,
	}
}

// mustPublishEntry creates a publish entry through the production constructor
// so entry fixture tests keep the same validation rules as normal manifests.
func mustPublishEntry(t *testing.T, spec manifest.PublishEntrySpec) manifest.PublishEntry {
	t.Helper()

	entry, err := manifest.NewPublishEntry(spec)
	if err != nil {
		t.Fatalf("manifest.NewPublishEntry() error = %v", err)
	}

	return entry
}

// fakeGit is a deterministic Git reader used by source workflow tests.
type fakeGit struct {
	// head is returned by Head unless headErr is set.
	head portgit.CommitHash

	// headErr forces Head to fail.
	headErr error

	// branch is returned by CurrentBranch unless branchErr is set.
	branch portgit.BranchName

	// branchErr forces CurrentBranch to fail.
	branchErr error

	// status is returned by Status unless statusErr is set.
	status portgit.Status

	// statusErr forces Status to fail.
	statusErr error
}

// Head returns the configured commit hash or the injected HEAD error.
func (g fakeGit) Head(context.Context, string) (portgit.CommitHash, error) {
	return g.head, g.headErr
}

// CurrentBranch returns the configured branch or the injected branch error.
func (g fakeGit) CurrentBranch(context.Context, string) (portgit.BranchName, error) {
	return g.branch, g.branchErr
}

// Status returns a detached copy of the configured status or the injected
// status error.
func (g fakeGit) Status(context.Context, string) (portgit.Status, error) {
	return cloneStatus(g.status), g.statusErr
}

// ConfigGet is unused by source inspection and returns a stable missing value.
func (g fakeGit) ConfigGet(context.Context, string, string) (string, bool, error) {
	return "", false, nil
}

// RefExists is unused by source inspection and returns a stable false result.
func (g fakeGit) RefExists(context.Context, string, string) (bool, error) { return false, nil }

// RemoteRefExists is unused by source inspection and returns a stable false
// result.
func (g fakeGit) RemoteRefExists(context.Context, string, string, string) (bool, error) {
	return false, nil
}

// RemoteRefHash is unused by source inspection and returns a stable missing ref.
func (g fakeGit) RemoteRefHash(context.Context, string, string, string) (portgit.CommitHash, bool, error) {
	return "", false, nil
}

// CommitMessage is unused by source inspection and returns an empty message.
func (g fakeGit) CommitMessage(context.Context, string, string) (string, error) { return "", nil }

// fakeFS is an in-memory filesystem reader/hasher with optional injected
// operation failures.
type fakeFS struct {
	// dirs stores known directory paths.
	dirs map[string]bool

	// files stores known regular-file contents.
	files map[string][]byte

	// existsErr forces Exists to fail for a specific path.
	existsErr map[string]error

	// isDirErr forces IsDir to fail for a specific path.
	isDirErr map[string]error

	// readErr forces ReadFile to fail for a specific path.
	readErr map[string]error

	// treeHashErr forces TreeHash to fail for a specific root path.
	treeHashErr map[string]error
}

// newFakeFS returns an empty fake filesystem.
func newFakeFS() *fakeFS {
	return &fakeFS{
		dirs:        map[string]bool{},
		files:       map[string][]byte{},
		existsErr:   map[string]error{},
		isDirErr:    map[string]error{},
		readErr:     map[string]error{},
		treeHashErr: map[string]error{},
	}
}

// addDir registers path as a directory in the fake filesystem.
func (fs *fakeFS) addDir(path string) { fs.dirs[path] = true }

// addFile registers path as a regular file with deterministic contents.
func (fs *fakeFS) addFile(path string, data string) { fs.files[path] = []byte(data) }

// Exists reports whether path is registered as a directory or file.
func (fs *fakeFS) Exists(_ context.Context, path string) (bool, error) {
	if err := fs.existsErr[path]; err != nil {
		return false, err
	}

	return fs.dirs[path] || fs.files[path] != nil, nil
}

// IsDir reports whether path is registered as a directory.
func (fs *fakeFS) IsDir(_ context.Context, path string) (bool, error) {
	if err := fs.isDirErr[path]; err != nil {
		return false, err
	}

	return fs.dirs[path], nil
}

// ReadFile returns a detached copy of registered file contents.
func (fs *fakeFS) ReadFile(_ context.Context, path string) ([]byte, error) {
	if err := fs.readErr[path]; err != nil {
		return nil, err
	}

	data, ok := fs.files[path]
	if !ok {
		return nil, fmt.Errorf("file %s not found", path)
	}
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

// WriteFile stores a detached file payload for tests that require the full
// filesystem port.
func (fs *fakeFS) WriteFile(
	_ context.Context,
	path string,
	data []byte,
	opts portfs.WriteFileOptions,
) error {
	if fs.files[path] != nil && !opts.Overwrite {
		return fmt.Errorf("file %s already exists", path)
	}
	if opts.CreateDirs {
		fs.addParents(filepath.Dir(path))
	}

	out := make([]byte, len(data))
	copy(out, data)
	fs.files[path] = out
	return nil
}

// MkdirAll records a directory and all missing parents.
func (fs *fakeFS) MkdirAll(_ context.Context, path string, _ portfs.MkdirOptions) error {
	fs.addParents(path)
	return nil
}

// RemoveAll removes matching fake files and directories.
func (fs *fakeFS) RemoveAll(_ context.Context, path string, _ portfs.RemoveOptions) error {
	delete(fs.dirs, path)
	delete(fs.files, path)
	prefix := path + "/"
	for dir := range fs.dirs {
		if strings.HasPrefix(dir, prefix) {
			delete(fs.dirs, dir)
		}
	}
	for file := range fs.files {
		if strings.HasPrefix(file, prefix) {
			delete(fs.files, file)
		}
	}
	return nil
}

// CleanDir removes fake directory contents while keeping the directory itself.
func (fs *fakeFS) CleanDir(_ context.Context, dir string, _ portfs.CleanDirOptions) error {
	if !fs.dirs[dir] {
		return fmt.Errorf("dir %s not found", dir)
	}
	prefix := dir + "/"
	for child := range fs.dirs {
		if strings.HasPrefix(child, prefix) {
			delete(fs.dirs, child)
		}
	}
	for file := range fs.files {
		if strings.HasPrefix(file, prefix) {
			delete(fs.files, file)
		}
	}
	return nil
}

// CopyTree mirrors fake files from src to dst.
func (fs *fakeFS) CopyTree(
	_ context.Context,
	src string,
	dst string,
	_ portfs.CopyTreeOptions,
) (portfs.CopyTreeResult, error) {
	if !fs.dirs[src] {
		return portfs.CopyTreeResult{}, fmt.Errorf("tree %s not found", src)
	}

	fs.addParents(dst)
	result := portfs.CopyTreeResult{DirectoriesCopied: 1}
	prefix := src + "/"
	for path, data := range fs.files {
		if !strings.HasPrefix(path, prefix) {
			continue
		}
		rel := strings.TrimPrefix(path, prefix)
		target := filepath.Join(dst, filepath.FromSlash(rel))
		if err := fs.WriteFile(context.Background(), target, data, portfs.WriteFileOptions{
			CreateDirs: true,
			Overwrite:  true,
		}); err != nil {
			return portfs.CopyTreeResult{}, err
		}
		result.FilesCopied++
		result.BytesCopied += int64(len(data))
	}
	return result, nil
}

// addParents registers path and its parents as directories.
func (fs *fakeFS) addParents(path string) {
	for path != "." && path != "/" && path != "" {
		fs.dirs[path] = true
		path = filepath.Dir(path)
	}
	if path == "/" {
		fs.dirs[path] = true
	}
}

// TreeHash returns a stable synthetic hash for files under root.
func (fs *fakeFS) TreeHash(
	_ context.Context,
	root string,
	_ portfs.TreeHashOptions,
) (portfs.TreeHash, error) {
	if err := fs.treeHashErr[root]; err != nil {
		return "", err
	}

	if !fs.dirs[root] {
		return "", fmt.Errorf("tree %s not found", root)
	}
	var names []string
	for path := range fs.files {
		if strings.HasPrefix(path, root+"/") {
			names = append(names, path)
		}
	}
	return portfs.TreeHash("sha256:tree:" + root + fmt.Sprintf(":%d", len(names))), nil
}

// standardFS returns a repository tree that satisfies standardPlan.
func standardFS() *fakeFS {
	fs := newFakeFS()
	fs.addDir("/repo")
	fs.addDir("/repo/staging")
	for _, name := range []string{"foundation", "control"} {
		root := "/repo/staging/src/arcoris.dev/" + name
		fs.addDir("/repo/staging/src")
		fs.addDir("/repo/staging/src/arcoris.dev")
		fs.addDir(root)
		fs.addDir(root + "/contracts")
		fs.addFile(root+"/go.mod", "module arcoris.dev/"+name+"\n")
		fs.addFile(root+"/contracts/doc.go", "package contracts\n")
	}
	return fs
}

// standardService wires source service dependencies for tests.
func standardService(fs *fakeFS, git fakeGit, opts Options) Service {
	return New(Dependencies{Git: git, FS: fs}, opts)
}

// cleanGit returns a clean checked-out source branch.
func cleanGit() fakeGit {
	return fakeGit{head: "abcdef1234567890", branch: "main", status: portgit.Status{Clean: true}}
}

// dirtyGit returns a checked-out source branch with one modified file.
func dirtyGit() fakeGit {
	git := cleanGit()
	git.status = portgit.Status{
		Clean: false,
		Entries: []portgit.StatusEntry{{
			Path: "file.go",
			Code: " M",
		}},
	}
	return git
}

// assertValidationHas checks that err is a source ValidationError containing
// code.
func assertValidationHas(t *testing.T, err error, code IssueCode) {
	t.Helper()
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T, want *ValidationError: %v", err, err)
	}
	if !ve.Has(code) {
		t.Fatalf("validation error does not contain %s: %v", code, ve)
	}
}

// assertIssuesHave checks an issue slice directly when the inspected helper
// returns raw issues instead of a ValidationError.
func assertIssuesHave(t *testing.T, issues []Issue, code IssueCode) {
	t.Helper()

	for _, issue := range issues {
		if issue.Code == code {
			return
		}
	}

	t.Fatalf("issues %v do not contain %s", issues, code)
}
