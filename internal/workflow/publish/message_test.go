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
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/buildinfo"
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
)

func TestCommitMessageContainsArcorisTrailers(t *testing.T) {
	oldVersion := buildinfo.Version
	buildinfo.Version = "v1.2.3"
	t.Cleanup(func() { buildinfo.Version = oldVersion })

	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}
	mod, _ := p.ModuleByName("foundation")

	message := commitMessage(mod, Request{Plan: p})

	for _, required := range []string{
		"Arcoris-Source-Repository:",
		"Arcoris-Module:",
		"Arcoris-Version:",
		"Arcoris-Target-Repository:",
		"Arcoris-Publish-Mode:",
		"Arcoris-Publisher-Version: v1.2.3",
	} {
		if !strings.Contains(message, required) {
			t.Fatalf("commit message missing %q:\n%s", required, message)
		}
	}
	if strings.Contains(message, "/repo") || strings.Contains(message, "/target") {
		t.Fatalf("commit message leaks local path:\n%s", message)
	}
}

func TestCommitMessageHonorsDisabledCommitTrailers(t *testing.T) {
	commitTrailers := false
	p, err := publishertest.Plan(
		publishertest.PlanOptions{
			Publish: manifest.PublishSpec{
				Provenance: manifest.ProvenanceSpec{CommitTrailers: &commitTrailers},
			},
		},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}
	mod, _ := p.ModuleByName("foundation")

	message := commitMessage(mod, Request{Plan: p})

	if message != "sync: publish foundation v0.3.0\n\n" {
		t.Fatalf("commitMessage() = %q", message)
	}
}
