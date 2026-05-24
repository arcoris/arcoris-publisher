#!/usr/bin/env bash
# Copyright 2026 The ARCORIS Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "smoke-exit-codes: build arcpub"
go build -o "$tmp/arcpub" ./cmd/arcpub

expect_code() {
	local want="$1"
	shift
	echo "smoke-exit-codes: expect $want: arcpub $*"
	set +e
	"$tmp/arcpub" "$@" >"$tmp/stdout.txt" 2>"$tmp/stderr.txt"
	local got="$?"
	set -e
	if [[ "$got" != "$want" ]]; then
		echo "expected exit $want, got $got for: arcpub $*" >&2
		echo "stdout:" >&2
		cat "$tmp/stdout.txt" >&2
		echo "stderr:" >&2
		cat "$tmp/stderr.txt" >&2
		exit 1
	fi
}

manifest="internal/testdata/e2e/minimal/arcpub.yaml"
expect_code 64 unknown
expect_code 64 plan --manifest "$manifest"
expect_code 64 plan --manifest "$manifest" --version not-a-version
expect_code 64 completion unknown
expect_code 64 version --output json --pretty --compact
expect_code 1 plan --manifest "$tmp/missing/arcpub.yaml" --version v0.1.0

echo "smoke-exit-codes: ok"
