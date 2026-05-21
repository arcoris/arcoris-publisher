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
	"fmt"
	"strconv"
	"strings"
)

// parseYAMLScalar converts one scalar token into the JSON-compatible value used
// by the strict struct decoder.
func parseYAMLScalar(raw string) (any, error) {
	switch raw {
	case "":
		return "", nil
	case "[]":
		return []any{}, nil
	case "{}":
		return map[string]any{}, nil
	case "null", "~":
		return nil, nil
	case "true", "True", "TRUE":
		return true, nil
	case "false", "False", "FALSE":
		return false, nil
	}

	if strings.HasPrefix(raw, "'") || strings.HasPrefix(raw, "\"") {
		return parseQuotedYAMLScalar(raw)
	}
	if strings.HasPrefix(raw, "[") || strings.HasPrefix(raw, "{") {
		return nil, unsupportedInlineYAMLCollectionError()
	}
	if number, ok := parseYAMLNumber(raw); ok {
		return number, nil
	}

	return raw, nil
}

// parseQuotedYAMLScalar handles both YAML single-quoted strings and JSON-style
// double-quoted strings.
func parseQuotedYAMLScalar(raw string) (string, error) {
	if strings.HasPrefix(raw, "'") {
		return parseSingleQuotedYAMLScalar(raw)
	}

	value, err := strconv.Unquote(raw)
	if err != nil {
		return "", fmt.Errorf("invalid quoted string: %w", err)
	}

	return value, nil
}

// parseSingleQuotedYAMLScalar implements YAML's doubled single-quote escape.
func parseSingleQuotedYAMLScalar(raw string) (string, error) {
	if !strings.HasSuffix(raw, "'") || len(raw) == 1 {
		return "", fmt.Errorf("unterminated single-quoted string")
	}

	return strings.ReplaceAll(raw[1:len(raw)-1], "''", "'"), nil
}

// parseYAMLNumber returns int64 or float64 values for simple numeric literals.
func parseYAMLNumber(raw string) (any, bool) {
	if raw == "" {
		return nil, false
	}
	if !isYAMLNumberLiteral(raw) {
		return nil, false
	}
	if strings.Contains(raw, ".") {
		return parseYAMLFloat(raw)
	}

	return parseYAMLInt(raw)
}

// isYAMLNumberLiteral performs the cheap character-level check before numeric
// parsing is attempted.
func isYAMLNumberLiteral(raw string) bool {
	for _, r := range raw {
		if (r < '0' || r > '9') && r != '-' && r != '+' && r != '.' {
			return false
		}
	}

	return true
}

// parseYAMLFloat parses decimal literals into float64.
func parseYAMLFloat(raw string) (any, bool) {
	value, err := strconv.ParseFloat(raw, 64)
	return value, err == nil
}

// parseYAMLInt parses integer literals into int64.
func parseYAMLInt(raw string) (any, bool) {
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, false
	}

	return value, true
}

// unsupportedInlineYAMLCollectionError explains that only empty inline
// collections are accepted by this small parser.
func unsupportedInlineYAMLCollectionError() error {
	return fmt.Errorf("inline YAML collections other than [] and {} are not supported")
}
