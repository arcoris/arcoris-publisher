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
	"encoding/json"
	"fmt"
	"strings"
)

// yamlToJSON converts the conservative YAML subset used by arcpub manifests to
// JSON so the same strict struct decoder can validate YAML and JSON inputs.
//
// The parser intentionally supports only the simple, deterministic subset needed
// by manifest files: indentation-based maps, lists, booleans, nulls, strings,
// numbers, and []/{} literals. Unsupported YAML features fail closed rather than
// being interpreted ambiguously.
func yamlToJSON(data []byte) ([]byte, error) {
	parser, err := newYAMLParser(string(data))
	if err != nil {
		return nil, err
	}
	value, err := parser.parseDocument()
	if err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

type yamlLine struct {
	number int
	indent int
	text   string
}

type yamlParser struct {
	lines []yamlLine
	pos   int
}

// newYAMLParser normalizes line endings, strips comments, rejects tabs, and
// stores only meaningful YAML lines with their original line numbers.
func newYAMLParser(input string) (*yamlParser, error) {
	rawLines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	lines := make([]yamlLine, 0, len(rawLines))
	for i, raw := range rawLines {
		if strings.Contains(raw, "\t") {
			return nil, fmt.Errorf("line %d: tabs are not supported in arcpub YAML", i+1)
		}
		withoutComment := stripYAMLComment(raw)
		trimmed := strings.TrimSpace(withoutComment)
		if trimmed == "" {
			continue
		}
		indent := len(withoutComment) - len(strings.TrimLeft(withoutComment, " "))
		lines = append(lines, yamlLine{number: i + 1, indent: indent, text: trimmed})
	}
	return &yamlParser{lines: lines}, nil
}

// stripYAMLComment removes comments that begin outside quoted strings while
// preserving # characters embedded in single- or double-quoted scalars.
func stripYAMLComment(line string) string {
	inSingle := false
	inDouble := false
	escaped := false
	for i, r := range line {
		if inDouble && escaped {
			escaped = false
			continue
		}
		if inDouble && r == '\\' {
			escaped = true
			continue
		}
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble && (i == 0 || line[i-1] == ' ') {
				return strings.TrimRight(line[:i], " ")
			}
		}
	}
	return line
}

// parseDocument parses the full YAML document and rejects any trailing content
// left after the top-level block has been consumed.
func (p *yamlParser) parseDocument() (any, error) {
	if len(p.lines) == 0 {
		return map[string]any{}, nil
	}
	value, err := p.parseBlock(p.lines[p.pos].indent)
	if err != nil {
		return nil, err
	}
	if p.pos != len(p.lines) {
		line := p.lines[p.pos]
		return nil, fmt.Errorf("line %d: unexpected content %q", line.number, line.text)
	}
	return value, nil
}

// parseBlock dispatches an indentation-aligned YAML block to either map or list
// parsing based on the first token at that indentation level.
func (p *yamlParser) parseBlock(indent int) (any, error) {
	if p.pos >= len(p.lines) {
		return map[string]any{}, nil
	}

	line := p.lines[p.pos]
	if line.indent < indent {
		return nil, expectedIndentationError(line, indent)
	}
	if line.indent > indent {
		return nil, unexpectedIndentationError(line)
	}

	if isYAMLListItem(line) {
		return p.parseList(indent)
	}

	return p.parseMap(indent)
}

// parseMap consumes consecutive key-value entries at one indentation level.
func (p *yamlParser) parseMap(indent int) (map[string]any, error) {
	out := make(map[string]any)
	for p.pos < len(p.lines) {
		line := p.lines[p.pos]
		if line.indent < indent {
			break
		}
		if line.indent > indent {
			return nil, unexpectedIndentationError(line)
		}
		if isYAMLListItem(line) {
			break
		}

		key, raw, hasValue, err := splitYAMLKeyValue(line)
		if err != nil {
			return nil, err
		}
		if _, exists := out[key]; exists {
			return nil, fmt.Errorf("line %d: duplicate key %q", line.number, key)
		}
		p.pos++
		if hasValue {
			value, err := parseYAMLScalar(raw)
			if err != nil {
				return nil, lineError(line, err)
			}

			out[key] = value
			continue
		}

		if p.pos >= len(p.lines) || p.lines[p.pos].indent <= indent {
			out[key] = map[string]any{}
			continue
		}
		value, err := p.parseBlock(p.lines[p.pos].indent)
		if err != nil {
			return nil, err
		}
		out[key] = value
	}
	return out, nil
}

// parseList consumes consecutive dash-prefixed items at one indentation level.
func (p *yamlParser) parseList(indent int) ([]any, error) {
	out := make([]any, 0)

	for p.pos < len(p.lines) {
		line := p.lines[p.pos]

		ok, err := validateListItemLine(line, indent)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		value, err := p.parseListItem(line, indent)
		if err != nil {
			return nil, err
		}

		out = append(out, value)
	}

	return out, nil
}

// parseListItem parses one list item after the dash has been recognized.
func (p *yamlParser) parseListItem(line yamlLine, indent int) (any, error) {
	item := strings.TrimSpace(strings.TrimPrefix(line.text, "-"))
	p.pos++

	if item == "" {
		return p.parseNestedListItem(indent)
	}

	if looksLikeInlineMapItem(item) {
		return p.parseInlineMapItem(line, indent, item)
	}

	value, err := parseYAMLScalar(item)
	if err != nil {
		return nil, lineError(line, err)
	}

	return value, nil
}

// parseNestedListItem parses a dash item whose value starts on following
// indented lines.
func (p *yamlParser) parseNestedListItem(indent int) (any, error) {
	if p.pos >= len(p.lines) || p.lines[p.pos].indent <= indent {
		return nil, nil
	}

	return p.parseBlock(p.lines[p.pos].indent)
}

// parseInlineMapItem parses list entries such as "- name: module" and then
// merges subsequent aligned continuation keys into the same map.
func (p *yamlParser) parseInlineMapItem(
	line yamlLine,
	indent int,
	item string,
) (map[string]any, error) {
	entryLine := yamlLine{
		number: line.number,
		indent: indent + 2,
		text:   item,
	}

	out := map[string]any{}
	if err := p.addInlineMapEntry(out, entryLine, indent); err != nil {
		return nil, err
	}
	if err := p.addMapContinuations(out, indent); err != nil {
		return nil, err
	}

	return out, nil
}

// addMapContinuations consumes key-value lines that continue an inline list map
// item, for example "sourceDir" following "- name".
func (p *yamlParser) addMapContinuations(
	out map[string]any,
	indent int,
) error {
	for p.pos < len(p.lines) && p.lines[p.pos].indent > indent {
		line := p.lines[p.pos]
		if isYAMLListItem(line) {
			break
		}
		if line.indent != indent+2 {
			return expectedListContinuationIndentationError(line, indent+2)
		}
		if err := p.addMapEntry(out, line, indent); err != nil {
			return err
		}
	}

	return nil
}

// addInlineMapEntry adds a synthetic inline entry without advancing parser
// position because the surrounding list parser already consumed that line.
func (p *yamlParser) addInlineMapEntry(
	out map[string]any,
	line yamlLine,
	parentIndent int,
) error {
	return p.addMapEntryValue(out, line, parentIndent, false)
}

// addMapEntry adds a real map entry from the current parser line and advances
// the parser after the key is consumed.
func (p *yamlParser) addMapEntry(
	out map[string]any,
	line yamlLine,
	parentIndent int,
) error {
	return p.addMapEntryValue(out, line, parentIndent, true)
}

// addMapEntryValue validates duplicate keys, optionally advances the cursor,
// and stores either the scalar value or the nested block value.
func (p *yamlParser) addMapEntryValue(
	out map[string]any,
	line yamlLine,
	parentIndent int,
	advance bool,
) error {
	key, raw, hasValue, err := splitYAMLKeyValue(line)
	if err != nil {
		return err
	}
	if _, exists := out[key]; exists {
		return duplicateYAMLKeyError(line, key)
	}

	if advance {
		p.pos++
	}

	value, err := p.parseMapEntryValue(line, parentIndent, raw, hasValue)
	if err != nil {
		return err
	}

	out[key] = value
	return nil
}

// parseMapEntryValue converts a map entry's raw scalar or nested block into the
// generic representation later marshaled as JSON.
func (p *yamlParser) parseMapEntryValue(
	line yamlLine,
	parentIndent int,
	raw string,
	hasValue bool,
) (any, error) {
	if hasValue {
		value, err := parseYAMLScalar(raw)
		if err != nil {
			return nil, lineError(line, err)
		}

		return value, nil
	}

	if p.pos >= len(p.lines) || p.lines[p.pos].indent <= parentIndent {
		return map[string]any{}, nil
	}

	return p.parseBlock(p.lines[p.pos].indent)
}
