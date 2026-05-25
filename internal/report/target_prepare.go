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
	"io"
	"net/url"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

// TargetPrepareReport is the stable DTO for target worktree preparation.
type TargetPrepareReport struct {
	Kind        string                      `json:"kind"`
	Status      Status                      `json:"status"`
	TargetRoot  string                      `json:"targetRoot,omitempty"`
	ModuleCount int                         `json:"moduleCount"`
	Modules     []TargetPrepareModuleReport `json:"modules"`
}

// TargetPrepareModuleReport describes one prepared target worktree.
type TargetPrepareModuleReport struct {
	Name        string                      `json:"name"`
	Repository  string                      `json:"repository"`
	Status      Status                      `json:"status"`
	WorktreeDir string                      `json:"worktreeDir,omitempty"`
	RemoteName  string                      `json:"remoteName,omitempty"`
	RemoteURL   string                      `json:"remoteURL,omitempty"`
	Actions     []TargetPrepareActionReport `json:"actions"`
}

// TargetPrepareActionReport describes one target preparation action.
type TargetPrepareActionReport struct {
	Name      string `json:"name"`
	Status    Status `json:"status"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Path      string `json:"path,omitempty"`
	RemoteURL string `json:"remoteURL,omitempty"`
}

// BuildTargetPrepareReport converts workflow target preparation results into a
// path-safe report DTO.
func BuildTargetPrepareReport(result target.PrepareResult, opts Options) TargetPrepareReport {
	modules := result.Modules()
	out := TargetPrepareReport{
		Kind:        "target-prepare",
		Status:      targetPrepareStatus(result.Status()),
		TargetRoot:  includePath(result.TargetRoot(), opts),
		ModuleCount: len(modules),
		Modules:     make([]TargetPrepareModuleReport, 0, len(modules)),
	}
	for _, mod := range modules {
		report := TargetPrepareModuleReport{
			Name:        mod.Module().String(),
			Repository:  mod.Repository().String(),
			Status:      targetPrepareStatus(mod.Status()),
			WorktreeDir: includePath(mod.WorktreeDir(), opts),
			RemoteName:  mod.RemoteName(),
			RemoteURL:   includeRemoteURL(mod.RemoteURL(), opts),
			Actions:     make([]TargetPrepareActionReport, 0, len(mod.Actions())),
		}
		for _, action := range mod.Actions() {
			report.Actions = append(report.Actions, TargetPrepareActionReport{
				Name:      action.Name(),
				Status:    targetPrepareStatus(action.Status()),
				Code:      action.Code(),
				Message:   action.Message(),
				Path:      includePath(action.Path(), opts),
				RemoteURL: includeRemoteURL(action.RemoteURL(), opts),
			})
		}
		out.Modules = append(out.Modules, report)
	}
	return out
}

func targetPrepareStatus(status target.PrepareStatus) Status {
	switch status {
	case target.PrepareStatusPrepared:
		return StatusPrepared
	case target.PrepareStatusPassed:
		return StatusPassed
	case target.PrepareStatusFailed:
		return StatusFailed
	case target.PrepareStatusSkipped:
		return StatusSkipped
	default:
		return StatusEmpty
	}
}

func includeRemoteURL(value string, opts Options) string {
	value = redactRemoteURL(value)
	if value == "" {
		return ""
	}
	if isLocalRemoteURL(value) && !opts.IncludeLocalPaths {
		return ""
	}
	return value
}

func redactRemoteURL(value string) string {
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err == nil && parsed.User != nil {
		parsed.User = url.User("redacted")
		return parsed.String()
	}
	if strings.Contains(value, "://") {
		return value
	}
	if strings.HasPrefix(value, "git@") {
		return value
	}
	if strings.Contains(value, "@") {
		return ""
	}
	return value
}

func isLocalRemoteURL(value string) bool {
	if strings.HasPrefix(value, "file://") {
		return true
	}
	return isLocalAbsolutePath(value)
}

func writeTargetPrepareText(w io.Writer, report TargetPrepareReport) error {
	if err := writeLine(w, "Target prepare"); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if report.TargetRoot != "" {
		if err := writeLine(w, "  Target root: %s", report.TargetRoot); err != nil {
			return err
		}
	}
	if err := writeLine(w, ""); err != nil {
		return err
	}
	if err := writeLine(w, "Modules:"); err != nil {
		return err
	}
	for _, mod := range report.Modules {
		if err := writeLine(w, "  %s: %s", mod.Name, mod.Status); err != nil {
			return err
		}
		for _, action := range mod.Actions {
			if err := writeLine(w, "    %s: %s", action.Name, action.Status); err != nil {
				return err
			}
		}
	}
	return nil
}
