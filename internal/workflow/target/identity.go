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

package target

import (
	"context"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

const (
	CommitIdentityCodeMissingBoth  = "missing_commit_identity"
	CommitIdentityCodeMissingName  = "missing_commit_user_name"
	CommitIdentityCodeMissingEmail = "missing_commit_user_email"
	CommitIdentityCodeReadFailed   = "commit_identity_read_failed"
)

// CommitIdentity is the effective Git identity that git commit would use.
type CommitIdentity struct {
	Name  string
	Email string
}

// CommitIdentityCheck contains the sanitized result of reading Git identity.
type CommitIdentityCheck struct {
	identity CommitIdentity
	code     string
	message  string
	err      error
}

// CheckCommitIdentity reads and validates the effective Git commit identity.
//
// It intentionally checks only for non-blank user.name and user.email. Git
// accepts a broad range of identity strings, so the workflow only catches the
// missing identity case early instead of imposing repository policy.
func CheckCommitIdentity(ctx context.Context, client git.RepositoryReader, repoDir string) CommitIdentityCheck {
	name, nameOK, err := client.ConfigGet(ctx, repoDir, "user.name")
	if err != nil {
		return CommitIdentityCheck{code: CommitIdentityCodeReadFailed, message: "Git commit identity could not be read", err: err}
	}
	email, emailOK, err := client.ConfigGet(ctx, repoDir, "user.email")
	if err != nil {
		return CommitIdentityCheck{code: CommitIdentityCodeReadFailed, message: "Git commit identity could not be read", err: err}
	}

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	nameOK = nameOK && name != ""
	emailOK = emailOK && email != ""
	switch {
	case !nameOK && !emailOK:
		return CommitIdentityCheck{code: CommitIdentityCodeMissingBoth, message: "missing Git user.name and user.email"}
	case !nameOK:
		return CommitIdentityCheck{code: CommitIdentityCodeMissingName, message: "missing Git user.name"}
	case !emailOK:
		return CommitIdentityCheck{code: CommitIdentityCodeMissingEmail, message: "missing Git user.email"}
	default:
		return CommitIdentityCheck{identity: CommitIdentity{Name: name, Email: email}, message: "Git commit identity is configured"}
	}
}

// Passed reports whether Git has enough identity to create a commit.
func (c CommitIdentityCheck) Passed() bool { return c.code == "" }

// Identity returns the effective identity without exposing it in reports.
func (c CommitIdentityCheck) Identity() CommitIdentity { return c.identity }

// Code returns the stable machine-readable failure code.
func (c CommitIdentityCheck) Code() string { return c.code }

// Message returns a sanitized diagnostic safe for default reports.
func (c CommitIdentityCheck) Message() string { return c.message }

// Err returns an underlying Git config read failure, if any.
func (c CommitIdentityCheck) Err() error { return c.err }
