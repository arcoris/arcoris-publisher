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

package plan

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/versioning"
)

func TestFromPublicationSetBuildsCompletePlan(t *testing.T) {
	set := mustPublicationSet(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	p, err := FromPublicationSet(set, versioning.Must("v0.3.0"))
	if err != nil {
		t.Fatalf("FromPublicationSet() error = %v", err)
	}
	assertModuleNames(t, p.ModuleNames(), "foundation", "control")
}

func TestFromPublicationSetReturnsVersioningError(t *testing.T) {
	set := mustPublicationSet(t, testModule{name: "foundation"})
	_, err := FromPublicationSet(set, versioning.Must("v0.0.0-20260521123456-abcdefabcdef"))
	if err == nil {
		t.Fatal("error = nil, want versioning error")
	}
}
