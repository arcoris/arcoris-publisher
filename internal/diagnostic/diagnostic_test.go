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

package diagnostic

import "testing"

type testCode string

const (
	testCodeInvalid testCode = "invalid"
	testCodeMissing testCode = "missing"
)

func TestIssueErrorFormatsPathWhenPresent(t *testing.T) {
	issue := Issue[testCode]{Code: testCodeInvalid, Path: "modules[0]", Message: "bad"}

	if got := issue.Error(); got != "invalid: modules[0]: bad" {
		t.Fatalf("Issue.Error() = %q", got)
	}
}

func TestValidationErrorHelpers(t *testing.T) {
	var nilErr *ValidationError[testCode]
	if !nilErr.Empty() {
		t.Fatal("nil ValidationError should be empty")
	}
	if got := nilErr.Error(); got != "validation failed" {
		t.Fatalf("nil Error() = %q", got)
	}

	err := &ValidationError[testCode]{
		Scope: "test",
		Issues: []Issue[testCode]{
			{Code: testCodeInvalid, Path: "x", Message: "bad"},
			{Code: testCodeMissing, Message: "missing"},
		},
	}
	if err.Empty() {
		t.Fatal("non-empty ValidationError reported empty")
	}
	if !err.Has(testCodeMissing) || err.Has(testCode("other")) {
		t.Fatalf("Has() returned wrong result for %v", err.Issues)
	}
	if got := err.Error(); got != "invalid: x: bad; missing: missing" {
		t.Fatalf("Error() = %q", got)
	}
}

func TestCollectorReturnsDetachedIssues(t *testing.T) {
	collector := NewCollector[testCode]("test")
	collector.Add(testCodeInvalid, "field", "bad %s", "value")

	err := collector.Err().(*ValidationError[testCode])
	err.Issues[0].Code = testCodeMissing

	issues := collector.Issues()
	if issues[0].Code != testCodeInvalid {
		t.Fatalf("collector issue mutated: %v", issues)
	}
}
