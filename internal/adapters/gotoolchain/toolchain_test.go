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

import "testing"

func TestNewAppliesDefaultBinaryAndDetachesEnvironment(t *testing.T) {
	env := []string{"A=1"}
	tool := New(&fakeRunner{}, Options{Env: env})
	env[0] = "A=mutated"

	if tool.goBin != defaultGoBinary {
		t.Fatalf("goBin = %q, want %q", tool.goBin, defaultGoBinary)
	}
	if tool.env[0] != "A=1" {
		t.Fatalf("env = %#v, want detached copy", tool.env)
	}
}

func TestNewUsesConfiguredBinary(t *testing.T) {
	tool := New(&fakeRunner{}, Options{GoBinary: "custom-go"})
	if tool.goBin != "custom-go" {
		t.Fatalf("goBin = %q, want custom-go", tool.goBin)
	}
}
