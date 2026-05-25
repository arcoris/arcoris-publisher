#!/usr/bin/env bash
set -euo pipefail

LAB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$LAB_DIR/common.sh"

require_cmd go
run_step "building cmd/arcpub"
mkdir -p "$(lab_bin)"
(cd "$(lab_repo_root)" && go build -o "$(lab_arcpub)" ./cmd/arcpub)
log "binary: $(lab_arcpub)"
