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

package versioning

import (
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestAssignmentAccessors(t *testing.T) {
	req := mustRequest(t, "v0.3.0", string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
	)
	assignments, err := Assign(req)
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	if assignments.Empty() {
		t.Fatalf("assignments unexpectedly empty")
	}
	if assignments.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", assignments.Len())
	}
	mv, ok := assignments.ModuleVersion(manifest.ModuleName("foundation"))
	if !ok {
		t.Fatalf("ModuleVersion(foundation) missing")
	}
	assertModuleVersion(t, mv, "foundation", "arcoris.dev/foundation", "v0.3.0")
	if _, ok := assignments.ModuleVersion(manifest.ModuleName("missing")); ok {
		t.Fatalf("ModuleVersion(missing) found")
	}
	if _, ok := assignments.RequirementsFor(manifest.ModuleName("missing")); ok {
		t.Fatalf("RequirementsFor(missing) found")
	}
	if _, ok := assignments.RequirementMapFor(manifest.ModuleName("missing")); ok {
		t.Fatalf("RequirementMapFor(missing) found")
	}
}

func TestZeroAssignments(t *testing.T) {
	var assignments Assignments
	if !assignments.Empty() {
		t.Fatalf("zero Assignments should be empty")
	}
	if assignments.Len() != 0 {
		t.Fatalf("zero Len() = %d", assignments.Len())
	}
	if _, ok := assignments.VersionOf(manifest.ModuleName("missing")); ok {
		t.Fatalf("zero VersionOf returned ok")
	}
}

func TestValidationErrorAccessors(t *testing.T) {
	issue := Issue{Code: IssueInvalidRequest, Path: "request", Message: "bad"}
	if got := issue.Error(); !strings.Contains(got, "request") || !strings.Contains(got, "bad") {
		t.Fatalf("Issue.Error() = %q", got)
	}
	issue.Path = ""
	if got := issue.Error(); !strings.Contains(got, "bad") {
		t.Fatalf("Issue.Error() without path = %q", got)
	}
	var nilErr *ValidationError
	if !nilErr.Empty() || nilErr.Has(IssueInvalidRequest) {
		t.Fatalf("nil ValidationError accessors are wrong")
	}
	err := &ValidationError{Issues: []Issue{{Code: IssueInvalidVersion, Message: "invalid"}}}
	if err.Empty() || !err.Has(IssueInvalidVersion) || err.Has(IssueInvalidRequest) {
		t.Fatalf("ValidationError accessors are wrong")
	}
	if got := err.Error(); !strings.Contains(got, "invalid") {
		t.Fatalf("ValidationError.Error() = %q", got)
	}
	empty := &ValidationError{}
	if got := empty.Error(); got == "" {
		t.Fatalf("empty ValidationError.Error() is empty")
	}
}

func TestMustPanicsOnInvalidVersion(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("Must did not panic")
		}
	}()
	_ = Must("bad")
}
