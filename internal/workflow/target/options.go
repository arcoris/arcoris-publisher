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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// RemoteURLFunc resolves a repository reference into a clone URL.
type RemoteURLFunc func(manifest.RepositoryRef) string

// Options configures target worktree preparation.
type Options struct {
	// RemoteName is the Git remote used for fetch and later publish stages.
	RemoteName string

	// RemoteURL resolves repositories when missing worktrees must be cloned.
	RemoteURL RemoteURLFunc

	// CheckoutBranch checks out the first target branch after worktree preparation.
	CheckoutBranch bool

	// CreateMissing creates an empty local worktree directory when no clone URL is configured.
	CreateMissing bool

	// Fetch updates existing target worktrees before they are reused.
	Fetch bool

	// FetchRequired turns fetch failures into target preparation errors when
	// Fetch is enabled. Leave it false for intentionally offline best-effort
	// preparation.
	FetchRequired bool

	// RequireClean rejects dirty target worktrees before construction.
	RequireClean bool
}

// DefaultOptions returns conservative target preparation defaults.
func DefaultOptions() Options {
	return Options{
		RemoteName:     "origin",
		CheckoutBranch: true,
		CreateMissing:  false,
		Fetch:          true,
		FetchRequired:  true,
		RequireClean:   true,
	}
}
