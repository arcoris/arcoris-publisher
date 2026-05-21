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
)

func TestValidationErrorFormattingWithPath(t *testing.T) {
	err := &ValidationError{Issues: []Issue{{Code: IssueInvalidAssignment, Path: "release", Message: "bad"}}}

	if got := err.Error(); got != "release: bad" {
		t.Fatalf("Error() = %q, want path-prefixed message", got)
	}
}

func TestValidationErrorFormatting(t *testing.T) {
	if got := (*ValidationError)(nil).Error(); got != "versioning validation failed" {
		t.Fatalf("unexpected nil error string: %q", got)
	}
	one := (&ValidationError{Issues: []Issue{{Code: IssueInvalidAssignment, Message: "bad"}}}).Error()
	if one != "bad" {
		t.Fatalf("unexpected one issue without path: %q", one)
	}
	many := (&ValidationError{Issues: []Issue{{Message: "a"}, {Message: "b"}}}).Error()
	if !strings.Contains(many, "2 issues") {
		t.Fatalf("unexpected many issue summary: %q", many)
	}
}
