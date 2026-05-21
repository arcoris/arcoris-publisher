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

package gitcli

import (
	"context"
	"errors"
	"os/exec"
	"strconv"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

type fakeRunner struct {
	specs   []processport.Spec
	results []processport.Result
	errs    []error
}

// Run records every Git process spec and returns the configured result/error.
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

// runGit executes a real git command for integration-style repository tests.
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

// must fails the current test when err is non-nil.
func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
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

// intsOf renders integer slices for reuse with assertStringSlice.
func intsOf(values []int) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = strconv.Itoa(value)
	}
	return out
}
