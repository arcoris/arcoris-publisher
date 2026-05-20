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
	"fmt"
	"net/http"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

// classifyHTTP maps GitHub HTTP status codes into remote-port error codes.
//
// GitHub sometimes reports rate limiting as 403 with X-RateLimit-Remaining: 0,
// so that header takes precedence over the generic access-denied mapping.
func classifyHTTP(status int, header http.Header, body string, path string) error {
	code := httpStatusCode(status, header)
	msg := fmt.Sprintf("github api request failed with status %d", status)
	return remoteError(code, msg, nil, porterr.Details{"path": path, "status": fmt.Sprint(status), "body": body})
}

// httpStatusCode translates provider HTTP status into the closest port code.
func httpStatusCode(status int, header http.Header) porterr.Code {
	switch status {
	case http.StatusNotFound:
		return remoteport.CodeRepositoryNotFound
	case http.StatusUnauthorized:
		return remoteport.CodeAuthenticationFailed
	case http.StatusForbidden:
		if header.Get("X-RateLimit-Remaining") == "0" {
			return remoteport.CodeRateLimited
		}
		return remoteport.CodeAccessDenied
	case http.StatusTooManyRequests:
		return remoteport.CodeRateLimited
	default:
		return remoteport.CodeReleaseFailed
	}
}

// requestFailedError maps lower-level client failures into a remote error.
func requestFailedError(cause error) error {
	return remoteError(remoteport.CodeReleaseFailed, "github request failed", cause, nil)
}

// responseReadError maps response body read failures into a remote error.
func responseReadError(path string, cause error) error {
	return remoteError(remoteport.CodeReleaseFailed, "github response read failed", cause, porterr.Details{"path": path})
}
