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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// UnmarshalJSON rejects wire representations that this binary cannot decode
// safely while keeping enum-like status values opaque. schemaVersion defines
// the structural compatibility boundary; unknown fields or versions fail
// closed, while unknown status values remain inspectable and are conservatively
// treated as publish blockers by TransactionStatus policy.
func (j *TransactionJournal) UnmarshalJSON(data []byte) error {
	type journalWire TransactionJournal

	var envelope struct {
		SchemaVersion *int `json:"schemaVersion"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}
	if envelope.SchemaVersion == nil {
		return fmt.Errorf("transaction journal schemaVersion is missing")
	}
	if *envelope.SchemaVersion != transactionSchemaVersion {
		return fmt.Errorf("unsupported transaction journal schemaVersion %d", *envelope.SchemaVersion)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var decoded journalWire
	if err := decoder.Decode(&decoded); err != nil {
		return err
	}
	if err := requireJSONEOF(decoder); err != nil {
		return err
	}

	*j = TransactionJournal(decoded)
	return nil
}

func requireJSONEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("transaction journal contains multiple JSON values")
		}
		return err
	}
	return nil
}

// normalizeTransactionJournalForWrite produces the canonical supported wire
// version. The storage layer intentionally does not validate status transitions:
// transactionRunner owns that state-machine policy, while storage preserves
// opaque enum values so unknown-but-structurally-compatible state remains
// diagnosable and fail-closed.
func normalizeTransactionJournalForWrite(journal TransactionJournal) (TransactionJournal, error) {
	if journal.SchemaVersion == 0 {
		journal.SchemaVersion = transactionSchemaVersion
	}
	if journal.SchemaVersion != transactionSchemaVersion {
		return TransactionJournal{}, fmt.Errorf("unsupported transaction journal schemaVersion %d", journal.SchemaVersion)
	}
	return journal, nil
}

func knownTransactionStatus(status TransactionStatus) bool {
	switch status {
	case TransactionStatusPending,
		TransactionStatusPreflighted,
		TransactionStatusSnapshotted,
		TransactionStatusCommittedLocally,
		TransactionStatusCandidatesPushed,
		TransactionStatusPromoting,
		TransactionStatusBranchesPromoted,
		TransactionStatusTagging,
		TransactionStatusCommitted,
		TransactionStatusFailed,
		TransactionStatusRollingBack,
		TransactionStatusRolledBack,
		TransactionStatusRollbackFailed:
		return true
	default:
		return false
	}
}
