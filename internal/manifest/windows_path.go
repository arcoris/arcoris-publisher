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

// looksWindowsAbsolute reports whether value starts like a Windows drive path.
//
// Manifest paths are always slash-separated and relative, but the check keeps
// "C:/work" from being accepted on non-Windows hosts where path.Clean would not
// treat it as absolute.
func looksWindowsAbsolute(value string) bool {
	if len(value) < 3 {
		return false
	}

	return isASCIILetter(value[0]) &&
		value[1] == ':' &&
		(value[2] == '/' || value[2] == '\\')
}

// isASCIILetter reports whether b is an English alphabetic byte.
func isASCIILetter(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}
