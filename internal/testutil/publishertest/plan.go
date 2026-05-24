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

// Package publishertest provides small production-constructor fixtures for
// package tests.
package publishertest

import (
	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

// Module describes one module in a realistic test publication set.
type Module struct {
	// Name is the staging and module manifest name.
	Name string

	// Dependencies are direct internal dependency names.
	Dependencies []string

	// Entries overrides the default explicit publication entries.
	Entries []manifest.PublishEntrySpec
}

// PlanOptions configures a realistic test publication plan.
type PlanOptions struct {
	// Version is the release-train version assigned to public modules.
	Version string

	// Publish customizes the top-level publish policy.
	Publish manifest.PublishSpec
}

// Plan builds a publication plan through manifest constructors, resolution,
// registry, graph, versioning, and plan builders.
func Plan(opts PlanOptions, modules ...Module) (plan.Plan, error) {
	if opts.Version == "" {
		opts.Version = "v0.3.0"
	}

	set, err := PublicationSet(opts, modules...)
	if err != nil {
		return plan.Plan{}, err
	}

	version, err := versioning.Parse(opts.Version)
	if err != nil {
		return plan.Plan{}, err
	}

	return plan.FromPublicationSet(set, version)
}

// PublicationSet builds a resolved publication set from realistic test
// manifests.
func PublicationSet(opts PlanOptions, modules ...Module) (resolved.PublicationSet, error) {
	stagingManifest, err := staging.New(stagingSpec(opts, modules))
	if err != nil {
		return resolved.PublicationSet{}, err
	}

	moduleManifests := make([]modulemanifest.Manifest, 0, len(modules))
	for _, mod := range modules {
		moduleManifest, err := modulemanifest.New(moduleSpec(mod))
		if err != nil {
			return resolved.PublicationSet{}, err
		}
		moduleManifests = append(moduleManifests, moduleManifest)
	}

	return resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest,
		Modules: moduleManifests,
	})
}

func stagingSpec(opts PlanOptions, modules []Module) staging.Spec {
	stagingModules := make([]staging.ModuleSpec, 0, len(modules))
	for _, mod := range modules {
		stagingModules = append(stagingModules, staging.ModuleSpec{
			Name:       mod.Name,
			SourceDir:  "src/arcoris.dev/" + mod.Name,
			Repository: "arcoris/" + mod.Name,
		})
	}

	return staging.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindStagingManifest),
		Metadata:   manifest.MetadataSpec{Name: "arcoris"},
		Source: manifest.SourceSpec{
			Repository:    "arcoris/arcoris",
			DefaultBranch: "main",
		},
		Publish: opts.Publish,
		Modules: stagingModules,
	}
}

func moduleSpec(mod Module) modulemanifest.Spec {
	entries := mod.Entries
	if entries == nil {
		entries = DefaultEntries()
	}

	return modulemanifest.Spec{
		APIVersion:   string(manifest.APIVersionV1Alpha1),
		Kind:         string(manifest.KindModuleManifest),
		Metadata:     manifest.MetadataSpec{Name: mod.Name},
		Module:       manifest.ModuleIdentitySpec{Path: "arcoris.dev/" + mod.Name},
		Dependencies: modulemanifest.DependenciesSpec{Internal: mod.Dependencies},
		Publish:      modulemanifest.PublishSpec{Entries: entries},
	}
}

// DefaultEntries returns the standard explicit publication entries used by
// workflow tests.
func DefaultEntries() []manifest.PublishEntrySpec {
	return []manifest.PublishEntrySpec{
		FileEntry("go.mod"),
		DirectoryEntry("contracts"),
	}
}

// FileEntry returns an explicit file publication entry whose source and target
// path are identical.
func FileEntry(path string) manifest.PublishEntrySpec {
	return manifest.PublishEntrySpec{
		Type: string(manifest.PublishEntryFile),
		From: path,
		To:   path,
	}
}

// DirectoryEntry returns an explicit directory publication entry whose source
// and target path are identical.
func DirectoryEntry(path string) manifest.PublishEntrySpec {
	return manifest.PublishEntrySpec{
		Type: string(manifest.PublishEntryDirectory),
		From: path,
		To:   path,
	}
}
