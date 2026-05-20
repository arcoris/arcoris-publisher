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

package registry

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// New builds a deterministic lookup registry from validated manifest modules.
//
// The constructor copies the module slice before indexing so callers can reuse
// or reshuffle their input slice without changing the registry's declaration
// order. Duplicate names, module paths, source directories, repositories, and
// branch mappings are reported together as a ValidationError.
func New(modules []manifest.Module) (Registry, error) {
	builder := newBuilder(modules)
	return builder.build()
}

// FromManifest builds a registry from every module declared in manifest.
//
// The manifest package already returns detached module slices. Calling through
// New keeps all registry validation in one path and preserves the same error
// shape whether the source was a manifest or an explicit module list.
func FromManifest(manifestValue manifest.Manifest) (Registry, error) {
	return New(manifestValue.Modules())
}

// Must constructs a registry and panics when modules are invalid.
//
// Must is intended for tests, fixtures, and static wiring where an invalid
// registry is a programmer error. Runtime publication paths should call New and
// return the validation diagnostics to the caller.
func Must(modules []manifest.Module) Registry {
	registry, err := New(modules)
	if err != nil {
		panic(err)
	}
	return registry
}
