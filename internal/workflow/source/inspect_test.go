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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	portgit "arcoris.dev/arcoris-publisher/internal/ports/git"
)

func TestInspectBuildsSourceSnapshot(t *testing.T) {
	svc := standardService(standardFS(), cleanGit(), DefaultOptions())

	snap, err := svc.Inspect(context.Background(), standardRequest(t))

	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}

	repository := snap.Repository()
	if repository.Head() != "abcdef1234567890" {
		t.Fatalf("Head() = %q", repository.Head())
	}
	if repository.Branch() != "main" {
		t.Fatalf("Branch() = %q", repository.Branch())
	}

	names := snap.ModuleNames()
	if len(names) != 2 || names[0] != "foundation" || names[1] != "control" {
		t.Fatalf("ModuleNames() = %v", names)
	}

	modules := snap.Modules()
	if len(modules[0].Entries()) != 2 {
		t.Fatalf("Entries() len = %d", len(modules[0].Entries()))
	}
	if modules[0].Hash().IsZero() {
		t.Fatal("Hash() is zero")
	}
}

func TestInspectRejectsEmptyPlanAndBlankPaths(t *testing.T) {
	svc := standardService(standardFS(), cleanGit(), DefaultOptions())

	_, err := svc.Inspect(context.Background(), Request{})

	assertValidationHas(t, err, IssueInvalidRequest)
}

func TestInspectRejectsMissingRootDirectories(t *testing.T) {
	t.Run("repository", func(t *testing.T) {
		fs := standardFS()
		delete(fs.dirs, "/repo")
		svc := standardService(fs, cleanGit(), DefaultOptions())

		_, err := svc.Inspect(context.Background(), standardRequest(t))

		assertValidationHas(t, err, IssueRepositoryMissing)
	})

	t.Run("staging", func(t *testing.T) {
		fs := standardFS()
		delete(fs.dirs, "/repo/staging")
		svc := standardService(fs, cleanGit(), DefaultOptions())

		_, err := svc.Inspect(context.Background(), standardRequest(t))

		assertValidationHas(t, err, IssueStagingMissing)
	})
}

func TestInspectRejectsRootPathsThatAreNotDirectories(t *testing.T) {
	fs := standardFS()
	delete(fs.dirs, "/repo")
	fs.addFile("/repo", "not a directory")
	delete(fs.dirs, "/repo/staging")
	fs.addFile("/repo/staging", "not a directory")
	svc := standardService(fs, cleanGit(), DefaultOptions())

	_, err := svc.Inspect(context.Background(), standardRequest(t))

	assertValidationHas(t, err, IssueRepositoryNotDirectory)
	assertValidationHas(t, err, IssueStagingNotDirectory)
}

func TestInspectWrapsFilesystemRootErrors(t *testing.T) {
	fs := standardFS()
	fs.existsErr["/repo"] = fmt.Errorf("exists failed")
	fs.isDirErr["/repo/staging"] = fmt.Errorf("isdir failed")
	svc := standardService(fs, cleanGit(), DefaultOptions())

	_, err := svc.Inspect(context.Background(), standardRequest(t))

	assertValidationHas(t, err, IssueInvalidRequest)
}

func TestInspectRejectsStagingOutsideRepository(t *testing.T) {
	fs := standardFS()
	fs.addDir("/elsewhere")
	svc := standardService(fs, cleanGit(), DefaultOptions())
	request := Request{
		Plan:          standardPlan(t),
		RepositoryDir: "/repo",
		StagingDir:    "/elsewhere",
	}

	_, err := svc.Inspect(context.Background(), request)

	assertValidationHas(t, err, IssueStagingOutsideRepo)
}

func TestInspectReportsGitReadErrors(t *testing.T) {
	cases := []struct {
		name string
		git  fakeGit
	}{
		{
			name: "head",
			git:  fakeGit{headErr: fmt.Errorf("head failed")},
		},
		{
			name: "branch",
			git: fakeGit{
				head:      "abcdef",
				branchErr: fmt.Errorf("branch failed"),
			},
		},
		{
			name: "status",
			git: fakeGit{
				head:      "abcdef",
				branch:    "main",
				statusErr: fmt.Errorf("status failed"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := standardService(standardFS(), tc.git, DefaultOptions())

			_, err := svc.Inspect(context.Background(), standardRequest(t))

			if err == nil {
				t.Fatal("Inspect() error = nil")
			}
		})
	}
}

func TestInspectDetachedHeadPolicy(t *testing.T) {
	t.Run("reject-by-default", func(t *testing.T) {
		git := cleanGit()
		git.branch = ""
		svc := standardService(standardFS(), git, DefaultOptions())

		_, err := svc.Inspect(context.Background(), standardRequest(t))

		assertValidationHas(t, err, IssueDetachedHead)
	})

	t.Run("allow-when-configured", func(t *testing.T) {
		git := cleanGit()
		git.branch = ""
		opts := Options{AllowDetachedHEAD: true, DisableHashes: true}
		svc := standardService(standardFS(), git, opts)

		snap, err := svc.Inspect(context.Background(), standardRequest(t))

		if err != nil {
			t.Fatalf("Inspect() error = %v", err)
		}
		if snap.Repository().Branch() != "" {
			t.Fatalf("Branch() = %q", snap.Repository().Branch())
		}
	})
}

func TestInspectDirtySourcePolicy(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		svc := standardService(standardFS(), dirtyGit(), DefaultOptions())

		_, err := svc.Inspect(context.Background(), standardRequest(t))

		assertValidationHas(t, err, IssueDirtySource)
	})

	t.Run("warn", func(t *testing.T) {
		svc := standardService(standardFS(), dirtyGit(), DefaultOptions())
		request := Request{
			Plan: mustPlan(
				t,
				string(manifest.DirtyPolicyWarn),
				testModule{name: "foundation"},
			),
			RepositoryDir: "/repo",
			StagingDir:    "/repo/staging",
		}

		snap, err := svc.Inspect(context.Background(), request)

		if err != nil {
			t.Fatalf("Inspect() error = %v", err)
		}
		if len(snap.Warnings()) != 1 || snap.Warnings()[0].Code != IssueDirtySource {
			t.Fatalf("Warnings() = %v", snap.Warnings())
		}
	})

	t.Run("allow", func(t *testing.T) {
		svc := standardService(standardFS(), dirtyGit(), DefaultOptions())
		request := Request{
			Plan: mustPlan(
				t,
				string(manifest.DirtyPolicyAllow),
				testModule{name: "foundation"},
			),
			RepositoryDir: "/repo",
			StagingDir:    "/repo/staging",
		}

		snap, err := svc.Inspect(context.Background(), request)

		if err != nil {
			t.Fatalf("Inspect() error = %v", err)
		}
		if len(snap.Warnings()) != 0 {
			t.Fatalf("Warnings() = %v, want none", snap.Warnings())
		}
	})
}

func TestDirtyPolicyRejectsUnsupportedValue(t *testing.T) {
	inspector := inspector{}
	repository := RepositorySnapshot{
		status: portgit.Status{
			Clean: false,
			Entries: []portgit.StatusEntry{{
				Path: "file.go",
				Code: " M",
			}},
		},
	}

	err := inspector.enforceDirtyPolicy(repository)

	assertValidationHas(t, err, IssueInvalidRequest)
}

func TestInspectCanSkipHashesWhenDisabled(t *testing.T) {
	git := cleanGit()
	git.branch = ""
	opts := Options{AllowDetachedHEAD: true, DisableHashes: true}
	svc := standardService(standardFS(), git, opts)

	snap, err := svc.Inspect(context.Background(), standardRequest(t))

	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	if !snap.Modules()[0].Hash().IsZero() {
		t.Fatal("Hash() is not zero")
	}
}

func TestInspectRejectsModuleSourceProblems(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		fs := standardFS()
		delete(fs.dirs, "/repo/staging/src/arcoris.dev/foundation")
		svc := standardService(fs, cleanGit(), DefaultOptions())

		_, err := svc.Inspect(context.Background(), standardRequest(t))

		assertValidationHas(t, err, IssueModuleSourceMissing)
	})

	t.Run("not-directory", func(t *testing.T) {
		fs := standardFS()
		delete(fs.dirs, "/repo/staging/src/arcoris.dev/foundation")
		fs.addFile("/repo/staging/src/arcoris.dev/foundation", "not a directory")
		svc := standardService(fs, cleanGit(), DefaultOptions())

		_, err := svc.Inspect(context.Background(), standardRequest(t))

		assertValidationHas(t, err, IssueModuleSourceNotDir)
	})
}

func TestValidateModuleRootsReportsEscapes(t *testing.T) {
	var issues issueCollector
	inspector := inspector{
		deps:       Dependencies{FS: standardFS()},
		stagingDir: "/repo/staging",
	}

	inspector.validateModuleRoots(
		context.Background(),
		&issues,
		manifest.ModuleName("foundation"),
		"modules[0]",
		"/repo/outside/foundation",
		"/repo/outside/foundation",
	)

	err := issues.err()
	assertValidationHas(t, err, IssueEntryPathEscape)
}

func TestInspectAllowsMissingOptionalEntry(t *testing.T) {
	optional := true
	entries := []manifest.PublishEntrySpec{
		fileEntrySpec("go.mod"),
		{
			Type:     string(manifest.PublishEntryFile),
			From:     "go.sum",
			To:       "go.sum",
			Optional: &optional,
		},
	}
	p := mustPlan(t, "", testModule{name: "foundation", entries: entries})
	svc := standardService(standardFS(), cleanGit(), DefaultOptions())

	snap, err := svc.Inspect(context.Background(), requestForPlan(p))

	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}

	entriesSnapshot := snap.Modules()[0].Entries()
	if len(entriesSnapshot) != 2 {
		t.Fatalf("Entries() len = %d", len(entriesSnapshot))
	}
	if entriesSnapshot[1].Present() {
		t.Fatal("missing optional entry is present")
	}
}

func TestInspectRejectsMissingRequiredEntry(t *testing.T) {
	p := mustPlan(
		t,
		"",
		testModule{name: "foundation", entries: []manifest.PublishEntrySpec{
			fileEntrySpec("missing.go"),
		}},
	)
	svc := standardService(standardFS(), cleanGit(), DefaultOptions())

	_, err := svc.Inspect(context.Background(), requestForPlan(p))

	assertValidationHas(t, err, IssueEntryMissing)
}

func TestInspectRejectsEntryTypeMismatches(t *testing.T) {
	cases := []struct {
		name string
		spec manifest.PublishEntrySpec
	}{
		{
			name: "file-entry-points-to-directory",
			spec: fileEntrySpec("contracts"),
		},
		{
			name: "directory-entry-points-to-file",
			spec: directoryEntrySpec("go.mod"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := mustPlan(
				t,
				"",
				testModule{name: "foundation", entries: []manifest.PublishEntrySpec{
					tc.spec,
				}},
			)
			svc := standardService(standardFS(), cleanGit(), DefaultOptions())

			_, err := svc.Inspect(context.Background(), requestForPlan(p))

			assertValidationHas(t, err, IssueEntryTypeMismatch)
		})
	}
}

func TestInspectRejectsEntryReadAndHashFailures(t *testing.T) {
	t.Run("file-read", func(t *testing.T) {
		fs := standardFS()
		fs.readErr["/repo/staging/src/arcoris.dev/foundation/go.mod"] = fmt.Errorf("read failed")
		svc := standardService(fs, cleanGit(), DefaultOptions())

		_, err := svc.Inspect(context.Background(), standardRequest(t))

		assertValidationHas(t, err, IssueEntryHashFailed)
	})

	t.Run("tree-hash", func(t *testing.T) {
		fs := standardFS()
		fs.treeHashErr["/repo/staging/src/arcoris.dev/foundation/contracts"] = fmt.Errorf("tree failed")
		svc := standardService(fs, cleanGit(), DefaultOptions())

		_, err := svc.Inspect(context.Background(), standardRequest(t))

		assertValidationHas(t, err, IssueEntryHashFailed)
	})
}

func TestInspectEntryReportsFilesystemErrors(t *testing.T) {
	entry := mustPublishEntry(t, fileEntrySpec("go.mod"))

	t.Run("exists", func(t *testing.T) {
		fs := standardFS()
		fs.existsErr["/repo/staging/src/arcoris.dev/foundation/go.mod"] = fmt.Errorf("exists failed")
		inspector := inspector{deps: Dependencies{FS: fs}}

		_, issues := inspector.inspectEntry(
			context.Background(),
			manifest.ModuleName("foundation"),
			"modules[0].publish.entries[0]",
			"/repo/staging/src/arcoris.dev/foundation",
			entry,
		)

		assertIssuesHave(t, issues, IssueInvalidRequest)
	})

	t.Run("is-dir", func(t *testing.T) {
		fs := standardFS()
		fs.isDirErr["/repo/staging/src/arcoris.dev/foundation/go.mod"] = fmt.Errorf("isdir failed")
		inspector := inspector{deps: Dependencies{FS: fs}}

		_, issues := inspector.inspectEntry(
			context.Background(),
			manifest.ModuleName("foundation"),
			"modules[0].publish.entries[0]",
			"/repo/staging/src/arcoris.dev/foundation",
			entry,
		)

		assertIssuesHave(t, issues, IssueInvalidRequest)
	})
}
