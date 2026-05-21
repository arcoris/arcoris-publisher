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

package manifest

// stringValue unwraps an optional string field or returns the supplied fallback.
//
// Manifest specs use pointers for optional scalar fields so an omitted value can
// be distinguished from an explicitly provided empty string when a field needs
// that distinction.
func stringValue(value *string, fallback string) string {
	if value == nil {
		return fallback
	}
	return *value
}

// boolValue unwraps an optional boolean field or returns the supplied fallback.
func boolValue(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}
