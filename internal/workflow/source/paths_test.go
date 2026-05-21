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

import "testing"

func TestPathHelpersResolveManifestPaths(t *testing.T) {
	moduleDir := resolveModuleDir("/repo/staging", "src/arcoris.dev/foundation")
	if moduleDir != "/repo/staging/src/arcoris.dev/foundation" {
		t.Fatalf("resolveModuleDir() = %q", moduleDir)
	}

	moduleRootDir := resolveModuleRootDir(moduleDir, "pkg")
	if moduleRootDir != "/repo/staging/src/arcoris.dev/foundation/pkg" {
		t.Fatalf("resolveModuleRootDir() = %q", moduleRootDir)
	}

	entry := mustPublishEntry(t, fileEntrySpec("contracts/api.go"))
	sourcePath := resolveEntrySource(moduleRootDir, entry)
	if sourcePath != "/repo/staging/src/arcoris.dev/foundation/pkg/contracts/api.go" {
		t.Fatalf("resolveEntrySource() = %q", sourcePath)
	}
}
