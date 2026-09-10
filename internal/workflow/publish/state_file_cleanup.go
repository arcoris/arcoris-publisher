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

package publish

// stateFileCleanupOutcome records which irreversible cleanup steps completed.
// Removed=true with Synced=false means the directory entry was unlinked but
// the durability of that unlink could not be confirmed.
type stateFileCleanupOutcome struct {
	Removed bool
	Synced  bool
}

// cleanupCreatedStateFile removes a state file that this process just created.
// Unlike normal lock release it deliberately does not parse or identity-check
// the file: acquisition may have failed while writing a partial representation.
func cleanupCreatedStateFile(
	path string,
	remove func(string) error,
	syncParent func(string) error,
) (stateFileCleanupOutcome, error) {
	if err := remove(path); err != nil {
		return stateFileCleanupOutcome{}, err
	}
	outcome := stateFileCleanupOutcome{Removed: true}
	if err := syncParent(path); err != nil {
		return outcome, err
	}
	outcome.Synced = true
	return outcome, nil
}
