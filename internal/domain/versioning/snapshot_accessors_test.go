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

func TestSnapshotAccessors(t *testing.T) {
	when := time.Date(2026, 5, 21, 1, 2, 3, 0, time.UTC)
	snapshot := MustSnapshot(SnapshotSpec{Base: "v1.2.3+build", Time: when, Commit: "abcdef123456"})
	if snapshot.Base() != "v1.2.3" {
		t.Fatalf("build metadata was not stripped from snapshot base: %q", snapshot.Base())
	}
	if !snapshot.Time().Equal(when) {
		t.Fatalf("unexpected snapshot time: %v", snapshot.Time())
	}
}
