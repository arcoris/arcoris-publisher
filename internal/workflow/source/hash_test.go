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

package source

import "testing"

func TestHashBasics(t *testing.T) {
	if !Hash("").IsZero() {
		t.Fatal("empty hash is not zero")
	}
	if Hash("sha256:abc").String() != "sha256:abc" {
		t.Fatalf("String() = %q", Hash("sha256:abc").String())
	}
}

func TestHashBytesIsStableAndTyped(t *testing.T) {
	first := hashBytes("file", "a", "bc")
	second := hashBytes("file", "a", "bc")
	if first != second {
		t.Fatalf("hashBytes() = %q and %q", first, second)
	}

	otherKind := hashBytes("dir", "a", "bc")
	if first == otherKind {
		t.Fatalf("hashBytes() ignored kind: %q", first)
	}
}

func TestCombineHashesSkipsZeroValues(t *testing.T) {
	combined := combineHashes("module", []Hash{"a", "", "b"})
	withoutZero := combineHashes("module", []Hash{"a", "b"})
	if combined != withoutZero {
		t.Fatalf("combineHashes() = %q, want %q", combined, withoutZero)
	}

	if got := combineHashes("module", []Hash{"", ""}); got != "" {
		t.Fatalf("combineHashes(empty) = %q", got)
	}
}
