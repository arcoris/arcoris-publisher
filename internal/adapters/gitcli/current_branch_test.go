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

package gitcli

import (
	"context"
	"testing"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestCurrentBranchReadsBranchName(t *testing.T) {
	client := New(&fakeRunner{results: []processport.Result{{Stdout: []byte("main\n")}}}, Options{})

	branch, err := client.CurrentBranch(context.Background(), "/repo")
	if err != nil || branch != "main" {
		t.Fatalf("CurrentBranch() = %q, %v", branch, err)
	}
}

func TestCurrentBranchDetachedHead(t *testing.T) {
	client := New(&fakeRunner{results: []processport.Result{{Stdout: []byte("\n")}}}, Options{})
	_, err := client.CurrentBranch(context.Background(), "/repo")
	assertPortCode(t, err, gitport.CodeRefNotFound)
}
