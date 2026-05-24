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

import "arcoris.dev/arcoris-publisher/internal/workflow/target"

// TargetReport summarizes prepared target workspaces.
type TargetReport struct {
	Present        bool                    `json:"present"`
	Status         Status                  `json:"status"`
	WorkspaceCount int                     `json:"workspaceCount"`
	Workspaces     []TargetWorkspaceReport `json:"workspaces"`
}

// TargetWorkspaceReport describes one target workspace without local paths by
// default.
type TargetWorkspaceReport struct {
	Module      string         `json:"module"`
	Repository  string         `json:"repository"`
	WorktreeDir string         `json:"worktreeDir,omitempty"`
	Branches    []BranchReport `json:"branches"`
}

// BuildTargetReport converts target workspaces into a stable report DTO.
func BuildTargetReport(workspaces target.WorkspaceSet, opts Options) TargetReport {
	items := workspaces.Workspaces()
	out := TargetReport{
		Present:        len(items) > 0,
		Status:         statusEmptyPresent(len(items) > 0),
		WorkspaceCount: len(items),
		Workspaces:     make([]TargetWorkspaceReport, 0, len(items)),
	}
	for _, workspace := range items {
		report := TargetWorkspaceReport{
			Module:      workspace.Module().String(),
			Repository:  workspace.Repository().String(),
			WorktreeDir: includePath(workspace.WorktreeDir(), opts),
			Branches:    make([]BranchReport, 0, len(workspace.Branches())),
		}
		for _, branch := range workspace.Branches() {
			report.Branches = append(report.Branches, BranchReport{
				Source: branch.Source().String(),
				Target: branch.Target().String(),
			})
		}
		out.Workspaces = append(out.Workspaces, report)
	}
	return out
}
