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

package plan

import "testing"

func TestValidationErrorHelpers(t *testing.T) {
	var nilErr *ValidationError
	if !nilErr.Empty() {
		t.Fatal("nil ValidationError should be empty")
	}
	if nilErr.Has(IssueEmptyPlan) {
		t.Fatal("nil ValidationError unexpectedly has issue")
	}
	if nilErr.Error() != "plan validation failed" {
		t.Fatalf("nil Error() = %q", nilErr.Error())
	}
	err := &ValidationError{
		Issues: []Issue{{
			Code:    IssueEmptyPlan,
			Path:    "modules",
			Message: "empty",
		}},
	}
	if err.Empty() {
		t.Fatal("non-empty ValidationError reported empty")
	}
	if !err.Has(IssueEmptyPlan) {
		t.Fatal("ValidationError.Has(empty_plan) = false")
	}
	if got := err.Issues[0].Error(); got != "empty_plan: modules: empty" {
		t.Fatalf("Issue.Error() = %q", got)
	}
	issue := Issue{Code: IssueInvalidRequest, Message: "bad"}
	if got := issue.Error(); got != "invalid_request: bad" {
		t.Fatalf("Issue.Error() without path = %q", got)
	}
	if err.Has(IssueGraphOrder) {
		t.Fatal("ValidationError.Has(graph_order) = true")
	}
	if got := err.Error(); got != "empty_plan: modules: empty" {
		t.Fatalf("ValidationError.Error() = %q", got)
	}
	empty := &ValidationError{}
	if got := empty.Error(); got != "plan validation failed" {
		t.Fatalf("empty ValidationError.Error() = %q", got)
	}
}
