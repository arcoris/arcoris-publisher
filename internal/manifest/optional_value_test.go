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

func TestStringValueReturnsFallbackOnlyWhenNil(t *testing.T) {
	if got := stringValue(nil, "fallback"); got != "fallback" {
		t.Fatalf("stringValue(nil) = %q", got)
	}
	value := ""
	if got := stringValue(&value, "fallback"); got != "" {
		t.Fatalf("stringValue(empty) = %q", got)
	}
}

func TestBoolValueReturnsFallbackOnlyWhenNil(t *testing.T) {
	if got := boolValue(nil, true); !got {
		t.Fatalf("boolValue(nil) = false")
	}
	value := false
	if got := boolValue(&value, true); got {
		t.Fatalf("boolValue(false) = true")
	}
}
