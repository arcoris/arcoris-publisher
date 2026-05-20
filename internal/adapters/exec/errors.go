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
	"strconv"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

// processError creates the structured process-port error used by this adapter.
func processError(code porterr.Code, message string, cause error, details porterr.Details) error {
	return porterr.New(porterr.KindProcess, code, message, cause).WithDetails(details)
}

// commandDetails returns non-secret command context suitable for diagnostics.
//
// Callers must pass already-redacted values. Keeping this helper small and
// explicit prevents accidental inclusion of stdout/stderr or environment values,
// which often carry credentials.
func commandDetails(name string, args []string, dir string, exitCode int) porterr.Details {
	details := porterr.Details{}
	if name != "" {
		details["name"] = name
	}
	if dir != "" {
		details["dir"] = dir
	}
	if len(args) > 0 {
		details["args"] = strings.Join(args, " ")
	}
	if exitCode != 0 {
		details["exit_code"] = strconv.Itoa(exitCode)
	}
	return details
}
