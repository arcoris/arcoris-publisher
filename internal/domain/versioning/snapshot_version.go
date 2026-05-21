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

package versioning

import "fmt"

// Version returns the Go pseudo-version derived from the snapshot.
func (s Snapshot) Version() Version {
	return Version(fmt.Sprintf("%s-%s-%s", s.base, s.timestamp(), s.shortCommit()))
}

// timestamp returns the pseudo-version UTC timestamp component.
func (s Snapshot) timestamp() string {
	return s.time.UTC().Format(snapshotTimestampLayout)
}

// shortCommit returns the pseudo-version commit component.
func (s Snapshot) shortCommit() string {
	if len(s.commit) > minimumSnapshotCommitLength {
		return s.commit[:minimumSnapshotCommitLength]
	}
	return s.commit
}
