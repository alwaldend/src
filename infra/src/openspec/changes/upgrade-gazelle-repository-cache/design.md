## Context

The host has shared Bazel action and repository-contents caches. Gazelle
0.47.0 keeps its Go download cache beneath each output base; the approved
0.53.0 release supports an optional private module cache for each fetch.

## Decisions

Update the existing pin in `third_party/include.MODULE.bazel`, then reconcile
the root module lock with Bazel. Use 0.53.0 because the user selected that
exact version. Keep independent nested-module pins under their owners.

Measure completed-repository reuse from fresh output bases with the existing
shared cache. Ephemeral mode is a disk-space policy, not a cache-sharing
mechanism. If the shared repository cache cannot reliably reuse Go repositories,
the user's follow-up authorizes a shared host Go module cache instead; keep the
host-owned decision and configuration in its companion change.

## Risks / Trade-offs

An upgraded Gazelle can change generated dependency BUILD files. Validate
representative consumers and the repository's Gazelle binary. Review all
module-lock resolution changes before delivery.

Fresh output bases can miss the repository cache because rule inputs include
output-base-dependent environment files. Preserve fetch diagnostics and
inspect actual cached repositories rather than treating command success or
elapsed time as proof of reuse.

## Verification

Use isolated task-owned worktrees, the same dependency declarations, and
repository cache hit evidence. Keep raw logs under `out/gazelle-cache`.

On 2026-10-13, equivalent root-workspace declarations in fresh linked
worktrees reused the same completed trees for `github.com/google/uuid`,
`github.com/BurntSushi/toml`, and `golang.org/x/net`. The source and
generated BUILD files matched; each external repository linked into the
configured shared cache. The first UUID fetch appeared in the workspace-rule
log, while its subsequent fetch execution was absent. A third fresh output
base loaded the host setting automatically and reused all three repositories
without their fetch execution arguments in its log. No temporary module-cache
directory remained inside the completed repositories.

The Gazelle binary, Ansible package, and module generator built successfully.
The lock reconciliation changed only Gazelle's registry metadata. Source,
cache-path, generated-output, and workspace-log evidence is retained in
`completed-repository-reuse.json` and `live-reuse-verification.json` beneath
task scratch. These observations cover matching pinned inputs, not arbitrary
repository cache misses.
