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

// UnmarshalJSON rejects journal representations that this binary cannot
// interpret safely. Durable recovery state is versioned state, not a
// best-effort configuration file: an unknown schema, status, rollback status,
// or JSON field must fail closed instead of being silently projected onto the
// current structs.
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
	if !knownTransactionStatus(decoded.Status) {
		return fmt.Errorf("unsupported transaction status %q", decoded.Status)
	}
	if !knownRollbackStatus(decoded.Rollback) {
		return fmt.Errorf("unsupported rollback status %q", decoded.Rollback)
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

// normalizeTransactionJournalForWrite produces the canonical v1 form used by
// the file store. A zero schema version is accepted only for in-process
// callers and normalized before persistence; non-zero unknown versions are
// rejected so old binaries cannot manufacture future state accidentally.
func normalizeTransactionJournalForWrite(journal TransactionJournal) (TransactionJournal, error) {
	if journal.SchemaVersion == 0 {
		journal.SchemaVersion = transactionSchemaVersion
	}
	if journal.SchemaVersion != transactionSchemaVersion {
		return TransactionJournal{}, fmt.Errorf("unsupported transaction journal schemaVersion %d", journal.SchemaVersion)
	}
	if !knownTransactionStatus(journal.Status) {
		return TransactionJournal{}, fmt.Errorf("unsupported transaction status %q", journal.Status)
	}
	if !knownRollbackStatus(journal.Rollback) {
		return TransactionJournal{}, fmt.Errorf("unsupported rollback status %q", journal.Rollback)
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

func knownRollbackStatus(status RollbackStatus) bool {
	switch status {
	case RollbackStatusEmpty,
		RollbackStatusPending,
		RollbackStatusSucceeded,
		RollbackStatusFailed:
		return true
	default:
		return false
	}
}
