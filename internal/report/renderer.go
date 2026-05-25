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

	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/workflow"
	"arcoris.dev/arcoris-publisher/internal/workflow/preflight"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

// Renderer renders publisher reports in a configured output format.
type Renderer struct{ opts Options }

// New returns a report renderer with normalized options.
func New(opts Options) Renderer { return Renderer{opts: normalizeOptions(opts)} }

// Plan renders a publication plan report.
func (r Renderer) Plan(w io.Writer, p plan.Plan) error {
	report := BuildPlanReport(p, r.opts)
	return renderFormatted(w, report, r.opts, writePlanText)
}

// Workflow renders an aggregate workflow report.
func (r Renderer) Workflow(w io.Writer, result workflow.Result) error {
	report := BuildWorkflowReport(result, r.opts)
	return renderFormatted(w, report, r.opts, writeWorkflowText)
}

// Verify renders a verification result report.
func (r Renderer) Verify(w io.Writer, result verify.Result) error {
	report := BuildVerifyReport(result, r.opts)
	return renderFormatted(w, report, r.opts, writeVerifyText)
}

// Publish renders a publication result report.
func (r Renderer) Publish(w io.Writer, result publish.Result) error {
	report := BuildPublishReport(result, r.opts)
	return renderFormatted(w, report, r.opts, writePublishText)
}

// Preflight renders read-only publish readiness checks.
func (r Renderer) Preflight(w io.Writer, result preflight.Result) error {
	report := BuildPreflightReport(result, r.opts)
	return renderFormatted(w, report, r.opts, writePreflightText)
}

// TransactionList renders a publish transaction list report.
func (r Renderer) TransactionList(w io.Writer, summaries []publish.TransactionSummary) error {
	report := BuildTransactionListReport(summaries)
	return renderFormatted(w, report, r.opts, writeTransactionListText)
}

// Transaction renders one publish transaction journal report.
func (r Renderer) Transaction(w io.Writer, journal publish.TransactionJournal) error {
	report := BuildTransactionReport(journal, r.opts)
	return renderFormatted(w, report, r.opts, writeTransactionText)
}

func renderFormatted[T any](
	w io.Writer,
	value T,
	opts Options,
	writeText func(io.Writer, T) error,
) error {
	switch opts.Format {
	case FormatText:
		return writeText(w, value)
	case FormatJSON:
		return writeJSON(w, value, opts.Pretty)
	default:
		return unsupportedFormatError(opts.Format)
	}
}

func unsupportedFormatError(format Format) error {
	return &Error{
		Code:    CodeUnsupportedFormat,
		Message: "unsupported report format: " + format.String(),
	}
}
