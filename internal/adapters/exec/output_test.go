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

func TestLimitedBufferUnlimited(t *testing.T) {
	buffer := newLimitedBuffer(0)
	n, err := buffer.Write([]byte("abcdef"))
	if err != nil || n != 6 {
		t.Fatalf("Write() = %d, %v; want 6, nil", n, err)
	}
	if got := string(buffer.Bytes()); got != "abcdef" {
		t.Fatalf("Bytes() = %q, want abcdef", got)
	}
	if buffer.Truncated() {
		t.Fatalf("unlimited buffer must not report truncation")
	}
}

func TestLimitedBufferTruncatesButReportsFullWrite(t *testing.T) {
	buffer := newLimitedBuffer(4)
	n, err := buffer.Write([]byte("abcdef"))
	if err != nil || n != 6 {
		t.Fatalf("Write() = %d, %v; want 6, nil", n, err)
	}
	if got := string(buffer.Bytes()); got != "abcd" {
		t.Fatalf("Bytes() = %q, want abcd", got)
	}
	if !buffer.Truncated() {
		t.Fatalf("limited buffer should report truncation")
	}
}

func TestLimitedBufferBytesAreDetached(t *testing.T) {
	buffer := newLimitedBuffer(0)
	_, _ = buffer.Write([]byte("abc"))
	data := buffer.Bytes()
	data[0] = 'x'

	if got := string(buffer.Bytes()); got != "abc" {
		t.Fatalf("Bytes() returned shared storage: %q", got)
	}
}
