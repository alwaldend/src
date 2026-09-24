## Context

See [proposal.md](proposal.md). `ReleaseDeployment` currently contains only
OCI configuration, `deploy` invokes Soras for OCI files, and the CLI always
requires `soras_path`. Local generated releases store payload under `files/`.
That local package layout is distinct from the requested public server layout.

## Goals / Non-Goals

**Goals:** Reuse the existing manifest and deploy command for file publication
and optional website activation over authenticated SSH.

**Non-Goals:** Ansible publication, an upload daemon, rollback commands,
automatic replication, automatic version deletion, or global version-policy
changes. Existing OCI metadata and deployment remain supported.

## Decisions

### Deployment metadata and execution

Add an SSH deployment variant with explicit environment/host, publisher
identity or SSH configuration reference, destination root, public project
name, and optional site archive designation. Reuse release version and file
names from their existing manifest owners. `projects/alwaldend.com` maps to
the public project key `alwaldend.com`; avoid accidentally producing
`projects/projects/alwaldend.com` on the server.

Keep connection values outside artifacts that would expose credentials. Use
standard SSH authentication and host verification, with argument-safe process
execution and bounded timeouts. Do not disable host-key checks. Reuse available
packaged SSH/SFTP facilities after inspecting repository dependencies.

The existing `deploy` command selects the requested environment explicitly.
An SSH-only manifest requires SSH tooling but no Soras executable. A manifest
with OCI deployment still requires Soras. Keep native build/deploy wrappers
aligned with the selected backend and preserve existing OCI behavior.

### Ordinary uploads

Copy each local payload file to an unpublished temporary destination, verify
its expected bytes/checksum, then rename it into
`<root>/projects/<project>/releases/<version>/<filename>`. There is no public
`files/` layer and no archive extraction for ordinary artifacts. A release is
a directory of independent files, not an all-files transaction.

Reject unsafe path components. Repeating an identical upload is idempotent;
different bytes at an existing filename/version produce a conflict, preserving
the published file. Interrupted uploads never appear under their final names.
Report partial success clearly if a release contains several files.

### Site activation and redeployment

Only an explicitly designated site archive is extracted. The archive remains
an ordinary public release file. Initially support `.tar.gz` and `.zip` with
prebuilt HTML/assets at archive root; do not run a website build on the VM.
Treat these format and root conventions as implementation defaults to document.

Validate archive entries and extract into a unique staging directory on the
same filesystem as `<root>/sites/<project>/releases/<version>`. Reject absolute
paths, traversal, link/device entries, and unsupported structure. Do not
follow archive-provided symlinks into other projects or staging paths.
Verify a readable `index.html` and expected content before publishing the
extracted directory and atomically replacing `sites/<project>/current`.

Use a per-project activation lock and bind extraction to the uploaded archive
digest so concurrent or conflicting deployments cannot corrupt a release.
If an already extracted release matches its recorded archive identity,
redeployment only reselects that release via the same deploy command. A failed
upload, extraction, validation, or activation leaves the old link intact.
Selecting an older release is ordinary redeployment, not a rollback command.

Keep staging and extraction bookkeeping outside the public release tree and
the site document root. Failed attempts clean up only their own temporary
paths. Disk-full errors preserve published files and the active site.

## Risks / Trade-offs

- Shell metacharacters and unusual filenames -> validate structural names,
  quote remote arguments safely, and test spaces, Unicode, and punctuation.
- Repeated publication to mutable names such as `head` -> require a concrete
  deployment version for differing content; do not silently overwrite a
  retained release. The website can redeploy any retained concrete version.
- Source packaging differs from public layout -> test exact remote paths and
  release-page links with a real generated manifest.
- Independent targets -> report each target's result without pretending that
  publication across hosts is atomic.

## Migration Plan

Add the metadata variant and regenerate bindings through the owning workflow,
then implement the backend and wire the site consumer. Prepare the behavioral
failure matrix and isolated SSH/HTTP fixture before implementation. Introduce
the new OpenSpec owner's packaging/validation registration in implementation.
Existing OCI consumers require no metadata migration; the site consumer moves
to explicit versioned SSH deployment through its linked plan.
