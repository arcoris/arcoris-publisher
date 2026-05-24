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

package provenance

import (
	"strings"
	"testing"
)

func TestBuildTrailersReturnsStableRequiredKeys(t *testing.T) {
	input := testInput(t)

	trailers := BuildTrailers(input)

	expected := []string{
		"Arcoris-Source-Repository",
		"Arcoris-Source-Commit",
		"Arcoris-Source-Branch",
		"Arcoris-Module",
		"Arcoris-Module-Path",
		"Arcoris-Version",
		"Arcoris-Target-Repository",
		"Arcoris-Target-Branches",
		"Arcoris-Publish-Mode",
		"Arcoris-Push-Policy",
		"Arcoris-Tag-Policy",
		"Arcoris-Publisher-Version",
		"Arcoris-Source-Dir",
		"Arcoris-Source-Hash",
		"Arcoris-Projection-Hash",
	}
	if len(trailers) != len(expected) {
		t.Fatalf("len(trailers) = %d", len(trailers))
	}
	for i, key := range expected {
		if trailers[i].Key != key {
			t.Fatalf("trailers[%d].Key = %q, want %q", i, trailers[i].Key, key)
		}
	}
}

func TestTrailersRenderDeterministicallyWithoutLocalPaths(t *testing.T) {
	input := testInput(t)

	first := BuildTrailers(input).Render()
	second := BuildTrailers(input).Render()

	if first != second {
		t.Fatalf("trailers are not deterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	for _, required := range []string{
		"Arcoris-Module: foundation",
		"Arcoris-Module-Path: arcoris.dev/foundation",
		"Arcoris-Publisher-Version: v9.8.7",
		"Arcoris-Source-Dir: src/arcoris.dev/foundation",
		"Arcoris-Projection-Hash: sha256:",
	} {
		if !strings.Contains(first, required) {
			t.Fatalf("trailers missing %q:\n%s", required, first)
		}
	}
	for _, localPath := range []string{"/repo", "/target", "/tmp"} {
		if strings.Contains(first, localPath) {
			t.Fatalf("trailers leak local path %q:\n%s", localPath, first)
		}
	}
}

func TestTrailersRenderSanitizesNewlines(t *testing.T) {
	trailers := Trailers{
		{Key: " Arcoris-Module\r\nInjected ", Value: " foundation\ncontrol\r\nnext "},
	}

	rendered := trailers.Render()

	body := strings.TrimSuffix(rendered, "\n")
	if strings.ContainsAny(body, "\r\n") {
		t.Fatalf("trailer contains raw newline inside key/value:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Arcoris-Module Injected: foundation control next") {
		t.Fatalf("trailer key/value not sanitized:\n%s", rendered)
	}
}

func TestTrailerSanitizationPreservesKnownKeys(t *testing.T) {
	key := "Arcoris-Publisher-Version"

	if sanitizeTrailerKey(key) != key {
		t.Fatalf("sanitizeTrailerKey(%q) = %q", key, sanitizeTrailerKey(key))
	}
}
