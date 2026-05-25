#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

run_step "resetting lab for dry-run"
bash "$LAB_DIR/setup.sh"

report="$(lab_reports)/dry-run.json"
stderr="$(lab_logs)/dry-run.stderr"

run_step "running publish --dry-run"
code="$(run_common_arcpub_json "$report" "$stderr" publish --dry-run)"
assert_code "$code" 0 "$stderr"
assert_json_if_python "$report"
assert_no_path_leak "$report"
assert_no_final_refs
log "dry-run constructs and verifies target worktrees, then skips commit, tag, and push"
