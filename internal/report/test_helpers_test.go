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

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/workflow"
	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

type workflowReportFixture struct {
	publish      bool
	dirtyModules []string
	verifyFails  bool
}

func reportPlan(t *testing.T) workflow.Request {
	t.Helper()

	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
		publishertest.Module{Name: "control", Dependencies: []string{"foundation"}},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	return workflow.Request{
		Plan:                p,
		SourceRepositoryDir: "/repo",
		StagingDir:          "/repo/staging",
		TargetRootDir:       "/target",
	}
}

func reportWorkflowResult(t *testing.T, cfg workflowReportFixture) workflow.Result {
	t.Helper()

	req := reportPlan(t)
	req.Publish = cfg.publish

	fs := reportWorkflowFS()
	fakeGit := porttest.NewGit()
	for _, name := range cfg.dirtyModules {
		fakeGit.Statuses[reportWorktreeDir(name)] = reportDirtyStatus()
	}

	goToolchain := porttest.GoToolchain{}
	if cfg.verifyFails {
		goToolchain.ModTidyError = errors.New("tidy failed")
	}

	runner := workflow.New(
		workflow.Dependencies{
			Source:     source.Dependencies{FS: fs, Git: fakeGit},
			Target:     target.Dependencies{FS: fs, Git: fakeGit},
			Construct:  construct.Dependencies{FS: fs},
			ModuleFile: modulefile.Dependencies{FS: fs},
			Verify:     verify.Dependencies{FS: fs, Go: goToolchain},
			Publish:    publish.Dependencies{Git: fakeGit},
		},
		workflow.Options{
			Target: target.Options{
				CreateMissing: true,
				RequireClean:  false,
			},
			Construct: construct.Options{PreserveGitDir: true},
			Verify:    verify.Options{RequireClean: false},
		},
	)

	result, err := runner.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("workflow.Run() error = %v", err)
	}

	return result
}

func reportWorkflowFS() *porttest.FileSystem {
	fs := porttest.NewFileSystem()
	for _, name := range []string{"foundation", "control"} {
		fs.AddFile(
			"/repo/staging/src/arcoris.dev/"+name+"/go.mod",
			[]byte("module arcoris.dev/"+name+"\n"),
		)
		fs.AddFile(
			"/repo/staging/src/arcoris.dev/"+name+"/contracts/doc.go",
			[]byte("package contracts\n"),
		)
	}
	return fs
}

func reportWorktreeDir(name string) string {
	return "/target/arcoris__" + name
}

func reportDirtyStatus() git.Status {
	return git.Status{
		Clean:   false,
		Entries: []git.StatusEntry{{Path: "go.mod", Code: " M"}},
	}
}

func renderReport(t *testing.T, fn func(*bytes.Buffer) error) string {
	t.Helper()

	var buf bytes.Buffer
	if err := fn(&buf); err != nil {
		t.Fatalf("render error = %v", err)
	}

	return buf.String()
}

func assertNoLocalPaths(t *testing.T, output string) {
	t.Helper()

	for _, path := range []string{"/repo", "/target", "/tmp", "/home/user"} {
		if strings.Contains(output, path) {
			t.Fatalf("report leaked local path %q:\n%s", path, output)
		}
	}
}

func assertContains(t *testing.T, output string, values ...string) {
	t.Helper()

	for _, value := range values {
		if !strings.Contains(output, value) {
			t.Fatalf("output missing %q:\n%s", value, output)
		}
	}
}
