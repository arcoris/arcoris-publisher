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

package github

import (
	"errors"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

// remoteError creates the structured remote-port error used by this adapter.
func remoteError(code porterr.Code, message string, cause error, details porterr.Details) error {
	return porterr.New(porterr.KindRemote, code, message, cause).WithDetails(details)
}

// remoteCode extracts a remote-port error code from a wrapped error.
func remoteCode(err error) (porterr.Code, bool) {
	var perr *porterr.Error
	if !errors.As(err, &perr) || perr.Kind != porterr.KindRemote {
		return "", false
	}
	return perr.Code, true
}

// isRemoteCode reports whether err already carries code.
func isRemoteCode(err error, code porterr.Code) bool {
	got, ok := remoteCode(err)
	return ok && got == code
}

// wrapRemoteOperationError preserves transport/classification errors when they
// are more specific than an endpoint-level fallback.
//
// For example, CreateRelease should return repository_not_found for a missing
// repo, not the less helpful release_failed fallback.
func wrapRemoteOperationError(fallback porterr.Code, message string, cause error, details porterr.Details) error {
	if code, ok := remoteCode(cause); ok {
		switch code {
		case remoteport.CodeRepositoryNotFound,
			remoteport.CodeAuthenticationFailed,
			remoteport.CodeAccessDenied,
			remoteport.CodeRateLimited,
			remoteport.CodeBranchProtected:
			return cause
		}
	}
	return remoteError(fallback, message, cause, details)
}
