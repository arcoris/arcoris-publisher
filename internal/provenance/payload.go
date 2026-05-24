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

package provenance

import "encoding/json"

// SchemaVersion identifies the stable provenance payload schema.
const SchemaVersion = "arcoris.provenance/v1"

// FilePayload is the deterministic JSON provenance document written to target
// repositories when file provenance is enabled.
type FilePayload struct {
	SchemaVersion string             `json:"schemaVersion"`
	Publisher     PublisherPayload   `json:"publisher"`
	Module        ModulePayload      `json:"module"`
	Source        SourcePayload      `json:"source"`
	Target        TargetPayload      `json:"target"`
	Publication   PublicationPayload `json:"publication"`
	Projection    ProjectionPayload  `json:"projection"`
}

// PublisherPayload describes the publisher binary that produced the output.
type PublisherPayload struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Dirty   string `json:"dirty"`
}

// ModulePayload identifies the public module being published.
type ModulePayload struct {
	Name       string `json:"name"`
	ModulePath string `json:"modulePath"`
	Version    string `json:"version"`
	SourceDir  string `json:"sourceDir"`
}

// SourcePayload describes repository-level source state without local paths.
type SourcePayload struct {
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
	Branch     string `json:"branch"`
	Hash       string `json:"hash"`
}

// TargetPayload describes the target repository and branch mappings.
type TargetPayload struct {
	Repository string   `json:"repository"`
	Branches   []string `json:"branches"`
}

// PublicationPayload records global publication policies relevant to this
// output.
type PublicationPayload struct {
	PublishMode string `json:"publishMode"`
	PushPolicy  string `json:"pushPolicy"`
	TagPolicy   string `json:"tagPolicy"`
}

// ProjectionPayload summarizes explicit published entries.
type ProjectionPayload struct {
	EntryCount     int     `json:"entryCount"`
	ProjectionHash string  `json:"projectionHash"`
	Entries        []Entry `json:"entries"`
}

// BuildFilePayload builds a stable provenance document from resolved workflow
// state.
func BuildFilePayload(input Input) FilePayload {
	entries := EntriesFromSourceModule(input.SourceModule)
	build := input.Build

	return FilePayload{
		SchemaVersion: SchemaVersion,
		Publisher: PublisherPayload{
			Version: build.Version(),
			Commit:  build.Commit(),
			Date:    build.Date(),
			Dirty:   build.Dirty(),
		},
		Module: ModulePayload{
			Name:       input.Module.Name().String(),
			ModulePath: input.Module.ModulePath().String(),
			Version:    input.Module.Version().String(),
			SourceDir:  input.Module.SourceDir().String(),
		},
		Source: SourcePayload{
			Repository: input.Plan.Source().Repository().String(),
			Commit:     string(input.Source.Repository().Head()),
			Branch:     string(input.Source.Repository().Branch()),
			Hash:       input.SourceModule.Hash().String(),
		},
		Target: TargetPayload{
			Repository: input.Module.Repository().String(),
			Branches:   input.targetBranches(),
		},
		Publication: PublicationPayload{
			PublishMode: string(input.Plan.PublishPolicy().Mode()),
			PushPolicy:  string(input.Plan.PublishPolicy().PushPolicy()),
			TagPolicy:   string(input.Plan.PublishPolicy().Tags().Mode()),
		},
		Projection: buildProjectionPayload(entries),
	}
}

// RenderFilePayload renders indented JSON with a trailing newline.
func RenderFilePayload(input Input) ([]byte, error) {
	data, err := json.MarshalIndent(BuildFilePayload(input), "", "  ")
	if err != nil {
		return nil, err
	}

	return append(data, '\n'), nil
}

func buildProjectionPayload(entries []Entry) ProjectionPayload {
	entries = normalizeEntries(entries)

	return ProjectionPayload{
		EntryCount:     len(entries),
		ProjectionHash: ProjectionHash(entries),
		Entries:        entries,
	}
}
