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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// WorkspaceSet contains target worktrees in plan module order.
type WorkspaceSet struct{ workspaces []ModuleWorkspace }

// Workspaces returns detached module target workspaces.
func (s WorkspaceSet) Workspaces() []ModuleWorkspace {
	out := make([]ModuleWorkspace, len(s.workspaces))
	copy(out, s.workspaces)
	return out
}

// ModuleNames returns workspace module names in plan order.
func (s WorkspaceSet) ModuleNames() []manifest.ModuleName {
	out := make([]manifest.ModuleName, 0, len(s.workspaces))
	for _, ws := range s.workspaces {
		out = append(out, ws.Module())
	}
	return out
}

// WorkspaceByModule returns the workspace for module.
func (s WorkspaceSet) WorkspaceByModule(name manifest.ModuleName) (ModuleWorkspace, bool) {
	for _, ws := range s.workspaces {
		if ws.Module() == name {
			return ws, true
		}
	}
	return ModuleWorkspace{}, false
}

// Len returns the number of prepared workspaces.
func (s WorkspaceSet) Len() int { return len(s.workspaces) }

// Empty reports whether no target workspaces were prepared.
func (s WorkspaceSet) Empty() bool { return len(s.workspaces) == 0 }

// ModuleWorkspace describes the target worktree prepared for one module.
type ModuleWorkspace struct {
	// module is the planned module name.
	module manifest.ModuleName

	// repository is the target repository reference.
	repository manifest.RepositoryRef

	// worktreeDir is the local prepared target repository directory.
	worktreeDir string

	// branches contains source-to-target branch mappings in plan order.
	branches []BranchWorkspace
}

// Module returns the planned module name.
func (w ModuleWorkspace) Module() manifest.ModuleName { return w.module }

// Repository returns the target repository.
func (w ModuleWorkspace) Repository() manifest.RepositoryRef { return w.repository }

// WorktreeDir returns the local target repository worktree directory.
func (w ModuleWorkspace) WorktreeDir() string { return w.worktreeDir }

// Branches returns detached branch workspaces.
func (w ModuleWorkspace) Branches() []BranchWorkspace {
	out := make([]BranchWorkspace, len(w.branches))
	copy(out, w.branches)
	return out
}

// BranchWorkspace records one source-to-target branch mapping prepared for a worktree.
type BranchWorkspace struct {
	// source is the source branch used by the plan.
	source manifest.BranchName

	// target is the target repository branch to publish.
	target manifest.BranchName
}

// newBranchWorkspace stores one branch mapping.
func newBranchWorkspace(source, target manifest.BranchName) BranchWorkspace {
	return BranchWorkspace{source: source, target: target}
}

// Source returns the source branch used by the plan.
func (b BranchWorkspace) Source() manifest.BranchName { return b.source }

// Target returns the target repository branch to publish.
func (b BranchWorkspace) Target() manifest.BranchName { return b.target }
