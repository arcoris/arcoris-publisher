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

	if err := validateUniqueJSONKeys(data); err != nil {
		return err
	}

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

// validateUniqueJSONKeys rejects duplicate object member names at every depth.
// encoding/json otherwise accepts duplicates and keeps the last value, which is
// unsafe for durable recovery state because two readers could assign different
// meaning to the same bytes. Decoder.Token returns decoded member names, so
// escaped spellings such as "id" and "\u0069d" are treated as duplicates.
func validateUniqueJSONKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := validateUniqueJSONValue(decoder); err != nil {
		return err
	}
	return requireJSONEOF(decoder)
}

func validateUniqueJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}

	switch delim {
	case '{':
		seen := make(map[string]struct{})
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("transaction journal JSON object contains a non-string key")
			}
			if _, duplicate := seen[key]; duplicate {
				return fmt.Errorf("transaction journal contains duplicate JSON key %q", key)
			}
			seen[key] = struct{}{}
			if err := validateUniqueJSONValue(decoder); err != nil {
				return err
			}
		}
		return consumeJSONDelimiter(decoder, '}')
	case '[':
		for decoder.More() {
			if err := validateUniqueJSONValue(decoder); err != nil {
				return err
			}
		}
		return consumeJSONDelimiter(decoder, ']')
	default:
		return fmt.Errorf("transaction journal contains unexpected JSON delimiter %q", delim)
	}
}

func consumeJSONDelimiter(decoder *json.Decoder, want json.Delim) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	got, ok := token.(json.Delim)
	if !ok || got != want {
		return fmt.Errorf("transaction journal JSON delimiter mismatch")
	}
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
