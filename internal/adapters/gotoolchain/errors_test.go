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
	"errors"
	"testing"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestWrapGoErrorMapsProcessNotFound(t *testing.T) {
	cause := porterr.New(porterr.KindProcess, processport.CodeNotFound, "missing", errors.New("missing"))

	err := wrapGoError(goport.CodeCommandFailed, "go failed", processport.Result{}, cause)
	assertPortCode(t, err, goport.CodeBinaryNotFound)
}

func TestWrapGoErrorMapsStderrNotFound(t *testing.T) {
	err := wrapGoError(
		goport.CodeCommandFailed,
		"go failed",
		processport.Result{Stderr: []byte("executable file not found")},
		errors.New("missing"),
	)
	assertPortCode(t, err, goport.CodeBinaryNotFound)
}
