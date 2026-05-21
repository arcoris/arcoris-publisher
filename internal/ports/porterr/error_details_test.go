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

import "testing"

func TestErrorWithDetailsClonesInput(t *testing.T) {
	details := Details{"repo": "target"}
	err := New(KindGit, Code("git_failed"), "git failed", nil).WithDetails(details)
	details["repo"] = "mutated"

	if got := err.Details["repo"]; got != "target" {
		t.Fatalf("WithDetails() repo = %q, want target", got)
	}
}

func TestErrorWithDetailsNilReceiver(t *testing.T) {
	if (*Error)(nil).WithDetails(Details{"key": "value"}) != nil {
		t.Fatalf("WithDetails() on nil receiver should return nil")
	}
}
