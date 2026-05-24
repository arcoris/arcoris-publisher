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

package construct

import (
	"context"
	"errors"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/buildinfo"
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/provenance"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

func TestConstructRejectsInvalidRequest(t *testing.T) {
	_, err := New(Dependencies{}, Options{}).Construct(context.Background(), Request{})

	validation, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if !validation.Has(IssueInvalidRequest) {
		t.Fatalf("validation issues = %v", validation.Issues)
	}
}

func TestConstructWritesStableProvenanceWithoutLocalPaths(t *testing.T) {
	oldVersion := buildinfo.Version
	buildinfo.Version = "v1.2.3"
	t.Cleanup(func() { buildinfo.Version = oldVersion })

	provenanceFile := "ARCPUB.json"
	p, err := publishertest.Plan(
		publishertest.PlanOptions{
			Publish: manifest.PublishSpec{
				Provenance: manifest.ProvenanceSpec{File: &provenanceFile},
			},
		},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fs := constructFixtureFS()
	git := porttest.NewGit()
	snapshot, err := source.New(
		source.Dependencies{FS: fs, Git: git},
		source.Options{},
	).Inspect(context.Background(), source.Request{
		Plan:          p,
		RepositoryDir: "/repo",
		StagingDir:    "/repo/staging",
	})
	if err != nil {
		t.Fatalf("source.Inspect() error = %v", err)
	}
	targets, err := target.New(
		target.Dependencies{FS: fs, Git: git},
		target.Options{CreateMissing: true},
	).Prepare(context.Background(), target.Request{
		Plan:    p,
		RootDir: "/target",
	})
	if err != nil {
		t.Fatalf("target.Prepare() error = %v", err)
	}

	result, err := New(
		Dependencies{FS: fs},
		Options{PreserveGitDir: true, GenerateProvenanceFile: true},
	).Construct(context.Background(), Request{
		Plan:    p,
		Source:  snapshot,
		Targets: targets,
	})
	if err != nil {
		t.Fatalf("Construct() error = %v", err)
	}

	module := result.Modules()[0]
	data, err := fs.ReadFile(context.Background(), module.WorktreeDir()+"/ARCPUB.json")
	if err != nil {
		t.Fatalf("ReadFile(provenance) error = %v", err)
	}
	text := string(data)
	for _, required := range []string{
		`"schemaVersion": "arcoris.provenance/v1"`,
		`"version": "v1.2.3"`,
		`"name": "foundation"`,
		`"repository": "arcoris/arcoris"`,
		`"repository": "arcoris/foundation"`,
		`"publishMode": "explicit-projection"`,
		`"projectionHash": "sha256:`,
	} {
		if !strings.Contains(text, required) {
			t.Fatalf("provenance missing %q:\n%s", required, text)
		}
	}
	if strings.Contains(text, "/repo") || strings.Contains(text, "/target") {
		t.Fatalf("provenance leaks local paths:\n%s", text)
	}
}

func TestAppendProvenanceFileFailsOnRenderError(t *testing.T) {
	req, module, fs := constructProvenanceContext(t, "ARCPUB.json")
	oldRender := renderProvenanceFilePayload
	renderProvenanceFilePayload = func(provenance.Input) ([]byte, error) {
		return nil, errors.New("render failed")
	}
	t.Cleanup(func() { renderProvenanceFilePayload = oldRender })

	var issues issueCollector
	operations := []Operation{newOperation(OperationClean, "", module.workspace.WorktreeDir())}
	ok := New(
		Dependencies{FS: fs},
		Options{GenerateProvenanceFile: true},
	).appendProvenanceFile(context.Background(), req, module, &operations, &issues)

	if ok {
		t.Fatal("appendProvenanceFile() = true")
	}
	assertConstructIssue(t, issues.Err(), IssueEntryCopyFailed)
	assertNoGeneratedOperation(t, operations)
}

func TestAppendProvenanceFileFailsOnWriteError(t *testing.T) {
	req, module, fs := constructProvenanceContext(t, "ARCPUB.json")
	failingFS := failingWriteFS{FileSystem: fs, failSuffix: "ARCPUB.json"}

	var issues issueCollector
	operations := []Operation{newOperation(OperationClean, "", module.workspace.WorktreeDir())}
	ok := New(
		Dependencies{FS: failingFS},
		Options{GenerateProvenanceFile: true},
	).appendProvenanceFile(context.Background(), req, module, &operations, &issues)

	if ok {
		t.Fatal("appendProvenanceFile() = true")
	}
	assertConstructIssue(t, issues.Err(), IssueEntryCopyFailed)
	assertNoGeneratedOperation(t, operations)
	if _, err := fs.ReadFile(context.Background(), module.workspace.WorktreeDir()+"/ARCPUB.json"); err == nil {
		t.Fatal("provenance file was written despite write error")
	}
}

func TestConstructFailsWhenProvenanceWriteFails(t *testing.T) {
	req, _, fs := constructProvenanceContext(t, "ARCPUB.json")
	failingFS := failingWriteFS{FileSystem: fs, failSuffix: "ARCPUB.json"}

	result, err := New(
		Dependencies{FS: failingFS},
		Options{GenerateProvenanceFile: true},
	).Construct(context.Background(), req)

	if err == nil {
		t.Fatal("Construct() error = nil")
	}
	if len(result.Modules()) != 0 {
		t.Fatalf("Construct() reported modules after write failure: %+v", result.Modules())
	}
	assertConstructIssue(t, err, IssueEntryCopyFailed)
}

func TestAppendProvenanceFileAppendsGeneratedOperationOnSuccess(t *testing.T) {
	req, module, fs := constructProvenanceContext(t, "ARCPUB.json")

	var issues issueCollector
	operations := []Operation{newOperation(OperationClean, "", module.workspace.WorktreeDir())}
	ok := New(
		Dependencies{FS: fs},
		Options{GenerateProvenanceFile: true},
	).appendProvenanceFile(context.Background(), req, module, &operations, &issues)

	if !ok {
		t.Fatal("appendProvenanceFile() = false")
	}
	if err := issues.Err(); err != nil {
		t.Fatalf("issues.Err() = %v", err)
	}
	if operations[len(operations)-1].Kind() != OperationWriteGenerated {
		t.Fatalf("last operation = %q", operations[len(operations)-1].Kind())
	}
}

func TestConstructSkipsProvenanceWhenFilePolicyDisabled(t *testing.T) {
	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fs := constructFixtureFS()
	git := porttest.NewGit()
	snapshot := inspectConstructSource(t, p, fs, git)
	targets := prepareConstructTargets(t, p, fs, git)

	result, err := New(
		Dependencies{FS: fs},
		Options{PreserveGitDir: true, GenerateProvenanceFile: true},
	).Construct(context.Background(), Request{
		Plan:    p,
		Source:  snapshot,
		Targets: targets,
	})
	if err != nil {
		t.Fatalf("Construct() error = %v", err)
	}

	for _, operation := range result.Modules()[0].Operations() {
		if operation.Kind() == OperationWriteGenerated {
			t.Fatalf("unexpected generated operation: %+v", operation)
		}
	}
}

func TestConstructWritesProvenanceWhenCommitTrailersDisabled(t *testing.T) {
	commitTrailers := false
	provenanceFile := "ARCPUB.json"
	p, err := publishertest.Plan(
		publishertest.PlanOptions{
			Publish: manifest.PublishSpec{
				Provenance: manifest.ProvenanceSpec{
					CommitTrailers: &commitTrailers,
					File:           &provenanceFile,
				},
			},
		},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fs := constructFixtureFS()
	git := porttest.NewGit()
	snapshot := inspectConstructSource(t, p, fs, git)
	targets := prepareConstructTargets(t, p, fs, git)

	result, err := New(
		Dependencies{FS: fs},
		Options{PreserveGitDir: true, GenerateProvenanceFile: true},
	).Construct(context.Background(), Request{
		Plan:    p,
		Source:  snapshot,
		Targets: targets,
	})
	if err != nil {
		t.Fatalf("Construct() error = %v", err)
	}

	module := result.Modules()[0]
	if _, err := fs.ReadFile(context.Background(), module.WorktreeDir()+"/ARCPUB.json"); err != nil {
		t.Fatalf("ReadFile(provenance) error = %v", err)
	}
}

func constructFixtureFS() *porttest.FileSystem {
	fs := porttest.NewFileSystem()
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/foundation/go.mod",
		[]byte("module arcoris.dev/foundation\n"),
	)
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/foundation/contracts/doc.go",
		[]byte("package contracts\n"),
	)
	return fs
}

func constructProvenanceContext(
	t *testing.T,
	provenanceFile string,
) (Request, moduleContext, *porttest.FileSystem) {
	t.Helper()

	p, err := publishertest.Plan(
		publishertest.PlanOptions{
			Publish: manifest.PublishSpec{
				Provenance: manifest.ProvenanceSpec{File: &provenanceFile},
			},
		},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fs := constructFixtureFS()
	git := porttest.NewGit()
	snapshot := inspectConstructSource(t, p, fs, git)
	targets := prepareConstructTargets(t, p, fs, git)

	req := Request{Plan: p, Source: snapshot, Targets: targets}
	module, ok := resolveModuleContext(req, "foundation", &issueCollector{})
	if !ok {
		t.Fatal("resolveModuleContext() failed")
	}

	return req, module, fs
}

func inspectConstructSource(
	t *testing.T,
	p plan.Plan,
	fs *porttest.FileSystem,
	git *porttest.Git,
) source.Snapshot {
	t.Helper()

	snapshot, err := source.New(
		source.Dependencies{FS: fs, Git: git},
		source.Options{},
	).Inspect(context.Background(), source.Request{
		Plan:          p,
		RepositoryDir: "/repo",
		StagingDir:    "/repo/staging",
	})
	if err != nil {
		t.Fatalf("source.Inspect() error = %v", err)
	}
	return snapshot
}

func prepareConstructTargets(
	t *testing.T,
	p plan.Plan,
	fs *porttest.FileSystem,
	git *porttest.Git,
) target.WorkspaceSet {
	t.Helper()

	targets, err := target.New(
		target.Dependencies{FS: fs, Git: git},
		target.Options{CreateMissing: true},
	).Prepare(context.Background(), target.Request{
		Plan:    p,
		RootDir: "/target",
	})
	if err != nil {
		t.Fatalf("target.Prepare() error = %v", err)
	}
	return targets
}

func assertConstructIssue(t *testing.T, err error, code IssueCode) {
	t.Helper()

	validation, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if !validation.Has(code) {
		t.Fatalf("validation issues = %v", validation.Issues)
	}
}

func assertNoGeneratedOperation(t *testing.T, operations []Operation) {
	t.Helper()

	for _, operation := range operations {
		if operation.Kind() == OperationWriteGenerated {
			t.Fatalf("unexpected generated operation: %+v", operation)
		}
	}
}

type failingWriteFS struct {
	*porttest.FileSystem
	failSuffix string
}

func (fs failingWriteFS) WriteFile(
	ctx context.Context,
	path string,
	data []byte,
	opts filesystem.WriteFileOptions,
) error {
	if strings.HasSuffix(path, fs.failSuffix) {
		return errors.New("write failed")
	}
	return fs.FileSystem.WriteFile(ctx, path, data, opts)
}
