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

// MetadataSpec is the raw metadata block shared by all manifest resources.
type MetadataSpec struct {
	Name string `json:"name" yaml:"name"`
}

// Metadata is the validated metadata block shared by all manifest resources.
type Metadata struct {
	name ManifestName
}

// NewMetadata validates spec and returns a Metadata value.
func NewMetadata(spec MetadataSpec) (Metadata, error) {
	name, err := ParseManifestName(spec.Name)
	if err != nil {
		return Metadata{}, err
	}
	return Metadata{name: name}, nil
}

// Name returns the validated manifest resource name.
func (m Metadata) Name() ManifestName { return m.name }

// Spec returns a serializable metadata representation.
func (m Metadata) Spec() MetadataSpec { return MetadataSpec{Name: string(m.name)} }
