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
	"errors"
	"testing"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestWrapGitCommandErrorClassifiesStderr(t *testing.T) {
	tests := []struct {
		name   string
		stderr string
		want   string
	}{
		{name: "auth", stderr: "permission denied", want: gitport.CodeAuthenticationFailed.String()},
		{name: "repo not found", stderr: "repository not found", want: gitport.CodeRepositoryNotFound.String()},
		{name: "push rejected", stderr: "non-fast-forward", want: gitport.CodePushRejected.String()},
		{name: "tag exists", stderr: "fatal: tag already exists", want: gitport.CodeTagAlreadyExists.String()},
		{name: "generic", stderr: "boom", want: gitport.CodeCommandFailed.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wrapGitCommandError("failed", processport.Result{Dir: "/repo", Stderr: []byte(tt.stderr)}, errors.New("cause"))
			assertPortCode(t, err, porterr.Code(tt.want))
		})
	}
}
