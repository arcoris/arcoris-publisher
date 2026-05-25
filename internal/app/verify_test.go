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

package app

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/versioning"
	"arcoris.dev/arcoris-publisher/internal/workflow"
	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/preflight"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

func TestVerifyRunsWorkflow(t *testing.T) {
	app, _ := appFixture(t)

	result, err := app.Verify(context.Background(), appRequest())

	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if result.Workflow().Verify().Failed() {
		t.Fatal("verification failed")
	}
	if result.Workflow().Publish().Published() {
		t.Fatal("verify use case published")
	}
}

func appFixture(t *testing.T) (App, *porttest.Git) {
	t.Helper()

	fs := appFS()
	fakeGit := porttest.NewGit()
	fakeGit.Refs["/target/arcoris__foundation\x00refs/heads/main"] = true
	fakeGit.Refs["/target/arcoris__control\x00refs/heads/stable"] = true
	fakeGit.RemoteURLs["/target/arcoris__foundation\x00origin"] = "file:///remotes/foundation.git"
	fakeGit.RemoteURLs["/target/arcoris__control\x00origin"] = "file:///remotes/control.git"
	fakeGit.ConfigValues["/target/arcoris__foundation\x00user.name"] = "ARCORIS Test"
	fakeGit.ConfigValues["/target/arcoris__foundation\x00user.email"] = "arcoris-test@example.invalid"
	fakeGit.ConfigValues["/target/arcoris__control\x00user.name"] = "ARCORIS Test"
	fakeGit.ConfigValues["/target/arcoris__control\x00user.email"] = "arcoris-test@example.invalid"
	deps := Dependencies{
		Workflow: workflow.Dependencies{
			Source:     source.Dependencies{FS: fs, Git: fakeGit},
			Target:     target.Dependencies{FS: fs, Git: fakeGit},
			Construct:  construct.Dependencies{FS: fs},
			ModuleFile: modulefile.Dependencies{FS: fs},
			Verify:     verify.Dependencies{FS: fs, Go: porttest.GoToolchain{}},
			Preflight:  preflight.Dependencies{FS: fs, Git: fakeGit},
			Publish:    publish.Dependencies{Git: fakeGit},
		},
	}
	opts := Options{
		Workflow: workflow.Options{
			Target: target.Options{
				CreateMissing: true,
				RequireClean:  false,
			},
			Construct: construct.Options{PreserveGitDir: true},
			Verify:    verify.Options{RequireClean: false},
			Publish: publish.Options{
				StateDir:          t.TempDir(),
				TransactionIDFunc: func(publish.TransactionInput) publish.TransactionID { return "tx-app" },
			},
		},
	}

	return New(deps, opts), fakeGit
}

func appRequest() Request {
	return Request{
		ManifestPath:        "../config/testdata/minimal/arcpub.yaml",
		Version:             versioning.Must("v0.3.0"),
		SourceRepositoryDir: "/repo",
		StagingDir:          "/repo/staging",
		TargetRootDir:       "/target",
	}
}

func appFS() *porttest.FileSystem {
	fs := porttest.NewFileSystem()
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/foundation/go.mod",
		[]byte("module arcoris.dev/foundation\n"),
	)
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/foundation/contracts/doc.go",
		[]byte("package contracts\n"),
	)
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/control/go.mod",
		[]byte("module arcoris.dev/control\n"),
	)
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/control/runtime/doc.go",
		[]byte("package runtime\n"),
	)
	fs.AddDir("/target")
	fs.AddDir("/target/arcoris__foundation")
	fs.AddDir("/target/arcoris__control")
	return fs
}
