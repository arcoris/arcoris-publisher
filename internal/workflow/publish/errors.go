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

import "fmt"

// ErrorCode identifies a publication-stage error.
type ErrorCode string

const (
	// CodeInvalidRequest indicates malformed publication input or missing ports.
	CodeInvalidRequest ErrorCode = "invalid_request"

	// CodeVerificationFailed indicates publication was blocked by failed checks.
	CodeVerificationFailed ErrorCode = "verification_failed"

	// CodeMissingSourceSnapshot indicates source provenance for a planned module
	// is unavailable.
	CodeMissingSourceSnapshot ErrorCode = "missing_source_snapshot"

	// CodePreflightFailed indicates publication was blocked before Git mutation.
	CodePreflightFailed ErrorCode = "preflight_failed"

	// CodePublishFailed indicates a Git publication operation failed.
	CodePublishFailed ErrorCode = "publish_failed"

	// CodePendingTransaction indicates an unresolved transaction already exists.
	CodePendingTransaction ErrorCode = "pending_transaction"

	// CodeJournalFailed indicates the durable transaction journal failed.
	CodeJournalFailed ErrorCode = "journal_failed"

	// CodeLockFailed indicates the transaction lock could not be acquired.
	CodeLockFailed ErrorCode = "lock_failed"

	// CodeTransactionFailed indicates the saga failed after transaction start.
	CodeTransactionFailed ErrorCode = "transaction_failed"

	// CodeRollbackFailed indicates automatic or explicit rollback failed.
	CodeRollbackFailed ErrorCode = "rollback_failed"

	// CodeCandidatePushFailed indicates a candidate ref push failed.
	CodeCandidatePushFailed ErrorCode = "candidate_push_failed"

	// CodePromotionFailed indicates final branch promotion failed.
	CodePromotionFailed ErrorCode = "promotion_failed"

	// CodeTagPublishFailed indicates final tag publication failed.
	CodeTagPublishFailed ErrorCode = "tag_publish_failed"

	// CodeRecoveryFailed indicates transaction recovery failed.
	CodeRecoveryFailed ErrorCode = "recovery_failed"

	// CodePruneFailed indicates transaction journal pruning failed.
	CodePruneFailed ErrorCode = "prune_failed"
)

// Error describes why publication could not complete.
type Error struct {
	// Code is the stable machine-readable error reason.
	Code ErrorCode

	// Message is the human-readable diagnostic text.
	Message string

	// Cause wraps the infrastructure error when available.
	Cause error
}

// Error returns a compact publication error string.
func (e *Error) Error() string {
	if e == nil {
		return "publish error"
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying infrastructure error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
