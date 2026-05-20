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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

// do sends one GitHub API request and decodes a JSON response.
//
// The helper centralizes headers, authentication, body encoding, response size
// limiting, and HTTP status classification so individual endpoint methods stay
// focused on request/response mapping.
func (p *Provider) do(ctx context.Context, method string, path string, body any, out any) error {
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return remoteError(remoteport.CodeReleaseFailed, "github request encode failed", err, nil)
		}
		r = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, p.baseURL+path, r)
	if err != nil {
		return remoteError(remoteport.CodeReleaseFailed, "github request creation failed", err, nil)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return remoteError(remoteport.CodeReleaseFailed, "github request failed", err, nil)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return remoteError(remoteport.CodeReleaseFailed, "github response read failed", err, porterr.Details{"path": path})
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return classifyHTTP(resp.StatusCode, resp.Header, string(data), path)
	}
	if out != nil && len(data) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return remoteError(remoteport.CodeReleaseFailed, "github response decode failed", err, porterr.Details{"path": path})
		}
	}
	return nil
}

// classifyHTTP maps GitHub HTTP status codes into remote-port error codes.
//
// GitHub sometimes reports rate limiting as 403 with X-RateLimit-Remaining: 0,
// so that header takes precedence over the generic access-denied mapping.
func classifyHTTP(status int, header http.Header, body string, path string) error {
	code := remoteport.CodeReleaseFailed
	msg := fmt.Sprintf("github api request failed with status %d", status)
	switch status {
	case http.StatusNotFound:
		code = remoteport.CodeRepositoryNotFound
	case http.StatusUnauthorized:
		code = remoteport.CodeAuthenticationFailed
	case http.StatusForbidden:
		if header.Get("X-RateLimit-Remaining") == "0" {
			code = remoteport.CodeRateLimited
		} else {
			code = remoteport.CodeAccessDenied
		}
	case http.StatusTooManyRequests:
		code = remoteport.CodeRateLimited
	}
	return remoteError(code, msg, nil, porterr.Details{"path": path, "status": fmt.Sprint(status), "body": body})
}
