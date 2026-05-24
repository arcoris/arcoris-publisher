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

	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
)

// PublishReport is the stable report DTO for publication results.
type PublishReport struct {
	Kind           string                `json:"kind"`
	Present        bool                  `json:"present"`
	Status         Status                `json:"status"`
	ModuleCount    int                   `json:"moduleCount"`
	PublishedCount int                   `json:"publishedCount"`
	SkippedCount   int                   `json:"skippedCount"`
	Modules        []PublishModuleReport `json:"modules"`
}

// PublishModuleReport describes publication outcome for one module.
type PublishModuleReport struct {
	Name    string   `json:"name"`
	Status  Status   `json:"status"`
	Commit  string   `json:"commit,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Pushed  bool     `json:"pushed"`
	Skipped bool     `json:"skipped"`
}

// BuildPublishReport converts publication results into a stable report DTO.
func BuildPublishReport(result publish.Result, _ Options) PublishReport {
	modules := result.Modules()
	out := PublishReport{
		Kind:        "publish",
		Present:     len(modules) > 0,
		Status:      StatusEmpty,
		ModuleCount: len(modules),
		Modules:     make([]PublishModuleReport, 0, len(modules)),
	}
	for _, mod := range modules {
		moduleReport := PublishModuleReport{
			Name:    mod.Module().String(),
			Commit:  string(mod.Commit()),
			Pushed:  mod.Pushed(),
			Skipped: mod.Skipped(),
			Tags:    gitTags(mod.Tags()),
		}
		switch {
		case mod.Published():
			moduleReport.Status = StatusPublished
			out.PublishedCount++
		case mod.Skipped():
			moduleReport.Status = StatusSkipped
			out.SkippedCount++
		default:
			moduleReport.Status = StatusPending
		}
		out.Modules = append(out.Modules, moduleReport)
	}
	switch {
	case out.PublishedCount > 0:
		out.Status = StatusPublished
	case out.ModuleCount == 0:
		out.Status = StatusEmpty
	case out.SkippedCount == out.ModuleCount:
		out.Status = StatusSkipped
	default:
		out.Status = StatusPending
	}
	return out
}

func gitTags[T ~string](tags []T) []string {
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		out = append(out, string(tag))
	}
	return out
}

func writePublishText(w io.Writer, report PublishReport) error {
	if err := writeLine(w, "Publication"); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if err := writeLine(w, "  Modules: %d", report.ModuleCount); err != nil {
		return err
	}
	if err := writeLine(w, "  Published: %d", report.PublishedCount); err != nil {
		return err
	}
	if err := writeLine(w, "  Skipped: %d", report.SkippedCount); err != nil {
		return err
	}
	if err := writeLine(w, ""); err != nil {
		return err
	}
	for _, mod := range report.Modules {
		if err := writeLine(w, "  %s: %s", mod.Name, mod.Status); err != nil {
			return err
		}
		if mod.Commit != "" {
			if err := writeLine(w, "    commit: %s", mod.Commit); err != nil {
				return err
			}
		}
		if err := writeLine(w, "    tags:   %s", commaOrDash(mod.Tags)); err != nil {
			return err
		}
		if err := writeLine(w, ""); err != nil {
			return err
		}
	}
	return nil
}
