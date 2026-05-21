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
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestAssignmentsValidateRejectsUnsupportedPolicy(t *testing.T) {
	assignments := Assignments{
		policy: manifest.VersionPolicy("unknown"),
		items: []ModuleVersion{{
			module:     testModuleName(t, "foundation"),
			modulePath: testModulePath(t, "arcoris.dev/foundation"),
			version:    MustVersion("v0.1.0"),
		}},
	}

	validationErr := mustValidationError(t, assignments.Validate())
	if !hasIssueCode(validationErr.Issues, IssueUnsupportedPolicy) {
		t.Fatalf("issues = %#v, want unsupported policy", validationErr.Issues)
	}
}

func TestAssignmentsValidatorWrapsPlainItemError(t *testing.T) {
	validator := newAssignmentsValidator(Assignments{})
	validator.addItemValidationError("items[0]", errors.New("plain failure"))

	if len(validator.issues) != 1 {
		t.Fatalf("issues len = %d, want 1", len(validator.issues))
	}
	if validator.issues[0].Path != "items[0]" || validator.issues[0].Message != "plain failure" {
		t.Fatalf("issue = %#v", validator.issues[0])
	}
}
