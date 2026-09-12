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
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

// validateRollbackJournal verifies persisted recovery input before it is used
// by the rollback workflow. Structural decoding keeps unknown status values
// inspectable, while rollback accepts only states and identifiers understood by
// the current binary.
func validateRollbackJournal(journal TransactionJournal) error {
	if err := validateTransactionID(journal.ID); err != nil {
		return fmt.Errorf("invalid transaction id: %w", err)
	}
	if !knownTransactionStatus(journal.Status) {
		return fmt.Errorf("unsupported transaction status %q", journal.Status)
	}
	if err := validateRecoveryRemote(journal.Remote); err != nil {
		return err
	}

	seenModules := make(map[string]struct{}, len(journal.Modules))
	for i := range journal.Modules {
		if err := validateRollbackModule(journal, i, seenModules); err != nil {
			return err
		}
	}
	return nil
}

func validateRollbackModule(journal TransactionJournal, index int, seenModules map[string]struct{}) error {
	mod := journal.Modules[index]
	moduleName := mod.Module.String()
	if moduleName == "" {
		return fmt.Errorf("module[%d] has empty module name", index)
	}
	if _, duplicate := seenModules[moduleName]; duplicate {
		return fmt.Errorf("module[%d] duplicates module %q", index, moduleName)
	}
	seenModules[moduleName] = struct{}{}

	if mod.WorktreeDir == "" {
		return fmt.Errorf("module %s has empty worktree path", mod.Module)
	}
	if err := validateRecoveryWorktree(mod.WorktreeDir); err != nil {
		return fmt.Errorf("module %s worktree path is invalid: %w", mod.Module, err)
	}

	if mod.TargetBranch == "" {
		return fmt.Errorf("module %s has empty target branch", mod.Module)
	}
	if err := validateGitRef(mod.FinalBranchRef); err != nil {
		return fmt.Errorf("module %s final branch ref is invalid: %w", mod.Module, err)
	}
	if !strings.HasPrefix(mod.FinalBranchRef, "refs/heads/") {
		return fmt.Errorf("module %s final branch ref is outside refs/heads", mod.Module)
	}
	if want := branchRef(mod.TargetBranch); mod.FinalBranchRef != want {
		return fmt.Errorf("module %s final branch ref %q does not match target branch %q", mod.Module, mod.FinalBranchRef, mod.TargetBranch)
	}

	if err := validateGitRef(mod.CandidateBranchRef); err != nil {
		return fmt.Errorf("module %s candidate ref is invalid: %w", mod.Module, err)
	}
	if want := candidateRef(journal.ID, mod.Module); mod.CandidateBranchRef != want {
		return fmt.Errorf("module %s candidate ref %q does not match transaction namespace", mod.Module, mod.CandidateBranchRef)
	}

	if mod.FinalTagRef != "" {
		if err := validateGitRef(mod.FinalTagRef); err != nil {
			return fmt.Errorf("module %s final tag ref is invalid: %w", mod.Module, err)
		}
		if !strings.HasPrefix(mod.FinalTagRef, "refs/tags/") {
			return fmt.Errorf("module %s final tag ref is outside refs/tags", mod.Module)
		}
	}

	for _, object := range []struct {
		name string
		hash git.CommitHash
	}{
		{name: "local base", hash: mod.LocalBaseHead},
		{name: "remote base", hash: mod.RemoteBaseCommit},
		{name: "created commit", hash: mod.CreatedCommit},
		{name: "remote tag", hash: mod.RemoteTagHash},
	} {
		if err := validateRecoveryObjectName(object.hash); err != nil {
			return fmt.Errorf("module %s %s object is invalid: %w", mod.Module, object.name, err)
		}
	}

	if mod.RemoteBaseExists && mod.RemoteBaseCommit == "" {
		return fmt.Errorf("module %s records an existing remote base without its object id", mod.Module)
	}
	if (mod.CandidatePushed || mod.FinalBranchPromoted || mod.LocalTagCreated || mod.RemoteTagPushed) && mod.CreatedCommit == "" {
		return fmt.Errorf("module %s records published state without its created commit", mod.Module)
	}
	if (mod.LocalTagCreated || mod.RemoteTagPushed) && mod.FinalTagRef == "" {
		return fmt.Errorf("module %s records tag state without a final tag ref", mod.Module)
	}
	if mod.RemoteTagHash != "" && !mod.RemoteTagPushed {
		return fmt.Errorf("module %s records a remote tag object without a pushed remote tag", mod.Module)
	}
	return nil
}

func validateRecoveryWorktree(path string) error {
	clean := filepath.Clean(path)
	if !filepath.IsAbs(clean) {
		return fmt.Errorf("path is not absolute")
	}
	if filepath.Dir(clean) == clean {
		return fmt.Errorf("path is a filesystem root")
	}
	if clean != path {
		return fmt.Errorf("path is not canonical")
	}
	return nil
}

func validateRecoveryRemote(remote string) error {
	if remote == "" {
		return nil
	}
	if strings.HasPrefix(remote, "-") {
		return fmt.Errorf("remote name %q starts with an option prefix", remote)
	}
	for _, r := range remote {
		if r == 0 || unicode.IsControl(r) || unicode.IsSpace(r) {
			return fmt.Errorf("remote name %q contains whitespace or control characters", remote)
		}
	}
	return nil
}

func validateRecoveryObjectName(hash git.CommitHash) error {
	if hash == "" {
		return nil
	}
	value := hash.String()
	if len(value) < 7 || len(value) > 64 {
		return fmt.Errorf("object id length %d is outside the supported range", len(value))
	}
	for _, r := range value {
		if r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F' {
			continue
		}
		return fmt.Errorf("object id contains non-hexadecimal character %q", r)
	}
	return nil
}
