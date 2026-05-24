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

import "arcoris.dev/arcoris-publisher/internal/workflow/construct"

// ConstructReport summarizes explicit-projection construction results.
type ConstructReport struct {
	Present        bool                    `json:"present"`
	Status         string                  `json:"status"`
	Changed        bool                    `json:"changed"`
	ModuleCount    int                     `json:"moduleCount"`
	OperationCount int                     `json:"operationCount"`
	Modules        []ConstructModuleReport `json:"modules"`
}

// ConstructModuleReport describes construction actions for one module.
type ConstructModuleReport struct {
	Name           string                     `json:"name"`
	Changed        bool                       `json:"changed"`
	WorktreeDir    string                     `json:"worktreeDir,omitempty"`
	OperationCount int                        `json:"operationCount"`
	Operations     []ConstructOperationReport `json:"operations"`
}

// ConstructOperationReport describes one construction operation.
type ConstructOperationReport struct {
	Kind       string `json:"kind"`
	SourcePath string `json:"sourcePath,omitempty"`
	TargetPath string `json:"targetPath,omitempty"`
}

// BuildConstructReport converts construction results into a report DTO.
func BuildConstructReport(result construct.Result, opts Options) ConstructReport {
	modules := result.Modules()
	out := ConstructReport{
		Present:     len(modules) > 0,
		Status:      statusEmptyPresent(len(modules) > 0),
		Changed:     result.Changed(),
		ModuleCount: len(modules),
		Modules:     make([]ConstructModuleReport, 0, len(modules)),
	}
	for _, mod := range modules {
		operations := mod.Operations()
		moduleReport := ConstructModuleReport{
			Name:           mod.Module().String(),
			Changed:        mod.Changed(),
			WorktreeDir:    includePath(mod.WorktreeDir(), opts),
			OperationCount: len(operations),
			Operations:     make([]ConstructOperationReport, 0, len(operations)),
		}
		for _, op := range operations {
			operationReport := ConstructOperationReport{
				Kind:       string(op.Kind()),
				SourcePath: includePath(op.SourcePath(), opts),
				TargetPath: includePath(op.TargetPath(), opts),
			}
			moduleReport.Operations = append(moduleReport.Operations, operationReport)
		}
		out.OperationCount += len(operations)
		out.Modules = append(out.Modules, moduleReport)
	}
	return out
}
