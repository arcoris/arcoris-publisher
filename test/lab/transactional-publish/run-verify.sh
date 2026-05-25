#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

ensure_lab_ready
mkdir -p "$(lab_reports)"
report="$(lab_reports)/verify.json"
stderr="$(lab_logs)/verify.stderr"

run_step "running verify"
code="$(run_common_arcpub_json "$report" "$stderr" verify)"
assert_code "$code" 0 "$stderr"
assert_json_if_python "$report"
assert_no_path_leak "$report"

run_step "target worktree status after verify"
git -C "$(repository_worktree arcoris/foundation)" status --porcelain || true
git -C "$(repository_worktree arcoris/control)" status --porcelain || true
log "verify constructs target trees; worktrees may be dirty until reset or publish"
