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

package git

// CheckoutOptions configures a Git checkout operation.
type CheckoutOptions struct {
	// Create creates a new branch at the requested ref.
	Create bool
	// Force permits checkout when local changes would otherwise block it.
	Force bool
	// Detach checks out the target commit without selecting a branch.
	Detach bool
	// Orphan creates a branch with no parent commit.
	Orphan bool
}

// CreateBranchOptions configures local branch creation.
type CreateBranchOptions struct {
	// Force permits replacing an existing branch reference.
	Force bool
}
