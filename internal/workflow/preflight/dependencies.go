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

package preflight

import (
	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
)

// Dependencies contains ports needed for read-only readiness checks.
type Dependencies struct {
	// Source contains source inspection ports. The service fills missing fields
	// from FS and Git so runtime wiring can stay compact.
	Source source.Dependencies

	// FS reads source and target filesystem state.
	FS filesystem.FileSystem

	// Git reads local and remote Git state without mutating refs.
	Git git.Client
}
