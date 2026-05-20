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

package gotoolchain

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

type fakeRunner struct {
	specs   []processport.Spec
	results []processport.Result
	errs    []error
}

// Run records every Go process spec and returns the configured result/error.
func (f *fakeRunner) Run(ctx context.Context, spec processport.Spec) (processport.Result, error) {
	f.specs = append(f.specs, spec)
	index := len(f.specs) - 1
	var result processport.Result
	if index < len(f.results) {
		result = f.results[index]
	}
	result.Name = spec.Name
	result.Args = append([]string(nil), spec.Args...)
	result.Dir = spec.Dir
	var err error
	if index < len(f.errs) {
		err = f.errs[index]
	}
	return result, err
}

// assertPortCode verifies that err is a structured port error with code.
func assertPortCode(t *testing.T, err error, code porterr.Code) {
	t.Helper()
	var perr *porterr.Error
	if !errors.As(err, &perr) {
		t.Fatalf("expected port error %s, got %T %v", code, err, err)
	}
	if perr.Code != code {
		t.Fatalf("expected code %s, got %s", code, perr.Code)
	}
}

// assertStringSlice compares slices without pulling in a test dependency.
func assertStringSlice(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("item %d = %q, want %q in %#v", i, got[i], want[i], got)
		}
	}
}

// assertContains verifies that values contains want.
func assertContains(t *testing.T, values []string, want string) {
	t.Helper()
	for _, value := range values {
		if value == want {
			return
		}
	}
	t.Fatalf("%q not found in %#v", want, values)
}
