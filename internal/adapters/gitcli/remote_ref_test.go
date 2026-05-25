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

func TestRemoteRefExistsUsesOriginDefault(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{ExitCode: 0}, {ExitCode: 2}}}
	client := New(runner, Options{})

	exists, err := client.RemoteRefExists(context.Background(), "/repo", "", "main")
	if err != nil || !exists {
		t.Fatalf("RemoteRefExists() = %v, %v", exists, err)
	}
	exists, err = client.RemoteRefExists(context.Background(), "/repo", "", "missing")
	if err != nil || exists {
		t.Fatalf("RemoteRefExists() missing = %v, %v", exists, err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"ls-remote", "--exit-code", "origin", "main"})
}

func TestDefaultRemote(t *testing.T) {
	if got := defaultRemote(""); got != "origin" {
		t.Fatalf("defaultRemote(empty) = %q, want origin", got)
	}
	if got := defaultRemote("upstream"); got != "upstream" {
		t.Fatalf("defaultRemote(upstream) = %q, want upstream", got)
	}
}

func TestRemoteRefHashParsesLsRemote(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{
		ExitCode: 0,
		Stdout:   []byte("abc123\trefs/heads/main\n"),
	}}}
	client := New(runner, Options{})

	hash, ok, err := client.RemoteRefHash(context.Background(), "/repo", "origin", "refs/heads/main")

	if err != nil || !ok || hash != "abc123" {
		t.Fatalf("RemoteRefHash() = %q, %v, %v", hash, ok, err)
	}
}
