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

package versioning

import (
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/registry"
)

func TestAssignRejectsMissingVersion(t *testing.T) {
	req := mustRequest(t, "v0.3.0", "", testModule{name: "foundation"})
	req.Version = ""
	_, err := Assign(req)
	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(IssueInvalidVersion) {
		t.Fatalf("Assign() error = %v, want missing version", err)
	}
}

func TestAssignRejectsZeroRequest(t *testing.T) {
	_, err := Assign(Request{Version: Must("v0.3.0")})
	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(IssueInvalidRequest) {
		t.Fatalf("Assign() error = %v, want invalid request", err)
	}
}

func TestAssignRejectsMissingRequestIndexes(t *testing.T) {
	set := mustPublicationSet(t, string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
	)
	_, err := Assign(Request{
		Set:     set,
		Version: Must("v0.3.0"),
	})

	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(IssueUnknownModule) {
		t.Fatalf("Assign() error = %v, want missing index validation", err)
	}
}

func TestAssignRejectsRegistryPathMismatch(t *testing.T) {
	set := mustPublicationSet(t, string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation", modulePath: "arcoris.dev/foundation"},
	)
	indexSet := mustPublicationSet(t, string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation", modulePath: "arcoris.dev/foundation-alt"},
	)
	reg, err := registry.New(indexSet)
	if err != nil {
		t.Fatalf("registry.New() error = %v", err)
	}
	g, err := graph.New(reg)
	if err != nil {
		t.Fatalf("graph.New() error = %v", err)
	}

	_, err = Assign(Request{
		Set:      set,
		Registry: reg,
		Graph:    g,
		Version:  Must("v0.3.0"),
	})

	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(IssueInvalidRequest) {
		t.Fatalf("Assign() error = %v, want registry mismatch validation", err)
	}
}
