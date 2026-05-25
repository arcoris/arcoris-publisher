#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"

echo "smoke-transactions-prune: e2e prune scenarios"
go test ./test/e2e -run TransactionsPrune -count=1 -v

echo "smoke-transactions-prune: lab prune scenario"
LAB_ROOT="${TMPDIR:-/tmp}/arcpub-smoke-transactions-prune"
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/setup.sh
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/run-prune.sh
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/cleanup.sh

echo "smoke-transactions-prune: ok"
