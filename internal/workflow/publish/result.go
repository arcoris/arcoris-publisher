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
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

// Result describes publication outcomes in plan order.
type Result struct {
	modules     []ModuleResult
	transaction TransactionJournal
}

// Modules returns detached module publication results.
func (r Result) Modules() []ModuleResult {
	out := make([]ModuleResult, len(r.modules))
	copy(out, r.modules)
	return out
}

// Transaction returns the durable publish transaction journal when publication
// created one.
func (r Result) Transaction() TransactionJournal { return r.transaction }

// HasTransaction reports whether a transaction journal is present.
func (r Result) HasTransaction() bool { return r.transaction.ID != "" }

// Published reports whether at least one module was pushed.
func (r Result) Published() bool {
	for _, m := range r.modules {
		if m.Published() {
			return true
		}
	}
	return false
}

// ModuleResult describes publication outcome for one module.
type ModuleResult struct {
	// module is the planned module name.
	module manifest.ModuleName

	// commit is the created publication commit.
	commit git.CommitHash

	// tags contains pushed release tags.
	tags []git.TagName

	// pushed reports whether refs were pushed.
	pushed bool

	// skipped reports whether publication skipped this module due to no changes.
	skipped bool
}

// Module returns the planned module name.
func (m ModuleResult) Module() manifest.ModuleName { return m.module }

// Commit returns the publication commit.
func (m ModuleResult) Commit() git.CommitHash { return m.commit }

// Tags returns detached pushed release tags.
func (m ModuleResult) Tags() []git.TagName {
	out := make([]git.TagName, len(m.tags))
	copy(out, m.tags)
	return out
}

// Pushed reports whether refs were pushed.
func (m ModuleResult) Pushed() bool { return m.pushed }

// Skipped reports whether publication skipped this module due to no changes.
func (m ModuleResult) Skipped() bool { return m.skipped }

// Published reports whether this module was actually published.
func (m ModuleResult) Published() bool { return m.pushed && !m.skipped }
