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

package source

import (
	"strings"
	"testing"
)

func TestValidationErrorFormatting(t *testing.T) {
	var nilErr *ValidationError
	if got := nilErr.Error(); got != "workflow validation failed" {
		t.Fatalf("nil Error() = %q", got)
	}
	if nilErr.Has(IssueInvalidRequest) {
		t.Fatal("nil Has(invalid_request) = true")
	}

	empty := &ValidationError{Scope: "source"}
	if got := empty.Error(); got != "source validation failed" {
		t.Fatalf("empty Error() = %q", got)
	}

	single := &ValidationError{Scope: "source", Issues: []Issue{{
		Code:    IssueInvalidRequest,
		Message: "bad request",
	}}}
	if got := single.Error(); got != "bad request" {
		t.Fatalf("single Error() = %q", got)
	}

	singleWithPath := &ValidationError{Scope: "source", Issues: []Issue{{
		Code:    IssueInvalidRequest,
		Path:    "request",
		Message: "bad request",
	}}}
	if got := singleWithPath.Error(); got != "request: bad request" {
		t.Fatalf("single path Error() = %q", got)
	}

	many := &ValidationError{Scope: "source", Issues: []Issue{
		{Code: IssueInvalidRequest, Path: "request", Message: "bad request"},
		{Code: IssueDirtySource, Message: "dirty"},
	}}
	if got := many.Error(); !strings.Contains(got, "2 issues") {
		t.Fatalf("many Error() = %q", got)
	}
	if many.Has(IssueDetachedHead) {
		t.Fatal("Has(detached_head) = true")
	}
}

func TestIssueCollectorReturnsDetachedValidationError(t *testing.T) {
	collector := newIssueCollector()
	collector.Add(IssueInvalidRequest, "", "request", "bad %s", "request")

	err := collector.Err()
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("err() type = %T", err)
	}

	ve.Issues[0].Code = IssueDetachedHead
	issues := collector.Issues()
	if issues[0].Code != IssueInvalidRequest {
		t.Fatalf("collector issue was mutated: %v", issues)
	}
}
