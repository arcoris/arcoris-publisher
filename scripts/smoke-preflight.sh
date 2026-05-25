#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"
LAB_ROOT="${TMPDIR:-/tmp}/arcpub-smoke-preflight"

echo "smoke-preflight: go test ./test/e2e -run Preflight"
go test ./test/e2e -run Preflight -count=1 -v

echo "smoke-preflight: lab preflight"
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/setup.sh
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/run-preflight.sh
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/cleanup.sh

echo "smoke-preflight: ok"
