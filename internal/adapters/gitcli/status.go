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

package gitcli

import (
	"context"
	"strings"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

// Status reads porcelain-v1 status and converts it into port status entries.
func (c *Client) Status(ctx context.Context, repoDir string) (gitport.Status, error) {
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"status", "--porcelain=v1", "-z"}, nil, true, true))
	if err != nil {
		return gitport.Status{}, wrapGitCommandError("git status failed", result, err)
	}
	entries := parseStatus(result.Stdout)
	return gitport.Status{Clean: len(entries) == 0, Entries: entries}, nil
}

// parseStatus parses NUL-delimited porcelain-v1 output.
//
// Rename and copy records include an extra path field after the primary entry;
// the parser skips that extra field because the port currently models only one
// path per status entry.
func parseStatus(out []byte) []gitport.StatusEntry {
	if len(out) == 0 {
		return nil
	}
	parts := strings.Split(string(out), "\x00")
	entries := make([]gitport.StatusEntry, 0, len(parts))
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		if part == "" {
			continue
		}
		if len(part) < 3 {
			entries = append(entries, gitport.StatusEntry{Path: part})
			continue
		}
		code := part[:2]
		path := strings.TrimPrefix(part[2:], " ")
		entries = append(entries, gitport.StatusEntry{Code: code, Path: path})
		if strings.HasPrefix(code, "R") || strings.HasPrefix(code, "C") {
			i++
		}
	}
	return entries
}
