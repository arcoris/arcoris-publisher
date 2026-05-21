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
	"testing"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestFinishSuccessfulProcessRejectsExplicitlyDisallowedZero(t *testing.T) {
	result, err := finishSuccessfulProcess(processport.Result{}, processport.Spec{Name: "cmd", AllowedExitCodes: []int{7}}, Redactor{})

	if result.ExitCode != 0 {
		t.Fatalf("ExitCode = %d, want 0", result.ExitCode)
	}
	assertPortCode(t, err, processport.CodeFailed)
}

func TestExitErrorHelpersUseStableCodes(t *testing.T) {
	spec := processport.Spec{Name: "secret-cmd", SensitiveValues: []string{"secret"}}
	redactor := NewRedactor("secret")

	assertPortCode(t, emptyNameError(spec, redactor), processport.CodeStartFailed)
	assertPortCode(t, rejectedExitCodeError(spec, redactor, context.Canceled, 7), processport.CodeFailed)
	assertPortCode(t, timedOutError(spec, redactor, context.DeadlineExceeded, -1), processport.CodeTimedOut)
	assertPortCode(t, cancelledError(spec, redactor, context.Canceled, -1), processport.CodeCancelled)
	assertPortCode(t, executableNotFoundError(spec, redactor, context.Canceled), processport.CodeNotFound)
	assertPortCode(t, startFailedError(spec, redactor, context.Canceled), processport.CodeStartFailed)
}

func TestRedactedCommandDetailsHidesSensitiveValues(t *testing.T) {
	spec := processport.Spec{Name: "git-token", Args: []string{"push", "git-token"}, Dir: "/tmp/git-token"}
	details := redactedCommandDetails(spec, NewRedactor("git-token"), 3)

	if details["name"] != "<redacted>" || details["args"] != "push <redacted>" || details["dir"] != "/tmp/<redacted>" || details["exit_code"] != "3" {
		t.Fatalf("redactedCommandDetails() = %#v", details)
	}
}
