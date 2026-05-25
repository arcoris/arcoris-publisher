#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

run_step "resetting lab for lock conflict"
bash "$LAB_DIR/setup.sh"
mkdir -p "$(lab_state_dir)"
printf 'transaction=tx-other\npid=1\nstartedAt=manual-lab\n' >"$(lab_state_dir)/publish.lock"

report="$(lab_reports)/lock-conflict.json"
stderr="$(lab_logs)/lock-conflict.stderr"

run_step "running publish with existing lock"
code="$(run_common_arcpub_json "$report" "$stderr" publish)"
assert_code "$code" 1 "$stderr"
assert_contains "$stderr" "lock"
assert_no_final_refs
log "lock conflict blocked publish before mutation"
