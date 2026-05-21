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
	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

// Module returns the underlying manifest module declaration.
func (p ModulePlan) Module() manifest.Module { return p.module }

// Name returns the planned module name.
func (p ModulePlan) Name() manifest.ModuleName { return p.module.Name() }

// ModulePath returns the planned module path.
func (p ModulePlan) ModulePath() manifest.ModulePath { return p.module.ModulePath() }

// SourceDir returns the staged source directory for the planned module.
func (p ModulePlan) SourceDir() manifest.SourceDir { return p.module.SourceDir() }

// Repository returns the target repository for the planned module.
func (p ModulePlan) Repository() manifest.RepositoryRef { return p.module.Repository() }

// Version returns the version assigned to the planned module.
func (p ModulePlan) Version() versioning.Version { return p.version }

// Action returns the planned action for the module.
func (p ModulePlan) Action() Action { return p.action }

// OrderIndex returns the dependency-order index of the module in the plan.
func (p ModulePlan) OrderIndex() int { return p.orderIndex }
