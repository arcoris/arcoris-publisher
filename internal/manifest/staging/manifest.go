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

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// Manifest is a validated top-level arcpub.yaml resource.
type Manifest struct {
	apiVersion manifest.APIVersion
	kind       manifest.Kind
	metadata   manifest.Metadata
	source     manifest.Source
	target     manifest.TargetPolicy
	publish    manifest.PublishPolicy
	defaults   Defaults
	modules    []Module
}

// New validates spec and returns a staging manifest.
func New(spec Spec) (Manifest, error) {
	var collector manifest.IssueCollector

	apiVersion, err := manifest.ValidateAPIVersion(spec.APIVersion)
	collector.AddError("apiVersion", err)

	kind, err := manifest.ValidateKind(spec.Kind, manifest.KindStagingManifest)
	collector.AddError("kind", err)

	metadata, err := manifest.NewMetadata(spec.Metadata)
	collector.AddError("metadata", err)

	source, err := manifest.NewSource(spec.Source)
	collector.AddError("source", err)

	target, err := manifest.NewTargetPolicy(spec.Target)
	collector.AddError("target", err)

	publish, err := manifest.NewPublishPolicy(spec.Publish)
	collector.AddError("publish", err)

	defaults, err := NewDefaults(spec.Defaults)
	collector.AddError("defaults", err)

	modules := make([]Module, 0, len(spec.Modules))
	for i, moduleSpec := range spec.Modules {
		mod, err := NewModule(moduleSpec)
		if err != nil {
			collector.AddError(fmt.Sprintf("modules[%d]", i), err)
			continue
		}
		modules = append(modules, mod)
	}
	if len(spec.Modules) == 0 {
		collector.Add(manifest.IssueMissingField, "modules", "at least one module is required")
	}

	if err := collector.Err(); err != nil {
		return Manifest{}, err
	}

	out := Manifest{
		apiVersion: apiVersion,
		kind:       kind,
		metadata:   metadata,
		source:     source,
		target:     target,
		publish:    publish,
		defaults:   defaults,
		modules:    modules,
	}
	if err := out.Validate(); err != nil {
		return Manifest{}, err
	}

	return out, nil
}

// APIVersion returns the staging manifest API version.
func (m Manifest) APIVersion() manifest.APIVersion { return m.apiVersion }

// Kind returns the staging manifest kind.
func (m Manifest) Kind() manifest.Kind { return m.kind }

// Metadata returns the staging manifest metadata.
func (m Manifest) Metadata() manifest.Metadata { return m.metadata }

// Source returns the source repository declaration.
func (m Manifest) Source() manifest.Source { return m.source }

// Target returns the target worktree preparation policy.
func (m Manifest) Target() manifest.TargetPolicy { return m.target }

// Publish returns the global publication policy.
func (m Manifest) Publish() manifest.PublishPolicy { return m.publish }

// Defaults returns top-level defaults.
func (m Manifest) Defaults() Defaults { return m.defaults }

// Modules returns detached top-level module declarations.
func (m Manifest) Modules() []Module {
	out := make([]Module, len(m.modules))
	copy(out, m.modules)
	return out
}
