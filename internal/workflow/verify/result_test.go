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

package verify

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestResultReportsFailureAndDetachedSlices(t *testing.T) {
	failed := NewCheckResult("go-test", StatusFailed, SeverityError, "failed")
	result := Result{modules: []ModuleResult{{
		module: manifest.ModuleName("control"),
		checks: []CheckResult{failed},
	}}}

	if result.Passed() || !result.Failed() {
		t.Fatal("Result failure state is wrong")
	}
	if len(result.FailedChecks()) != 1 {
		t.Fatalf("FailedChecks() len = %d", len(result.FailedChecks()))
	}

	modules := result.Modules()
	modules[0].module = "mutated"
	if result.Modules()[0].Module() != "control" {
		t.Fatal("Modules() returned attached slice")
	}
}

func TestCheckResultAccessors(t *testing.T) {
	check := NewCheckResult("go-list", StatusPassed, SeverityInfo, "ok").withPath("go.mod")
	if check.Name() != "go-list" {
		t.Fatalf("Name() = %q", check.Name())
	}
	if check.Status() != StatusPassed {
		t.Fatalf("Status() = %q", check.Status())
	}
	if check.Severity() != SeverityInfo {
		t.Fatalf("check = %#v", check)
	}
	if check.Message() != "ok" || check.Path() != "go.mod" {
		t.Fatalf("check message/path = %#v", check)
	}
}
