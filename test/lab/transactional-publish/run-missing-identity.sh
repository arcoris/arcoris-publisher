#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

require_cmd git

run_step "resetting lab for missing commit identity scenario"
bash "$LAB_DIR/setup.sh"

report="$(lab_reports)/missing-identity-preflight.json"
stderr="$(lab_logs)/missing-identity-preflight.stderr"

run_step "clearing target worktree Git identity"
clear_target_identity
isolate_git_identity

run_step "running preflight; missing commit identity should fail readiness"
code="$(run_common_arcpub_json "$report" "$stderr" preflight)"
assert_code "$code" 1 "$stderr"
assert_json_if_python "$report"
assert_no_path_leak "$report"
assert_contains "$report" '"name": "commit-identity"'
assert_contains "$report" '"code": "missing_commit_identity"'

if [[ -e "$(lab_state_dir)/publish.lock" || -d "$(lab_state_dir)/transactions" ]]; then
  fail "missing identity preflight created publish transaction state"
fi

log "missing identity report: $report"
