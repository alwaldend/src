## Context

See `proposal.md`. The live bazelrc currently matches its Ansible owner,
including the existing shared repository-contents cache. Go module-cache
policy must be selected from observed reuse, not inferred from its name.

## Decisions and authority

The user requested: "make another pr updating them in the role and patch the
host config on this host", then selected Gazelle 0.53.0 and ephemeral caching.
The subsequent question about reliable reuse was followed by: "Is there really
no better way? If so, just use it". This authorizes the shared Go module-cache
fallback if completed-repository reuse proves insufficient.

Verify reuse before choosing the final host policy. For the live bazelrc,
apply only the selected cache setting after checking for concurrent changes;
retain ownership and mode. This records the bounded imperative patch requested
by the user and its reconciliation into the owning Ansible file. Provision
any required shared module-cache directory in the same role with matching
live ownership and permissions. No storage or unrelated host deployment is
authorized.

## Risks / Trade-offs

Ephemeral caching saves disk space but redownloads modules on a repository
cache miss. The 0.53.0 default retains ZIPs but prunes extracted trees; that
cache remains scoped to an output base. A stable host module cache reuses
both across worktrees and survives Gazelle upgrades, at the cost of retained
disk space and corresponding managed write permissions.

## Migration and verification

Validate packaged desired state before touching live configuration. Preserve
a private task-local backup of the live bazelrc, verify the exact intended
diff, and compare the resulting file with its checked-in owner. If necessary,
validate any Codex permission change in an isolated opt-in profile before
changing a shared default. Cache experiments and receipts remain under ignored
task scratch.

## Selected policy and live reconciliation

The companion repository change's fresh-worktree experiment verified completed
Go repository reuse, so the fallback shared host Go module cache is unnecessary.
The selected policy remains `GO_REPOSITORY_EPHEMERAL_MODCACHE=1`.

The packaged Ansible manifest maps `files/bazelrc` to the exact owning source.
The bounded live patch added the same setting and comments after confirming
the previous live file matched its checked-in baseline. The live file matches
the new source byte-for-byte, retaining uid/gid 1000 and mode 0600. The shared
cache paths already existed, so this operation changed only the user bazelrc.
Detailed evidence is retained in task-local `package-verification.json` and
`live-patch-verification.json`.
