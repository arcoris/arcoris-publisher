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
	"fmt"
	"strings"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

// validateGitPositional rejects values that can be reinterpreted as Git CLI
// options or contain control characters. Process execution already passes an
// argv vector rather than a shell command, so whitespace inside a normal value
// is not an injection boundary and remains allowed.
func validateGitPositional(kind string, value string) error {
	if value == "" {
		return invalidGitArgument(kind, "value is empty")
	}
	if strings.HasPrefix(value, "-") {
		return invalidGitArgument(kind, "value starts with an option prefix")
	}
	for _, r := range value {
		if r == 0 || r < 0x20 || r == 0x7f {
			return invalidGitArgument(kind, fmt.Sprintf("value contains control character %q", r))
		}
	}
	return nil
}

func validatedRemote(remote string) (string, error) {
	remote = defaultRemote(remote)
	if err := validateGitPositional("remote", remote); err != nil {
		return "", err
	}
	return remote, nil
}

func invalidGitArgument(kind string, reason string) error {
	return gitError(
		gitport.CodeInvalidArgument,
		"invalid Git "+kind+" argument",
		fmt.Errorf("%s", reason),
		nil,
	)
}
