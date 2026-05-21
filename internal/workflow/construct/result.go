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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Result describes constructed target module trees in plan order.
type Result struct{ modules []ModuleResult }

// Modules returns detached module construction results in plan order.
func (r Result) Modules() []ModuleResult {
	out := make([]ModuleResult, len(r.modules))
	copy(out, r.modules)
	return out
}

// ModuleByName returns the construction result for name.
func (r Result) ModuleByName(name manifest.ModuleName) (ModuleResult, bool) {
	for _, m := range r.modules {
		if m.Module() == name {
			return m, true
		}
	}
	return ModuleResult{}, false
}

// Changed reports whether construction changed any module worktree.
func (r Result) Changed() bool {
	for _, m := range r.modules {
		if m.Changed() {
			return true
		}
	}
	return false
}

// ModuleResult describes construction work performed for one module.
type ModuleResult struct {
	// module is the planned module name.
	module manifest.ModuleName

	// worktreeDir is the target worktree that received copied entries.
	worktreeDir string

	// operations records deterministic construction actions.
	operations []Operation

	// changed reports whether any construction operation affected this module.
	changed bool
}

// Module returns the planned module name.
func (m ModuleResult) Module() manifest.ModuleName { return m.module }

// WorktreeDir returns the target worktree that received copied entries.
func (m ModuleResult) WorktreeDir() string { return m.worktreeDir }

// Operations returns detached construction operations.
func (m ModuleResult) Operations() []Operation {
	out := make([]Operation, len(m.operations))
	copy(out, m.operations)
	return out
}

// Changed reports whether construction affected this module.
func (m ModuleResult) Changed() bool { return m.changed }

// OperationKind identifies one construction operation kind.
type OperationKind string

const (
	// OperationClean records a target worktree cleanup.
	OperationClean OperationKind = "clean"

	// OperationCopyFile records copying one explicit file entry.
	OperationCopyFile OperationKind = "copy-file"

	// OperationCopyDirectory records copying one explicit directory entry.
	OperationCopyDirectory OperationKind = "copy-directory"

	// OperationSkipOptional records an absent optional source entry.
	OperationSkipOptional OperationKind = "skip-optional"

	// OperationWriteGenerated records writing generated metadata.
	OperationWriteGenerated OperationKind = "write-generated"
)

// Operation describes one deterministic construction action.
type Operation struct {
	// kind identifies the action.
	kind OperationKind

	// sourcePath is the source entry path, empty for generated operations.
	sourcePath string

	// targetPath is the target path affected by the operation.
	targetPath string
}

// newOperation creates one construction operation record.
func newOperation(kind OperationKind, sourcePath, targetPath string) Operation {
	return Operation{kind: kind, sourcePath: sourcePath, targetPath: targetPath}
}

// Kind returns the operation kind.
func (o Operation) Kind() OperationKind { return o.kind }

// SourcePath returns the source entry path, empty for generated operations.
func (o Operation) SourcePath() string { return o.sourcePath }

// TargetPath returns the target path affected by the operation.
func (o Operation) TargetPath() string { return o.targetPath }
