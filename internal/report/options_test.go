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

package report

import "testing"

func TestDefaultOptions(t *testing.T) {
	t.Parallel()

	opts := DefaultOptions()
	if opts.Format != FormatText || !opts.Pretty || opts.IncludeLocalPaths {
		t.Fatalf("DefaultOptions() = %+v", opts)
	}
}

func TestNormalizeOptionsDefaultsZeroValue(t *testing.T) {
	t.Parallel()

	opts := normalizeOptions(Options{})
	if opts.Format != FormatText || !opts.Pretty {
		t.Fatalf("normalizeOptions(zero) = %+v", opts)
	}
}

func TestNewUsesDefaultOptionsForZeroValue(t *testing.T) {
	t.Parallel()

	renderer := New(Options{})
	if renderer.opts != DefaultOptions() {
		t.Fatalf("New(Options{}) opts = %+v", renderer.opts)
	}
}

func TestNormalizeOptionsKeepsExplicitCompactJSON(t *testing.T) {
	t.Parallel()

	opts := normalizeOptions(Options{Format: FormatJSON})
	if opts.Format != FormatJSON || opts.Pretty {
		t.Fatalf("normalizeOptions(json) = %+v", opts)
	}
}
