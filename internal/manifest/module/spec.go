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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Spec is the raw serializable shape of arcpub.module.yaml.
type Spec struct {
	APIVersion   string                      `json:"apiVersion" yaml:"apiVersion"`
	Kind         string                      `json:"kind" yaml:"kind"`
	Metadata     manifest.MetadataSpec       `json:"metadata" yaml:"metadata"`
	Module       manifest.ModuleIdentitySpec `json:"module" yaml:"module"`
	Dependencies DependenciesSpec            `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Publish      PublishSpec                 `json:"publish" yaml:"publish"`
	Verification manifest.VerificationSpec   `json:"verification,omitempty" yaml:"verification,omitempty"`
}
