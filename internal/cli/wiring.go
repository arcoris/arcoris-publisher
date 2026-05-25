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
	"context"

	"arcoris.dev/arcoris-publisher/internal/app"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

// Application is the app-layer boundary used by CLI commands.
//
// app.App satisfies this interface. Tests may provide fakes without importing
// concrete infrastructure adapters.
type Application interface {
	BuildPlan(ctx context.Context, manifestPath string, version versioning.Version) (plan.Plan, error)
	Preflight(ctx context.Context, req app.Request) (app.Result, error)
	Verify(ctx context.Context, req app.Request) (app.Result, error)
	Publish(ctx context.Context, req app.Request) (app.Result, error)
	ListTransactions(ctx context.Context, req app.TransactionRequest) (app.TransactionListResult, error)
	ShowTransaction(ctx context.Context, req app.TransactionRequest) (app.TransactionResult, error)
	RollbackTransaction(ctx context.Context, req app.TransactionRequest) (app.TransactionResult, error)
}

// AppFactory constructs an application instance for one command invocation.
//
// The app options argument contains command-line overrides such as --dry-run.
type AppFactory func(app.Options) (Application, error)

// Dependencies contains CLI collaborators.
type Dependencies struct {
	// App is an optional preconstructed application used mostly by tests and
	// embedders. Command flags that affect app construction, such as --dry-run,
	// are only applied when AppFactory is used.
	App Application

	// AppFactory lazily constructs the application when a command needs it. It is
	// ignored when App is set.
	AppFactory AppFactory

	// BuildInfo returns publisher build metadata for the version command.
	BuildInfo BuildInfoFunc
}

func (d Dependencies) application(opts app.Options) (Application, error) {
	if d.App != nil {
		return d.App, nil
	}
	if d.AppFactory != nil {
		return d.AppFactory(opts)
	}
	return nil, &Error{
		Code:    CodeMissingApplication,
		Message: "application dependency is required for plan, verify, and publish commands",
	}
}

// ApplicationFromApp constructs the CLI application boundary from app
// dependencies and options.
//
// Concrete infrastructure adapters should be assembled by the caller and passed
// through app.Dependencies. This helper deliberately does not import adapters so
// workflow packages remain port-oriented.
func ApplicationFromApp(deps app.Dependencies, opts app.Options) Application {
	return app.New(deps, opts)
}

// ApplicationFactoryFromApp returns an AppFactory backed by app.New.
func ApplicationFactoryFromApp(deps app.Dependencies) AppFactory {
	return func(opts app.Options) (Application, error) {
		return ApplicationFromApp(deps, opts), nil
	}
}

var _ Application = app.App{}
