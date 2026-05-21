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

package registry

import (
	"strings"
	"testing"
)

func TestIssueErrorIncludesPath(t *testing.T) {
	issue := Issue{
		Code:    IssueDuplicateModuleName,
		Path:    "modules[1].name",
		Message: "duplicate",
	}

	if !strings.Contains(issue.Error(), "modules[1].name") {
		t.Fatalf("unexpected issue error: %s", issue.Error())
	}
}

func TestValidationErrorFormatsIssues(t *testing.T) {
	err := (&ValidationError{
		Issues: []Issue{
			{Code: IssueDuplicateModuleName, Message: "duplicate"},
		},
	}).Error()

	if !strings.Contains(err, string(IssueDuplicateModuleName)) {
		t.Fatalf("unexpected validation error: %s", err)
	}
}

func TestEmptyValidationErrorString(t *testing.T) {
	if (&ValidationError{Scope: "registry"}).Error() != "registry validation failed" {
		t.Fatal("unexpected empty validation error")
	}
}
