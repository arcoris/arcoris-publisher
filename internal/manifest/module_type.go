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

import "fmt"

// ModuleType identifies the kind of module artifact being published.
type ModuleType string

const (
	// ModuleTypeGo identifies a Go module publication unit.
	ModuleTypeGo ModuleType = "go"
)

// ParseModuleType validates a module type string.
func ParseModuleType(value string) (ModuleType, error) {
	switch ModuleType(value) {
	case ModuleTypeGo:
		return ModuleType(value), nil
	default:
		return "", fmt.Errorf("unsupported module.type %q", value)
	}
}
