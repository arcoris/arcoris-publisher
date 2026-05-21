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
	"arcoris.dev/arcoris-publisher/internal/domain/graph"
	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

// Manifest returns the validated manifest used to build the plan.
func (p Plan) Manifest() manifest.Manifest { return p.manifest }

// Registry returns the module registry used to build the plan.
func (p Plan) Registry() registry.Registry { return p.registry }

// Graph returns the dependency graph used to build the plan.
func (p Plan) Graph() graph.Graph { return p.graph }

// Versions returns version assignments used to build the plan.
func (p Plan) Versions() versioning.Assignments { return p.assignments }

// Source returns the authoritative source repository declaration.
func (p Plan) Source() manifest.Source { return p.manifest.Source() }

// Policy returns the global publication policy.
func (p Plan) Policy() manifest.Policy { return p.manifest.Policy() }

// Len returns the number of modules planned for publication.
func (p Plan) Len() int { return len(p.modules) }

// Empty reports whether the plan contains no publishable modules.
func (p Plan) Empty() bool { return len(p.modules) == 0 }

// Modules returns module plans in dependency-first publish order.
func (p Plan) Modules() []ModulePlan { return cloneModulePlans(p.modules) }

// SkippedModules returns known non-published modules in manifest declaration order.
func (p Plan) SkippedModules() []SkippedModule { return cloneSkippedModules(p.skipped) }
