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

	"arcoris.dev/arcoris-publisher/internal/workflow/preflight"
)

// PreflightReport is the stable DTO for read-only publish readiness checks.
type PreflightReport struct {
	Kind        string                  `json:"kind"`
	Status      Status                  `json:"status"`
	Version     string                  `json:"version,omitempty"`
	Checks      []PreflightCheckReport  `json:"checks"`
	Modules     []PreflightModuleReport `json:"modules"`
	ModuleCount int                     `json:"moduleCount"`
}

// PreflightCheckReport describes one global or module readiness check.
type PreflightCheckReport struct {
	Name     string `json:"name"`
	Status   Status `json:"status"`
	Severity string `json:"severity,omitempty"`
	Code     string `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	Path     string `json:"path,omitempty"`
}

// PreflightModuleReport describes checks for one planned module.
type PreflightModuleReport struct {
	Name        string                 `json:"name"`
	Repository  string                 `json:"repository"`
	WorktreeDir string                 `json:"worktreeDir,omitempty"`
	Status      Status                 `json:"status"`
	Checks      []PreflightCheckReport `json:"checks"`
}

// BuildPreflightReport converts workflow preflight results into a report DTO.
func BuildPreflightReport(result preflight.Result, opts Options) PreflightReport {
	modules := result.Modules()
	out := PreflightReport{
		Kind:        "preflight",
		Status:      preflightStatus(result.Status()),
		Version:     result.Version(),
		Checks:      preflightChecks(result.Checks(), opts),
		Modules:     make([]PreflightModuleReport, 0, len(modules)),
		ModuleCount: len(modules),
	}
	for _, mod := range modules {
		out.Modules = append(out.Modules, PreflightModuleReport{
			Name:        mod.Name().String(),
			Repository:  mod.Repository().String(),
			WorktreeDir: includePath(mod.WorktreeDir(), opts),
			Status:      preflightStatus(mod.Status()),
			Checks:      preflightChecks(mod.Checks(), opts),
		})
	}
	return out
}

func preflightChecks(checks []preflight.CheckResult, opts Options) []PreflightCheckReport {
	out := make([]PreflightCheckReport, 0, len(checks))
	for _, check := range checks {
		out = append(out, PreflightCheckReport{
			Name:     check.Name(),
			Status:   preflightStatus(check.Status()),
			Severity: string(check.Severity()),
			Code:     check.Code(),
			Message:  check.Message(),
			Path:     includePath(check.Path(), opts),
		})
	}
	return out
}

func preflightStatus(status preflight.Status) Status {
	switch status {
	case preflight.StatusPassed:
		return StatusPassed
	case preflight.StatusFailed:
		return StatusFailed
	case preflight.StatusSkipped:
		return StatusSkipped
	case preflight.StatusWarning:
		return StatusWarning
	default:
		return StatusEmpty
	}
}

func writePreflightText(w io.Writer, report PreflightReport) error {
	if err := writeLine(w, "Preflight"); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if report.Version != "" {
		if err := writeLine(w, "  Version: %s", report.Version); err != nil {
			return err
		}
	}
	if err := writeLine(w, ""); err != nil {
		return err
	}
	if err := writeLine(w, "Global checks:"); err != nil {
		return err
	}
	for _, check := range report.Checks {
		if err := writeLine(w, "  %s: %s", check.Name, check.Status); err != nil {
			return err
		}
		if check.Message != "" {
			if err := writeLine(w, "    %s", check.Message); err != nil {
				return err
			}
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
		for _, check := range mod.Checks {
			if err := writeLine(w, "    %s: %s", check.Name, check.Status); err != nil {
				return err
			}
		}
	}
	return nil
}
