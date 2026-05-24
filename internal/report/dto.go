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

package report

// BranchReport describes one source-to-target branch mapping.
type BranchReport struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// PublishPolicyReport describes global publication policy values relevant to
// reports.
type PublishPolicyReport struct {
	Mode                  string `json:"mode"`
	PushPolicy            string `json:"pushPolicy"`
	TagPolicy             string `json:"tagPolicy"`
	TagEnabled            bool   `json:"tagEnabled"`
	ProvenanceFileEnabled bool   `json:"provenanceFileEnabled"`
	ProvenanceFile        string `json:"provenanceFile,omitempty"`
	CommitTrailers        bool   `json:"commitTrailers"`
}

// PublishEntryReport describes one explicit publication entry.
type PublishEntryReport struct {
	Kind      string `json:"kind"`
	From      string `json:"from"`
	To        string `json:"to"`
	Optional  bool   `json:"optional"`
	Recursive bool   `json:"recursive"`
}

// DependencyRequirementReport describes one direct internal module requirement.
type DependencyRequirementReport struct {
	Module     string `json:"module"`
	ModulePath string `json:"modulePath"`
	Version    string `json:"version"`
}
