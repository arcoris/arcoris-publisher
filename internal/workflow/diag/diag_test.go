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

package diag

import (
	"strings"
	"testing"
)

type testCode string

const (
	testInvalid testCode = "invalid"
	testMissing testCode = "missing"
)

func TestCollectorErrIsDetached(t *testing.T) {
	collector := NewCollector[testCode]("test")
	collector.Add(testInvalid, "", "path", "bad %s", "value")

	err := collector.Err().(*ValidationError[testCode])
	err.Issues[0].Code = testMissing

	if collector.Issues()[0].Code != testInvalid {
		t.Fatal("collector issue mutated through returned error")
	}
}

func TestValidationErrorFormattingAndHas(t *testing.T) {
	err := &ValidationError[testCode]{
		Scope: "test",
		Issues: []Issue[testCode]{
			{Code: testInvalid, Path: "a", Message: "bad"},
			{Code: testMissing, Message: "missing"},
		},
	}

	if !err.Has(testMissing) {
		t.Fatal("Has(missing) = false")
	}
	if got := err.Error(); !strings.Contains(got, "test validation failed with 2 issues") {
		t.Fatalf("Error() = %q", got)
	}
}
