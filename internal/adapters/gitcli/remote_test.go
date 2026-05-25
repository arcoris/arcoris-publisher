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

func TestRemoteURLBuildsCommand(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{Stdout: []byte("https://example/repo.git\n")}}}
	client := New(runner, Options{})

	got, ok, err := client.RemoteURL(context.Background(), "/repo", "")
	if err != nil || !ok || got != "https://example/repo.git" {
		t.Fatalf("RemoteURL() = %q, %v, %v", got, ok, err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"remote", "get-url", "origin"})
	assertStringSlice(t, intsOf(runner.specs[0].AllowedExitCodes), []string{"0", "2"})
}

func TestRemoteURLMissing(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{ExitCode: 2}}}
	client := New(runner, Options{})

	got, ok, err := client.RemoteURL(context.Background(), "/repo", "upstream")
	if err != nil || ok || got != "" {
		t.Fatalf("RemoteURL() missing = %q, %v, %v", got, ok, err)
	}
}

func TestAddRemoteBuildsCommandAndRedactsURL(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	url := "https://token@example/repo.git"
	if err := client.AddRemote(context.Background(), "/repo", "", url); err != nil {
		t.Fatalf("AddRemote() error = %v", err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"remote", "add", "origin", url})
	assertStringSlice(t, runner.specs[0].SensitiveValues, []string{url})
}
