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

package process

import "arcoris.dev/arcoris-publisher/internal/ports/porterr"

const (
	// CodeFailed identifies a process that completed with an unaccepted exit code.
	CodeFailed porterr.Code = "process_failed"
	// CodeNotFound identifies a missing executable.
	CodeNotFound porterr.Code = "process_not_found"
	// CodeTimedOut identifies a process timeout.
	CodeTimedOut porterr.Code = "process_timed_out"
	// CodeCancelled identifies a process cancelled through context cancellation.
	CodeCancelled porterr.Code = "process_cancelled"
	// CodeStartFailed identifies a process that could not be started.
	CodeStartFailed porterr.Code = "process_start_failed"
)
