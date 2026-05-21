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

import "fmt"

// expectedIndentationError reports a line that dedented before the current block
// could start.
func expectedIndentationError(line yamlLine, indent int) error {
	return fmt.Errorf(
		"line %d: expected indentation %d, got %d",
		line.number,
		indent,
		line.indent,
	)
}

// unexpectedIndentationError reports indentation that is deeper than the parser
// can accept at the current block boundary.
func unexpectedIndentationError(line yamlLine) error {
	return fmt.Errorf(
		"line %d: unexpected indentation %d",
		line.number,
		line.indent,
	)
}

// expectedListContinuationIndentationError reports malformed continuation lines
// for inline map items inside lists.
func expectedListContinuationIndentationError(
	line yamlLine,
	indent int,
) error {
	return fmt.Errorf(
		"line %d: expected list item continuation indentation %d",
		line.number,
		indent,
	)
}

// duplicateYAMLKeyError reports repeated keys within the same map object.
func duplicateYAMLKeyError(line yamlLine, key string) error {
	return fmt.Errorf("line %d: duplicate key %q", line.number, key)
}

// lineError annotates scalar parse errors with the original YAML line number.
func lineError(line yamlLine, err error) error {
	return fmt.Errorf("line %d: %w", line.number, err)
}
