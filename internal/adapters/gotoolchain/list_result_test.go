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

func TestParseListResultKeepsRawOutputWhenJSONDisabled(t *testing.T) {
	result, err := parseListResult([]byte("stdout"), []byte("stderr"), goport.ListOptions{})
	if err != nil {
		t.Fatalf("parseListResult() error = %v", err)
	}
	if string(result.Stdout) != "stdout" || string(result.Stderr) != "stderr" || len(result.Packages) != 0 {
		t.Fatalf("parseListResult() = %#v", result)
	}
}

func TestParseListResultParsesJSONPackages(t *testing.T) {
	result, err := parseListResult([]byte(`{"ImportPath":"example.com/m"}`+"\n"), nil, goport.ListOptions{JSON: true})
	if err != nil {
		t.Fatalf("parseListResult() error = %v", err)
	}
	if len(result.Packages) != 1 || result.Packages[0].ImportPath != "example.com/m" {
		t.Fatalf("Packages = %#v", result.Packages)
	}
}
