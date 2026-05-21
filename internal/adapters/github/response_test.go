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
	"io"
	"net/http"
	"strings"
	"testing"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

func TestDecodeAPIResponseHandlesEmptySuccess(t *testing.T) {
	resp := &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}
	if err := New(Options{}).decodeAPIResponse(resp, "/path", nil); err != nil {
		t.Fatalf("decodeAPIResponse() error = %v", err)
	}
}

func TestDecodeJSONResponseMapsMalformedJSON(t *testing.T) {
	var out struct {
		OK bool `json:"ok"`
	}
	err := decodeJSONResponse("/path", []byte("{"), &out)
	assertPortCode(t, err, remoteport.CodeReleaseFailed)
}

func TestReadResponseBodyUsesBoundedReader(t *testing.T) {
	resp := &http.Response{Body: io.NopCloser(strings.NewReader("body"))}
	data, err := readResponseBody(resp)
	if err != nil || string(data) != "body" {
		t.Fatalf("readResponseBody() = %q, %v", data, err)
	}
}
