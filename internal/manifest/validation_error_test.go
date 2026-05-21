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
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestValidationErrorFormatsAllIssuesInOrder(t *testing.T) {
	err := &manifest.ValidationError{Issues: []manifest.Issue{
		{Code: manifest.IssueMissingField, Path: "metadata.name", Message: "required"},
		{Code: manifest.IssueInvalidValue, Path: "source.repository", Message: "bad"},
	}}
	got := err.Error()
	if !strings.Contains(got, "metadata.name") || !strings.Contains(got, "source.repository") {
		t.Fatalf("unexpected validation error: %q", got)
	}
}

func TestValidationErrorHandlesEmptyIssueList(t *testing.T) {
	if got := (&manifest.ValidationError{}).Error(); got != "manifest validation failed" {
		t.Fatalf("unexpected empty validation error: %q", got)
	}
}
