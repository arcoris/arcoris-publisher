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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Plan is an immutable-by-convention publication execution snapshot.
//
// Modules are stored in deterministic dependency-before-dependent publication
// order. Internal and disabled modules are not included in the executable module
// list; they may still have participated in earlier validation and topology.
type Plan struct {
	// metadata is the publication-set metadata copied from the resolved model.
	metadata manifest.Metadata

	// source is the effective source repository declaration.
	source manifest.Source

	// publish is the effective global publication policy.
	publish manifest.PublishPolicy

	// modules contains executable public modules in publication order.
	modules []ModulePlan

	// byName indexes modules by resolved manifest module name.
	byName map[manifest.ModuleName]int

	// byPath indexes modules by published Go module path.
	byPath map[manifest.ModulePath]int

	// byRepository indexes modules by target repository.
	byRepository map[manifest.RepositoryRef]int
}

// Metadata returns the publication-set metadata used to build the plan.
func (p Plan) Metadata() manifest.Metadata { return p.metadata }

// Source returns the effective source repository declaration.
func (p Plan) Source() manifest.Source { return p.source }

// PublishPolicy returns the effective global publication policy.
func (p Plan) PublishPolicy() manifest.PublishPolicy { return p.publish }

// Len returns the number of public modules in the plan.
func (p Plan) Len() int { return len(p.modules) }

// Empty reports whether the plan has no public modules.
func (p Plan) Empty() bool { return len(p.modules) == 0 }

// Modules returns detached module plans in publication order.
func (p Plan) Modules() []ModulePlan {
	out := make([]ModulePlan, len(p.modules))
	copy(out, p.modules)
	return out
}

// ModuleNames returns planned module names in publication order.
func (p Plan) ModuleNames() []manifest.ModuleName {
	out := make([]manifest.ModuleName, 0, len(p.modules))
	for _, mod := range p.modules {
		out = append(out, mod.Name())
	}
	return out
}
