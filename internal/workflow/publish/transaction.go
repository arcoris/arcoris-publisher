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
	"hash/fnv"
	"strings"
	"time"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

const transactionSchemaVersion = 1

// TransactionID identifies one durable publish transaction.
type TransactionID string

// String returns the transaction identifier text.
func (id TransactionID) String() string { return string(id) }

// TransactionStatus identifies the saga phase recorded in the publish journal.
type TransactionStatus string

const (
	// TransactionStatusPending is written before any Git mutation.
	TransactionStatusPending TransactionStatus = "pending"
	// TransactionStatusPreflighted means all module preflight checks passed.
	TransactionStatusPreflighted TransactionStatus = "preflighted"
	// TransactionStatusSnapshotted means rollback refs were captured.
	TransactionStatusSnapshotted TransactionStatus = "snapshotted"
	// TransactionStatusCommittedLocally means publication commits were created.
	TransactionStatusCommittedLocally TransactionStatus = "committed_locally"
	// TransactionStatusCandidatesPushed means candidate refs exist remotely.
	TransactionStatusCandidatesPushed TransactionStatus = "candidates_pushed"
	// TransactionStatusPromoting means final branch promotion is in progress.
	TransactionStatusPromoting TransactionStatus = "promoting"
	// TransactionStatusBranchesPromoted means final branches were updated.
	TransactionStatusBranchesPromoted TransactionStatus = "branches_promoted"
	// TransactionStatusTagging means final tag publication is in progress.
	TransactionStatusTagging TransactionStatus = "tagging"
	// TransactionStatusCommitted means final refs were published successfully.
	TransactionStatusCommitted TransactionStatus = "committed"
	// TransactionStatusFailed means publication failed before rollback began.
	TransactionStatusFailed TransactionStatus = "failed"
	// TransactionStatusRollingBack means compensating rollback is in progress.
	TransactionStatusRollingBack TransactionStatus = "rolling_back"
	// TransactionStatusRolledBack means rollback completed successfully.
	TransactionStatusRolledBack TransactionStatus = "rolled_back"
	// TransactionStatusRollbackFailed means rollback could not safely complete.
	TransactionStatusRollbackFailed TransactionStatus = "rollback_failed"
)

// Terminal reports whether a transaction no longer blocks a new publish.
func (s TransactionStatus) Terminal() bool {
	switch s {
	case TransactionStatusCommitted, TransactionStatusRolledBack:
		return true
	default:
		return false
	}
}

// RollbackStatus identifies rollback outcome.
type RollbackStatus string

const (
	// RollbackStatusEmpty means rollback was not needed or not started.
	RollbackStatusEmpty RollbackStatus = ""
	// RollbackStatusPending means rollback has not completed.
	RollbackStatusPending RollbackStatus = "pending"
	// RollbackStatusSucceeded means all compensating actions completed.
	RollbackStatusSucceeded RollbackStatus = "succeeded"
	// RollbackStatusFailed means at least one compensating action failed.
	RollbackStatusFailed RollbackStatus = "failed"
)

// TransactionIDFunc creates deterministic or production transaction IDs.
type TransactionIDFunc func(TransactionInput) TransactionID

// TransactionInput describes data available to transaction ID generators.
type TransactionInput struct {
	// Version is the publication version shared by planned modules.
	Version string
	// SourceCommit is the source repository HEAD when available.
	SourceCommit git.CommitHash
}

// TransactionJournal is the durable JSON state for a saga-style publish.
//
// It is intentionally honest: Git remotes do not provide distributed ACID
// transactions across repositories, so the journal records enough state for
// best-effort rollback and later manual recovery when rollback safety checks
// cannot be satisfied.
type TransactionJournal struct {
	SchemaVersion int                      `json:"schemaVersion"`
	ID            TransactionID            `json:"id"`
	Status        TransactionStatus        `json:"status"`
	Rollback      RollbackStatus           `json:"rollbackStatus,omitempty"`
	Version       string                   `json:"version,omitempty"`
	Remote        string                   `json:"remote,omitempty"`
	StartedAt     time.Time                `json:"startedAt"`
	UpdatedAt     time.Time                `json:"updatedAt"`
	Modules       []ModuleTransactionState `json:"modules"`
	Warnings      []string                 `json:"warnings,omitempty"`
	Failure       string                   `json:"failure,omitempty"`
	ManualActions []ManualRecoveryAction   `json:"manualRecoveryActions,omitempty"`
}

// Summary returns the compact journal identity used by list commands.
func (j TransactionJournal) Summary() TransactionSummary {
	return TransactionSummary{
		ID:        j.ID,
		Status:    j.Status,
		Rollback:  j.Rollback,
		Version:   j.Version,
		StartedAt: j.StartedAt,
		UpdatedAt: j.UpdatedAt,
	}
}

// TransactionSummary is a compact transaction journal listing item.
type TransactionSummary struct {
	ID        TransactionID     `json:"id"`
	Status    TransactionStatus `json:"status"`
	Rollback  RollbackStatus    `json:"rollbackStatus,omitempty"`
	Version   string            `json:"version,omitempty"`
	StartedAt time.Time         `json:"startedAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

// ModuleTransactionState records all refs needed to publish or roll back one
// module.
type ModuleTransactionState struct {
	Module              manifest.ModuleName    `json:"module"`
	Repository          manifest.RepositoryRef `json:"repository"`
	WorktreeDir         string                 `json:"worktreeDir,omitempty"`
	TargetBranch        manifest.BranchName    `json:"targetBranch,omitempty"`
	FinalBranchRef      string                 `json:"finalBranchRef,omitempty"`
	FinalTagRef         string                 `json:"finalTagRef,omitempty"`
	LocalBaseHead       git.CommitHash         `json:"localBaseHead,omitempty"`
	LocalBaseBranch     git.BranchName         `json:"localBaseBranch,omitempty"`
	RemoteBaseCommit    git.CommitHash         `json:"remoteBaseCommit,omitempty"`
	RemoteBaseExists    bool                   `json:"remoteBaseExists"`
	CreatedCommit       git.CommitHash         `json:"createdCommit,omitempty"`
	CandidateBranchRef  string                 `json:"candidateBranchRef,omitempty"`
	Skipped             bool                   `json:"skipped"`
	CandidatePushed     bool                   `json:"candidatePushed"`
	FinalBranchPromoted bool                   `json:"finalBranchPromoted"`
	LocalTagCreated     bool                   `json:"localTagCreated"`
	RemoteTagPushed     bool                   `json:"remoteTagPushed"`
	RemoteTagHash       git.CommitHash         `json:"remoteTagHash,omitempty"`
	Rollback            ModuleRollbackState    `json:"rollback,omitempty"`
}

// ModuleRollbackState records compensating actions already completed.
type ModuleRollbackState struct {
	RemoteTagDeleted    bool     `json:"remoteTagDeleted"`
	FinalBranchRestored bool     `json:"finalBranchRestored"`
	CandidateDeleted    bool     `json:"candidateDeleted"`
	LocalTagDeleted     bool     `json:"localTagDeleted"`
	LocalWorktreeReset  bool     `json:"localWorktreeReset"`
	FailedActions       []string `json:"failedActions,omitempty"`
}

// ManualRecoveryAction describes a rollback action the operator must review.
type ManualRecoveryAction struct {
	Module       manifest.ModuleName    `json:"module,omitempty"`
	Repository   manifest.RepositoryRef `json:"repository,omitempty"`
	Ref          string                 `json:"ref,omitempty"`
	ExpectedHash git.CommitHash         `json:"expectedHash,omitempty"`
	DesiredHash  git.CommitHash         `json:"desiredHash,omitempty"`
	Message      string                 `json:"message"`
	Command      string                 `json:"command,omitempty"`
}

func newTransactionJournal(id TransactionID, req Request, remote string, modules []modulePreflight, now time.Time) TransactionJournal {
	return TransactionJournal{
		SchemaVersion: transactionSchemaVersion,
		ID:            id,
		Status:        TransactionStatusPending,
		Version:       transactionVersion(req),
		Remote:        remote,
		StartedAt:     now,
		UpdatedAt:     now,
		Modules:       newModuleStates(id, req, modules),
	}
}

func newModuleStates(id TransactionID, req Request, modules []modulePreflight) []ModuleTransactionState {
	out := make([]ModuleTransactionState, 0, len(modules))
	for _, item := range modules {
		branch := item.mod.Branches()[0].Target()
		tagRef := ""
		if req.Plan.PublishPolicy().Tags().Enabled() {
			tagRef = "refs/tags/" + item.mod.Version().String()
		}
		out = append(out, ModuleTransactionState{
			Module:             item.mod.Name(),
			Repository:         item.mod.Repository(),
			WorktreeDir:        item.worktree,
			TargetBranch:       branch,
			FinalBranchRef:     branchRef(branch),
			FinalTagRef:        tagRef,
			CandidateBranchRef: candidateRef(id, item.mod.Name()),
			Skipped:            item.skip,
		})
	}
	return out
}

func transactionVersion(req Request) string {
	for _, mod := range req.Plan.Modules() {
		return mod.Version().String()
	}
	return ""
}

func defaultTransactionID(input TransactionInput) TransactionID {
	version := safeRefComponent(input.Version)
	if version == "" {
		version = "unknown"
	}
	return TransactionID(fmt.Sprintf("tx-%s-%s", time.Now().UTC().Format("20060102T150405.000000000Z"), version))
}

func candidateRef(id TransactionID, module manifest.ModuleName) string {
	return "refs/heads/arcpub/tx/" + safeRefComponent(id.String()) + "/" + safeRefComponent(module.String())
}

func branchRef(branch manifest.BranchName) string {
	return "refs/heads/" + branch.String()
}

func validateTransactionID(id TransactionID) error {
	if id == "" {
		return fmt.Errorf("empty transaction id")
	}
	if safeRefComponent(id.String()) != id.String() {
		return fmt.Errorf("unsafe transaction id")
	}
	return validateGitRefComponent(id.String())
}

func validateGitRef(ref string) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	if strings.HasPrefix(ref, "/") || strings.HasSuffix(ref, "/") {
		return fmt.Errorf("ref has leading or trailing slash")
	}
	if strings.Contains(ref, "//") {
		return fmt.Errorf("ref contains empty component")
	}
	if strings.Contains(ref, "..") || strings.Contains(ref, "@{") {
		return fmt.Errorf("ref contains forbidden sequence")
	}
	if strings.HasSuffix(ref, ".") {
		return fmt.Errorf("ref ends with dot")
	}
	for _, r := range ref {
		if r <= 0x20 || r == 0x7f {
			return fmt.Errorf("ref contains control or space")
		}
		switch r {
		case '\\', '^', '~', ':', '?', '*', '[':
			return fmt.Errorf("ref contains forbidden character %q", r)
		}
	}
	for _, component := range strings.Split(ref, "/") {
		if err := validateGitRefComponent(component); err != nil {
			return err
		}
	}
	return nil
}

func validateGitRefComponent(component string) error {
	if component == "" {
		return fmt.Errorf("empty ref component")
	}
	if strings.HasPrefix(component, ".") {
		return fmt.Errorf("ref component %q starts with dot", component)
	}
	if strings.HasSuffix(strings.ToLower(component), ".lock") {
		return fmt.Errorf("ref component %q ends with .lock", component)
	}
	return nil
}

func safeRefComponent(value string) string {
	original := strings.TrimSpace(value)
	value = strings.TrimSpace(value)
	if value == "" {
		return "value-" + shortHash(original)
	}
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_', r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	out := b.String()
	for strings.Contains(out, "..") {
		out = strings.ReplaceAll(out, "..", ".")
	}
	out = strings.Trim(out, "-.")
	if strings.HasSuffix(strings.ToLower(out), ".lock") {
		out = strings.TrimSuffix(out[:len(out)-5], ".-") + "-lock"
	}
	if out == "" {
		return "value-" + shortHash(original)
	}
	return out
}

func shortHash(value string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(value))
	return fmt.Sprintf("%08x", h.Sum32())
}
