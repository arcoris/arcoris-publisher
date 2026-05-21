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

// ModuleIdentitySpec is the raw module identity declaration from arcpub.module.yaml.
type ModuleIdentitySpec struct {
	Type  *string `json:"type,omitempty" yaml:"type,omitempty"`
	Path  string  `json:"path" yaml:"path"`
	Root  *string `json:"root,omitempty" yaml:"root,omitempty"`
	GoMod *string `json:"goMod,omitempty" yaml:"goMod,omitempty"`
}

// ModuleIdentity is the validated module identity declaration.
type ModuleIdentity struct {
	moduleType ModuleType
	path       ModulePath
	root       RelativePath
	goMod      RelativePath
}

// NewModuleIdentity validates spec and applies module identity defaults.
func NewModuleIdentity(spec ModuleIdentitySpec) (ModuleIdentity, error) {
	moduleType, err := ParseModuleType(stringValue(spec.Type, string(ModuleTypeGo)))
	if err != nil {
		return ModuleIdentity{}, err
	}
	modulePath, err := ParseModulePath(spec.Path)
	if err != nil {
		return ModuleIdentity{}, err
	}
	root, err := ParseRelativePath("module.root", stringValue(spec.Root, "."), true)
	if err != nil {
		return ModuleIdentity{}, err
	}
	goMod, err := ParseRelativePath("module.goMod", stringValue(spec.GoMod, "go.mod"), false)
	if err != nil {
		return ModuleIdentity{}, err
	}
	return ModuleIdentity{moduleType: moduleType, path: modulePath, root: root, goMod: goMod}, nil
}

// Type returns the module type.
func (i ModuleIdentity) Type() ModuleType { return i.moduleType }

// Path returns the published module path.
func (i ModuleIdentity) Path() ModulePath { return i.path }

// Root returns the module root relative to arcpub.module.yaml.
func (i ModuleIdentity) Root() RelativePath { return i.root }

// GoMod returns the go.mod path relative to the module root.
func (i ModuleIdentity) GoMod() RelativePath { return i.goMod }
