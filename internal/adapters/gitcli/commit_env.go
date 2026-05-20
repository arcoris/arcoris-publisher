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

package gitcli

import (
	"strconv"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

// commitEnv converts commit identity options into Git environment variables.
func commitEnv(opts gitport.CommitOptions) []string {
	env := make([]string, 0, 6)
	env = appendOptionalEnv(env, "GIT_AUTHOR_NAME", opts.AuthorName)
	env = appendOptionalEnv(env, "GIT_AUTHOR_EMAIL", opts.AuthorEmail)
	env = appendOptionalEnv(env, "GIT_COMMITTER_NAME", opts.CommitterName)
	env = appendOptionalEnv(env, "GIT_COMMITTER_EMAIL", opts.CommitterEmail)
	if !opts.AuthorDate.IsZero() {
		env = append(env, "GIT_AUTHOR_DATE="+gitTimestamp(opts.AuthorDate.Unix()))
	}
	if !opts.CommitterDate.IsZero() {
		env = append(env, "GIT_COMMITTER_DATE="+gitTimestamp(opts.CommitterDate.Unix()))
	}
	return env
}

// appendOptionalEnv appends KEY=VALUE only when value is not empty.
func appendOptionalEnv(env []string, key string, value string) []string {
	if value == "" {
		return env
	}
	return append(env, key+"="+value)
}

// gitTimestamp formats a Unix timestamp in the stable UTC form Git accepts.
func gitTimestamp(unix int64) string {
	return strconv.FormatInt(unix, 10) + " +0000"
}
