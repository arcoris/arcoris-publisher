#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

for repo in arcoris/foundation arcoris/control; do
  worktree="$(repository_worktree "$repo")"
  run_step "worktree $repo: $worktree"
  echo "status:"
  git -C "$worktree" status --porcelain || true
  echo "branch:"
  git -C "$worktree" branch --show-current || true
  echo "remote:"
  git -C "$worktree" remote -v || true
  echo "last commit:"
  git -C "$worktree" log -1 --format=%B || true
  echo "files:"
  find "$worktree" -maxdepth 3 -type f ! -path '*/.git/*' | sort || true
done
