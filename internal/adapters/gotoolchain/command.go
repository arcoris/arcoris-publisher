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
	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// command builds the process spec shared by go list, go test, and go mod tidy.
//
// Environment variables are applied in adapter-default, call-specific, and
// derived-option order. Derived options such as WorkspaceMode and GOPROXY win so
// the typed option fields have predictable precedence.
func (t *Toolchain) command(moduleDir string, args []string, common goport.CommonOptions) processport.Spec {
	return processport.Spec{
		Name:            binary(t.goBin, common.GoBinary),
		Args:            append([]string(nil), args...),
		Dir:             moduleDir,
		Env:             t.commandEnv(common),
		Timeout:         common.Timeout,
		CaptureStdout:   true,
		CaptureStderr:   true,
		SensitiveValues: append([]string(nil), common.SensitiveValues...),
	}
}
