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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/app"
	"arcoris.dev/arcoris-publisher/internal/workflow"
)

func TestApplicationFactoryFromAppReturnsApplication(t *testing.T) {
	t.Parallel()

	application, err := ApplicationFactoryFromApp(app.Dependencies{})(app.Options{})
	if err != nil {
		t.Fatalf("ApplicationFactoryFromApp() error = %v", err)
	}
	if application == nil {
		t.Fatal("ApplicationFactoryFromApp() returned nil application")
	}
}

func TestDependenciesApplicationRequiresApp(t *testing.T) {
	t.Parallel()

	_, err := (Dependencies{}).application(app.Options{})
	if err == nil {
		t.Fatal("Dependencies.application() error = nil")
	}
}

func TestDependenciesApplicationUsesFactoryOptions(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	var got app.Options
	deps := Dependencies{
		AppFactory: func(opts app.Options) (Application, error) {
			got = opts
			return fake, nil
		},
	}

	application, err := deps.application(app.Options{Workflow: workflow.Options{DryRun: true}})
	if err != nil {
		t.Fatalf("Dependencies.application() error = %v", err)
	}
	if application != fake {
		t.Fatalf("application = %#v", application)
	}
	if !got.Workflow.DryRun {
		t.Fatalf("factory options were not captured")
	}
}
