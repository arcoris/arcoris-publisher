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
	"errors"
	"testing"
)

func TestNewRejectsDuplicateIndexes(t *testing.T) {
	_, err := New(duplicatePublicationSet(t))
	if err == nil {
		t.Fatal("expected duplicate validation error")
	}

	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	want := map[IssueCode]bool{
		IssueDuplicateModuleName: false,
		IssueDuplicateModulePath: false,
		IssueDuplicateSourceDir:  false,
		IssueDuplicateRepository: false,
	}
	for _, issue := range validation.Issues {
		if _, ok := want[issue.Code]; ok {
			want[issue.Code] = true
		}
	}

	for code, seen := range want {
		if !seen {
			t.Fatalf("missing issue code %s in %#v", code, validation.Issues)
		}
	}
}

func TestInvalidPublicationSetBuildsValidationError(t *testing.T) {
	err := invalidPublicationSet("set", "bad set")
	if err == nil {
		t.Fatal("expected validation error")
	}

	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if validation.Issues[0].Code != IssueInvalidPublicationSet {
		t.Fatalf("unexpected issue: %#v", validation.Issues[0])
	}
}
