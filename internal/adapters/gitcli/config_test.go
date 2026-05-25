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
	"errors"
	"testing"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestConfigGetBuildsCommand(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{Stdout: []byte("ARCORIS Test\n")}}}
	client := New(runner, Options{})

	got, ok, err := client.ConfigGet(context.Background(), "/repo", "user.name")
	if err != nil || !ok || got != "ARCORIS Test" {
		t.Fatalf("ConfigGet() = %q, %v, %v", got, ok, err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"config", "--get", "user.name"})
	assertStringSlice(t, intsOf(runner.specs[0].AllowedExitCodes), []string{"0", "1"})
}

func TestConfigGetMissing(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{ExitCode: 1}}}
	client := New(runner, Options{})

	got, ok, err := client.ConfigGet(context.Background(), "/repo", "user.email")
	if err != nil || ok || got != "" {
		t.Fatalf("ConfigGet() missing = %q, %v, %v", got, ok, err)
	}
}

func TestConfigGetCommandFailure(t *testing.T) {
	runner := &fakeRunner{
		results: []processport.Result{{Stderr: []byte("fatal: bad config\n")}},
		errs:    []error{errors.New("exit 2")},
	}
	client := New(runner, Options{})

	_, _, err := client.ConfigGet(context.Background(), "/repo", "user.name")
	assertPortCode(t, err, gitport.CodeCommandFailed)
}
