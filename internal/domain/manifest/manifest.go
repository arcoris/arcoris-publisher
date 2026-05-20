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

// Package manifest defines the validated domain model for ARCORIS Publisher
// manifest files.
//
// The package describes static publication policy: source repository, publishable
// modules, target repositories, branch mappings, dependency declarations, and
// publication policy. It deliberately does not load YAML, access the filesystem,
// call Git, invoke the Go toolchain, or execute publication workflows.
package manifest

// Spec is the raw, serializable representation of an ARCORIS Publisher manifest.
//
// Spec is intentionally DTO-like. Configuration loaders may decode YAML, JSON, or
// another format into Spec and then call New to obtain a validated Manifest value.
type Spec struct {
	// Version is the manifest schema version, for example "v1".
	Version string `json:"version" yaml:"version"`
	// Source describes the authoritative source repository.
	Source SourceSpec `json:"source" yaml:"source"`
	// Policy describes global publication policy.
	Policy PolicySpec `json:"policy" yaml:"policy"`
	// Modules declares all modules known to the publisher.
	Modules []ModuleSpec `json:"modules" yaml:"modules"`
}

// Manifest is a validated publication-policy aggregate.
//
// Manifest values are immutable by convention. Accessors return detached slices
// so caller mutations cannot alter a validated aggregate after construction.
type Manifest struct {
	version Version
	source  Source
	policy  Policy
	modules []Module
}
