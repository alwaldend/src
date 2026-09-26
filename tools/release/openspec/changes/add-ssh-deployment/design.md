## Context

The release tool already generates manifests/pages and deploys OCI files.
The user requested SSH support plus `infra/download` publishing targets, then
explicitly selected ordinary rsync or rclone instead of a custom upload
protocol. The implementation uses rsync over OpenSSH.

## Goals / Non-Goals

Reuse the existing `deploy` command, manifest version and filenames. Keep
publication independent of Ansible. Preserve the download account as content
owner without configuring authorized keys. Do not add a custom receiver,
upload service, rollback command, replication, or automatic deletion.

## Decisions

### Authentication and transfer

Use administrator SSH and `sudo -n -u download rsync` on download hosts.
OpenSSH retains authentication and strict host trust. Connection settings are
public metadata; identity/known-hosts file paths are runtime options only.
The download role installs native rsync alongside the filesystem tools.

Select `local` or `yandex` explicitly. `infra/download:publish.<environment>`
packages a destination and accepts a release directory. The public project
removes one leading `projects/` from the manifest project or uses an explicit
override. Ordinary files land directly in the version directory. Rsync owns
checksums, temporary files, retry/update behavior, and error reporting. It
updates existing filenames without deleting other release files. This replaces
the earlier custom immutable-file/conflict protocol at the user's request.

### Independent publication steps

The user confirmed independent upload, extraction, and linking during review.
`--steps` is an ordered subset of `upload,extract,link`, defaulting to upload.
The deploy method creates a publisher and dispatches each selected operation;
configuration and validation stay separate from the transfer operations.

Upload copies manifest files using rsync and its normal temporary-file rename.
Extraction reads the already-published archive, reuses the rooted Go extractor,
and transfers the completed tree into private remote staging. Promotion keeps
an existing completed version unchanged. Linking checks the retained site and
atomically replaces current, without an archive or any rsync operation.

The extraction decision is to reuse the existing validated Go implementation.
Running remote tar/unzip would avoid transferring the archive back, but would
introduce host dependencies and a different archive-safety contract. Reusing
the extractor preserves traversal/link/device checks and needs no installed
receiver. The cost is local temporary space and an extra archive transfer.
Reconsider that choice if measured archive sizes make it a material bottleneck.

A short per-project flock protects promotion and linking. Staging uses the
content filesystem and retains its noexec mount. Invalid step sequences fail
before remote mutation; upload failures can leave completed ordinary files,
while extraction failures preserve the currently selected site.

### Quoting and error context

Preserve literal remote arguments and rsync's SSH command with the existing
POSIX single-quote escaping approach. The existing pinned anmitsu/go-shlex
package only tokenizes strings; its missing quoting API does not justify adding
another dependency without approval. The user rejected the added quoting
library and required an explicit approval gate in AGENTS.md. No new external
dependency is needed for publication. Propagated errors wrap their causes with
operation and path context; repo-go owns that repository procedure.

### Build and backend selection

The protobuf and Bazel deployment rule represent OCI or SSH. SSH-only wrapper
targets omit Soras. A supplied SSH deployment/host publishes all manifest
files to that destination; otherwise an explicit environment selects attached
SSH metadata. Existing OCI operation remains the default without an environment.
Generated download URLs use the same public project and version mapping.

## Risks / Trade-offs

Rsync and OpenSSH must be installed locally and rsync remotely. Website
extraction also needs local free space. Multi-file rsync is not a transaction;
completed ordinary files can remain after failure. Website activation waits
for the complete site transfer. Administrators retain their existing sudo
capability; the content process runs as the dedicated non-sudo account.

## Acceptance

See [acceptance.md](acceptance.md) for the failure matrix and retained evidence.
Tests use an isolated SSH/Nginx environment, not either production VM.
