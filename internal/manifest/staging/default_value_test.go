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

package staging

import "testing"

func TestStringOrDefaultUsesFallbackOnlyWhenNil(t *testing.T) {
	if got := stringOrDefault(nil, "fallback"); got != "fallback" {
		t.Fatalf("stringOrDefault(nil) = %q", got)
	}
	value := ""
	if got := stringOrDefault(&value, "fallback"); got != "" {
		t.Fatalf("stringOrDefault(empty) = %q", got)
	}
}

func TestBoolOrDefaultUsesFallbackOnlyWhenNil(t *testing.T) {
	if got := boolOrDefault(nil, true); !got {
		t.Fatalf("boolOrDefault(nil) = false")
	}
	value := false
	if got := boolOrDefault(&value, true); got {
		t.Fatalf("boolOrDefault(false) = true")
	}
}
