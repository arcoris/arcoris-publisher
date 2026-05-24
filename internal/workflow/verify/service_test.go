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

package verify

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

func TestVerifyRejectsInvalidRequest(t *testing.T) {
	_, err := New(Dependencies{}, Options{}).Verify(context.Background(), Request{})

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodeInvalidRequest {
		t.Fatalf("Code = %q", got.Code)
	}
}

func TestTargetWorktreeCheck(t *testing.T) {
	tests := []struct {
		name   string
		fs     fakeReader
		status Status
	}{
		{
			name:   "missing",
			fs:     fakeReader{},
			status: StatusFailed,
		},
		{
			name: "not directory",
			fs: fakeReader{
				paths: map[string]fakePath{"/target": {exists: true}},
			},
			status: StatusFailed,
		},
		{
			name: "directory",
			fs: fakeReader{
				paths: map[string]fakePath{"/target": {exists: true, dir: true}},
			},
			status: StatusPassed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(Dependencies{FS: tt.fs}, Options{})

			check := service.targetWorktreeCheck(context.Background(), "/target")

			if check.Status() != tt.status {
				t.Fatalf("Status() = %s, want %s", check.Status(), tt.status)
			}
			if check.Path() != "/target" {
				t.Fatalf("Path() = %q", check.Path())
			}
		})
	}
}

func TestGoModTidyPassesWhenFilesAreUnchanged(t *testing.T) {
	req, fs, _ := verifyRequest(t)
	result, err := New(
		Dependencies{FS: fs, Go: porttest.GoToolchain{}},
		Options{},
	).Verify(context.Background(), req)

	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	assertCheckStatus(t, result, "go-mod-tidy", StatusPassed)
}

func TestGoModTidyFailsWhenGoModChangesWithoutGit(t *testing.T) {
	req, fs, moduleRoot := verifyRequest(t)
	goTool := porttest.GoToolchain{
		ModTidyHook: func(ctx context.Context, dir string) error {
			return fs.WriteFile(
				ctx,
				dir+"/go.mod",
				[]byte("module arcoris.dev/foundation\n\nrequire example.com/new v1.0.0\n"),
				filesystem.WriteFileOptions{Overwrite: true},
			)
		},
	}

	result, err := New(
		Dependencies{FS: fs, Go: goTool},
		Options{},
	).Verify(context.Background(), req)

	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if moduleRoot == "" {
		t.Fatal("moduleRoot is empty")
	}
	assertCheckStatus(t, result, "go-mod-tidy", StatusFailed)
	if !result.Failed() {
		t.Fatal("Failed() = false")
	}
}

func TestGoModTidyFailsWhenGoSumChangesWithoutGit(t *testing.T) {
	req, fs, _ := verifyRequest(t)
	goTool := porttest.GoToolchain{
		ModTidyHook: func(ctx context.Context, dir string) error {
			return fs.WriteFile(
				ctx,
				dir+"/go.sum",
				[]byte("example.com/new v1.0.0 h1:test\n"),
				filesystem.WriteFileOptions{CreateDirs: true, Overwrite: true},
			)
		},
	}

	result, err := New(
		Dependencies{FS: fs, Go: goTool},
		Options{},
	).Verify(context.Background(), req)

	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	assertCheckStatus(t, result, "go-mod-tidy", StatusFailed)
}

func TestGoModTidyFailureIsVerificationFailure(t *testing.T) {
	req, fs, _ := verifyRequest(t)
	goTool := porttest.GoToolchain{ModTidyError: errors.New("tidy failed")}

	result, err := New(
		Dependencies{FS: fs, Go: goTool},
		Options{},
	).Verify(context.Background(), req)

	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	assertCheckStatus(t, result, "go-mod-tidy", StatusFailed)
	if !result.Failed() {
		t.Fatal("Failed() = false")
	}
}

func verifyRequest(t *testing.T) (Request, *porttest.FileSystem, string) {
	t.Helper()

	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fs := porttest.NewFileSystem()
	fakeGit := porttest.NewGit()
	targets, err := target.New(
		target.Dependencies{FS: fs, Git: fakeGit},
		target.Options{CreateMissing: true},
	).Prepare(context.Background(), target.Request{
		Plan:    p,
		RootDir: "/target",
	})
	if err != nil {
		t.Fatalf("target.Prepare() error = %v", err)
	}

	ws, ok := targets.WorkspaceByModule("foundation")
	if !ok {
		t.Fatal("workspace for foundation not found")
	}

	moduleRoot := ws.WorktreeDir()
	fs.AddFile(moduleRoot+"/go.mod", []byte("module arcoris.dev/foundation\n"))

	return Request{Plan: p, Targets: targets}, fs, moduleRoot
}

func assertCheckStatus(t *testing.T, result Result, checkName CheckName, status Status) {
	t.Helper()

	for _, module := range result.Modules() {
		for _, check := range module.Checks() {
			if check.Name() == checkName {
				if check.Status() != status {
					t.Fatalf("%s status = %s, want %s", checkName, check.Status(), status)
				}
				return
			}
		}
	}

	t.Fatalf("check %q not found in %#v", checkName, result.Modules())
}

type fakePath struct {
	exists bool
	dir    bool
	data   []byte
}

type fakeReader struct {
	paths map[string]fakePath
}

func (fs fakeReader) Exists(_ context.Context, name string) (bool, error) {
	path := fs.paths[name]
	return path.exists, nil
}

func (fs fakeReader) IsDir(_ context.Context, name string) (bool, error) {
	path := fs.paths[name]
	return path.dir, nil
}

func (fs fakeReader) ReadFile(_ context.Context, path string) ([]byte, error) {
	data := fs.paths[path].data
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}
