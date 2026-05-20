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

package filesystem

import "context"

// Cleaner describes safe directory cleanup operations.
type Cleaner interface {
	CleanDir(ctx context.Context, dir string, opts CleanDirOptions) error
}

// CleanDirOptions configures safe recursive cleanup of one directory.
type CleanDirOptions struct {
	// Preserve contains relative paths or patterns that must survive cleanup.
	Preserve []string
	// RequireGitDir makes cleanup fail unless dir contains repository metadata.
	RequireGitDir bool
	// AllowMissing treats a missing target directory as a successful no-op.
	AllowMissing bool
	// SafetyRoot confines deletion to a known parent directory.
	SafetyRoot string
}
