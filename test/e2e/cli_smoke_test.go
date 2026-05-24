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

package e2e_test

import "testing"

func TestArcpubHelp(t *testing.T) {
	result := runArcpub(t, "help")
	assertExitCode(t, result, 0)
	for _, command := range []string{"plan", "verify", "publish", "version", "completion"} {
		assertContains(t, result.Stdout, command)
	}
}

func TestArcpubRootHelp(t *testing.T) {
	result := runArcpub(t, "--help")
	assertExitCode(t, result, 0)
	for _, command := range []string{"plan", "verify", "publish", "version", "completion"} {
		assertContains(t, result.Stdout, command)
	}
}

func TestArcpubVersionText(t *testing.T) {
	result := runArcpub(t, "version")
	assertExitCode(t, result, 0)
	for _, want := range []string{"arcpub", "commit:", "date:", "dirty:"} {
		assertContains(t, result.Stdout, want)
	}
}

func TestArcpubVersionJSON(t *testing.T) {
	result := runArcpub(t, "version", "--output", "json")
	assertExitCode(t, result, 0)
	decoded := assertJSON(t, result.Stdout)
	for _, key := range []string{"version", "commit", "date", "dirty"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("version JSON missing key %q: %#v", key, decoded)
		}
	}
}

func TestArcpubCompletionBash(t *testing.T) {
	result := runArcpub(t, "completion", "bash")
	assertExitCode(t, result, 0)
	assertContains(t, result.Stdout, "__start_arcpub")
}

func TestArcpubCompletionZsh(t *testing.T) {
	result := runArcpub(t, "completion", "zsh")
	assertExitCode(t, result, 0)
	assertContains(t, result.Stdout, "#compdef arcpub")
}

func TestArcpubCompletionFish(t *testing.T) {
	result := runArcpub(t, "completion", "fish")
	assertExitCode(t, result, 0)
	assertContains(t, result.Stdout, "complete -c arcpub")
}

func TestArcpubCompletionPowerShell(t *testing.T) {
	result := runArcpub(t, "completion", "powershell")
	assertExitCode(t, result, 0)
	assertContains(t, result.Stdout, "Register-ArgumentCompleter")
}
