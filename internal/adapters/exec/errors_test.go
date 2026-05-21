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

package exec

import (
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestProcessErrorCarriesKindCodeCauseAndDetails(t *testing.T) {
	cause := errors.New("failed")
	err := processError(processport.CodeFailed, "process failed", cause, porterr.Details{"name": "go"})

	var perr *porterr.Error
	if !errors.As(err, &perr) {
		t.Fatalf("processError() type = %T", err)
	}
	if perr.Kind != porterr.KindProcess || perr.Code != processport.CodeFailed || perr.Cause != cause {
		t.Fatalf("processError() = %#v", perr)
	}
	if perr.Details["name"] != "go" {
		t.Fatalf("processError() details = %#v", perr.Details)
	}
}

func TestCommandDetailsOmitsEmptyValues(t *testing.T) {
	details := commandDetails("go", []string{"test", "./..."}, "", 2)

	if details["name"] != "go" || details["args"] != "test ./..." || details["exit_code"] != "2" {
		t.Fatalf("commandDetails() = %#v", details)
	}
	if _, ok := details["dir"]; ok {
		t.Fatalf("commandDetails() should omit empty dir: %#v", details)
	}
}
