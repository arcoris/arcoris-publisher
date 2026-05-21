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

package module

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// DependenciesSpec is the raw module dependency declaration.
type DependenciesSpec struct {
	Internal []string `json:"internal,omitempty" yaml:"internal,omitempty"`
}

// Dependencies is the validated module dependency declaration.
type Dependencies struct {
	internal []manifest.ModuleName
}

// NewDependencies validates spec and returns Dependencies.
func NewDependencies(spec DependenciesSpec) (Dependencies, error) {
	var collector manifest.IssueCollector
	internal := make([]manifest.ModuleName, 0, len(spec.Internal))
	seen := make(map[manifest.ModuleName]int, len(spec.Internal))
	for i, raw := range spec.Internal {
		name, err := manifest.ParseModuleName(raw)
		if err != nil {
			collector.AddError(fmt.Sprintf("internal[%d]", i), err)
			continue
		}
		if prev, exists := seen[name]; exists {
			collector.Add(manifest.IssueDuplicateValue, fmt.Sprintf("internal[%d]", i), "duplicate dependency %q previously declared at internal[%d]", name, prev)
			continue
		}
		seen[name] = i
		internal = append(internal, name)
	}
	if err := collector.Err(); err != nil {
		return Dependencies{}, err
	}
	return Dependencies{internal: internal}, nil
}

// Internal returns detached internal module dependency names.
func (d Dependencies) Internal() []manifest.ModuleName { return manifest.CloneModuleNames(d.internal) }
