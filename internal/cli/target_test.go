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

package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunTargetPreparePassesRequest(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	cli := New(Dependencies{App: fake}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"target", "prepare",
		"--manifest", "staging/arcpub.yaml",
		"--version", "v0.3.0",
		"--target-root", "/target",
		"--remote-template", "file:///remotes/{name}.git",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(target prepare) code = %d stderr = %s", code, stderr.String())
	}
	if !fake.prepareTargetsCalled {
		t.Fatal("PrepareTargets was not called")
	}
	if fake.prepareTargetsRequest.ManifestPath != "staging/arcpub.yaml" ||
		fake.prepareTargetsRequest.TargetRootDir != "/target" ||
		fake.prepareTargetsRequest.TargetRemoteTemplate != "file:///remotes/{name}.git" {
		t.Fatalf("target prepare request = %+v", fake.prepareTargetsRequest)
	}
}

func TestRunTargetPrepareMissingVersionIsUsage(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	cli := New(Dependencies{App: fake}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"target", "prepare"}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(target prepare) code = %d stderr = %s", code, stderr.String())
	}
	if fake.prepareTargetsCalled {
		t.Fatal("PrepareTargets called for missing version")
	}
}

func TestRunTargetPrepareInvalidRemoteTemplateIsUsage(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	cli := New(Dependencies{App: fake}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"target", "prepare",
		"--version", "v0.3.0",
		"--remote-template", "file:///remotes/{bogus}.git",
	}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(target prepare invalid template) code = %d stderr = %s", code, stderr.String())
	}
	if fake.prepareTargetsCalled {
		t.Fatal("PrepareTargets called for invalid remote template")
	}
}

func TestRunTargetPrepareHelpDescribesMutationBoundary(t *testing.T) {
	t.Parallel()

	cli := New(Dependencies{}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"target", "prepare", "--help"}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(target prepare --help) code = %d stderr = %s", code, stderr.String())
	}
	output := stdout.String()
	for _, want := range []string{"clone missing worktrees", "fetch remote state", "does not construct", "commit, tag, push"} {
		if !strings.Contains(output, want) {
			t.Fatalf("target prepare help missing %q:\n%s", want, output)
		}
	}
}
