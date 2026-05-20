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

// CloneOptions configures a Git clone operation.
type CloneOptions struct {
	// NoTags disables tag discovery during clone.
	NoTags bool
	// Depth requests a shallow clone when greater than zero.
	Depth int
	// Bare creates a repository without a working tree.
	Bare bool
	// Mirror creates a full mirror suitable for ref replication.
	Mirror bool
	// SensitiveValues are raw values that adapters must redact in diagnostics.
	SensitiveValues []string
}
