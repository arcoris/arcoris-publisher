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

package config

// LoaderOptions configures a manifest loader.
type LoaderOptions struct {
	Reader  Reader
	Decoder Decoder
	Locator Locator
}

// Loader reads top-level and module-level ARCORIS Publisher manifests and
// resolves them into an effective publication set.
type Loader struct {
	reader  Reader
	decoder Decoder
	locator Locator
}

// NewLoader creates a Loader with safe defaults for omitted dependencies.
func NewLoader(opts LoaderOptions) Loader {
	return Loader{
		reader:  defaultReader(opts.Reader),
		decoder: defaultDecoder(opts.Decoder),
		locator: defaultLocator(opts.Locator),
	}
}

// defaultReader chooses the filesystem-backed reader when tests or callers do
// not inject a virtual reader.
func defaultReader(reader Reader) Reader {
	if reader != nil {
		return reader
	}
	return OSReader{}
}

// defaultDecoder chooses the strict manifest decoder when callers do not need a
// custom syntax layer.
func defaultDecoder(decoder Decoder) Decoder {
	if decoder != nil {
		return decoder
	}
	return StrictDecoder{}
}

// defaultLocator keeps discovery deterministic by restoring conventional names
// when the caller supplies an empty locator.
func defaultLocator(locator Locator) Locator {
	if len(locator.Names) != 0 {
		return locator
	}
	return DefaultLocator()
}
