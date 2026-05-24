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

// Options controls report rendering.
type Options struct {
	// Format selects text or JSON rendering.
	Format Format

	// Pretty controls JSON indentation. Text reports ignore this option.
	Pretty bool

	// IncludeLocalPaths allows reports to include local absolute paths. The
	// default is false because reports are often uploaded to CI artifacts, pasted
	// into issues, or committed as diagnostics.
	IncludeLocalPaths bool
}

// DefaultOptions returns safe report defaults.
func DefaultOptions() Options {
	return Options{Format: FormatText, Pretty: true}
}

func normalizeOptions(opts Options) Options {
	if opts.Format == "" {
		if opts == (Options{}) {
			return DefaultOptions()
		}
		opts.Format = FormatText
	}
	return opts
}
