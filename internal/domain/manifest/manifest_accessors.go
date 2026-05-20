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

// Version returns the manifest schema version.
func (m Manifest) Version() Version { return m.version }

// Source returns the authoritative source repository declaration.
func (m Manifest) Source() Source { return m.source }

// Policy returns the global publication policy.
func (m Manifest) Policy() Policy { return m.policy }

// Modules returns a detached copy of validated modules.
func (m Manifest) Modules() []Module {
	return append([]Module(nil), m.modules...)
}

// ModuleByName returns the module with name and whether it was found.
func (m Manifest) ModuleByName(name ModuleName) (Module, bool) {
	for _, module := range m.modules {
		if module.Name() == name {
			return module, true
		}
	}
	return Module{}, false
}
