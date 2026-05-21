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

package process

import "testing"

func TestSpec_WithSensitiveValuesClonesInput(t *testing.T) {
	values := []string{"secret"}
	spec := Spec{Name: "git"}.WithSensitiveValues(values...)
	values[0] = "mutated"

	if got := spec.SensitiveValues[0]; got != "secret" {
		t.Fatalf("expected detached sensitive values, got %q", got)
	}
}

func TestSpec_WithAllowedExitCodesClonesInput(t *testing.T) {
	codes := []int{0, 2}
	spec := Spec{Name: "grep"}.WithAllowedExitCodes(codes...)
	codes[1] = 99

	if got := spec.AllowedExitCodes[1]; got != 2 {
		t.Fatalf("expected detached allowed exit codes, got %d", got)
	}
}

func TestSpec_WithMethodsKeepOtherFields(t *testing.T) {
	spec := Spec{Name: "git", Args: []string{"status"}, Dir: "/repo"}.
		WithSensitiveValues("token").
		WithAllowedExitCodes(0, 1)

	if spec.Name != "git" || spec.Args[0] != "status" || spec.Dir != "/repo" {
		t.Fatalf("with methods should preserve existing fields: %#v", spec)
	}
	if len(spec.SensitiveValues) != 1 || spec.SensitiveValues[0] != "token" {
		t.Fatalf("unexpected sensitive values: %#v", spec.SensitiveValues)
	}
	if len(spec.AllowedExitCodes) != 2 || spec.AllowedExitCodes[1] != 1 {
		t.Fatalf("unexpected allowed exit codes: %#v", spec.AllowedExitCodes)
	}
}
