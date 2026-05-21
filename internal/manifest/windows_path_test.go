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

func TestLooksWindowsAbsoluteDetectsDrivePaths(t *testing.T) {
	for _, value := range []string{"C:/work", "d:\\work"} {
		if !looksWindowsAbsolute(value) {
			t.Fatalf("looksWindowsAbsolute(%q) = false", value)
		}
	}
}

func TestLooksWindowsAbsoluteIgnoresRelativeColonPaths(t *testing.T) {
	for _, value := range []string{"release:v1", "1:/not-drive", "C:relative"} {
		if looksWindowsAbsolute(value) {
			t.Fatalf("looksWindowsAbsolute(%q) = true", value)
		}
	}
}

func TestIsASCIILetterDetectsOnlyLetters(t *testing.T) {
	if !isASCIILetter('A') || !isASCIILetter('z') || isASCIILetter('1') {
		t.Fatalf("unexpected ASCII letter classification")
	}
}
