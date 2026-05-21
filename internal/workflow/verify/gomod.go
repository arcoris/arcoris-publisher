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

package verify

import "strings"

// goModInfo contains only go.mod data needed by verification.
type goModInfo struct {
	// module is the declared module path.
	module string

	// requires maps direct requirement module paths to versions.
	requires map[string]string

	// localReplaces contains modules replaced with local filesystem paths.
	localReplaces []string
}

// parseGoMod extracts module, require, and local replace directives.
func parseGoMod(data []byte) goModInfo {
	info := goModInfo{requires: map[string]string{}}
	lines := strings.Split(string(data), "\n")
	inRequire := false
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}
		fields := strings.Fields(t)
		if len(fields) >= 2 && fields[0] == "module" {
			info.module = fields[1]
			continue
		}
		if t == "require (" {
			inRequire = true
			continue
		}
		if inRequire && t == ")" {
			inRequire = false
			continue
		}
		if inRequire {
			if len(fields) >= 2 {
				info.requires[fields[0]] = fields[1]
			}
			continue
		}
		if len(fields) >= 3 && fields[0] == "require" {
			info.requires[fields[1]] = fields[2]
			continue
		}
		if strings.HasPrefix(t, "replace ") && strings.Contains(t, "=>") {
			parts := strings.SplitN(strings.TrimPrefix(t, "replace "), "=>", 2)
			if len(parts) == 2 {
				newPath := strings.TrimSpace(parts[1])
				if strings.HasPrefix(newPath, ".") || strings.HasPrefix(newPath, "/") {
					oldPath := strings.Fields(strings.TrimSpace(parts[0]))[0]
					info.localReplaces = append(info.localReplaces, oldPath)
				}
			}
		}
	}
	return info
}
