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

	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

func TestPublishRunsWorkflow(t *testing.T) {
	app, fakeGit := appFixture(t)
	fakeGit.Statuses["/target/arcoris__foundation"] = dirtyStatus()
	fakeGit.Statuses["/target/arcoris__control"] = dirtyStatus()

	result, err := app.Publish(context.Background(), appRequest())

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Workflow().Publish().Published() {
		t.Fatal("publish use case did not publish")
	}
}

func dirtyStatus() git.Status {
	return git.Status{
		Clean:   false,
		Entries: []git.StatusEntry{{Path: "go.mod", Code: " M"}},
	}
}
