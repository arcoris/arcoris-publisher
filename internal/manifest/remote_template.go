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

import (
	"fmt"
	"strings"
	"unicode"
)

// RemoteTemplate renders logical repository refs into concrete Git transport
// URLs. It keeps repository identity separate from clone/fetch transport.
type RemoteTemplate string

// ParseRemoteTemplate validates the placeholders supported by target
// preparation.
func ParseRemoteTemplate(value string) (RemoteTemplate, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("remoteTemplate is required")
	}
	if strings.TrimSpace(value) != value {
		return "", fmt.Errorf("remoteTemplate must not have surrounding whitespace")
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return "", fmt.Errorf("remoteTemplate must not contain control characters")
		}
	}
	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '{':
			end := strings.IndexByte(value[i+1:], '}')
			if end < 0 {
				return "", fmt.Errorf("remoteTemplate has an unterminated placeholder")
			}
			name := value[i+1 : i+1+end]
			switch name {
			case "repository", "owner", "name":
			default:
				return "", fmt.Errorf("remoteTemplate has unsupported placeholder %q", name)
			}
			i += end + 1
		case '}':
			return "", fmt.Errorf("remoteTemplate has an unmatched closing brace")
		}
	}
	return RemoteTemplate(value), nil
}

// Resolve renders a concrete Git remote URL for repository.
func (t RemoteTemplate) Resolve(repository RepositoryRef, module ModuleName) (string, error) {
	parts := strings.Split(repository.String(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("repository %q is not in owner/name form", repository)
	}
	name := parts[1]
	if name == "" && module != "" {
		name = module.String()
	}
	out := strings.NewReplacer(
		"{repository}", repository.String(),
		"{owner}", parts[0],
		"{name}", name,
	).Replace(t.String())
	if strings.TrimSpace(out) == "" {
		return "", fmt.Errorf("remoteTemplate resolved to an empty URL")
	}
	for _, r := range out {
		if unicode.IsControl(r) {
			return "", fmt.Errorf("remoteTemplate resolved to a URL with control characters")
		}
	}
	return out, nil
}

// String returns the template in manifest form.
func (t RemoteTemplate) String() string { return string(t) }
