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
	"io"
	"net/http"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

// newAPIRequest builds one authenticated GitHub REST request.
func (p *Provider) newAPIRequest(ctx context.Context, method string, path string, body any) (*http.Request, error) {
	reader, err := encodeRequestBody(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, p.baseURL+path, reader)
	if err != nil {
		return nil, remoteError(remoteport.CodeReleaseFailed, "github request creation failed", err, nil)
	}
	p.attachRequestHeaders(req, body != nil)
	return req, nil
}

// encodeRequestBody marshals optional JSON payloads for endpoint methods.
func encodeRequestBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, remoteError(remoteport.CodeReleaseFailed, "github request encode failed", err, nil)
	}
	return bytes.NewReader(data), nil
}

// attachRequestHeaders applies GitHub API versioning and optional auth.
func (p *Provider) attachRequestHeaders(req *http.Request, hasBody bool) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}
	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}
}
