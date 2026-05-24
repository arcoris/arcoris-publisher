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
)

// commitMessage renders deterministic publication provenance for one module.
func commitMessage(mod plan.ModulePlan, req Request) string {
	var b strings.Builder
	fmt.Fprintf(&b, "sync: publish %s %s\n\n", mod.Name(), mod.Version())
	if req.Plan.Source().Repository() != "" {
		fmt.Fprintf(&b, "Arcoris-Source-Repository: %s\n", req.Plan.Source().Repository())
	}
	if req.Source.Repository().Head() != "" {
		fmt.Fprintf(&b, "Arcoris-Source-Commit: %s\n", req.Source.Repository().Head())
	}
	if req.Source.Repository().Branch() != "" {
		fmt.Fprintf(&b, "Arcoris-Source-Branch: %s\n", req.Source.Repository().Branch())
	}
	fmt.Fprintf(&b, "Arcoris-Module: %s\n", mod.Name())
	fmt.Fprintf(&b, "Arcoris-Module-Path: %s\n", mod.ModulePath())
	fmt.Fprintf(&b, "Arcoris-Version: %s\n", mod.Version())
	fmt.Fprintf(&b, "Arcoris-Target-Repository: %s\n", mod.Repository())
	fmt.Fprintf(&b, "Arcoris-Target-Branches: %s\n", targetBranchTrailer(mod))
	fmt.Fprintf(&b, "Arcoris-Publish-Mode: %s\n", req.Plan.PublishPolicy().Mode())
	fmt.Fprintf(&b, "Arcoris-Push-Policy: %s\n", req.Plan.PublishPolicy().PushPolicy())
	fmt.Fprintf(&b, "Arcoris-Tag-Policy: %s\n", req.Plan.PublishPolicy().Tags().Mode())
	return b.String()
}

func targetBranchTrailer(mod plan.ModulePlan) string {
	branches := mod.Branches()
	values := make([]string, 0, len(branches))
	for _, branch := range branches {
		values = append(values, branch.Target().String())
	}
	return strings.Join(values, ",")
}
