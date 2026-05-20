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

// Package github implements the remote provider port for GitHub-compatible APIs.
//
// The adapter speaks GitHub's REST API and translates provider-specific response
// shapes into the remote port's repository, branch protection, pull request, and
// release types. It intentionally does not perform Git transport operations.
package github

import (
	"net/http"
	"strings"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

const defaultAPIBaseURL = "https://api.github.com"

// Provider implements remote.Provider for GitHub's REST API.
//
// Provider is safe to reuse. Tests can supply a custom http.Client transport,
// while production code can use the default client and the public GitHub API.
type Provider struct {
	client  *http.Client
	baseURL string
	token   string
}

// Options configures a GitHub provider.
type Options struct {
	// Client sends HTTP requests. Nil uses http.DefaultClient.
	Client *http.Client
	// BaseURL points at the GitHub API root. Empty uses https://api.github.com.
	BaseURL string
	// Token is sent as a bearer token when non-empty.
	Token string
}

// New creates a GitHub provider.
func New(opts Options) *Provider {
	client := opts.Client
	if client == nil {
		client = http.DefaultClient
	}
	base := opts.BaseURL
	if base == "" {
		base = defaultAPIBaseURL
	}
	return &Provider{client: client, baseURL: strings.TrimRight(base, "/"), token: opts.Token}
}

var _ remoteport.Provider = (*Provider)(nil)
