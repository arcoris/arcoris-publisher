#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

export GOTOOLCHAIN="${GOTOOLCHAIN:-local}"

echo "smoke-transactions-lock: e2e lock scenarios"
go test ./test/e2e -run TransactionsLock -count=1 -v

echo "smoke-transactions-lock: lab lock scenario"
LAB_ROOT="${TMPDIR:-/tmp}/arcpub-smoke-transactions-lock"
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/setup.sh
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/run-lock.sh
ARCPUB_LAB_ROOT="$LAB_ROOT" bash test/lab/transactional-publish/cleanup.sh

echo "smoke-transactions-lock: ok"
