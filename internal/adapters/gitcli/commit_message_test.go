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

func TestCommitMessageDefaultsToHead(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{Stdout: []byte("subject\n\nbody\n")}}}
	client := New(runner, Options{})

	message, err := client.CommitMessage(context.Background(), "/repo", "")
	if err != nil || message != "subject\n\nbody\n" {
		t.Fatalf("CommitMessage() = %q, %v", message, err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"log", "-1", "--format=%B", "HEAD"})
}

func TestDefaultRef(t *testing.T) {
	if got := defaultRef(""); got != "HEAD" {
		t.Fatalf("defaultRef(empty) = %q, want HEAD", got)
	}
	if got := defaultRef("main"); got != "main" {
		t.Fatalf("defaultRef(main) = %q, want main", got)
	}
}
