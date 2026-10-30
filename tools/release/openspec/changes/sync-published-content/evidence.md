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
