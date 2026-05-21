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

package porterr

import (
	"errors"
	"testing"
)

func TestNewPopulatesStructuredErrorFields(t *testing.T) {
	cause := errors.New("root cause")
	err := New(KindGit, Code("git_failed"), "git failed", cause)

	if err.Kind != KindGit || err.Code != Code("git_failed") || err.Message != "git failed" || err.Cause != cause {
		t.Fatalf("New() = %#v", err)
	}
	if err.Details != nil || err.Temporary {
		t.Fatalf("New() should not attach optional metadata: %#v", err)
	}
}
