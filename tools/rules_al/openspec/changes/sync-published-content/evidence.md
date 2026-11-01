# Published-content synchronization acceptance

The 2026-10-26 isolated fixture invokes the built release CLI against real
loopback OpenSSH and rsync, with an ephemeral signed client key. It reads and
writes only task-owned directories. The synchronization fixture does not touch
deployment hosts, Vault credentials, DNS records, or production content.
The separate public Git redirect change owns its authorized live deployment.

The fixture was written before implementation and first failed because the
`sync` command did not exist. The user's later clarification made local
authoritative; the earlier retention/conflict behavior was replaced and its
results are not acceptance evidence for the final semantics.

The source-authoritative fixture passes 16 cases: dry-run deletion reporting,
required-mount refusal, SSH host-trust refusal, source symlink rejection,
escaping current-link rejection, missing site index, initial mirror, unchanged
repeat, ordinary-file replacement, differing same-version site replacement,
failed staged write, expired operation deadline, destination symlink rejection,
removal of an absent source selection, empty-source mirroring, and missing
source-root refusal. It also checks destination-only release removal, unchanged
private state, payload bytes, permissions, unselected site versions, and
controller temporary-directory cleanup.

The destination root is non-writable to model the deployed administrator-owned
mount root. Only its projects/sites/staging directories are writable. This
disconfirmed swapping the top-level trees; the implementation exchanges project
directories within those existing publisher-owned parents instead. No host
permission change or additional privilege is required.

For the staged-write failure, `prlimit` restricts only the fixture sshd's child
processes to 4096-byte file writes while a larger artifact is transferred.
Rsync fails in private staging, and the complete public content hash/mode/link
inventory remains unchanged. The fixture restores its daemon's limit and
removes its temporary process and credentials afterward.

Repeat with:

```sh
bazel_agent bazel test //tools/release/main/go:sync_e2e_test //tools/release/main/go:go_test
bazel_agent bazel build //infra/download:sync.yandex
```

`bazel-testlogs/tools/release/main/go/sync_e2e_test/test.outputs/sync-receipt.json`
records all case names and representative content inventories. The generated
component launcher was inspected: it selects local as source, Yandex as
destination, `ssh=1` injection, strict shared SSH host trust, and required
mountpoints. Only credential file references appear in its arguments.

Limits: the fixture does not establish live Vault authentication, SELinux,
HTTP serving during exchanges, Btrfs capacity, or crash recovery. Source and
destination publishers must be idle; there is no global snapshot or transaction
across projects. The public command has not been run. Its dry-run output must
be reviewed before an operator authorizes production replacements and removals.

Validation receipts and command logs are retained privately under
`out/download-sync/`; final publication uses the repository delivery workflow.

## Rebase and cleanup review, 2026-11-01

The user approved preserving trunk's third-party replacement uploads and
explicit SSH-agent control during PR #202's rebase. Project publication retains
complete private staging, collision refusal and archive-preserving retention.
The namespace CLI fixture verifies replacement in third_party, refusal in
projects, and existing agent-enabled/disabled authentication.

The prior cleanup fragment deleted a mounted sentinel before reporting a busy
mount point. Linking now rejects mounted ancestry and directory mounts in the
release tree before selection or cleanup. The real SSH acceptance fixture binds
private state onto an obsolete release and verifies refusal, unchanged state
and unchanged selection, then unmounts and verifies ordinary cleanup. The
fixture runs inside a private user/mount namespace; only its controller retains
mount capabilities. The SSH server drops inherited/ambient capabilities so
restrictive-permission checks remain meaningful.

The moved release target and shared DOWNLOAD_SSH_* injector references are
used by the fixture. The passing focused commands are
`bazel_agent bazel test //tools/release/cmd/release:go_test` and
`bazel_agent bazel test //infra/download/test/sync:sync_test`. The latter
completed in 173.788 seconds. The operations render contains atomic-publication
and archive-based recovery labels; the DNS runbook matches the configured apex
CNAME. Exact-candidate validation/publication receipts are under
`out/rebase-review/`. No live infrastructure operation is part of this rebase.

The broad candidate check passed 51 of 52 tests; its sole failure was the
dedicated workspace Buildifier test on three already-merged BUILD files. The
configured formatter normalizes their layout/order without changing target
attributes, dependencies or pins. The root documentation/packaging change
records that required format-only reconciliation. Passing final gates are
recorded against the subsequent candidate by the delivery tool.
