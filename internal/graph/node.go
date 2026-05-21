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

package graph

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Node is the graph-level projection of one resolved publication module.
//
// Node intentionally stores only topology-relevant fields. Target repositories,
// publish entries, verification policy, and other execution details belong to
// the resolved manifest model and later planning layers, not to graph topology.
type Node struct {
	name       manifest.ModuleName
	modulePath manifest.ModulePath
	visibility manifest.Visibility
}

// Name returns the module name represented by the node.
func (n Node) Name() manifest.ModuleName { return n.name }

// ModulePath returns the published Go module path associated with the node.
func (n Node) ModulePath() manifest.ModulePath { return n.modulePath }

// Visibility returns the resolved module visibility.
func (n Node) Visibility() manifest.Visibility { return n.visibility }

// Publishable reports whether the node should participate in publication order.
func (n Node) Publishable() bool { return n.visibility == manifest.VisibilityPublic }

// Internal reports whether the node represents a known but non-published module.
func (n Node) Internal() bool { return n.visibility == manifest.VisibilityInternal }

// Disabled reports whether the node is disabled.
//
// Graph construction currently excludes disabled modules, so this is normally
// false for all nodes returned by Graph.Nodes. The method is retained to keep
// visibility handling explicit at call sites.
func (n Node) Disabled() bool { return n.visibility == manifest.VisibilityDisabled }
