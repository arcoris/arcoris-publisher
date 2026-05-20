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
	"context"
)

// do sends one GitHub API request and decodes a JSON response.
//
// The helper centralizes headers, authentication, body encoding, response size
// limiting, and HTTP status classification so individual endpoint methods stay
// focused on request/response mapping.
func (p *Provider) do(ctx context.Context, method string, path string, body any, out any) error {
	req, err := p.newAPIRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return requestFailedError(err)
	}
	return p.decodeAPIResponse(resp, path, out)
}
