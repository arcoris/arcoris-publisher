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
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
)

// readTransactionJournalFile is the single decoder for durable journal files.
// Structural state-file violations are classified as journal corruption so
// Load, List, preflight, prune, lock inspection, and diagnostics all fail
// closed with the same semantic category. Ordinary filesystem read failures
// remain distinguishable because diagnostics must report incomplete inspection
// rather than mislabeling an inaccessible file as corrupt.
func readTransactionJournalFile(path string, requested TransactionID) (TransactionJournal, error) {
	data, err := readBoundedStateFile(path, maxTransactionJournalBytes)
	if err != nil {
		if errors.Is(err, errStateFileTooLarge) ||
			errors.Is(err, errStateFileNotRegular) ||
			errors.Is(err, errStateFileChanged) {
			return TransactionJournal{}, journalCorruptf(
				"transaction journal %s has an unsafe file representation: %v",
				filepath.Base(path),
				err,
			)
		}
		return TransactionJournal{}, err
	}

	var journal TransactionJournal
	if err := json.Unmarshal(data, &journal); err != nil {
		return TransactionJournal{}, journalCorruptf(
			"transaction journal %s is corrupt: %w",
			filepath.Base(path),
			err,
		)
	}
	if err := validateJournalIdentity(filepath.Base(path), requested, journal.ID); err != nil {
		return TransactionJournal{}, err
	}
	return journal, nil
}

func journalReadErrorMessage(name string, err error) string {
	return fmt.Sprintf("read transaction journal %s failed: %v", name, err)
}
