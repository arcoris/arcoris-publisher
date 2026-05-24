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

package runtime

import (
	"context"
	"strings"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

type recordingRunner struct {
	specs []processport.Spec
}

func (r *recordingRunner) Run(_ context.Context, spec processport.Spec) (processport.Result, error) {
	r.specs = append(r.specs, cloneSpec(spec))
	return processport.Result{
		Name:     spec.Name,
		Args:     append([]string(nil), spec.Args...),
		Dir:      spec.Dir,
		ExitCode: 0,
		Stdout:   stdoutForSpec(spec),
	}, nil
}

func cloneSpec(spec processport.Spec) processport.Spec {
	spec.Args = append([]string(nil), spec.Args...)
	spec.Env = append([]string(nil), spec.Env...)
	spec.AllowedExitCodes = append([]int(nil), spec.AllowedExitCodes...)
	spec.SensitiveValues = append([]string(nil), spec.SensitiveValues...)
	return spec
}

func stdoutForSpec(spec processport.Spec) []byte {
	switch {
	case spec.Name == "test-git" && strings.Join(spec.Args, " ") == "rev-parse HEAD":
		return []byte("abcdef1234567890\n")
	case spec.Name == "test-go" && strings.Join(spec.Args, " ") == "env -json":
		return []byte("{}\n")
	default:
		return nil
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
