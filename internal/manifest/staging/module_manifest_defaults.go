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

package staging

import "arcoris.dev/arcoris-publisher/internal/manifest"

// ModuleManifestDefaultsSpec declares default arcpub.module.yaml location rules.
type ModuleManifestDefaultsSpec struct {
	Path     *string `json:"path,omitempty" yaml:"path,omitempty"`
	Required *bool   `json:"required,omitempty" yaml:"required,omitempty"`
}

// ModuleManifestDefaults is the validated default module manifest location policy.
type ModuleManifestDefaults struct {
	path     manifest.RelativePath
	required bool
}

// NewModuleManifestDefaults validates spec and applies safe defaults.
func NewModuleManifestDefaults(spec ModuleManifestDefaultsSpec) (ModuleManifestDefaults, error) {
	path, err := manifest.ParseRelativePath("moduleManifest.path", stringOrDefault(spec.Path, "arcpub.module.yaml"), false)
	if err != nil {
		return ModuleManifestDefaults{}, err
	}
	return ModuleManifestDefaults{path: path, required: boolOrDefault(spec.Required, true)}, nil
}

// Path returns the default arcpub.module.yaml path relative to each sourceDir.
func (d ModuleManifestDefaults) Path() manifest.RelativePath { return d.path }

// Required reports whether every staging module must have a module manifest.
func (d ModuleManifestDefaults) Required() bool { return d.required }
