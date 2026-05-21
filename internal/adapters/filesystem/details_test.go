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

package filesystem

import "testing"

func TestPorterrDetailsBuildsPairsAndIgnoresOddTail(t *testing.T) {
	details := porterrDetails("path", "/repo", "dangling")

	if details["path"] != "/repo" {
		t.Fatalf("porterrDetails() = %#v", details)
	}
	if _, ok := details["dangling"]; ok {
		t.Fatalf("porterrDetails() should ignore odd trailing values: %#v", details)
	}
}
