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

// Spec returns a detached serializable representation of the manifest.
func (m Manifest) Spec() Spec {
	return Spec{
		Version: string(m.version),
		Source:  m.source.Spec(),
		Policy:  m.policy.Spec(),
		Modules: manifestModuleSpecs(m.modules),
	}
}

// manifestModuleSpecs converts validated modules back into raw module specs.
func manifestModuleSpecs(modules []Module) []ModuleSpec {
	specs := make([]ModuleSpec, 0, len(modules))
	for _, module := range modules {
		specs = append(specs, module.Spec())
	}
	return specs
}
