#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

run_step "resetting lab for corrupted journal"
bash "$LAB_DIR/setup.sh"
mkdir -p "$(lab_state_dir)/transactions"
printf '{' >"$(lab_state_dir)/transactions/tx-corrupt.json"

report="$(lab_reports)/corrupted-journal.json"
stderr="$(lab_logs)/corrupted-journal.stderr"

run_step "running publish with corrupted journal"
code="$(run_common_arcpub_json "$report" "$stderr" publish)"
assert_code "$code" 1 "$stderr"
assert_contains "$stderr" "corrupt"
assert_no_final_refs
log "corrupted journal blocked publish before mutation"
