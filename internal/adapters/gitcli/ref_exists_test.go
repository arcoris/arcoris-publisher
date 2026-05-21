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

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestRefExistsUsesAllowedExitCodes(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{ExitCode: 0}, {ExitCode: 1}}}
	client := New(runner, Options{})

	exists, err := client.RefExists(context.Background(), "/repo", "HEAD")
	if err != nil || !exists {
		t.Fatalf("RefExists() = %v, %v", exists, err)
	}
	exists, err = client.RefExists(context.Background(), "/repo", "missing")
	if err != nil || exists {
		t.Fatalf("RefExists() missing = %v, %v", exists, err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"rev-parse", "--verify", "--quiet", "HEAD"})
	assertStringSlice(t, intsOf(runner.specs[0].AllowedExitCodes), []string{"0", "1"})
}
