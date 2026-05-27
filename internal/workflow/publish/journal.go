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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// JournalStore persists publish transactions for rollback and operator
// recovery.
type JournalStore interface {
	Create(ctx context.Context, journal TransactionJournal) error
	Update(ctx context.Context, journal TransactionJournal) error
	Load(ctx context.Context, id TransactionID) (TransactionJournal, error)
	List(ctx context.Context) ([]TransactionSummary, error)
	HasPending(ctx context.Context) (TransactionSummary, bool, error)
}

// FileJournalStore stores transaction journals as newline-terminated JSON files.
type FileJournalStore struct {
	stateDir string
}

var errTransactionJournalCorrupt = errors.New("transaction journal corrupt")

func journalCorruptf(format string, args ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{errTransactionJournalCorrupt}, args...)...)
}

// NewFileJournalStore returns a journal store rooted at stateDir.
func NewFileJournalStore(stateDir string) FileJournalStore {
	return FileJournalStore{stateDir: stateDir}
}

// StateDir returns the journal state directory.
func (s FileJournalStore) StateDir() string { return s.stateDir }

// Create writes a new transaction journal atomically.
func (s FileJournalStore) Create(ctx context.Context, journal TransactionJournal) error {
	return s.write(ctx, journal)
}

// Update replaces an existing transaction journal atomically.
func (s FileJournalStore) Update(ctx context.Context, journal TransactionJournal) error {
	return s.write(ctx, journal)
}

// Load reads one transaction journal by id.
func (s FileJournalStore) Load(ctx context.Context, id TransactionID) (TransactionJournal, error) {
	if err := ctx.Err(); err != nil {
		return TransactionJournal{}, err
	}
	path, err := s.journalPath(id)
	if err != nil {
		return TransactionJournal{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return TransactionJournal{}, err
	}
	var journal TransactionJournal
	if err := json.Unmarshal(data, &journal); err != nil {
		return TransactionJournal{}, journalCorruptf("transaction journal %s is corrupt: %w", filepath.Base(path), err)
	}
	if err := validateJournalIdentity(filepath.Base(path), id, journal.ID); err != nil {
		return TransactionJournal{}, err
	}
	return journal, nil
}

// List returns transaction summaries sorted by start time.
func (s FileJournalStore) List(ctx context.Context) ([]TransactionSummary, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	dir := s.transactionsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	out := make([]TransactionSummary, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		var journal TransactionJournal
		if err := json.Unmarshal(data, &journal); err != nil {
			return nil, journalCorruptf("transaction journal %s is corrupt: %w", entry.Name(), err)
		}
		if err := validateJournalIdentity(entry.Name(), "", journal.ID); err != nil {
			return nil, err
		}
		out = append(out, journal.Summary())
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].StartedAt.Equal(out[j].StartedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].StartedAt.Before(out[j].StartedAt)
	})
	return out, nil
}

// HasPending reports the oldest non-terminal transaction, if any.
func (s FileJournalStore) HasPending(ctx context.Context) (TransactionSummary, bool, error) {
	summaries, err := s.List(ctx)
	if err != nil {
		return TransactionSummary{}, false, err
	}
	for _, summary := range summaries {
		if summary.Status.BlocksNewPublish() {
			return summary, true, nil
		}
	}
	return TransactionSummary{}, false, nil
}

func (s FileJournalStore) write(ctx context.Context, journal TransactionJournal) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path, err := s.journalPath(journal.ID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(s.stateDir, 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(journal, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return atomicWriteFile(path, data, 0o600)
}

func (s FileJournalStore) journalPath(id TransactionID) (string, error) {
	if id == "" {
		return "", fmt.Errorf("transaction id is required")
	}
	if s.stateDir == "" {
		return "", fmt.Errorf("state dir is required")
	}
	name := id.String()
	if err := validateTransactionID(id); err != nil {
		return "", fmt.Errorf("transaction id %q is not a safe journal file name", id)
	}
	return filepath.Join(s.transactionsDir(), name+".json"), nil
}

func (s FileJournalStore) transactionsDir() string {
	return filepath.Join(s.stateDir, "transactions")
}

func validateJournalIdentity(filename string, requested TransactionID, actual TransactionID) error {
	if err := validateTransactionID(actual); err != nil {
		return journalCorruptf("transaction journal %s has unsafe transaction id %q: %w", filename, actual, err)
	}
	if expected := strings.TrimSuffix(filename, ".json"); expected != actual.String() {
		return journalCorruptf("transaction journal %s contains transaction id %q", filename, actual)
	}
	if requested != "" && actual != requested {
		return journalCorruptf("transaction journal %s contains transaction id %q, want %q", filename, actual, requested)
	}
	return nil
}

func deriveStateDir(explicit string, targets []modulePreflight) string {
	if explicit != "" {
		return explicit
	}
	for _, target := range targets {
		if target.worktree != "" {
			return filepath.Join(filepath.Dir(target.worktree), ".arcpub", "state")
		}
	}
	return ""
}

func atomicWriteFile(path string, data []byte, perm fs.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+"-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpName)
		}
	}()

	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	cleanup = false
	return syncParentDir(path)
}

func syncParentDir(path string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}
