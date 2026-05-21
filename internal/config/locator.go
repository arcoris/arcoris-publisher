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

package config

import (
	"context"
	"path/filepath"
)

// Locator discovers a top-level arcpub manifest in a directory.
type Locator struct {
	Names []string
}

// DefaultLocator returns a locator for conventional top-level manifest names.
func DefaultLocator() Locator {
	return Locator{
		Names: StagingManifestNames(),
	}
}

// Find searches dir for the first supported top-level manifest name.
func (l Locator) Find(ctx context.Context, reader Reader, dir string) (string, error) {
	if reader == nil {
		reader = OSReader{}
	}
	names := l.Names
	if len(names) == 0 {
		names = StagingManifestNames()
	}
	for _, name := range names {
		path := filepath.Join(dir, filepath.FromSlash(name))
		exists, err := reader.Exists(ctx, path)
		if err != nil {
			return "", inspectCandidateError(path, err)
		}
		if exists {
			return path, nil
		}
	}

	return "", configNotFoundError(dir)
}
