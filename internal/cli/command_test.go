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

package cli

import "testing"

func TestParseCommandDefaultsToHelp(t *testing.T) {
	t.Parallel()

	cmd, rest, err := parseCommand(nil)
	if err != nil {
		t.Fatalf("parseCommand() error = %v", err)
	}
	if cmd != commandHelp || len(rest) != 0 {
		t.Fatalf("parseCommand() = %q, %v", cmd, rest)
	}
}

func TestParseCommandKnownCommand(t *testing.T) {
	t.Parallel()

	cmd, rest, err := parseCommand([]string{"verify", "--version", "v0.1.0"})
	if err != nil {
		t.Fatalf("parseCommand() error = %v", err)
	}
	if cmd != commandVerify || len(rest) != 2 || rest[0] != "--version" {
		t.Fatalf("parseCommand() = %q, %v", cmd, rest)
	}
}

func TestParseCommandRejectsUnknown(t *testing.T) {
	t.Parallel()

	_, _, err := parseCommand([]string{"unknown"})
	if err == nil {
		t.Fatal("parseCommand() error = nil")
	}
}

func TestIsHelpRequest(t *testing.T) {
	t.Parallel()

	if !isHelpRequest([]string{"--help"}) || !isHelpRequest([]string{"-h"}) {
		t.Fatal("help request was not detected")
	}
	if isHelpRequest([]string{"--version", "v0.3.0"}) {
		t.Fatal("non-help args detected as help")
	}
}
