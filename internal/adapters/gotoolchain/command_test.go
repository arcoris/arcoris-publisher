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
	"time"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

func TestCommandBuildsProcessSpec(t *testing.T) {
	tool := New(&fakeRunner{}, Options{GoBinary: "go-default", Env: []string{"GOPROXY=base", "KEEP=1"}})
	spec := tool.command("/repo", []string{"list"}, goport.CommonOptions{
		GoBinary:        "go-custom",
		WorkspaceMode:   goport.WorkspaceOff,
		Env:             []string{"EXTRA=1"},
		Timeout:         time.Second,
		PrivateModules:  []string{"example.com/private", "example.com/other"},
		Proxy:           "direct",
		SumDB:           "off",
		SensitiveValues: []string{"token"},
	})

	if spec.Name != "go-custom" || spec.Dir != "/repo" || spec.Timeout != time.Second {
		t.Fatalf("unexpected spec identity: %#v", spec)
	}
	assertStringSlice(t, spec.Args, []string{"list"})
	assertContains(t, spec.Env, "KEEP=1")
	assertContains(t, spec.Env, "EXTRA=1")
	assertContains(t, spec.Env, "GOWORK=off")
	assertContains(t, spec.Env, "GOPROXY=direct")
	assertContains(t, spec.Env, "GOSUMDB=off")
	assertContains(t, spec.Env, "GOPRIVATE=example.com/private,example.com/other")
	assertStringSlice(t, spec.SensitiveValues, []string{"token"})
	if !spec.CaptureStdout || !spec.CaptureStderr {
		t.Fatalf("expected stdout/stderr capture")
	}
}
