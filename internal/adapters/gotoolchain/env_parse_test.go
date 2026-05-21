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

func TestParseEnvResult(t *testing.T) {
	result, err := parseEnvResult([]byte(`{"GOWORK":"off","GOFLAGS":""}`))
	if err != nil {
		t.Fatalf("parseEnvResult() error = %v", err)
	}
	if result.Value("GOWORK") != "off" || !result.HasValue("GOFLAGS") {
		t.Fatalf("parseEnvResult() = %#v", result)
	}
}

func TestParseEnvResultRejectsMalformedJSON(t *testing.T) {
	_, err := parseEnvResult([]byte(`{`))
	assertPortCode(t, err, goport.CodeCommandFailed)
}
