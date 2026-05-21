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
	"net/http"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

func TestClassifyHTTP(t *testing.T) {
	tests := []struct {
		name   string
		status int
		header http.Header
		want   string
	}{
		{name: "not found", status: http.StatusNotFound, want: remoteport.CodeRepositoryNotFound.String()},
		{name: "unauthorized", status: http.StatusUnauthorized, want: remoteport.CodeAuthenticationFailed.String()},
		{name: "forbidden", status: http.StatusForbidden, want: remoteport.CodeAccessDenied.String()},
		{name: "rate limit 403", status: http.StatusForbidden, header: http.Header{"X-Ratelimit-Remaining": []string{"0"}}, want: remoteport.CodeRateLimited.String()},
		{name: "rate limit 429", status: http.StatusTooManyRequests, want: remoteport.CodeRateLimited.String()},
		{name: "generic", status: http.StatusUnprocessableEntity, want: remoteport.CodeReleaseFailed.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := classifyHTTP(tt.status, tt.header, "{}", "/path")
			assertPortCode(t, err, porterr.Code(tt.want))
		})
	}
}

func TestHTTPTransportErrorHelpers(t *testing.T) {
	cause := errors.New("network")
	assertPortCode(t, requestFailedError(cause), remoteport.CodeReleaseFailed)
	assertPortCode(t, responseReadError("/path", cause), remoteport.CodeReleaseFailed)
}
