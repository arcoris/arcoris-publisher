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

package versioning

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// Module returns the dependency module name.
func (v DependencyVersion) Module() manifest.ModuleName { return v.module }

// ModulePath returns the dependency public module path.
func (v DependencyVersion) ModulePath() manifest.ModulePath { return v.modulePath }

// Version returns the dependency version requirement.
func (v DependencyVersion) Version() Version { return v.version }
