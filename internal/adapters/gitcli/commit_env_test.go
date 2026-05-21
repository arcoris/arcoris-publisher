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
	"testing"
	"time"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

func TestCommitEnvSkipsEmptyValuesAndFormatsDates(t *testing.T) {
	env := commitEnv(gitport.CommitOptions{
		AuthorName: "Author",
		AuthorDate: time.Unix(123, 0).UTC(),
	})

	assertStringSlice(t, env, []string{"GIT_AUTHOR_NAME=Author", "GIT_AUTHOR_DATE=123 +0000"})
}

func TestAppendOptionalEnvAndGitTimestamp(t *testing.T) {
	env := appendOptionalEnv(nil, "A", "")
	if len(env) != 0 {
		t.Fatalf("appendOptionalEnv(empty) = %#v, want empty", env)
	}
	env = appendOptionalEnv(env, "A", "1")
	assertStringSlice(t, env, []string{"A=1"})
	if got := gitTimestamp(42); got != "42 +0000" {
		t.Fatalf("gitTimestamp() = %q, want 42 +0000", got)
	}
}
