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
	"path/filepath"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// repositoryWorktree maps a repository reference to a deterministic local path.
func repositoryWorktree(root string, repo manifest.RepositoryRef) string {
	name := strings.ReplaceAll(repo.String(), "/", "__")
	return filepath.Join(root, name)
}
