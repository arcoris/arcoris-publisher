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
	"net/http"
	"testing"
)

func TestDoAddsHeadersAndDecodesResponse(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Fatalf("Authorization = %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Accept") != "application/vnd.github+json" {
			t.Fatalf("Accept = %q", r.Header.Get("Accept"))
		}
		if r.Header.Get("X-GitHub-Api-Version") == "" {
			t.Fatalf("missing GitHub API version")
		}
		return jsonResponse(200, `{"ok":true}`)
	})
	provider.token = "secret"

	var out struct {
		OK bool `json:"ok"`
	}
	if err := provider.do(context.Background(), "GET", "/ok", nil, &out); err != nil {
		t.Fatalf("do() error = %v", err)
	}
	if !out.OK {
		t.Fatalf("decoded response = %#v", out)
	}
}
