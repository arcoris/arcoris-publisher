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

import "arcoris.dev/arcoris-publisher/internal/ports/porterr"

const (
	// CodeCommandFailed identifies a generic Go command failure.
	CodeCommandFailed porterr.Code = "go_command_failed"
	// CodeBinaryNotFound identifies a missing Go executable.
	CodeBinaryNotFound porterr.Code = "go_binary_not_found"
	// CodeListFailed identifies a go list failure.
	CodeListFailed porterr.Code = "go_list_failed"
	// CodeTestFailed identifies a go test failure.
	CodeTestFailed porterr.Code = "go_test_failed"
	// CodeModTidyFailed identifies a go mod tidy failure.
	CodeModTidyFailed porterr.Code = "go_mod_tidy_failed"
	// CodeWorkspaceLeak identifies workspace state leaking into module mode.
	CodeWorkspaceLeak porterr.Code = "go_workspace_leak"
	// CodeVersionUnsupported identifies an unsupported Go version.
	CodeVersionUnsupported porterr.Code = "go_version_unsupported"
)
