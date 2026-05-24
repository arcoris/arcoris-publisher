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

package provenance

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/buildinfo"
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
)

func testInput(t *testing.T) Input {
	t.Helper()

	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	mod, ok := p.ModuleByName(manifest.ModuleName("foundation"))
	if !ok {
		t.Fatal("foundation module missing from plan")
	}

	fs := porttest.NewFileSystem()
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/foundation/go.mod",
		[]byte("module arcoris.dev/foundation\n"),
	)
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/foundation/contracts/doc.go",
		[]byte("package contracts\n"),
	)

	snapshot, err := source.New(
		source.Dependencies{FS: fs, Git: porttest.NewGit()},
		source.Options{},
	).Inspect(context.Background(), source.Request{
		Plan:          p,
		RepositoryDir: "/repo",
		StagingDir:    "/repo/staging",
	})
	if err != nil {
		t.Fatalf("source.Inspect() error = %v", err)
	}

	return Input{
		Plan:         p,
		Module:       mod,
		Source:       snapshot,
		SourceModule: sourceModuleByName(t, snapshot, mod.Name()),
		Build:        testBuildInfo(t),
	}
}

func sourceModuleByName(
	t *testing.T,
	snapshot source.Snapshot,
	name manifest.ModuleName,
) source.ModuleSnapshot {
	t.Helper()

	for _, module := range snapshot.Modules() {
		if module.Name() == name {
			return module
		}
	}
	t.Fatalf("source module %q missing", name)
	return source.ModuleSnapshot{}
}

func testBuildInfo(t *testing.T) buildinfo.Info {
	t.Helper()

	oldVersion := buildinfo.Version
	oldCommit := buildinfo.Commit
	oldDate := buildinfo.Date
	oldDirty := buildinfo.Dirty

	buildinfo.Version = "v9.8.7"
	buildinfo.Commit = "feedface"
	buildinfo.Date = "2026-05-24T12:00:00Z"
	buildinfo.Dirty = "false"

	t.Cleanup(func() {
		buildinfo.Version = oldVersion
		buildinfo.Commit = oldCommit
		buildinfo.Date = oldDate
		buildinfo.Dirty = oldDirty
	})

	return buildinfo.Current()
}
