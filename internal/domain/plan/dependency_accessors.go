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

// Module returns the dependency module name.
func (p DependencyPlan) Module() manifest.ModuleName { return p.module }

// ModulePath returns the dependency module path written into go.mod.
func (p DependencyPlan) ModulePath() manifest.ModulePath { return p.modulePath }

// Version returns the dependency version requirement.
func (p DependencyPlan) Version() versioning.Version { return p.version }
