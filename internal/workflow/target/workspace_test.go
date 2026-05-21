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

package target

import (
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"testing"
)

func TestWorkspaceSetAccessors(t *testing.T) {
	name := manifest.ModuleName("control")
	set := WorkspaceSet{workspaces: []ModuleWorkspace{{module: name, worktreeDir: "/tmp/control"}}}
	if set.Len() != 1 || set.Empty() {
		t.Fatalf("unexpected set size")
	}
	if _, ok := set.WorkspaceByModule(name); !ok {
		t.Fatalf("workspace missing")
	}
	ws := set.Workspaces()
	ws[0].module = "mutated"
	if set.ModuleNames()[0] != name {
		t.Fatalf("workspace accessor did not detach")
	}
}
