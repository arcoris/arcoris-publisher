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

package cli

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/app"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/versioning"
	"arcoris.dev/arcoris-publisher/internal/workflow"
	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("writer failed")
}

type fakeApplication struct {
	plan                  plan.Plan
	buildPlanCalled       bool
	prepareTargetsCalled  bool
	preflightCalled       bool
	verifyCalled          bool
	publishCalled         bool
	listCalled            bool
	showCalled            bool
	rollbackCalled        bool
	pruneCalled           bool
	buildPlanManifest     string
	buildPlanVersion      versioning.Version
	prepareTargetsRequest app.Request
	preflightRequest      app.Request
	verifyRequest         app.Request
	publishRequest        app.Request
	transactionRequest    app.TransactionRequest
	transactionPruneReq   app.TransactionPruneRequest
	transactionPruneRes   app.TransactionPruneResult
	buildPlanError        error
	prepareTargetsError   error
	preflightError        error
	verifyError           error
	publishError          error
	transactionError      error
}

func newFakeApplication(t *testing.T) *fakeApplication {
	t.Helper()
	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}
	return &fakeApplication{plan: p}
}

func (f *fakeApplication) BuildPlan(_ context.Context, manifestPath string, version versioning.Version) (plan.Plan, error) {
	f.buildPlanCalled = true
	f.buildPlanManifest = manifestPath
	f.buildPlanVersion = version
	return f.plan, f.buildPlanError
}

func (f *fakeApplication) PrepareTargets(_ context.Context, req app.Request) (app.Result, error) {
	f.prepareTargetsCalled = true
	f.prepareTargetsRequest = req
	return app.Result{}, f.prepareTargetsError
}

func (f *fakeApplication) Verify(_ context.Context, req app.Request) (app.Result, error) {
	f.verifyCalled = true
	f.verifyRequest = req
	return app.Result{}, f.verifyError
}

func (f *fakeApplication) Publish(_ context.Context, req app.Request) (app.Result, error) {
	f.publishCalled = true
	f.publishRequest = req
	return app.Result{}, f.publishError
}

func (f *fakeApplication) Preflight(_ context.Context, req app.Request) (app.Result, error) {
	f.preflightCalled = true
	f.preflightRequest = req
	return app.Result{}, f.preflightError
}

func (f *fakeApplication) ListTransactions(_ context.Context, req app.TransactionRequest) (app.TransactionListResult, error) {
	f.listCalled = true
	f.transactionRequest = req
	return app.TransactionListResult{}, f.transactionError
}

func (f *fakeApplication) ShowTransaction(_ context.Context, req app.TransactionRequest) (app.TransactionResult, error) {
	f.showCalled = true
	f.transactionRequest = req
	return app.TransactionResult{}, f.transactionError
}

func (f *fakeApplication) RollbackTransaction(_ context.Context, req app.TransactionRequest) (app.TransactionResult, error) {
	f.rollbackCalled = true
	f.transactionRequest = req
	return app.TransactionResult{}, f.transactionError
}

func (f *fakeApplication) PruneTransactions(_ context.Context, req app.TransactionPruneRequest) (app.TransactionPruneResult, error) {
	f.pruneCalled = true
	f.transactionPruneReq = req
	return f.transactionPruneRes, f.transactionError
}

type realApplicationOptions struct {
	tidyError error
	dirty     bool
}

func newRealApplication(t *testing.T, opts realApplicationOptions) app.App {
	t.Helper()

	fs := porttest.NewFileSystem()
	for _, module := range []struct {
		name string
		dir  string
	}{
		{name: "foundation", dir: "contracts"},
		{name: "control", dir: "runtime"},
	} {
		fs.AddFile(
			"/repo/staging/src/arcoris.dev/"+module.name+"/go.mod",
			[]byte("module arcoris.dev/"+module.name+"\n"),
		)
		fs.AddFile(
			"/repo/staging/src/arcoris.dev/"+module.name+"/"+module.dir+"/doc.go",
			[]byte("package "+module.dir+"\n"),
		)
	}

	fakeGit := porttest.NewGit()
	for _, dir := range []string{"/target/arcoris__foundation", "/target/arcoris__control"} {
		fakeGit.ConfigValues[dir+"\x00user.name"] = "ARCORIS Test"
		fakeGit.ConfigValues[dir+"\x00user.email"] = "arcoris-test@example.invalid"
	}
	if opts.dirty {
		for _, dir := range []string{"/target/arcoris__foundation", "/target/arcoris__control"} {
			fakeGit.Statuses[dir] = git.Status{
				Clean:   false,
				Entries: []git.StatusEntry{{Path: "go.mod", Code: " M"}},
			}
		}
	}

	return app.New(
		app.Dependencies{
			Workflow: workflow.Dependencies{
				Source:     source.Dependencies{FS: fs, Git: fakeGit},
				Target:     target.Dependencies{FS: fs, Git: fakeGit},
				Construct:  construct.Dependencies{FS: fs},
				ModuleFile: modulefile.Dependencies{FS: fs},
				Verify: verify.Dependencies{
					FS: fs,
					Go: porttest.GoToolchain{ModTidyError: opts.tidyError},
				},
				Publish: publish.Dependencies{Git: fakeGit},
			},
		},
		app.Options{
			Workflow: workflow.Options{
				Target: target.Options{
					CreateMissing: true,
					RequireClean:  false,
				},
				Construct: construct.Options{PreserveGitDir: true},
				Verify:    verify.Options{RequireClean: false},
				Publish: publish.Options{
					StateDir:          t.TempDir(),
					TransactionIDFunc: func(publish.TransactionInput) publish.TransactionID { return "tx-cli" },
				},
			},
		},
	)
}

func workflowCommandArgs(command string) []string {
	return []string{
		command,
		"--manifest", "../config/testdata/minimal/arcpub.yaml",
		"--version", "v0.3.0",
		"--source-repo", "/repo",
		"--staging-dir", "/repo/staging",
		"--target-root", "/target",
	}
}
