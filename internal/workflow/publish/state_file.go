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
	"io"
	"os"
	"path/filepath"
)

const (
	maxTransactionLockBytes = 16 << 10
	maxOperationLockBytes   = 16 << 10
	maxTransactionJournalBytes = 8 << 20
)

var errStateFileTooLarge = errors.New("transaction state file too large")

// readBoundedStateFile prevents malformed recovery state from turning a
// read-only inspection path into an unbounded allocation. Errors intentionally
// mention only the base name so default diagnostics do not leak local paths.
func readBoundedStateFile(path string, maxBytes int64) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("%w: %s exceeds %d bytes", errStateFileTooLarge, filepath.Base(path), maxBytes)
	}
	return data, nil
}
