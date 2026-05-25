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

package target

import "arcoris.dev/arcoris-publisher/internal/manifest"

// PrepareStatus is a stable target-prepare state string.
type PrepareStatus string

const (
	PrepareStatusPrepared PrepareStatus = "prepared"
	PrepareStatusPassed   PrepareStatus = "passed"
	PrepareStatusFailed   PrepareStatus = "failed"
	PrepareStatusSkipped  PrepareStatus = "skipped"
)

// PrepareResult contains action-level target worktree preparation outcomes.
type PrepareResult struct {
	status     PrepareStatus
	targetRoot string
	modules    []PrepareModuleResult
}

// Status returns the aggregate target preparation status.
func (r PrepareResult) Status() PrepareStatus { return r.status }

// TargetRoot returns the local root containing prepared target worktrees.
func (r PrepareResult) TargetRoot() string { return r.targetRoot }

// Modules returns detached module results.
func (r PrepareResult) Modules() []PrepareModuleResult {
	out := make([]PrepareModuleResult, len(r.modules))
	copy(out, r.modules)
	return out
}

// Failed reports whether any module preparation failed.
func (r PrepareResult) Failed() bool { return r.status == PrepareStatusFailed }

// Empty reports whether no preparation checks have run.
func (r PrepareResult) Empty() bool { return r.targetRoot == "" && len(r.modules) == 0 }

// PrepareModuleResult contains preparation actions for one planned module.
type PrepareModuleResult struct {
	module      manifest.ModuleName
	repository  manifest.RepositoryRef
	worktreeDir string
	remoteName  string
	remoteURL   string
	status      PrepareStatus
	actions     []PrepareActionResult
}

// Module returns the planned module name.
func (m PrepareModuleResult) Module() manifest.ModuleName { return m.module }

// Repository returns the target repository identity.
func (m PrepareModuleResult) Repository() manifest.RepositoryRef { return m.repository }

// WorktreeDir returns the local target worktree directory.
func (m PrepareModuleResult) WorktreeDir() string { return m.worktreeDir }

// RemoteName returns the Git remote used for fetches.
func (m PrepareModuleResult) RemoteName() string { return m.remoteName }

// RemoteURL returns the resolved or configured remote transport URL.
func (m PrepareModuleResult) RemoteURL() string { return m.remoteURL }

// Status returns the aggregate module preparation status.
func (m PrepareModuleResult) Status() PrepareStatus { return m.status }

// Actions returns detached action results.
func (m PrepareModuleResult) Actions() []PrepareActionResult {
	out := make([]PrepareActionResult, len(m.actions))
	copy(out, m.actions)
	return out
}

// PrepareActionResult describes one stable preparation action.
type PrepareActionResult struct {
	name    string
	status  PrepareStatus
	code    string
	message string
	path    string
	remote  string
}

// Name returns the stable action identifier.
func (a PrepareActionResult) Name() string { return a.name }

// Status returns the action status.
func (a PrepareActionResult) Status() PrepareStatus { return a.status }

// Code returns a machine-readable failure code.
func (a PrepareActionResult) Code() string { return a.code }

// Message returns a human-readable diagnostic.
func (a PrepareActionResult) Message() string { return a.message }

// Path returns the local filesystem path related to the action, if any.
func (a PrepareActionResult) Path() string { return a.path }

// RemoteURL returns the remote URL related to the action, if any.
func (a PrepareActionResult) RemoteURL() string { return a.remote }

func preparePassed(name, message string) PrepareActionResult {
	return PrepareActionResult{name: name, status: PrepareStatusPassed, message: message}
}

func prepareSkipped(name, message string) PrepareActionResult {
	return PrepareActionResult{name: name, status: PrepareStatusSkipped, message: message}
}

func prepareFailed(name, code, message string) PrepareActionResult {
	return PrepareActionResult{name: name, status: PrepareStatusFailed, code: code, message: message}
}

func preparePath(action PrepareActionResult, path string) PrepareActionResult {
	action.path = path
	return action
}

func prepareRemote(action PrepareActionResult, remote string) PrepareActionResult {
	action.remote = remote
	return action
}

func aggregatePrepareActions(actions []PrepareActionResult) PrepareStatus {
	for _, action := range actions {
		if action.Status() == PrepareStatusFailed {
			return PrepareStatusFailed
		}
	}
	return PrepareStatusPrepared
}

func aggregatePrepareModules(modules []PrepareModuleResult) PrepareStatus {
	for _, module := range modules {
		if module.Status() == PrepareStatusFailed {
			return PrepareStatusFailed
		}
	}
	return PrepareStatusPrepared
}
