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

import (
	"context"
	"testing"
	"time"
)

func TestContextWithOptionalTimeoutReusesParentWithoutTimeout(t *testing.T) {
	parent := context.Background()
	ctx, cancel := contextWithOptionalTimeout(parent, 0)
	defer cancel()

	if ctx != parent {
		t.Fatalf("contextWithOptionalTimeout() should reuse parent when timeout is zero")
	}
}

func TestContextWithOptionalTimeoutCreatesDeadline(t *testing.T) {
	ctx, cancel := contextWithOptionalTimeout(context.Background(), time.Second)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		t.Fatalf("contextWithOptionalTimeout() should create a deadline")
	}
}
