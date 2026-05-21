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
	"testing"
	"time"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestNewResultCopiesCommandIdentity(t *testing.T) {
	args := []string{"status"}
	startedAt := time.Unix(1, 0)
	result := newResult(processport.Spec{Name: "git", Args: args, Dir: "/repo"}, startedAt)
	args[0] = "mutated"

	if result.Name != "git" || result.Args[0] != "status" || result.Dir != "/repo" || !result.StartedAt.Equal(startedAt) {
		t.Fatalf("newResult() = %#v", result)
	}
}

func TestFinishResultRedactsCapturedOutput(t *testing.T) {
	startedAt := time.Now().Add(-time.Second)
	capture := &outputCapture{stdout: newLimitedBuffer(0), stderr: newLimitedBuffer(0)}
	_, _ = capture.stdout.Write([]byte("secret out"))
	_, _ = capture.stderr.Write([]byte("secret err"))
	result := processport.Result{}

	finishResult(&result, startedAt, capture, NewRedactor("secret"))

	if string(result.Stdout) != "<redacted> out" || string(result.Stderr) != "<redacted> err" {
		t.Fatalf("finishResult() stdout=%q stderr=%q", result.Stdout, result.Stderr)
	}
	if result.FinishedAt.IsZero() || result.Duration <= 0 {
		t.Fatalf("finishResult() did not record timing: %#v", result)
	}
}
