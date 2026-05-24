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

package app

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/versioning"
)

func TestBuildPlanLoadsConfigPipeline(t *testing.T) {
	p, err := New(Dependencies{}, Options{}).BuildPlan(
		context.Background(),
		"../config/testdata/minimal/arcpub.yaml",
		versioning.Must("v0.3.0"),
	)

	if err != nil {
		t.Fatalf("BuildPlan() error = %v", err)
	}
	names := p.ModuleNames()
	if len(names) != 2 || names[0] != "foundation" || names[1] != "control" {
		t.Fatalf("ModuleNames() = %v", names)
	}

	control, ok := p.ModuleByName("control")
	if !ok {
		t.Fatal("control module missing")
	}
	requirements := control.Requirements()
	if len(requirements) != 1 || requirements[0].ModulePath() != "arcoris.dev/foundation" {
		t.Fatalf("control requirements = %#v", requirements)
	}
}

func TestBuildPlanRejectsInvalidConfig(t *testing.T) {
	_, err := New(Dependencies{}, Options{}).BuildPlan(
		context.Background(),
		"../config/testdata/unknown-dependency/arcpub.yaml",
		versioning.Must("v0.3.0"),
	)

	if err == nil {
		t.Fatal("BuildPlan() error = nil")
	}
}

func TestBuildPlanRejectsMissingVersion(t *testing.T) {
	_, err := New(Dependencies{}, Options{}).BuildPlan(
		context.Background(),
		"../config/testdata/minimal/arcpub.yaml",
		versioning.Version(""),
	)

	if err == nil {
		t.Fatal("BuildPlan() error = nil")
	}
}
