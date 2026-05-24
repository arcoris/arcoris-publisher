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

package buildinfo

import "testing"

func TestCurrentReturnsDefaultBuildInfo(t *testing.T) {
	restoreBuildVariables(t)

	info := Current()

	if info.Version() != "dev" {
		t.Fatalf("Version() = %q", info.Version())
	}
	if info.Commit() != "unknown" {
		t.Fatalf("Commit() = %q", info.Commit())
	}
	if info.Date() != "unknown" {
		t.Fatalf("Date() = %q", info.Date())
	}
	if info.Dirty() != "unknown" {
		t.Fatalf("Dirty() = %q", info.Dirty())
	}
	if !info.IsDev() {
		t.Fatal("IsDev() = false")
	}
}

func TestCurrentNormalizesEmptyBuildVariables(t *testing.T) {
	restoreBuildVariables(t)
	Version = ""
	Commit = ""
	Date = ""
	Dirty = ""

	info := Current()

	if info.Version() != "dev" {
		t.Fatalf("Version() = %q", info.Version())
	}
	if info.Commit() != "unknown" {
		t.Fatalf("Commit() = %q", info.Commit())
	}
	if info.Date() != "unknown" {
		t.Fatalf("Date() = %q", info.Date())
	}
	if info.Dirty() != "unknown" {
		t.Fatalf("Dirty() = %q", info.Dirty())
	}
}

func TestCurrentNormalizesWhitespaceBuildVariables(t *testing.T) {
	restoreBuildVariables(t)
	Version = "   "
	Commit = "\n"
	Date = "\t"
	Dirty = "\r\n"

	info := Current()

	if info.Version() != "dev" {
		t.Fatalf("Version() = %q", info.Version())
	}
	if info.Commit() != "unknown" {
		t.Fatalf("Commit() = %q", info.Commit())
	}
	if info.Date() != "unknown" {
		t.Fatalf("Date() = %q", info.Date())
	}
	if info.Dirty() != "unknown" {
		t.Fatalf("Dirty() = %q", info.Dirty())
	}
}

func TestCurrentReturnsCustomBuildInfo(t *testing.T) {
	restoreBuildVariables(t)
	Version = "v1.2.3"
	Commit = "abc123"
	Date = "2026-05-24T00:00:00Z"
	Dirty = "false"

	info := Current()

	if info.Version() != "v1.2.3" {
		t.Fatalf("Version() = %q", info.Version())
	}
	if info.Commit() != "abc123" {
		t.Fatalf("Commit() = %q", info.Commit())
	}
	if info.Date() != "2026-05-24T00:00:00Z" {
		t.Fatalf("Date() = %q", info.Date())
	}
	if info.Dirty() != "false" {
		t.Fatalf("Dirty() = %q", info.Dirty())
	}
	if info.IsDev() {
		t.Fatal("IsDev() = true")
	}
}

func TestCurrentTrimsCustomBuildInfo(t *testing.T) {
	restoreBuildVariables(t)
	Version = " v1.2.3 "
	Commit = "\tabc123\n"
	Date = " 2026-05-24T00:00:00Z "
	Dirty = " false "

	info := Current()

	if info.Version() != "v1.2.3" {
		t.Fatalf("Version() = %q", info.Version())
	}
	if info.Commit() != "abc123" {
		t.Fatalf("Commit() = %q", info.Commit())
	}
	if info.Date() != "2026-05-24T00:00:00Z" {
		t.Fatalf("Date() = %q", info.Date())
	}
	if info.Dirty() != "false" {
		t.Fatalf("Dirty() = %q", info.Dirty())
	}
}

func TestInfoZeroValueUsesDefaults(t *testing.T) {
	var info Info

	if info.Version() != "dev" {
		t.Fatalf("Version() = %q", info.Version())
	}
	if info.Commit() != "unknown" {
		t.Fatalf("Commit() = %q", info.Commit())
	}
	if info.Date() != "unknown" {
		t.Fatalf("Date() = %q", info.Date())
	}
	if info.Dirty() != "unknown" {
		t.Fatalf("Dirty() = %q", info.Dirty())
	}
}

func TestInfoMapReturnsDetachedValues(t *testing.T) {
	restoreBuildVariables(t)
	Version = "v1.2.3"

	info := Current()
	first := info.Map()
	first["version"] = "changed"

	second := info.Map()
	if second["version"] != "v1.2.3" {
		t.Fatalf("Map() was not detached: %v", second)
	}
}

func restoreBuildVariables(t *testing.T) {
	t.Helper()

	oldVersion := Version
	oldCommit := Commit
	oldDate := Date
	oldDirty := Dirty

	Version = defaultVersion
	Commit = defaultCommit
	Date = defaultDate
	Dirty = defaultDirty

	t.Cleanup(func() {
		Version = oldVersion
		Commit = oldCommit
		Date = oldDate
		Dirty = oldDirty
	})
}
