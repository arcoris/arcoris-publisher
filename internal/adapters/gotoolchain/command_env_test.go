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
	"testing"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

func TestCommandEnvOverlaysDerivedOptions(t *testing.T) {
	tool := New(&fakeRunner{}, Options{Env: []string{"GOPROXY=base", "KEEP=1"}})
	env := tool.commandEnv(goport.CommonOptions{
		WorkspaceMode:  goport.WorkspaceOff,
		Env:            []string{"EXTRA=1"},
		PrivateModules: []string{"example.com/private"},
		Proxy:          "direct",
		SumDB:          "off",
	})

	assertContains(t, env, "KEEP=1")
	assertContains(t, env, "EXTRA=1")
	assertContains(t, env, "GOWORK=off")
	assertContains(t, env, "GOPROXY=direct")
	assertContains(t, env, "GOSUMDB=off")
	assertContains(t, env, "GOPRIVATE=example.com/private")
}

func TestSetEnvReplacesExistingAssignment(t *testing.T) {
	env := setEnv([]string{"A=1", "B=2"}, "A", "3")
	assertStringSlice(t, env, []string{"A=3", "B=2"})
}
