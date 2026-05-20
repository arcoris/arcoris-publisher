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

package git

// FetchTagsMode describes how tags should be fetched.
type FetchTagsMode string

const (
	// FetchTagsDefault uses the Git implementation's default tag fetching behavior.
	FetchTagsDefault FetchTagsMode = "default"
	// FetchTagsAll fetches all tags.
	FetchTagsAll FetchTagsMode = "all"
	// FetchTagsNone disables tag fetching.
	FetchTagsNone FetchTagsMode = "none"
)

// String returns the stable string representation of the fetch tag mode.
func (m FetchTagsMode) String() string {
	return string(m)
}

// Valid reports whether the mode is one of the supported fetch tag modes.
func (m FetchTagsMode) Valid() bool {
	switch m {
	case FetchTagsDefault, FetchTagsAll, FetchTagsNone:
		return true
	default:
		return false
	}
}

// FetchOptions configures a Git fetch operation.
type FetchOptions struct {
	// Prune removes remote-tracking refs that no longer exist upstream.
	Prune bool
	// Tags controls tag fetching behavior.
	Tags FetchTagsMode
	// RefSpecs limits or expands the refs fetched from the remote.
	RefSpecs []RefSpec
	// SensitiveValues are raw values that adapters must redact in diagnostics.
	SensitiveValues []string
}
