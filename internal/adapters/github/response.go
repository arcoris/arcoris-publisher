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
	"encoding/json"
	"io"
	"net/http"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

// maxResponseBodyBytes bounds GitHub response bodies kept in memory.
//
// API errors can include verbose validation payloads. The limit preserves useful
// diagnostics without allowing an endpoint to consume unbounded memory.
const maxResponseBodyBytes int64 = 4 << 20

// decodeAPIResponse reads, classifies, and decodes one GitHub response.
func (p *Provider) decodeAPIResponse(resp *http.Response, path string, out any) error {
	defer resp.Body.Close()
	data, err := readResponseBody(resp)
	if err != nil {
		return responseReadError(path, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return classifyHTTP(resp.StatusCode, resp.Header, string(data), path)
	}
	return decodeJSONResponse(path, data, out)
}

// readResponseBody reads a bounded snapshot of the HTTP response body.
func readResponseBody(resp *http.Response) ([]byte, error) {
	return io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
}

// decodeJSONResponse unmarshals successful JSON responses when a target is supplied.
func decodeJSONResponse(path string, data []byte, out any) error {
	if out == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return remoteError(remoteport.CodeReleaseFailed, "github response decode failed", err, porterr.Details{"path": path})
	}
	return nil
}
