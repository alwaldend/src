# Dependency extraction validation

Observed on 2026-10-30 in the isolated `t3/extract-third-party-deps` worktree.
The source is security commit `5e4bf7ccb7e7948e12c12ae12b92d048bbba06e3`,
extracted onto master `a9ac58edab4a0f98aad8b2e87b442c6d248ad887`.
Task evidence is under ignored `out/dependency-extraction-20261030/`.

## Source and module boundaries

The Lynis BUILD and module fragment are byte-identical to the source commit.
The README only drops the link to the consumer role excluded from this pull
request. The new owner record describes extraction without a behavior change.

`bazel_agent bazel run //tools/bazel_module:update` preserves exactly the two
scanner includes from the source commit. Running
`bazel_agent bazel mod deps --lockfile_mode=update` produces no lockfile change.
The companion ClamAV dependency retains its own OpenSpec record. No consumer
role, scan configuration, host change, or live operation is included here.

## Artifact verification

`bazel_agent bazel build //third_party/com_github_cisofy_lynis:all` succeeds.
The archive resolved through `cquery --output=files` is 355590 bytes and has
SHA-256 `b5314a07fd85fa3ffc7da57b508f0108ec3280d84e4af823f805d95cbbc2428c`,
which matches the unchanged SRI pin. Extraction retains `lynis/lynis`,
`lynis/default.prf`, `lynis/db`, and `lynis/LICENSE`.

From the extracted `lynis/` directory, `./lynis --version` exits successfully,
prints `3.1.7`, and has empty stderr. `lynis-version-cwd.txt` and
`lynis-version-cwd.stderr` retain this result. Use this early-exit option:
`show version` attempts to create default log, report, and PID files. The
initial invocation's writes were blocked by the filesystem sandbox and are
retained in `artifact-verification.txt`. Invoking from outside the extracted
directory also fails include discovery; `lynis-version.txt` records that
diagnosis. No audit was requested or performed.

Strict validation of `extract-lynis-dependency` passes through
`OPENSPEC_PROJECT=third_party/com_github_cisofy_lynis bazel_agent bazel run
//tools/openspec -- validate extract-lynis-dependency --strict`.
Delivery receipts retain the exact committed candidate and mandatory gate
results separately from these source and artifact observations.
