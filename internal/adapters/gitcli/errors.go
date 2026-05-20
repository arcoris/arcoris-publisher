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
	"strings"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// gitError creates the structured Git-port error used by this adapter.
func gitError(code porterr.Code, message string, cause error, details porterr.Details) error {
	return porterr.New(porterr.KindGit, code, message, cause).WithDetails(details)
}

// wrapGitCommandError maps stderr text from the Git CLI into stable Git codes.
//
// Git's command-line API often reports important failure classes only in stderr,
// so this classifier intentionally recognizes common phrases while preserving
// the raw stderr in Details for diagnostics.
func wrapGitCommandError(message string, result processport.Result, cause error) error {
	stderr := strings.ToLower(string(result.Stderr))
	code := gitport.CodeCommandFailed
	switch {
	case strings.Contains(stderr, "authentication failed") || strings.Contains(stderr, "permission denied"):
		code = gitport.CodeAuthenticationFailed
	case strings.Contains(stderr, "repository not found") || strings.Contains(stderr, "not found"):
		code = gitport.CodeRepositoryNotFound
	case strings.Contains(stderr, "non-fast-forward") || strings.Contains(stderr, "fetch first") || strings.Contains(stderr, "protected branch") || strings.Contains(stderr, "pre-receive hook declined"):
		code = gitport.CodePushRejected
	case strings.Contains(stderr, "already exists") && strings.Contains(stderr, "tag"):
		code = gitport.CodeTagAlreadyExists
	}
	return gitError(code, message, cause, porterr.Details{"repo": result.Dir, "stderr": strings.TrimSpace(string(result.Stderr))})
}

// trimOutput removes Git's trailing newlines from scalar command output.
func trimOutput(b []byte) string { return strings.TrimSpace(string(b)) }
