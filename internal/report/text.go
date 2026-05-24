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

package report

import (
	"fmt"
	"io"
	"strings"
)

func writeLine(w io.Writer, format string, args ...any) error {
	_, err := fmt.Fprintf(w, format+"\n", args...)
	if err != nil {
		return &Error{Code: CodeWriteFailed, Message: "text write failed", Cause: err}
	}
	return nil
}

func commaOrDash(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}

func statusPassedFailed(failed bool) string {
	if failed {
		return "failed"
	}
	return "passed"
}

func statusEmptyPresent(present bool) string {
	if present {
		return "present"
	}
	return "empty"
}
