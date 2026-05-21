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
	"strings"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// goError creates the structured Go-toolchain error used by this adapter.
func goError(code porterr.Code, message string, cause error, details porterr.Details) error {
	return porterr.New(porterr.KindGo, code, message, cause).WithDetails(details)
}

// wrapGoError maps process and stderr failures into Go-toolchain error codes.
//
// The process adapter reports missing executables as a process error. This
// wrapper converts that boundary-specific failure into go_binary_not_found so
// callers do not need to understand how the Go adapter is implemented.
func wrapGoError(code porterr.Code, message string, result processport.Result, cause error) error {
	stderr := strings.TrimSpace(string(result.Stderr))

	if isProcessNotFoundError(cause) || mentionsMissingExecutable(stderr) {
		code = goport.CodeBinaryNotFound
	}

	return goError(code, message, cause, porterr.Details{"dir": result.Dir, "stderr": stderr})
}

func isProcessNotFoundError(cause error) bool {
	var perr *porterr.Error
	return errors.As(cause, &perr) && perr.Kind == porterr.KindProcess && perr.Code == processport.CodeNotFound
}

func mentionsMissingExecutable(stderr string) bool {
	return strings.Contains(strings.ToLower(stderr), "executable file not found")
}
