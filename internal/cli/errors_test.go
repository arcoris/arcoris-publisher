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

package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorWrapsCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("boom")
	err := &Error{Code: CodeReportFailed, Message: "render failed", Cause: cause}

	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is() = false for %v", err)
	}
	if !strings.Contains(err.Error(), string(CodeReportFailed)) || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("Error() = %q", err.Error())
	}
}

func TestNilErrorIsSafe(t *testing.T) {
	t.Parallel()

	var err *Error
	if err.Error() != "" || err.Unwrap() != nil {
		t.Fatalf("nil Error behavior changed")
	}
}
