# Transactional Publish Local Lab

This lab is a controlled, no-network way to inspect transactional `arcpub publish`
behavior before using real repositories. It creates only local temporary Git
repositories: one source repository, two target worktrees, and two local bare
remotes.

The publish transaction is saga-style, not ACID. `arcpub` stages candidate refs,
promotes final branches, pushes tags, and records a durable journal. If a phase
fails, it runs compensating rollback. If rollback safety checks cannot prove a
ref still belongs to the transaction, the journal remains for manual recovery.

## Requirements

- Bash
- Git
- Go
- No network access is required

The default runtime directory is:

```sh
${TMPDIR:-/tmp}/arcpub-transaction-lab
```

Override it with:

```sh
ARCPUB_LAB_ROOT=/custom/path bash test/lab/transactional-publish/setup.sh
```

Do not point `ARCPUB_LAB_ROOT` at important repositories.

## What It Creates

```text
$ARCPUB_LAB_ROOT/
  bin/arcpub
  source/
  targets/
    arcoris__foundation/
    arcoris__control/
    .arcpub/state/
  remotes/
    foundation.git
    control.git
  reports/
  logs/
```

The fixture publishes only explicit entries. `secret.txt` and `private/` files
must never appear in pushed target trees. Provenance is enabled at
`.arcoris/provenance.json`.

## Target Preparation

`run-target-prepare.sh` demonstrates `arcpub target prepare`. The script seeds
the local bare remotes, removes the target worktrees, then lets `arcpub` clone,
fetch, checkout, and validate them:

```sh
bash test/lab/transactional-publish/setup.sh
bash test/lab/transactional-publish/run-target-prepare.sh
```

The remote template used by the lab is:

```text
file://$ARCPUB_LAB_ROOT/remotes/{name}.git
```

The command supports `{repository}` for the full logical repository ref,
`{owner}` for the first component, and `{name}` for the final component.
`target prepare` is not read-only: it may create local target worktrees, add
`origin`, fetch, and checkout branches. It does not construct publication files,
rewrite `go.mod`, write provenance, create journals or locks, commit, tag, or
push.

## Happy Path

```sh
bash test/lab/transactional-publish/setup.sh
bash test/lab/transactional-publish/run-plan.sh
bash test/lab/transactional-publish/run-target-prepare.sh
bash test/lab/transactional-publish/run-preflight.sh
bash test/lab/transactional-publish/run-verify.sh
bash test/lab/transactional-publish/run-dry-run.sh
bash test/lab/transactional-publish/run-happy-path.sh
bash test/lab/transactional-publish/show-transactions.sh
```

`run-happy-path.sh` recreates the lab first so dirty worktrees from verify or
dry-run cannot affect publish preflight.

`run-preflight.sh` is the read-only safety check. It does not construct target
files, rewrite `go.mod`, create journals, create locks, or push refs. `verify`
constructs and validates target trees. `publish --dry-run` also constructs and
verifies target trees, then skips commit, tag, and push. `publish` performs the
transactional mutation.

`target prepare` sits before preflight in an operator workflow. It prepares the
local Git worktrees that preflight and publish inspect, but it does not publish
anything.

## Failure Scenarios

Each failure script recreates the lab and installs a local bare-repo hook:

```sh
bash test/lab/transactional-publish/run-candidate-failure.sh
bash test/lab/transactional-publish/run-promotion-failure.sh
bash test/lab/transactional-publish/run-tag-failure.sh
bash test/lab/transactional-publish/run-lock-conflict.sh
bash test/lab/transactional-publish/run-corrupted-journal.sh
```

These scenarios demonstrate rollback after:

- candidate ref push rejection;
- final branch promotion rejection;
- tag push rejection.

They also demonstrate that an existing publish lock or a corrupted transaction
journal blocks a new publish before remote mutation.

The scripts save JSON reports, inspect transaction status, verify candidate refs
are cleaned, verify final refs are absent after rollback, and run rollback twice
to demonstrate idempotency.

## Inspect And Recover

```sh
bash test/lab/transactional-publish/show-transactions.sh
bash test/lab/transactional-publish/show-transactions.sh <transaction-id>
bash test/lab/transactional-publish/run-rollback.sh <transaction-id>
bash test/lab/transactional-publish/inspect-remotes.sh
bash test/lab/transactional-publish/inspect-worktrees.sh
```

Reports are stored under `$ARCPUB_LAB_ROOT/reports`. Journals are stored under:

```text
$ARCPUB_LAB_ROOT/targets/.arcpub/state/transactions
```

`rollback_failed` means automatic compensation could not safely complete. Inspect
the transaction report for manual recovery actions before starting another
publish.

## Cleanup

```sh
bash test/lab/transactional-publish/cleanup.sh
```

Cleanup refuses to remove `/`, the repository root, or the user home directory.
