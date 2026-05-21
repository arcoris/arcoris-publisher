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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestMustReturnsAssignmentsForValidSpec(t *testing.T) {
	assignments := Must(testRegistry(t), AssignmentSpec{Release: "v0.5.0"})

	if version, ok := assignments.VersionOfModule(testModuleName(t, "foundation")); !ok || version != "v0.5.0" {
		t.Fatalf("VersionOfModule() = %q, %v", version, ok)
	}
}

func TestAssignVersionRejectsInvalidModule(t *testing.T) {
	_, err := assignVersion([]manifest.Module{{}}, manifest.VersionPolicyReleaseTrain, MustVersion("v0.1.0"))
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueInvalidAssignment) {
		t.Fatalf("issues = %#v, want invalid assignment", validationErr.Issues)
	}
}

func TestAssignVersionPropagatesAssignmentValidationError(t *testing.T) {
	_, err := assignVersion(nil, manifest.VersionPolicy("unknown"), MustVersion("v0.1.0"))
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueUnsupportedPolicy) {
		t.Fatalf("issues = %#v, want unsupported policy", validationErr.Issues)
	}
}
