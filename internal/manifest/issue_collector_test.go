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
	"errors"
	"fmt"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestIssueCollectorCollectsDirectAndNestedErrors(t *testing.T) {
	var collector manifest.IssueCollector
	collector.Add(manifest.IssueInvalidValue, "field", "bad %s", "value")
	collector.AddError("nested", fmt.Errorf("wrapped: %w", manifest.NewFieldError(manifest.IssueMissingField, "name", "required")))
	err := collector.Err()
	if err == nil || !strings.Contains(err.Error(), "nested.name") {
		t.Fatalf("unexpected collector error: %v", err)
	}
}

func TestIssueCollectorWrapsPlainErrorsAtPath(t *testing.T) {
	var collector manifest.IssueCollector
	collector.AddError("source", errors.New("boom"))
	issues := collector.Issues()
	if len(issues) != 1 || issues[0].Path != "source" || issues[0].Code != manifest.IssueInvalidValue {
		t.Fatalf("unexpected collected issues: %#v", issues)
	}
}

func TestIssueCollectorHandlesNilAndRootValidationErrors(t *testing.T) {
	var collector manifest.IssueCollector
	collector.AddError("ignored", nil)
	collector.AddError("metadata", &manifest.ValidationError{Issues: []manifest.Issue{{Code: manifest.IssueMissingField, Message: "required"}}})
	collector.AddError("", manifest.NewFieldError(manifest.IssueInvalidValue, "kind", "bad"))
	issues := collector.Issues()
	if len(issues) != 2 || issues[0].Path != "metadata" || issues[1].Path != "kind" {
		t.Fatalf("unexpected collected validation issues: %#v", issues)
	}
}

func TestIssueCollectorIssuesAreDetached(t *testing.T) {
	var collector manifest.IssueCollector
	collector.Add(manifest.IssueInvalidValue, "field", "bad")
	issues := collector.Issues()
	issues[0].Message = "mutated"
	if collector.Issues()[0].Message == "mutated" {
		t.Fatalf("issues accessor leaked internal slice")
	}
}

func TestIssueCollectorErrReturnsNilWhenEmpty(t *testing.T) {
	var collector manifest.IssueCollector
	if err := collector.Err(); err != nil {
		t.Fatalf("empty collector returned error: %v", err)
	}
}
