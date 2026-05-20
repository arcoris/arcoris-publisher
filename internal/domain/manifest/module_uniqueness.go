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

// validateLocalUniqueness checks duplicate declarations inside one module.
func (m Module) validateLocalUniqueness() error {
	if err := validateUniqueBranches(m.branches); err != nil {
		return err
	}
	return validateUniqueDependencies(m.dependencies)
}

// validateUniqueBranches rejects duplicate source branch mappings.
func validateUniqueBranches(branches []BranchMapping) error {
	seen := map[BranchName]struct{}{}
	for _, branch := range branches {
		key := branch.Source()
		if _, exists := seen[key]; exists {
			return fmt.Errorf("branches: duplicate source branch %q", key)
		}
		seen[key] = struct{}{}
	}
	return nil
}

// validateUniqueDependencies rejects duplicate direct dependencies.
func validateUniqueDependencies(dependencies []Dependency) error {
	seen := map[ModuleName]struct{}{}
	for _, dependency := range dependencies {
		key := dependency.Module()
		if _, exists := seen[key]; exists {
			return fmt.Errorf("dependencies: duplicate dependency %q", key)
		}
		seen[key] = struct{}{}
	}
	return nil
}
