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

package publish

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

func TestPublishRejectsInvalidRequest(t *testing.T) {
	_, err := New(Dependencies{}, Options{}).Publish(context.Background(), Request{})

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodeInvalidRequest {
		t.Fatalf("Code = %q", got.Code)
	}
}

func TestBranchRefspecUsesTargetBranch(t *testing.T) {
	got := branchRefspec(manifest.BranchName("release/v1"))
	want := git.RefSpec("refs/heads/release/v1:refs/heads/release/v1")

	if got != want {
		t.Fatalf("branchRefspec() = %q, want %q", got, want)
	}
}
