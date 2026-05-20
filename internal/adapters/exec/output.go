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

import "bytes"

// limitedBuffer captures at most limit bytes.
//
// A zero or negative limit means unlimited capture. Write always reports that it
// consumed the full input slice even when the internal buffer is already full;
// this keeps child process writes from seeing artificial short writes caused by
// diagnostic capture limits.
type limitedBuffer struct {
	limit     int64
	buf       bytes.Buffer
	truncated bool
}

// newLimitedBuffer creates a capture buffer with the requested byte limit.
func newLimitedBuffer(limit int64) *limitedBuffer { return &limitedBuffer{limit: limit} }

// Write appends p until the configured limit is reached.
func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.limit <= 0 {
		_, err := b.buf.Write(p)
		return len(p), err
	}
	remaining := b.limit - int64(b.buf.Len())
	if remaining <= 0 {
		b.truncated = b.truncated || len(p) > 0
		return len(p), nil
	}
	if int64(len(p)) > remaining {
		_, err := b.buf.Write(p[:remaining])
		b.truncated = true
		return len(p), err
	}
	_, err := b.buf.Write(p)
	return len(p), err
}

// Bytes returns a detached snapshot of captured bytes.
func (b *limitedBuffer) Bytes() []byte { return append([]byte(nil), b.buf.Bytes()...) }

// Truncated reports whether any bytes were omitted because of the limit.
func (b *limitedBuffer) Truncated() bool { return b.truncated }
