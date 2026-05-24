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
	"errors"
	"strings"
	"testing"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("writer failed")
}

func TestRendererUnsupportedFormatReturnsTypedError(t *testing.T) {
	t.Parallel()

	err := New(Options{Format: Format("yaml")}).Plan(failingWriter{}, reportPlan(t).Plan)
	if err == nil {
		t.Fatal("Plan() error = nil")
	}

	var reportErr *Error
	if !errors.As(err, &reportErr) || reportErr.Code != CodeUnsupportedFormat {
		t.Fatalf("Plan() error = %v", err)
	}
}

func TestRendererWriterErrorReturnsTypedError(t *testing.T) {
	t.Parallel()

	result := reportWorkflowResult(t, workflowReportFixture{
		publish:      true,
		dirtyModules: []string{"foundation", "control"},
	})
	failedVerify := reportWorkflowResult(t, workflowReportFixture{verifyFails: true}).Verify()

	tests := []struct {
		name string
		run  func() error
	}{
		{name: "plan", run: func() error {
			return New(Options{Format: FormatText}).Plan(failingWriter{}, reportPlan(t).Plan)
		}},
		{name: "workflow", run: func() error {
			return New(Options{Format: FormatText}).Workflow(failingWriter{}, result)
		}},
		{name: "verify failed", run: func() error {
			return New(Options{Format: FormatText}).Verify(failingWriter{}, failedVerify)
		}},
		{name: "publish skipped", run: func() error {
			return New(Options{Format: FormatText}).Publish(
				failingWriter{},
				reportWorkflowResult(t, workflowReportFixture{publish: true}).Publish(),
			)
		}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.run()
			if err == nil {
				t.Fatal("render error = nil")
			}

			var reportErr *Error
			if !errors.As(err, &reportErr) || reportErr.Code != CodeWriteFailed {
				t.Fatalf("render error = %v", err)
			}
			if !strings.Contains(err.Error(), "writer failed") {
				t.Fatalf("writer error not wrapped: %v", err)
			}
		})
	}
}

func TestRendererPublishSkippedIsNotAnError(t *testing.T) {
	t.Parallel()

	result := reportWorkflowResult(t, workflowReportFixture{publish: true}).Publish()
	if err := New(Options{Format: FormatText}).Publish(noopWriter{}, result); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
}

type noopWriter struct{}

func (noopWriter) Write(data []byte) (int, error) {
	return len(data), nil
}
