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

package publish

import (
	"fmt"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
)

// commitMessage renders deterministic publication provenance for one module.
func commitMessage(mod plan.ModulePlan, src source.Snapshot) string {
	var b strings.Builder
	fmt.Fprintf(&b, "sync: publish %s from %s\n\n", mod.Name(), src.Repository().RepositoryDir())
	if src.Repository().Head() != "" {
		fmt.Fprintf(&b, "Arcoris-Source-Commit: %s\n", src.Repository().Head())
	}
	if src.Repository().Branch() != "" {
		fmt.Fprintf(&b, "Arcoris-Source-Branch: %s\n", src.Repository().Branch())
	}
	fmt.Fprintf(&b, "Arcoris-Module: %s\n", mod.Name())
	fmt.Fprintf(&b, "Arcoris-Module-Path: %s\n", mod.ModulePath())
	fmt.Fprintf(&b, "Arcoris-Version: %s\n", mod.Version())
	return b.String()
}
