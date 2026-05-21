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

package manifest_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestIssueErrorIncludesPathWhenPresent(t *testing.T) {
	issue := manifest.Issue{
		Code:    manifest.IssueInvalidValue,
		Path:    "metadata.name",
		Message: "bad value",
	}
	if got := issue.Error(); got != "invalid_value: metadata.name: bad value" {
		t.Fatalf("unexpected issue error: %q", got)
	}
}

func TestIssueErrorOmitsEmptyPath(t *testing.T) {
	issue := manifest.Issue{Code: manifest.IssueInvalidValue, Message: "bad value"}
	if got := issue.Error(); got != "invalid_value: bad value" {
		t.Fatalf("unexpected issue error: %q", got)
	}
}
