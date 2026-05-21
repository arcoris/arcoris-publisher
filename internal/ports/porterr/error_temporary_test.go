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

func TestErrorWithTemporaryReturnsUpdatedCopy(t *testing.T) {
	base := New(KindProcess, Code("process_failed"), "failed", nil)
	err := base.WithTemporary(true)

	if !err.Temporary {
		t.Fatalf("WithTemporary(true) did not set temporary flag")
	}
	if base.Temporary {
		t.Fatalf("WithTemporary should not mutate receiver")
	}
}

func TestErrorWithTemporaryNilReceiver(t *testing.T) {
	if (*Error)(nil).WithTemporary(true) != nil {
		t.Fatalf("WithTemporary() on nil receiver should return nil")
	}
}
