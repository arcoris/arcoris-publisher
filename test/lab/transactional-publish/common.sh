#!/usr/bin/env bash
set -euo pipefail

LAB_COMMON_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LAB_REPO_ROOT="$(cd "$LAB_COMMON_DIR/../../.." && pwd)"
ARCPUB_LAB_ROOT="${ARCPUB_LAB_ROOT:-${TMPDIR:-/tmp}/arcpub-transaction-lab}"
GOTOOLCHAIN="${GOTOOLCHAIN:-local}"
GOCACHE="${GOCACHE:-$ARCPUB_LAB_ROOT/cache/go-build}"
GIT_TERMINAL_PROMPT=0
export GOTOOLCHAIN GOCACHE GIT_TERMINAL_PROMPT

lab_repo_root() { printf '%s\n' "$LAB_REPO_ROOT"; }
lab_root() { printf '%s\n' "$ARCPUB_LAB_ROOT"; }
lab_bin() { printf '%s/bin\n' "$(lab_root)"; }
lab_arcpub() { printf '%s/arcpub\n' "$(lab_bin)"; }
lab_source() { printf '%s/source\n' "$(lab_root)"; }
lab_targets() { printf '%s/targets\n' "$(lab_root)"; }
lab_state_dir() { printf '%s/.arcpub/state\n' "$(lab_targets)"; }
lab_remotes() { printf '%s/remotes\n' "$(lab_root)"; }
lab_reports() { printf '%s/reports\n' "$(lab_root)"; }
lab_logs() { printf '%s/logs\n' "$(lab_root)"; }

fail() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

log() {
  printf '[arcpub-lab] %s\n' "$*"
}

run_step() {
  log "$*"
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"
}

guard_lab_root() {
  local root
  root="$(lab_root)"
  [[ -n "$root" ]] || fail "ARCPUB_LAB_ROOT is empty"
  [[ "$root" != "/" ]] || fail "refusing to remove /"
  [[ "$root" != "$HOME" ]] || fail "refusing to remove HOME"
  [[ "$root" != "$(lab_repo_root)" ]] || fail "refusing to remove repository root"
}

reset_lab_root() {
  guard_lab_root
  rm -rf "$(lab_root)"
  mkdir -p "$(lab_bin)" "$(lab_source)" "$(lab_targets)" "$(lab_remotes)" "$(lab_reports)" "$(lab_logs)"
}

ensure_lab_ready() {
  [[ -x "$(lab_arcpub)" ]] || fail "arcpub binary missing; run setup.sh or build.sh"
  [[ -f "$(lab_source)/arcpub.yaml" ]] || fail "lab source fixture missing; run setup.sh"
}

run_git() {
  git "$@"
}

run_arcpub() {
  "$(lab_arcpub)" "$@"
}

capture_arcpub() {
  local stdout="$1"
  local stderr="$2"
  shift 2
  set +e
  "$(lab_arcpub)" "$@" >"$stdout" 2>"$stderr"
  local code=$?
  set -e
  printf '%s\n' "$code"
}

assert_code() {
  local got="$1"
  local want="$2"
  local stderr_file="${3:-}"
  if [[ "$got" != "$want" ]]; then
    [[ -n "$stderr_file" && -f "$stderr_file" ]] && sed -n '1,120p' "$stderr_file" >&2
    fail "exit code $got, want $want"
  fi
}

assert_file_exists() {
  [[ -f "$1" ]] || fail "file missing: $1"
}

assert_json_if_python() {
  local file="$1"
  if command -v python3 >/dev/null 2>&1; then
    python3 -m json.tool "$file" >/dev/null
  else
    log "python3 not found; skip JSON validation for $file"
  fi
}

assert_no_path_leak() {
  local file="$1"
  local root
  root="$(lab_root)"
  if grep -F "$root" "$file" >/dev/null 2>&1; then
    fail "report leaks lab path $root: $file"
  fi
}

assert_contains() {
  local file="$1"
  local text="$2"
  grep -F "$text" "$file" >/dev/null || fail "expected '$text' in $file"
}

assert_ref_exists() {
  local bare="$1"
  local ref="$2"
  git --git-dir "$bare" show-ref --verify --quiet "$ref" || fail "missing ref $ref in $bare"
}

assert_ref_missing() {
  local bare="$1"
  local ref="$2"
  if git --git-dir "$bare" show-ref --verify --quiet "$ref"; then
    fail "unexpected ref $ref in $bare"
  fi
}

assert_candidate_refs_absent() {
  local bare="$1"
  local refs
  refs="$(git --git-dir "$bare" for-each-ref --format='%(refname)' refs/heads/arcpub/tx || true)"
  [[ -z "$refs" ]] || fail "candidate refs remain in $bare: $refs"
}

assert_clean_worktree() {
  local worktree="$1"
  local status
  status="$(git -C "$worktree" status --porcelain)"
  [[ -z "$status" ]] || fail "dirty worktree $worktree: $status"
}

assert_tree_contains() {
  local bare="$1"
  local ref="$2"
  local path="$3"
  git --git-dir "$bare" cat-file -e "$ref:$path" || fail "tree $ref in $bare missing $path"
}

assert_tree_missing() {
  local bare="$1"
  local ref="$2"
  local path="$3"
  if git --git-dir "$bare" cat-file -e "$ref:$path" 2>/dev/null; then
    fail "tree $ref in $bare unexpectedly contains $path"
  fi
}

repository_worktree() {
  local repository="$1"
  printf '%s/%s\n' "$(lab_targets)" "${repository//\//__}"
}

bare_repo() {
  local repository="$1"
  case "$repository" in
    arcoris/foundation) printf '%s/foundation.git\n' "$(lab_remotes)" ;;
    arcoris/control) printf '%s/control.git\n' "$(lab_remotes)" ;;
    *) fail "unknown repository: $repository" ;;
  esac
}

common_workflow_args() {
  printf '%s\0' \
    --manifest "$(lab_source)/arcpub.yaml" \
    --version v0.1.0 \
    --source-repo "$(lab_source)" \
    --staging-dir "$(lab_source)" \
    --target-root "$(lab_targets)"
}

common_publish_args() {
  common_workflow_args
  printf '%s\0' --state-dir "$(lab_state_dir)"
}

run_common_arcpub_json() {
  local report="$1"
  local err="$2"
  local command="$3"
  shift 3
  local args=("$command")
  local common_func=common_workflow_args
  if [[ "$command" == "publish" ]]; then
    common_func=common_publish_args
  fi
  while IFS= read -r -d '' arg; do
    args+=("$arg")
  done < <($common_func)
  args+=("$@" --output json)
  capture_arcpub "$report" "$err" "${args[@]}"
}

configure_git_identity() {
  local dir="$1"
  git -C "$dir" config user.name "ARCORIS Lab"
  git -C "$dir" config user.email "arcoris-lab@example.invalid"
  git -C "$dir" config commit.gpgsign false
  git -C "$dir" config core.autocrlf false
  git -C "$dir" config init.defaultBranch main
}

init_git_repo() {
  local dir="$1"
  git -C "$dir" init -b main >/dev/null 2>&1 || {
    git -C "$dir" init >/dev/null
    git -C "$dir" checkout -B main >/dev/null
  }
  configure_git_identity "$dir"
  git -C "$dir" add .
  git -C "$dir" commit -m "lab: seed fixture" >/dev/null
}

init_bare_repo() {
  local dir="$1"
  mkdir -p "$(dirname "$dir")"
  git init --bare "$dir" >/dev/null
  git --git-dir "$dir" config receive.denyDeleteCurrent ignore
}

init_target_worktree() {
  local worktree="$1"
  local remote="$2"
  mkdir -p "$worktree"
  git -C "$worktree" init -b main >/dev/null 2>&1 || {
    git -C "$worktree" init >/dev/null
    git -C "$worktree" checkout -B main >/dev/null
  }
  configure_git_identity "$worktree"
  printf 'seed\n' >"$worktree/.seed"
  git -C "$worktree" add .seed
  git -C "$worktree" commit -m "lab: seed target" >/dev/null
  git -C "$worktree" remote add origin "$remote"
}

seed_bare_main() {
  local repository="$1"
  local bare
  local worktree
  bare="$(bare_repo "$repository")"
  worktree="$(repository_worktree "$repository")"

  if git --git-dir "$bare" show-ref --verify --quiet refs/heads/main; then
    return 0
  fi

  if [[ -d "$worktree/.git" ]]; then
    git -C "$worktree" push origin HEAD:refs/heads/main >/dev/null 2>&1
  else
    local seed
    seed="$(lab_root)/seed/${repository//\//__}"
    rm -rf "$seed"
    mkdir -p "$seed"
    init_target_worktree "$seed" "$bare"
    git -C "$seed" push origin HEAD:refs/heads/main >/dev/null 2>&1
  fi
  git --git-dir "$bare" symbolic-ref HEAD refs/heads/main
}

target_prepare_remote_template() {
  printf 'file://%s/{name}.git\n' "$(lab_remotes)"
}

write_fixture_source() {
  cp -R "$(lab_repo_root)/internal/testdata/e2e/local-publish/." "$(lab_source)/"
}

write_reject_hook() {
  local bare="$1"
  local mode="$2"
  local hook="$bare/hooks/pre-receive"
  case "$mode" in
    candidate-control)
      cat >"$hook" <<'HOOK'
#!/usr/bin/env sh
set -eu
while read old new ref
do
  case "$ref" in
    refs/heads/arcpub/tx/*/control)
      if [ "$new" != "0000000000000000000000000000000000000000" ]; then
        echo "rejecting candidate $ref" >&2
        exit 1
      fi
      ;;
  esac
done
exit 0
HOOK
      ;;
    promotion-control)
      cat >"$hook" <<'HOOK'
#!/usr/bin/env sh
set -eu
while read old new ref
do
  case "$ref" in
    refs/heads/main)
      if [ "$new" != "0000000000000000000000000000000000000000" ]; then
        echo "rejecting promotion $ref" >&2
        exit 1
      fi
      ;;
  esac
done
exit 0
HOOK
      ;;
    tags)
      cat >"$hook" <<'HOOK'
#!/usr/bin/env sh
set -eu
while read old new ref
do
  case "$ref" in
    refs/tags/*)
      if [ "$new" != "0000000000000000000000000000000000000000" ]; then
        echo "rejecting tag $ref" >&2
        exit 1
      fi
      ;;
  esac
done
exit 0
HOOK
      ;;
    *) fail "unknown hook mode: $mode" ;;
  esac
  chmod +x "$hook"
}

transaction_id_from_json() {
  local file="$1"
  if command -v python3 >/dev/null 2>&1; then
    python3 - "$file" <<'PY'
import json, sys
data = json.load(open(sys.argv[1], encoding="utf-8"))
tx = data.get("publish", {}).get("transaction", {})
print(tx.get("id", ""))
PY
  else
    grep -m1 '"id"[[:space:]]*:' "$file" | sed 's/.*"id"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/'
  fi
}

latest_transaction_id() {
  local tx_dir
  tx_dir="$(lab_state_dir)/transactions"
  [[ -d "$tx_dir" ]] || return 0
  find "$tx_dir" -name 'tx-*.json' -type f -exec basename {} .json \; 2>/dev/null | sort | tail -n 1
}

assert_transaction_status() {
  local file="$1"
  local want="$2"
  assert_contains "$file" "\"status\": \"$want\""
}

assert_no_final_refs() {
  for repo in arcoris/foundation arcoris/control; do
    local bare
    bare="$(bare_repo "$repo")"
    assert_ref_missing "$bare" refs/heads/main
    assert_ref_missing "$bare" refs/tags/v0.1.0
    assert_candidate_refs_absent "$bare"
  done
}

assert_published_repo() {
  local repo="$1"
  local bare
  bare="$(bare_repo "$repo")"
  assert_ref_exists "$bare" refs/heads/main
  assert_ref_exists "$bare" refs/tags/v0.1.0
  assert_candidate_refs_absent "$bare"
  assert_tree_contains "$bare" refs/heads/main go.mod
  assert_tree_contains "$bare" refs/heads/main README.md
  assert_tree_contains "$bare" refs/heads/main contracts/doc.go
  assert_tree_contains "$bare" refs/heads/main .arcoris/provenance.json
  assert_tree_missing "$bare" refs/heads/main secret.txt
  assert_tree_missing "$bare" refs/heads/main private/secret.go
  assert_tree_missing "$bare" refs/heads/main private
}
