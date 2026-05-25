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

// PushOptions configures a Git push operation.
type PushOptions struct {
	// Force permits overwriting remote refs without lease protection.
	Force bool
	// ForceWithLease permits overwriting only when the remote ref is unchanged.
	ForceWithLease bool
	// ForceWithLeaseRef names the exact ref protected by ForceWithLeaseExpect.
	ForceWithLeaseRef string
	// ForceWithLeaseExpect is the object expected at ForceWithLeaseRef.
	ForceWithLeaseExpect CommitHash
	// Atomic requests all ref updates to succeed or fail together.
	Atomic bool
	// SensitiveValues are raw values that adapters must redact in diagnostics.
	SensitiveValues []string
}
