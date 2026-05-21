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
	portfs "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	portgit "arcoris.dev/arcoris-publisher/internal/ports/git"
	portgo "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// Dependencies contains infrastructure ports used by Service.
type Dependencies struct {
	// FS reads target worktree files and directories.
	FS portfs.Reader

	// Git checks whether verification changed target worktrees.
	Git portgit.RepositoryReader

	// Go runs go list, go test, and go mod tidy checks when enabled.
	Go portgo.Toolchain
}
