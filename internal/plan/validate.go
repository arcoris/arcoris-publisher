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

package plan

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// rebuildIndexes reconstructs plan lookup maps and rejects duplicate public
// routing keys.
func (p *Plan) rebuildIndexes() error {
	p.byName = make(map[manifest.ModuleName]int, len(p.modules))
	p.byPath = make(map[manifest.ModulePath]int, len(p.modules))
	p.byRepository = make(map[manifest.RepositoryRef]int, len(p.modules))

	collector := newIssueCollector()
	for i, mod := range p.modules {
		p.indexModule(&collector, i, mod)
	}

	return collector.Err()
}

// indexModule stores one module in every lookup index.
func (p *Plan) indexModule(
	collector *issueCollector,
	index int,
	mod ModulePlan,
) {
	path := fmt.Sprintf("modules[%d]", index)

	p.indexModuleName(collector, path, index, mod)
	p.indexModulePath(collector, path, index, mod)
	p.indexRepository(collector, path, index, mod)
}

// indexModuleName stores the module-name index or reports a duplicate.
func (p *Plan) indexModuleName(
	collector *issueCollector,
	path string,
	index int,
	mod ModulePlan,
) {
	if prev, exists := p.byName[mod.Name()]; exists {
		collector.Add(
			IssueDuplicateModuleName,
			path+".name",
			"duplicate module name %q previously planned at modules[%d]",
			mod.Name(),
			prev,
		)
		return
	}

	p.byName[mod.Name()] = index
}

// indexModulePath stores the Go module path index or reports a duplicate.
func (p *Plan) indexModulePath(
	collector *issueCollector,
	path string,
	index int,
	mod ModulePlan,
) {
	if prev, exists := p.byPath[mod.ModulePath()]; exists {
		collector.Add(
			IssueDuplicateModulePath,
			path+".modulePath",
			"duplicate module path %q previously planned at modules[%d]",
			mod.ModulePath(),
			prev,
		)
		return
	}

	p.byPath[mod.ModulePath()] = index
}

// indexRepository stores the target repository index or reports a duplicate.
func (p *Plan) indexRepository(
	collector *issueCollector,
	path string,
	index int,
	mod ModulePlan,
) {
	if prev, exists := p.byRepository[mod.Repository()]; exists {
		collector.Add(
			IssueDuplicateRepository,
			path+".repository",
			"duplicate repository %q previously planned at modules[%d]",
			mod.Repository(),
			prev,
		)
		return
	}

	p.byRepository[mod.Repository()] = index
}

// validate checks the final executable plan shape.
func validate(p Plan) error {
	collector := newIssueCollector()
	if len(p.modules) == 0 {
		collector.Add(IssueEmptyPlan, "modules", "publication plan has no public modules")
	}
	for i, mod := range p.modules {
		path := fmt.Sprintf("modules[%d]", i)
		validateModulePlan(&collector, path, mod)
	}

	return collector.Err()
}

// validateModulePlan checks invariants expected by downstream workflow stages.
func validateModulePlan(collector *issueCollector, path string, mod ModulePlan) {
	if mod.Visibility() != manifest.VisibilityPublic {
		collector.Add(
			IssueNonPublishableModule,
			path+".visibility",
			"planned module %q is not public",
			mod.Name(),
		)
	}

	if mod.Version().IsZero() {
		collector.Add(
			IssueMissingAssignment,
			path+".version",
			"planned module %q has no version",
			mod.Name(),
		)
	}

	if len(mod.Branches()) == 0 {
		collector.Add(
			IssueEmptyBranches,
			path+".branches",
			"planned module %q has no branch mappings",
			mod.Name(),
		)
	}

	if len(mod.PublishEntries()) == 0 {
		collector.Add(
			IssueEmptyPublishEntries,
			path+".publish.entries",
			"planned module %q has no explicit publish entries",
			mod.Name(),
		)
	}
}
