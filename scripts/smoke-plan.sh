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

go build -o "$tmp/arcpub" ./cmd/arcpub

manifest="internal/testdata/e2e/minimal/arcpub.yaml"
"$tmp/arcpub" plan --manifest "$manifest" --version v0.1.0 >/dev/null
"$tmp/arcpub" plan --manifest "$manifest" --version v0.1.0 --output json >"$tmp/plan.json"
if command -v python3 >/dev/null 2>&1; then
	python3 -m json.tool "$tmp/plan.json" >/dev/null
else
	echo "python3 not found; skipping JSON validation"
fi

echo "smoke-plan: ok"
