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

package exec

import (
	"context"
	"strings"
	"testing"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestPrepareCommandBuildsOSCommandAndCapture(t *testing.T) {
	runner := New(Options{Env: []string{"BASE=1"}})
	spec := processport.Spec{
		Name:             "go",
		Args:             []string{"version"},
		Dir:              "/tmp",
		Env:              []string{"BASE=2", "EXTRA=1"},
		CaptureStdout:    true,
		CaptureStderr:    true,
		MaxStdoutBytes:   16,
		MaxStderrBytes:   32,
		AllowedExitCodes: []int{0},
	}

	cmd, capture := runner.prepareCommand(context.Background(), spec)

	if cmd.Args[0] != "go" || strings.Join(cmd.Args, " ") != "go version" || cmd.Dir != "/tmp" {
		t.Fatalf("prepareCommand() command = path %q args %#v dir %q", cmd.Path, cmd.Args, cmd.Dir)
	}
	if capture.stdout == nil || capture.stderr == nil || cmd.Stdout != capture.stdout || cmd.Stderr != capture.stderr {
		t.Fatalf("prepareCommand() did not wire capture buffers")
	}
	if got := outputLimit(0); got != defaultOutputLimit {
		t.Fatalf("outputLimit(0) = %d, want %d", got, defaultOutputLimit)
	}
	if got := outputLimit(12); got != 12 {
		t.Fatalf("outputLimit(12) = %d, want 12", got)
	}
}
