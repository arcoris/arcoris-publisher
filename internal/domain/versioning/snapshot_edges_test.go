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

import (
	"testing"
	"time"
)

func TestSnapshotVersionUsesExactLengthCommit(t *testing.T) {
	snapshot := MustSnapshot(SnapshotSpec{
		Time:   time.Date(2026, 5, 21, 12, 0, 0, 0, time.UTC),
		Commit: "abcdef123456",
	})

	if got, want := snapshot.Version(), Version("v0.0.0-20260521120000-abcdef123456"); got != want {
		t.Fatalf("Version() = %q, want %q", got, want)
	}
}

func TestNewSnapshotRejectsBlankCommit(t *testing.T) {
	_, err := NewSnapshot(SnapshotSpec{Time: time.Date(2026, 5, 21, 12, 0, 0, 0, time.UTC), Commit: "   "})
	if err == nil {
		t.Fatalf("NewSnapshot() succeeded, expected blank commit error")
	}
}

func TestSnapshotAssignmentsRejectsZeroSnapshot(t *testing.T) {
	_, err := SnapshotAssignments(testRegistry(t), Snapshot{})
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueInvalidAssignment) {
		t.Fatalf("issues = %#v, want invalid assignment", validationErr.Issues)
	}
}
