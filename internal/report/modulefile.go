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

import "arcoris.dev/arcoris-publisher/internal/workflow/modulefile"

// ModuleFileReport summarizes target go.mod rewrites.
type ModuleFileReport struct {
	Present     bool                     `json:"present"`
	Status      string                   `json:"status"`
	Changed     bool                     `json:"changed"`
	ModuleCount int                      `json:"moduleCount"`
	Modules     []ModuleFileModuleReport `json:"modules"`
}

// ModuleFileModuleReport describes one rewritten go.mod.
type ModuleFileModuleReport struct {
	Name         string                        `json:"name"`
	Changed      bool                          `json:"changed"`
	GoModPath    string                        `json:"goModPath,omitempty"`
	Requirements []DependencyRequirementReport `json:"requirements"`
}

// BuildModuleFileReport converts modulefile results into a report DTO.
func BuildModuleFileReport(result modulefile.Result, opts Options) ModuleFileReport {
	modules := result.Modules()
	out := ModuleFileReport{
		Present:     len(modules) > 0,
		Status:      statusEmptyPresent(len(modules) > 0),
		Changed:     result.Changed(),
		ModuleCount: len(modules),
		Modules:     make([]ModuleFileModuleReport, 0, len(modules)),
	}
	for _, mod := range modules {
		requirements := mod.Requirements()
		moduleReport := ModuleFileModuleReport{
			Name:         mod.Module().String(),
			Changed:      mod.Changed(),
			GoModPath:    includePath(mod.GoModPath(), opts),
			Requirements: make([]DependencyRequirementReport, 0, len(requirements)),
		}
		for _, req := range requirements {
			moduleReport.Requirements = append(moduleReport.Requirements, DependencyRequirementReport{
				ModulePath: req.ModulePath().String(),
				Version:    req.Version(),
			})
		}
		out.Modules = append(out.Modules, moduleReport)
	}
	return out
}
