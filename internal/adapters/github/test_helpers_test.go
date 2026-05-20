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
	"io"
	"net/http"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

type testResponse struct {
	status int
	header http.Header
	body   string
}

// roundTripFunc lets tests model GitHub HTTP responses without opening sockets.
type roundTripFunc func(*http.Request) testResponse

// RoundTrip implements http.RoundTripper for roundTripFunc.
func (fn roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	response := fn(r)
	header := response.header
	if header == nil {
		header = http.Header{}
	}
	return &http.Response{
		StatusCode: response.status,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(response.body)),
		Request:    r,
	}, nil
}

// newTestProvider constructs a Provider backed by an in-memory transport.
func newTestProvider(fn roundTripFunc) *Provider {
	return New(Options{
		BaseURL: "https://api.github.test",
		Client:  &http.Client{Transport: fn},
	})
}

// jsonResponse creates an HTTP JSON fixture response.
func jsonResponse(status int, body string) testResponse {
	return testResponse{
		status: status,
		header: http.Header{"Content-Type": []string{"application/json"}},
		body:   body,
	}
}

// readRequestBody reads and returns a request body in tests.
func readRequestBody(t *testing.T, r *http.Request) string {
	t.Helper()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("ReadAll(body) error = %v", err)
	}
	return string(data)
}

// assertPortCode verifies that err is a structured port error with code.
func assertPortCode(t *testing.T, err error, code porterr.Code) {
	t.Helper()
	var perr *porterr.Error
	if !errors.As(err, &perr) {
		t.Fatalf("expected port error %s, got %T %v", code, err, err)
	}
	if perr.Code != code {
		t.Fatalf("expected code %s, got %s", code, perr.Code)
	}
}
