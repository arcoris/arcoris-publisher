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

package verify

import "testing"

func TestParseGoMod(t *testing.T) {
	data := []byte(`module arcoris.dev/control

require arcoris.dev/foundation v0.1.0
replace arcoris.dev/foundation => ../foundation
`)
	info := parseGoMod(data)
	if info.module != "arcoris.dev/control" {
		t.Fatalf("module = %q", info.module)
	}
	if got := info.requires["arcoris.dev/foundation"]; got != "v0.1.0" {
		t.Fatalf("require = %q", got)
	}
	if len(info.localReplaces) != 1 {
		t.Fatalf("local replaces = %v", info.localReplaces)
	}
}
