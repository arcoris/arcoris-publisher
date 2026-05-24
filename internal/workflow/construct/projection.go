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

package construct

import (
	"fmt"
	"path/filepath"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
)

// projectionEntry reserves one target path in the explicit publication
// projection before the target worktree is mutated.
type projectionEntry struct {
	// index is the explicit publish entry index, or -1 for generated files.
	index int

	// kind distinguishes file and directory subtree reservations.
	kind manifest.PublishEntryKind

	// path is a clean slash-separated target-relative path.
	path string

	// generated reports whether the entry is publisher-generated metadata.
	generated bool
}

// projectionCollision describes two target reservations that cannot safely
// coexist in an explicit projection.
type projectionCollision struct {
	current  projectionEntry
	previous projectionEntry
	message  string
}

// validateProjection rejects target collisions before clean or copy operations
// can mutate the target worktree.
func validateProjection(
	req Request,
	module moduleContext,
	opts Options,
	issues *issueCollector,
) bool {
	before := issues.Len()
	entries := projectionEntries(module.plan)
	entries = appendProvenanceProjectionEntry(entries, req, opts)

	for _, entry := range entries {
		targetPath := pathutil.JoinRelative(
			module.workspace.WorktreeDir(),
			projectionPath(entry.path),
		)
		if err := pathutil.EnsureInside(module.workspace.WorktreeDir(), targetPath); err != nil {
			issues.AddMessage(IssueTargetPathEscape, module.plan.Name(), targetPath, err.Error())
		}
	}

	for _, collision := range detectProjectionCollisions(entries) {
		issues.AddMessage(
			IssueProjectionClash,
			module.plan.Name(),
			collision.current.issuePath(),
			collision.message,
		)
	}

	return issues.Len() == before
}

type projectionPath string

func (p projectionPath) String() string { return string(p) }

// projectionEntries reserves all explicit publish targets, including optional
// entries that were absent in the source snapshot. Optional entries still own
// their target paths so a later source file cannot silently collide with another
// projection rule.
func projectionEntries(mod plan.ModulePlan) []projectionEntry {
	entries := mod.PublishEntries()
	out := make([]projectionEntry, 0, len(entries))
	for i, entry := range entries {
		out = append(out, projectionEntry{
			index: i,
			kind:  entry.Kind(),
			path:  cleanProjectionPath(entry.To().String()),
		})
	}

	return out
}

// appendProvenanceProjectionEntry reserves the generated provenance file target
// when file provenance is enabled for construction.
func appendProvenanceProjectionEntry(
	entries []projectionEntry,
	req Request,
	opts Options,
) []projectionEntry {
	if !opts.GenerateProvenanceFile || !req.Plan.PublishPolicy().Provenance().FileEnabled() {
		return entries
	}

	return append(entries, projectionEntry{
		index:     -1,
		kind:      manifest.PublishEntryFile,
		path:      cleanProjectionPath(req.Plan.PublishPolicy().Provenance().File().String()),
		generated: true,
	})
}

// detectProjectionCollisions returns deterministic collisions by comparing each
// later reservation with the reservations declared before it.
func detectProjectionCollisions(entries []projectionEntry) []projectionCollision {
	collisions := []projectionCollision{}
	for i, current := range entries {
		for j := 0; j < i; j++ {
			previous := entries[j]
			collision, ok := newProjectionCollision(current, previous)
			if ok {
				collisions = append(collisions, collision)
			}
		}
	}

	return collisions
}

func newProjectionCollision(
	current projectionEntry,
	previous projectionEntry,
) (projectionCollision, bool) {
	if current.path == previous.path {
		return projectionCollision{
			current:  current,
			previous: previous,
			message: fmt.Sprintf(
				"target %q is already reserved by %s",
				current.path,
				previous.label(),
			),
		}, true
	}

	if previous.kind == manifest.PublishEntryDirectory && pathContains(previous.path, current.path) {
		return projectionCollision{
			current:  current,
			previous: previous,
			message: fmt.Sprintf(
				"target %q is inside directory target %q reserved by %s",
				current.path,
				previous.path,
				previous.label(),
			),
		}, true
	}

	if current.kind == manifest.PublishEntryDirectory && pathContains(current.path, previous.path) {
		return projectionCollision{
			current:  current,
			previous: previous,
			message: fmt.Sprintf(
				"directory target %q contains target %q reserved by %s",
				current.path,
				previous.path,
				previous.label(),
			),
		}, true
	}

	return projectionCollision{}, false
}

func (e projectionEntry) issuePath() string {
	if e.generated {
		return "publish.provenance.file"
	}

	return fmt.Sprintf("publish.entries[%d].to", e.index)
}

func (e projectionEntry) label() string {
	if e.generated {
		return "generated provenance file"
	}

	return fmt.Sprintf("publish.entries[%d]", e.index)
}

func cleanProjectionPath(path string) string {
	clean := filepath.ToSlash(filepath.Clean(path))
	if clean == "" {
		return "."
	}

	return clean
}

func pathContains(parent, child string) bool {
	if parent == "." {
		return child != "."
	}

	return strings.HasPrefix(child, parent+"/")
}
