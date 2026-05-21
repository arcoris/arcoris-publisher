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

import (
	"errors"
	"strings"
	"testing"
)

func TestYAMLErrorsIncludeLineNumber(t *testing.T) {
	line := yamlLine{number: 7, indent: 4}
	errs := []error{
		expectedIndentationError(line, 2),
		unexpectedIndentationError(line),
		expectedListContinuationIndentationError(line, 2),
		duplicateYAMLKeyError(line, "name"),
		lineError(line, errors.New("bad scalar")),
	}

	for _, err := range errs {
		if !strings.Contains(err.Error(), "line 7") {
			t.Fatalf("expected line number in %q", err)
		}
	}
}
