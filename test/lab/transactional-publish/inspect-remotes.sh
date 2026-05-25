#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

for repo in arcoris/foundation arcoris/control; do
  bare="$(bare_repo "$repo")"
  run_step "remote $repo: $bare"
  echo "heads:"
  git --git-dir "$bare" show-ref --heads || true
  echo "tags:"
  git --git-dir "$bare" show-ref --tags || true
  echo "candidate refs:"
  git --git-dir "$bare" for-each-ref refs/heads/arcpub/tx || true
  echo "latest main commit:"
  git --git-dir "$bare" log -1 --format=%B main || true
  echo "tree:"
  git --git-dir "$bare" ls-tree -r --name-only main || true
done
