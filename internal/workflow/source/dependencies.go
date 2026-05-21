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

package source

import (
	portfs "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	portgit "arcoris.dev/arcoris-publisher/internal/ports/git"
)

// Dependencies contains infrastructure ports used by Service.
//
// Dependencies MUST be satisfied by ports, not adapters directly. Application
// wiring chooses concrete implementations.
type Dependencies struct {
	// Git reads repository provenance and status information.
	Git portgit.RepositoryReader

	// FS reads source paths and computes content hashes.
	FS portfs.FileSystem
}
