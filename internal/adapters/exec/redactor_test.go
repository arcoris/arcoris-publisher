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

package exec

import "testing"

func TestRedactorRedactsStringsBytesAndSlices(t *testing.T) {
	redactor := NewRedactor("token", "")

	if got := redactor.RedactString("token value"); got != "<redacted> value" {
		t.Fatalf("RedactString() = %q", got)
	}
	if got := string(redactor.RedactBytes([]byte("token bytes"))); got != "<redacted> bytes" {
		t.Fatalf("RedactBytes() = %q", got)
	}
	values := redactor.RedactSlice([]string{"token", "safe"})
	assertStringSlice(t, values, []string{"<redacted>", "safe"})
}

func TestRedactorReturnsDetachedValues(t *testing.T) {
	redactor := NewRedactor("secret")
	in := []byte("secret")
	out := redactor.RedactBytes(in)
	in[0] = 'x'

	if got := string(out); got != "<redacted>" {
		t.Fatalf("RedactBytes() should detach output, got %q", got)
	}
}

func TestRedactorNilSlice(t *testing.T) {
	if got := NewRedactor("x").RedactSlice(nil); got != nil {
		t.Fatalf("nil RedactSlice() = %#v, want nil", got)
	}
}
