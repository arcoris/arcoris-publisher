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

// NewSnapshot validates spec and returns Snapshot.
//
// Empty Base means v0.0.0, matching Go pseudo-version conventions for modules
// without a previous release. Build metadata is stripped from the base because
// pseudo-versions encode identity through timestamp and commit hash.
func NewSnapshot(spec SnapshotSpec) (Snapshot, error) {
	base, err := parseSnapshotBase(spec.Base)
	if err != nil {
		return Snapshot{}, err
	}
	timestamp, err := normalizeSnapshotTime(spec.Time)
	if err != nil {
		return Snapshot{}, err
	}
	commit, err := normalizeCommit(spec.Commit)
	if err != nil {
		return Snapshot{}, fmt.Errorf("commit: %w", err)
	}
	return Snapshot{base: base, time: timestamp, commit: commit}, nil
}

// MustSnapshot constructs Snapshot and panics if spec is invalid.
func MustSnapshot(spec SnapshotSpec) Snapshot {
	snapshot, err := NewSnapshot(spec)
	if err != nil {
		panic(err)
	}
	return snapshot
}
