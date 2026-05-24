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

import (
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// goModInfo contains only go.mod data needed by verification.
type goModInfo struct {
	// module is the declared module path.
	module string

	// requires maps direct requirement module paths to versions.
	requires map[string]string

	// localReplaces contains modules replaced with local filesystem paths.
	localReplaces []string
}

// parseGoMod extracts module, require, and local replace directives using the
// Go module parser so block directives, versions, and comments are interpreted
// the same way as the Go toolchain.
func parseGoMod(data []byte) goModInfo {
	info := goModInfo{requires: map[string]string{}}
	file, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return info
	}

	if file.Module != nil {
		info.module = file.Module.Mod.Path
	}
	for _, req := range file.Require {
		info.requires[req.Mod.Path] = req.Mod.Version
	}
	for _, replacement := range file.Replace {
		if isLocalReplacePath(replacement.New.Path) {
			info.localReplaces = append(info.localReplaces, replacement.Old.Path)
		}
	}

	return info
}

func isLocalReplacePath(path string) bool {
	if path == "." || path == ".." || filepath.IsAbs(path) {
		return true
	}

	return strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") ||
		strings.HasPrefix(path, ".\\") ||
		strings.HasPrefix(path, "..\\")
}
