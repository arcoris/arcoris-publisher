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

package report

import (
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
)

// SourceReport summarizes inspected source repository state without local paths
// by default.
type SourceReport struct {
	Present     bool                   `json:"present"`
	Status      Status                 `json:"status"`
	ModuleCount int                    `json:"moduleCount"`
	Repository  SourceRepositoryReport `json:"repository"`
	Warnings    []IssueReport          `json:"warnings,omitempty"`
	Modules     []SourceModuleReport   `json:"modules"`
}

// SourceRepositoryReport describes source Git state.
type SourceRepositoryReport struct {
	Head          string `json:"head"`
	Branch        string `json:"branch"`
	Dirty         bool   `json:"dirty"`
	RepositoryDir string `json:"repositoryDir,omitempty"`
	StagingDir    string `json:"stagingDir,omitempty"`
}

// SourceModuleReport describes source-side state for one module.
type SourceModuleReport struct {
	Name          string `json:"name"`
	Hash          string `json:"hash"`
	EntryCount    int    `json:"entryCount"`
	SourceDir     string `json:"sourceDir,omitempty"`
	ModuleRootDir string `json:"moduleRootDir,omitempty"`
}

// IssueReport describes a non-fatal workflow diagnostic.
type IssueReport struct {
	Code    string `json:"code,omitempty"`
	Module  string `json:"module,omitempty"`
	Path    string `json:"path,omitempty"`
	Message string `json:"message,omitempty"`
}

// BuildSourceReport converts a source snapshot into a stable report DTO.
func BuildSourceReport(snapshot source.Snapshot, opts Options) SourceReport {
	modules := snapshot.Modules()
	repository := snapshot.Repository()
	warnings := snapshot.Warnings()
	present := len(modules) > 0 ||
		len(warnings) > 0 ||
		repository.Head() != "" ||
		repository.Branch() != "" ||
		repository.RepositoryDir() != "" ||
		repository.StagingDir() != ""

	dirty := false
	if present {
		dirty = repository.Dirty()
	}

	out := SourceReport{
		Present:     present,
		Status:      statusEmptyPresent(present),
		ModuleCount: len(modules),
		Repository: SourceRepositoryReport{
			Head:          string(repository.Head()),
			Branch:        string(repository.Branch()),
			Dirty:         dirty,
			RepositoryDir: includePath(repository.RepositoryDir(), opts),
			StagingDir:    includePath(repository.StagingDir(), opts),
		},
		Modules: make([]SourceModuleReport, 0, len(modules)),
	}
	for _, issue := range warnings {
		report := IssueReport{
			Code:    string(issue.Code),
			Module:  issue.Module.String(),
			Path:    includePath(issue.Path, opts),
			Message: issue.Message,
		}
		out.Warnings = append(out.Warnings, report)
	}
	for _, mod := range modules {
		out.Modules = append(out.Modules, SourceModuleReport{
			Name:          mod.Name().String(),
			Hash:          mod.Hash().String(),
			EntryCount:    len(mod.Entries()),
			SourceDir:     includePath(mod.SourceDir(), opts),
			ModuleRootDir: includePath(mod.ModuleRootDir(), opts),
		})
	}
	return out
}
