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

// Manifest is a validated module-level arcpub.module.yaml resource.
type Manifest struct {
	apiVersion   manifest.APIVersion
	kind         manifest.Kind
	metadata     manifest.Metadata
	module       manifest.ModuleIdentity
	dependencies Dependencies
	publish      Publish
	verification manifest.VerificationOverride
}

// New validates spec and returns a module manifest.
func New(spec Spec) (Manifest, error) {
	var collector manifest.IssueCollector
	apiVersion, err := manifest.ValidateAPIVersion(spec.APIVersion)
	collector.AddError("apiVersion", err)
	kind, err := manifest.ValidateKind(spec.Kind, manifest.KindModuleManifest)
	collector.AddError("kind", err)
	metadata, err := manifest.NewMetadata(spec.Metadata)
	collector.AddError("metadata", err)
	identity, err := manifest.NewModuleIdentity(spec.Module)
	collector.AddError("module", err)
	dependencies, err := NewDependencies(spec.Dependencies)
	collector.AddError("dependencies", err)
	publish, err := NewPublish(spec.Publish)
	collector.AddError("publish", err)
	verification, err := manifest.NewVerificationOverride(spec.Verification)
	collector.AddError("verification", err)
	if err := collector.Err(); err != nil {
		return Manifest{}, err
	}
	return Manifest{apiVersion: apiVersion, kind: kind, metadata: metadata, module: identity, dependencies: dependencies, publish: publish, verification: verification}, nil
}

// APIVersion returns the module manifest API version.
func (m Manifest) APIVersion() manifest.APIVersion { return m.apiVersion }

// Kind returns the module manifest kind.
func (m Manifest) Kind() manifest.Kind { return m.kind }

// Metadata returns the module manifest metadata.
func (m Manifest) Metadata() manifest.Metadata { return m.metadata }

// Module returns the module identity declaration.
func (m Manifest) Module() manifest.ModuleIdentity { return m.module }

// Dependencies returns the module dependency declaration.
func (m Manifest) Dependencies() Dependencies { return m.dependencies }

// Publish returns explicit module publication entries.
func (m Manifest) Publish() Publish { return m.publish }

// Verification returns the module-level verification override.
func (m Manifest) Verification() manifest.VerificationOverride { return m.verification }
