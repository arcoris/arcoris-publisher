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

package manifest

import "fmt"

// Dependency declares a direct dependency on another manifest module.
type Dependency struct {
	module ModuleName
}

// NewDependency validates a dependency module name and returns a Dependency.
func NewDependency(module string) (Dependency, error) {
	name, err := ParseModuleName(module)
	if err != nil {
		return Dependency{}, fmt.Errorf("module: %w", err)
	}
	return Dependency{module: name}, nil
}

// Module returns the referenced module name.
func (d Dependency) Module() ModuleName { return d.module }
