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
	portfs "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	portgit "arcoris.dev/arcoris-publisher/internal/ports/git"
)

// Dependencies contains infrastructure ports used by Service.
type Dependencies struct {
	// Git performs target repository clone, fetch, status, and checkout operations.
	Git portgit.WorktreeClient

	// FS reads and creates target worktree directories.
	FS portfs.FileSystem
}
