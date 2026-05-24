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

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunCompletionShells(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		shell string
		want  string
	}{
		{name: "bash", shell: "bash", want: "__start_arcpub"},
		{name: "zsh", shell: "zsh", want: "#compdef arcpub"},
		{name: "fish", shell: "fish", want: "complete -c arcpub"},
		{name: "powershell", shell: "powershell", want: "Register-ArgumentCompleter"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer
			code := New(Dependencies{}, Options{}).Run(
				context.Background(),
				[]string{"completion", tt.shell},
				&stdout,
				&stderr,
			)

			if code != ExitOK {
				t.Fatalf("Run(completion %s) code = %d stderr = %q", tt.shell, code, stderr.String())
			}
			if !strings.Contains(stdout.String(), tt.want) {
				t.Fatalf("completion output missing %q", tt.want)
			}
		})
	}
}

func TestRunCompletionRejectsInvalidShell(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := New(Dependencies{}, Options{}).Run(
		context.Background(),
		[]string{"completion", "xonsh"},
		&stdout,
		&stderr,
	)

	if code != ExitUsage {
		t.Fatalf("Run(completion invalid) code = %d", code)
	}
	if !strings.Contains(stderr.String(), "unsupported completion shell") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunCompletionRequiresShell(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := New(Dependencies{}, Options{}).Run(
		context.Background(),
		[]string{"completion"},
		&stdout,
		&stderr,
	)

	if code != ExitUsage {
		t.Fatalf("Run(completion missing shell) code = %d", code)
	}
}
