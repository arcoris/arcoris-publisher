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
	"errors"
	"fmt"
)

var errInvalidTransactionStatusTransition = errors.New("invalid transaction status transition")

// validateTransactionStatusTransition defines the durable publish state
// machine. Recovery is deliberately allowed to enter rolling_back from any
// known publish-blocking state, while committed and rolled_back are absorbing
// terminal states. Self-transitions are permitted so interrupted recovery can
// safely persist its current phase again.
func validateTransactionStatusTransition(from, to TransactionStatus) error {
	if !knownTransactionStatus(from) {
		return fmt.Errorf("%w: unknown current status %q", errInvalidTransactionStatusTransition, from)
	}
	if !knownTransactionStatus(to) {
		return fmt.Errorf("%w: unknown target status %q", errInvalidTransactionStatusTransition, to)
	}
	if from == to {
		return nil
	}
	if to == TransactionStatusRollingBack {
		if from.BlocksNewPublish() && from != TransactionStatusCommitted && from != TransactionStatusRolledBack {
			return nil
		}
	}
	if to == TransactionStatusFailed && publishPhaseCanFail(from) {
		return nil
	}

	switch from {
	case TransactionStatusPending:
		if to == TransactionStatusPreflighted {
			return nil
		}
	case TransactionStatusPreflighted:
		if to == TransactionStatusSnapshotted {
			return nil
		}
	case TransactionStatusSnapshotted:
		if to == TransactionStatusCommittedLocally {
			return nil
		}
	case TransactionStatusCommittedLocally:
		if to == TransactionStatusCandidatesPushed {
			return nil
		}
	case TransactionStatusCandidatesPushed:
		if to == TransactionStatusPromoting {
			return nil
		}
	case TransactionStatusPromoting:
		if to == TransactionStatusBranchesPromoted {
			return nil
		}
	case TransactionStatusBranchesPromoted:
		if to == TransactionStatusTagging || to == TransactionStatusCommitted {
			return nil
		}
	case TransactionStatusTagging:
		if to == TransactionStatusCommitted {
			return nil
		}
	case TransactionStatusRollingBack:
		if to == TransactionStatusRolledBack || to == TransactionStatusRollbackFailed {
			return nil
		}
	}
	return fmt.Errorf("%w: %s -> %s", errInvalidTransactionStatusTransition, from, to)
}

func publishPhaseCanFail(status TransactionStatus) bool {
	switch status {
	case TransactionStatusPending,
		TransactionStatusPreflighted,
		TransactionStatusSnapshotted,
		TransactionStatusCommittedLocally,
		TransactionStatusCandidatesPushed,
		TransactionStatusPromoting,
		TransactionStatusBranchesPromoted,
		TransactionStatusTagging:
		return true
	default:
		return false
	}
}
