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

package modulefile

import (
	"bytes"
	"path/filepath"
	"sort"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"golang.org/x/mod/modfile"
)

// rewriteGoMod edits go.mod through the Go module parser instead of rebuilding
// directive blocks by hand.
func rewriteGoMod(
	data []byte,
	mod plan.ModulePlan,
	removeLocalReplaces bool,
) ([]byte, []RequirementUpdate, bool, error) {
	file, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, nil, false, err
	}

	managed := managedRequirements(mod)
	if err := file.AddModuleStmt(mod.ModulePath().String()); err != nil {
		return nil, nil, false, err
	}

	dropManagedRequirements(file, managed)
	if removeLocalReplaces {
		dropManagedLocalReplaces(file, managed)
	}

	updates := make([]RequirementUpdate, 0, len(mod.Requirements()))
	for _, req := range sortedRequirements(mod.Requirements()) {
		if err := file.AddRequire(req.ModulePath().String(), req.Version().String()); err != nil {
			return nil, nil, false, err
		}
		updates = append(updates, RequirementUpdate{
			modulePath: req.ModulePath(),
			version:    req.Version().String(),
		})
	}

	file.Cleanup()
	out, err := file.Format()
	if err != nil {
		return nil, nil, false, err
	}

	return out, updates, !bytes.Equal(out, data), nil
}

// managedRequirements indexes direct internal requirements controlled by the
// plan so unrelated go.mod directives can be preserved untouched.
func managedRequirements(mod plan.ModulePlan) map[string]struct{} {
	managed := make(map[string]struct{}, len(mod.Requirements()))
	for _, req := range mod.Requirements() {
		managed[req.ModulePath().String()] = struct{}{}
	}

	return managed
}

// dropManagedRequirements removes stale internal requirements before the
// current planned versions are added back.
func dropManagedRequirements(file *modfile.File, managed map[string]struct{}) {
	for _, req := range file.Require {
		if _, ok := managed[req.Mod.Path]; ok {
			_ = file.DropRequire(req.Mod.Path)
		}
	}
}

// dropManagedLocalReplaces removes local development replaces for managed
// internal modules while leaving external or remote replacements intact.
func dropManagedLocalReplaces(file *modfile.File, managed map[string]struct{}) {
	for _, replacement := range file.Replace {
		if _, ok := managed[replacement.Old.Path]; !ok {
			continue
		}
		if !isLocalReplacePath(replacement.New.Path) {
			continue
		}

		_ = file.DropReplace(replacement.Old.Path, replacement.Old.Version)
	}
}

// isLocalReplacePath reports whether a replace target points at the local
// filesystem rather than another module path.
func isLocalReplacePath(path string) bool {
	if path == "." || path == ".." || filepath.IsAbs(path) {
		return true
	}

	return strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") ||
		strings.HasPrefix(path, ".\\") ||
		strings.HasPrefix(path, "..\\")
}

// sortedRequirements keeps added managed requirements deterministic even if a
// future plan implementation changes its internal requirement order.
func sortedRequirements(in []plan.DependencyRequirement) []plan.DependencyRequirement {
	out := make([]plan.DependencyRequirement, len(in))
	copy(out, in)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].ModulePath().String() < out[j].ModulePath().String()
	})

	return out
}

// parseModuleLine returns the module path from go.mod data when present. The
// verification workflow keeps this helper for lightweight module-path checks.
func parseModuleLine(data []byte) manifest.ModulePath {
	file, err := modfile.ParseLax("go.mod", data, nil)
	if err != nil || file.Module == nil {
		return ""
	}

	path, _ := manifest.ParseModulePath(file.Module.Mod.Path)
	return path
}
