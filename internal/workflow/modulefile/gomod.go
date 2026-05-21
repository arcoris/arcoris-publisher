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
	"sort"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/plan"
)

// rewriteGoMod rewrites module identity and direct internal requirements using
// the immutable plan. It preserves unrelated requirements and comments.
func rewriteGoMod(
	data []byte,
	mod plan.ModulePlan,
	removeLocalReplaces bool,
) ([]byte, []RequirementUpdate, bool) {
	original := string(data)
	lines := strings.Split(strings.ReplaceAll(original, "\r\n", "\n"), "\n")
	out := make([]string, 0, len(lines)+len(mod.Requirements())+4)
	moduleSet := false
	requirePaths := map[string]struct{}{}
	for _, req := range mod.Requirements() {
		requirePaths[req.ModulePath().String()] = struct{}{}
	}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "module ") {
			out = append(out, "module "+mod.ModulePath().String())
			moduleSet = true
			continue
		}
		if strings.HasPrefix(trimmed, "replace ") && removeLocalReplaces {
			old := strings.Fields(strings.TrimPrefix(trimmed, "replace "))
			if len(old) > 0 {
				if _, ok := requirePaths[old[0]]; ok {
					continue
				}
			}
		}
		if strings.HasPrefix(trimmed, "require ") && !strings.HasPrefix(trimmed, "require (") {
			fields := strings.Fields(strings.TrimPrefix(trimmed, "require "))
			if len(fields) >= 2 {
				if _, ok := requirePaths[fields[0]]; ok {
					continue
				}
			}
		}
		if trimmed == "require (" {
			out = append(out, line)
			for i++; i < len(lines); i++ {
				inner := lines[i]
				innerTrimmed := strings.TrimSpace(inner)
				if innerTrimmed == ")" {
					out = append(out, inner)
					break
				}
				fields := strings.Fields(innerTrimmed)
				if len(fields) >= 2 {
					if _, ok := requirePaths[fields[0]]; ok {
						continue
					}
				}
				out = append(out, inner)
			}
			continue
		}
		out = append(out, line)
	}
	if !moduleSet {
		out = append([]string{"module " + mod.ModulePath().String(), ""}, out...)
	}
	updates := make([]RequirementUpdate, 0, len(mod.Requirements()))
	reqs := mod.Requirements()
	sort.SliceStable(reqs, func(i, j int) bool {
		return reqs[i].ModulePath().String() < reqs[j].ModulePath().String()
	})
	if len(reqs) > 0 {
		out = append(out, "")
		out = append(out, "require (")
		for _, req := range reqs {
			out = append(out, "\t"+req.ModulePath().String()+" "+req.Version().String())
			updates = append(updates, RequirementUpdate{
				modulePath: req.ModulePath(),
				version:    req.Version().String(),
			})
		}
		out = append(out, ")")
	}
	result := strings.TrimRight(strings.Join(out, "\n"), "\n") + "\n"
	return []byte(result), updates, result != original
}

// parseModuleLine returns the module path from go.mod data when present.
func parseModuleLine(data []byte) manifest.ModulePath {
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 2 && fields[0] == "module" {
			p, _ := manifest.ParseModulePath(fields[1])
			return p
		}
	}
	return ""
}
