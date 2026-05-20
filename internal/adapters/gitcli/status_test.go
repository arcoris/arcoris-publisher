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

package gitcli

import (
	"context"
	"testing"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestStatusParsesPorcelain(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{Stdout: []byte(" M file.txt\x00R  new.txt\x00old.txt\x00")}}}
	client := New(runner, Options{})

	status, err := client.Status(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if status.Clean || len(status.Entries) != 2 {
		t.Fatalf("Status() = %#v, want 2 dirty entries", status)
	}
	if status.Entries[1].Code != "R " || status.Entries[1].Path != "new.txt" {
		t.Fatalf("rename entry = %#v", status.Entries[1])
	}
}

func TestParseStatusHandlesMalformedEntry(t *testing.T) {
	entries := parseStatus([]byte("?\x00"))
	if len(entries) != 1 || entries[0].Path != "?" {
		t.Fatalf("parseStatus() = %#v", entries)
	}
}
