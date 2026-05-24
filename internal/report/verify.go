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

	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

// VerifyReport is the stable report DTO for verification results.
type VerifyReport struct {
	Kind        string               `json:"kind"`
	Present     bool                 `json:"present"`
	Status      Status               `json:"status"`
	ModuleCount int                  `json:"moduleCount"`
	FailedCount int                  `json:"failedCount"`
	Modules     []VerifyModuleReport `json:"modules"`
}

// VerifyModuleReport describes verification checks for one module.
type VerifyModuleReport struct {
	Name        string              `json:"name"`
	Status      Status              `json:"status"`
	FailedCount int                 `json:"failedCount"`
	Checks      []VerifyCheckReport `json:"checks"`
}

// VerifyCheckReport describes one verification check outcome.
type VerifyCheckReport struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Severity string `json:"severity"`
	Message  string `json:"message,omitempty"`
	Path     string `json:"path,omitempty"`
}

// BuildVerifyReport converts verification results into a stable report DTO.
func BuildVerifyReport(result verify.Result, opts Options) VerifyReport {
	modules := result.Modules()
	out := VerifyReport{
		Kind:        "verify",
		Present:     len(modules) > 0,
		Status:      verifyStatus(result),
		ModuleCount: len(modules),
		Modules:     make([]VerifyModuleReport, 0, len(modules)),
	}
	for _, mod := range modules {
		checks := mod.Checks()
		moduleReport := VerifyModuleReport{
			Name:   mod.Module().String(),
			Status: statusPassedFailed(mod.Failed()),
			Checks: make([]VerifyCheckReport, 0, len(checks)),
		}
		for _, check := range checks {
			checkReport := VerifyCheckReport{
				Name:     string(check.Name()),
				Status:   string(check.Status()),
				Severity: string(check.Severity()),
				Message:  check.Message(),
				Path:     includePath(check.Path(), opts),
			}
			if check.Status() == verify.StatusFailed {
				moduleReport.FailedCount++
				out.FailedCount++
			}
			moduleReport.Checks = append(moduleReport.Checks, checkReport)
		}
		out.Modules = append(out.Modules, moduleReport)
	}
	return out
}

func verifyStatus(result verify.Result) Status {
	if len(result.Modules()) == 0 {
		return StatusEmpty
	}
	return statusPassedFailed(result.Failed())
}

func writeVerifyText(w io.Writer, report VerifyReport) error {
	if err := writeLine(w, "Verification"); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if err := writeLine(w, "  Modules: %d", report.ModuleCount); err != nil {
		return err
	}
	if err := writeLine(w, "  Failed checks: %d", report.FailedCount); err != nil {
		return err
	}
	if err := writeLine(w, ""); err != nil {
		return err
	}
	for _, mod := range report.Modules {
		if err := writeLine(w, "  %s: %s", mod.Name, mod.Status); err != nil {
			return err
		}
		for _, check := range mod.Checks {
			marker := "-"
			switch check.Status {
			case string(verify.StatusPassed):
				marker = "ok"
			case string(verify.StatusFailed):
				marker = "fail"
			case string(verify.StatusSkipped):
				marker = "skip"
			case string(verify.StatusWarning):
				marker = "warn"
			}
			if check.Message == "" {
				if err := writeLine(w, "    %s %s [%s]", marker, check.Name, check.Severity); err != nil {
					return err
				}
				continue
			}
			if err := writeLine(w, "    %s %s [%s]: %s", marker, check.Name, check.Severity, check.Message); err != nil {
				return err
			}
		}
		if err := writeLine(w, ""); err != nil {
			return err
		}
	}
	return nil
}
