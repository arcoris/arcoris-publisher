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
	"io"
	"net/http"
	"strings"
	"testing"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

func TestNewAPIRequestEncodesBodyAndHeaders(t *testing.T) {
	provider := New(Options{BaseURL: "https://api.example.test", Token: "token"})

	req, err := provider.newAPIRequest(context.Background(), http.MethodPost, "/repos", map[string]string{"name": "repo"})
	if err != nil {
		t.Fatalf("newAPIRequest() error = %v", err)
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if req.URL.String() != "https://api.example.test/repos" || !strings.Contains(string(data), `"name":"repo"`) {
		t.Fatalf("newAPIRequest() url/body = %s %s", req.URL, data)
	}
	if req.Header.Get("Authorization") != "Bearer token" || req.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("headers = %#v", req.Header)
	}
}

func TestEncodeRequestBodyMapsMarshalFailure(t *testing.T) {
	_, err := encodeRequestBody(map[string]any{"bad": func() {}})
	assertPortCode(t, err, remoteport.CodeReleaseFailed)
}
