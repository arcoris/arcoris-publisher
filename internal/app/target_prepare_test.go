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

package app

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
)

func TestPrepareTargetsRunsTargetWorkflow(t *testing.T) {
	app, fakeGit := appFixture(t)

	result, err := app.PrepareTargets(context.Background(), appRequest())

	if err != nil {
		t.Fatalf("PrepareTargets() error = %v", err)
	}
	if result.TargetPrepare().Failed() {
		t.Fatalf("target prepare failed: %#v", result.TargetPrepare())
	}
	if result.Plan().Empty() {
		t.Fatal("target prepare did not keep the built plan")
	}
	if !gitCallSeen(fakeGit.Calls, "fetch") {
		t.Fatalf("git calls = %#v, want fetch", fakeGit.Calls)
	}
}

func gitCallSeen(calls []porttest.GitCall, op string) bool {
	for _, call := range calls {
		if call.Op == op {
			return true
		}
	}
	return false
}
