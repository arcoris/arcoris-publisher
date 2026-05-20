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

func TestListBuildsCommandAndParsesPackages(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{
		Stdout: []byte(`{"ImportPath":"example.com/m","Module":{"Path":"example.com/m"}}
`),
	}}}
	tool := New(runner, Options{})

	result, err := tool.List(context.Background(), "/repo", goport.ListOptions{
		CommonOptions: goport.CommonOptions{WorkspaceMode: goport.WorkspaceOff, Tags: []string{"integration", "linux"}},
		JSON:          true,
		Deps:          true,
		Test:          true,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(result.Packages) != 1 || result.Packages[0].ImportPath != "example.com/m" {
		t.Fatalf("unexpected packages %#v", result.Packages)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"list", "-json", "-deps", "-test", "-tags", "integration,linux", "./..."})
	assertContains(t, runner.specs[0].Env, "GOWORK=off")
}

func TestListMalformedJSONReturnsError(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{Stdout: []byte(`{`)}}}
	tool := New(runner, Options{})

	result, err := tool.List(context.Background(), "/repo", goport.ListOptions{JSON: true})
	if len(result.Packages) != 0 {
		t.Fatalf("Packages = %#v, want empty on parse failure", result.Packages)
	}
	assertPortCode(t, err, goport.CodeListFailed)
}
