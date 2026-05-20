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

func TestDetailsCloneNilForEmptyMap(t *testing.T) {
	var nilDetails Details
	if clone := nilDetails.Clone(); clone != nil {
		t.Fatalf("nil Clone() = %#v, want nil", clone)
	}

	emptyDetails := Details{}
	if clone := emptyDetails.Clone(); clone != nil {
		t.Fatalf("empty Clone() = %#v, want nil", clone)
	}
}

func TestDetailsCloneDetachesMap(t *testing.T) {
	details := Details{"repo": "target", "branch": "main"}
	clone := details.Clone()
	details["repo"] = "mutated"

	if got := clone["repo"]; got != "target" {
		t.Fatalf("clone repo = %q, want target", got)
	}
	if got := clone["branch"]; got != "main" {
		t.Fatalf("clone branch = %q, want main", got)
	}
}
