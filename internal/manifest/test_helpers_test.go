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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func stringPtr(value string) *string { return &value }

func boolPtr(value bool) *bool { return &value }

func requireValidationIssuePaths(t *testing.T, err error, want ...string) {
	t.Helper()

	var validation *manifest.ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}
	if len(validation.Issues) != len(want) {
		t.Fatalf(
			"unexpected issue count: got %d want %d: %#v",
			len(validation.Issues),
			len(want),
			validation.Issues,
		)
	}
	for i, issue := range validation.Issues {
		if issue.Path != want[i] {
			t.Fatalf("unexpected issue path at %d: got %q want %q", i, issue.Path, want[i])
		}
	}
}
