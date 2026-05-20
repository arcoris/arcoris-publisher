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

// New validates spec and returns a Manifest aggregate.
func New(spec Spec) (Manifest, error) {
	version, err := parseManifestVersion(spec)
	if err != nil {
		return Manifest{}, err
	}
	source, err := parseManifestSource(spec)
	if err != nil {
		return Manifest{}, err
	}
	policy, err := parseManifestPolicy(spec)
	if err != nil {
		return Manifest{}, err
	}
	modules, err := parseManifestModules(spec.Modules)
	if err != nil {
		return Manifest{}, err
	}
	manifest := Manifest{version: version, source: source, policy: policy, modules: modules}
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

// Must constructs a Manifest and panics when spec is invalid.
//
// Must is intended for tests and static wiring. Runtime code should call New and
// return diagnostics to the caller.
func Must(spec Spec) Manifest {
	manifest, err := New(spec)
	if err != nil {
		panic(err)
	}
	return manifest
}
