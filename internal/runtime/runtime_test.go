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

import "testing"

func TestNewWithDependenciesFillsMissingDependencies(t *testing.T) {
	t.Parallel()

	provided := &recordingRunner{}
	rt := NewWithDependencies(Dependencies{Process: provided}, Options{})
	if rt.Dependencies().Process != provided {
		t.Fatal("NewWithDependencies replaced provided process runner")
	}
	if rt.Dependencies().Git == nil || rt.Dependencies().FileSystem == nil || rt.Dependencies().Go == nil {
		t.Fatalf("NewWithDependencies did not fill missing dependencies: %+v", rt.Dependencies())
	}
}

func TestOptionsReturnsDetachedEnvironment(t *testing.T) {
	t.Parallel()

	env := []string{"A=B"}
	rt := New(Options{Env: env})
	env[0] = "A=mutated"

	opts := rt.Options()
	if opts.Env[0] != "A=B" {
		t.Fatalf("Options().Env = %v", opts.Env)
	}
	opts.Env[0] = "A=changed"

	if got := rt.Options().Env[0]; got != "A=B" {
		t.Fatalf("Options() returned aliased Env, got %q", got)
	}
}
