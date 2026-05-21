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

func TestNewSnapshotBuildsPseudoVersion(t *testing.T) {
	snapshot, err := NewSnapshot(SnapshotSpec{
		Time:   time.Date(2026, 5, 21, 12, 34, 56, 0, time.FixedZone("UTC+3", 3*60*60)),
		Commit: "ABCDEF1234567890",
	})
	if err != nil {
		t.Fatalf("NewSnapshot returned error: %v", err)
	}
	if snapshot.Base() != "v0.0.0" {
		t.Fatalf("unexpected base: %q", snapshot.Base())
	}
	if snapshot.Commit() != "abcdef1234567890" {
		t.Fatalf("unexpected normalized commit: %q", snapshot.Commit())
	}
	if got, want := snapshot.Version(), Version("v0.0.0-20260521093456-abcdef123456"); got != want {
		t.Fatalf("unexpected snapshot version: got %q want %q", got, want)
	}
}

func TestNewSnapshotRejectsInvalidInput(t *testing.T) {
	cases := []SnapshotSpec{
		{Time: time.Now(), Commit: "abc"},
		{Time: time.Time{}, Commit: "abcdef123456"},
		{Time: time.Now(), Commit: "abcdef12345z"},
		{Base: "bad", Time: time.Now(), Commit: "abcdef123456"},
	}
	for _, tc := range cases {
		if _, err := NewSnapshot(tc); err == nil {
			t.Fatalf("NewSnapshot(%+v) succeeded, expected error", tc)
		}
	}
}

func TestMustSnapshotPanicsOnInvalidSnapshot(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("MustSnapshot did not panic")
		}
	}()
	_ = MustSnapshot(SnapshotSpec{})
}
