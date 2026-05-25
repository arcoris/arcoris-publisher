#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

ensure_lab_ready
mkdir -p "$(lab_reports)"
report="$(lab_reports)/preflight.json"
stderr="$(lab_logs)/preflight.stderr"

foundation_head="$(git -C "$(repository_worktree arcoris/foundation)" rev-parse HEAD)"
control_head="$(git -C "$(repository_worktree arcoris/control)" rev-parse HEAD)"
state_dir_existed=0
if [[ -e "$(lab_state_dir)" ]]; then
  state_dir_existed=1
fi
transaction_count_before=0
if [[ -d "$(lab_state_dir)/transactions" ]]; then
  transaction_count_before="$(find "$(lab_state_dir)/transactions" -type f -name '*.json' | wc -l | tr -d '[:space:]')"
fi
lock_existed_before=0
if [[ -e "$(lab_state_dir)/publish.lock" ]]; then
  lock_existed_before=1
fi

run_step "running read-only preflight"
code="$(run_common_arcpub_json "$report" "$stderr" preflight)"
assert_code "$code" 0 "$stderr"
assert_json_if_python "$report"
assert_no_path_leak "$report"

[[ "$(git -C "$(repository_worktree arcoris/foundation)" rev-parse HEAD)" == "$foundation_head" ]] || fail "foundation HEAD changed"
[[ "$(git -C "$(repository_worktree arcoris/control)" rev-parse HEAD)" == "$control_head" ]] || fail "control HEAD changed"
assert_clean_worktree "$(repository_worktree arcoris/foundation)"
assert_clean_worktree "$(repository_worktree arcoris/control)"

if [[ "$state_dir_existed" == "0" && -e "$(lab_state_dir)" ]]; then
  fail "preflight created transaction state dir: $(lab_state_dir)"
fi
transaction_count_after=0
if [[ -d "$(lab_state_dir)/transactions" ]]; then
  transaction_count_after="$(find "$(lab_state_dir)/transactions" -type f -name '*.json' | wc -l | tr -d '[:space:]')"
fi
[[ "$transaction_count_after" == "$transaction_count_before" ]] || fail "preflight changed transaction journal count"
lock_existed_after=0
if [[ -e "$(lab_state_dir)/publish.lock" ]]; then
  lock_existed_after=1
fi
[[ "$lock_existed_after" == "$lock_existed_before" ]] || fail "preflight changed publish lock state"

log "preflight report: $report"
