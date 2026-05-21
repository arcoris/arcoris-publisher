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

package target

import (
	"strings"
	"testing"
)

func TestValidationErrorFormattingAndHas(t *testing.T) {
	err := &ValidationError{Issues: []Issue{
		{Code: IssueInvalidRequest, Path: "plan", Message: "plan is empty"},
		{Code: IssueWorktreeDirty, Path: "/tmp/w", Message: "dirty"},
	}}

	if !err.Has(IssueWorktreeDirty) {
		t.Fatal("Has(worktree_dirty) = false")
	}
	if got := err.Error(); !strings.Contains(got, "2 issues") {
		t.Fatalf("Error() = %q", got)
	}
}

func TestIssueCollectorReturnsDetachedIssues(t *testing.T) {
	collector := newIssueCollector()
	collector.Add(IssueInvalidRequest, "", "plan", "plan is empty")

	err := collector.Err().(*ValidationError)
	err.Issues[0].Code = IssueWorktreeDirty

	if collector.Issues()[0].Code != IssueInvalidRequest {
		t.Fatal("collector issue mutated through error")
	}
}
