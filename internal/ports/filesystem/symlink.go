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

// SymlinkPolicy controls how filesystem operations treat symbolic links.
//
// The zero value is intentionally invalid so callers must choose a policy when
// an operation can encounter links. Adapters may provide their own default only
// when an option struct clearly documents that behavior.
type SymlinkPolicy string

const (
	// SymlinkReject rejects symbolic links.
	SymlinkReject SymlinkPolicy = "reject"
	// SymlinkPreserve preserves symbolic links as links.
	SymlinkPreserve SymlinkPolicy = "preserve"
	// SymlinkFollow follows symbolic links. Callers should use this mode only
	// when the adapter also enforces strict root-escape checks.
	SymlinkFollow SymlinkPolicy = "follow"
)

// String returns the stable string representation of the symlink policy.
func (p SymlinkPolicy) String() string {
	return string(p)
}

// Valid reports whether the policy is one of the supported symbolic-link modes.
func (p SymlinkPolicy) Valid() bool {
	switch p {
	case SymlinkReject, SymlinkPreserve, SymlinkFollow:
		return true
	default:
		return false
	}
}
