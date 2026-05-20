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
	"strings"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// commandEnv builds the environment for one Go command invocation.
//
// Environment variables are applied in adapter-default, call-specific, and
// derived-option order. Derived options such as WorkspaceMode and GOPROXY win so
// the typed option fields have predictable precedence.
func (t *Toolchain) commandEnv(common goport.CommonOptions) []string {
	env := append([]string(nil), t.env...)
	env = append(env, common.Env...)
	env = applyWorkspaceMode(env, common.WorkspaceMode)
	env = applyProxyEnv(env, common)
	return env
}

// applyWorkspaceMode turns the typed workspace option into GOWORK.
func applyWorkspaceMode(env []string, mode goport.WorkspaceMode) []string {
	if mode == goport.WorkspaceOff {
		return setEnv(env, "GOWORK", "off")
	}
	return env
}

// applyProxyEnv overlays module download related variables.
func applyProxyEnv(env []string, common goport.CommonOptions) []string {
	if common.Proxy != "" {
		env = setEnv(env, "GOPROXY", common.Proxy)
	}
	if common.SumDB != "" {
		env = setEnv(env, "GOSUMDB", common.SumDB)
	}
	if len(common.PrivateModules) > 0 {
		env = setEnv(env, "GOPRIVATE", strings.Join(common.PrivateModules, ","))
	}
	return env
}

// setEnv overlays one KEY=VALUE assignment in env.
func setEnv(env []string, key, value string) []string {
	entry := key + "=" + value
	for i, existing := range env {
		if strings.HasPrefix(existing, key+"=") {
			env[i] = entry
			return env
		}
	}
	return append(env, entry)
}
