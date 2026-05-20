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
	"testing"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestEnvParsesGoEnvJSON(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{Stdout: []byte(`{"GOWORK":"off"}`)}}}
	tool := New(runner, Options{GoBinary: "custom-go"})

	env, err := tool.Env(context.Background(), goport.EnvOptions{})
	if err != nil {
		t.Fatalf("Env() error = %v", err)
	}
	if env.Value("GOWORK") != "off" {
		t.Fatalf("GOWORK = %q, want off", env.Value("GOWORK"))
	}
	if runner.specs[0].Name != "custom-go" {
		t.Fatalf("spec.Name = %q, want custom-go", runner.specs[0].Name)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"env", "-json"})
}

func TestEnvRejectsMalformedJSON(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{Stdout: []byte(`{`)}}}
	tool := New(runner, Options{})

	_, err := tool.Env(context.Background(), goport.EnvOptions{})
	assertPortCode(t, err, goport.CodeCommandFailed)
}
