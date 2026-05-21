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

package resolved

// ValueSourceKind identifies where an effective value came from.
type ValueSourceKind string

const (
	// SourceBuiltInDefault marks a value supplied by code-level defaults.
	SourceBuiltInDefault ValueSourceKind = "built-in-default"
	// SourceStagingManifest marks a value supplied by arcpub.yaml top-level fields.
	SourceStagingManifest ValueSourceKind = "staging-manifest"
	// SourceStagingDefault marks a value supplied by arcpub.yaml defaults.
	SourceStagingDefault ValueSourceKind = "staging-default"
	// SourceStagingModule marks a value supplied by arcpub.yaml modules[].
	SourceStagingModule ValueSourceKind = "staging-module"
	// SourceModuleManifest marks a value supplied by arcpub.module.yaml.
	SourceModuleManifest ValueSourceKind = "module-manifest"
)

// ValueSource describes the origin of one effective value.
type ValueSource struct {
	Kind ValueSourceKind
	Path string
}

// ResolvedField records one effective value for diagnostics and explain output.
type ResolvedField struct {
	Path   string
	Value  string
	Source ValueSource
}

// ResolutionTrace records value origins without polluting PublicationSet types.
type ResolutionTrace struct {
	fields []ResolvedField
}

// Add appends one resolved field to the trace.
func (t *ResolutionTrace) Add(path string, value string, source ValueSource) {
	t.fields = append(t.fields, ResolvedField{Path: path, Value: value, Source: source})
}

// Fields returns detached resolved fields.
func (t ResolutionTrace) Fields() []ResolvedField {
	out := make([]ResolvedField, len(t.fields))
	copy(out, t.fields)
	return out
}
