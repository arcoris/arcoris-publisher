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

package manifest

import "testing"

func TestHasParentTraversalSegmentDetectsCompleteSegmentsOnly(t *testing.T) {
	if !hasParentTraversalSegment("safe/../secret") {
		t.Fatalf("expected parent traversal segment")
	}
	if hasParentTraversalSegment("safe/..hidden/secret") {
		t.Fatalf("did not expect partial segment to count as traversal")
	}
}
