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
	"errors"
	"runtime"
	"testing"
	"time"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestRunnerRunCapturesAndRedactsOutput(t *testing.T) {
	runner := New(Options{})
	spec := shellSpec("printf 'secret-value'")
	spec.CaptureStdout = true
	spec.SensitiveValues = []string{"secret-value"}

	result, err := runner.Run(context.Background(), spec)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if string(result.Stdout) != "<redacted>" {
		t.Fatalf("stdout was not redacted: %q", result.Stdout)
	}
}

func TestRunnerRunAcceptsExplicitExitCode(t *testing.T) {
	runner := New(Options{})
	spec := shellSpec("exit 7")
	spec.AllowedExitCodes = []int{7}

	result, err := runner.Run(context.Background(), spec)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.ExitCode != 7 {
		t.Fatalf("ExitCode = %d, want 7", result.ExitCode)
	}
}

func TestRunnerRunRejectsZeroWhenExplicitlyDisallowed(t *testing.T) {
	runner := New(Options{})
	spec := shellSpec("exit 0")
	spec.AllowedExitCodes = []int{7}

	result, err := runner.Run(context.Background(), spec)
	if result.ExitCode != 0 {
		t.Fatalf("ExitCode = %d, want 0", result.ExitCode)
	}
	assertPortCode(t, err, processport.CodeFailed)
}

func TestRunnerRunReportsRejectedExitCode(t *testing.T) {
	runner := New(Options{})
	_, err := runner.Run(context.Background(), shellSpec("exit 5"))
	assertPortCode(t, err, processport.CodeFailed)
}

func TestRunnerRunReportsTimeout(t *testing.T) {
	runner := New(Options{})
	spec := shellSpec("sleep 1")
	spec.Timeout = 10 * time.Millisecond

	_, err := runner.Run(context.Background(), spec)
	assertPortCode(t, err, processport.CodeTimedOut)
}

func TestRunnerRunReportsEmptyName(t *testing.T) {
	runner := New(Options{})

	result, err := runner.Run(context.Background(), processport.Spec{})
	if result.ExitCode != 0 {
		t.Fatalf("ExitCode = %d, want zero for startup failure", result.ExitCode)
	}
	assertPortCode(t, err, processport.CodeStartFailed)
}

func TestRunnerRunReportsMissingExecutable(t *testing.T) {
	runner := New(Options{})

	result, err := runner.Run(context.Background(), processport.Spec{Name: "arcoris-missing-command-for-test"})
	if result.ExitCode != 0 {
		t.Fatalf("ExitCode = %d, want zero for missing executable", result.ExitCode)
	}
	assertPortCode(t, err, processport.CodeNotFound)
}

func shellSpec(script string) processport.Spec {
	if runtime.GOOS == "windows" {
		return processport.Spec{Name: "cmd", Args: []string{"/C", script}, CaptureStdout: true, CaptureStderr: true}
	}
	return processport.Spec{Name: "sh", Args: []string{"-c", script}, CaptureStdout: true, CaptureStderr: true}
}

func assertPortCode(t *testing.T, err error, code porterr.Code) {
	t.Helper()
	var perr *porterr.Error
	if !errors.As(err, &perr) {
		t.Fatalf("expected port error, got %T %v", err, err)
	}
	if perr.Code != code {
		t.Fatalf("expected code %s, got %s", code, perr.Code)
	}
}
