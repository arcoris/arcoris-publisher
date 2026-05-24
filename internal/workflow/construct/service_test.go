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
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
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
		`"module": "foundation"`,
		`"sourceRepository": "arcoris/arcoris"`,
		`"targetRepository": "arcoris/foundation"`,
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
