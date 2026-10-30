---
title: Release
description: Release rules
languages:
  - bzl
  - go
tags:
  - proto
  - releases
  - bzl_rules
---

The release tool generates release manifests/pages and deploys their files.
OCI deployment uses Soras. SSH deployment uses the operator's OpenSSH and
rsync executables and requires rsync on the destination; it needs neither
Soras nor Ansible. No custom receiver or upload service is installed.

## Build a portable release bundle

`al_release_bundle` packages declared payloads as `release.json` and
`files/<filename>` in a Bazel directory output. It reuses the manifest
generator for file sizes and hashes and reads `STABLE_VERSION` from the
versioning owner's workspace status. Run its build through
`tools/versioning/cmd/versioning/versioning.sh bazel -- build --config=release`.
It does not contact deployment destinations or calculate another version.

The CLI equivalent is `generate --project projects/<name> --version_file
<status-file> --add_file <payload> --output_dir <empty-directory>`. Repeat
`--add_file` for multiple payloads. Missing or duplicate version values,
unsafe path components, missing payloads, colliding filenames, and empty
bundles fail. The output directory must be absent or empty. A completed
manifest is written only after payload copying finishes.

The website's `//projects/alwaldend.com:release` target consumes this rule;
its owner documents the complete build and SSH publication commands.

## Publish over SSH

A release directory contains `release.json` and `files/<filename>`. The
manifest owns the version, project, and payload names. Pass an explicit
SSH environment and destination, or a deployment JSON containing an `ssh`
object. The download component provides preconfigured targets. These default
to the existing `ansible` administrator and use `sudo -n -u download` for
content operations. Host configuration from merged
[PR #102](https://github.com/alwaldend/src/pull/102) provides the content account,
storage, and server rsync package. Configure the selected host before publishing.

```sh
bazel_agent bazel run //infra/download:publish.local -- \
  --release_dir "$PWD/out/release"

bazel_agent bazel run //infra/download:publish.yandex -- \
  --release_dir "$PWD/out/site-release" --site_archive website.tar.gz \
  --steps upload,extract,link
```

For other SSH destinations:

```sh
bazel_agent bazel run //tools/release -- deploy \
  --release_dir "$PWD/out/release" --environment local \
  --ssh_host download-host --ssh_user admin --publish_user download \
  --ssh_root /srv/download
```

`--ssh_deployment` or `--ssh_host` publishes all manifest files through that
SSH destination, regardless of their original OCI metadata. Otherwise,
`--environment` selects matching SSH metadata attached to individual items.
Omitting the environment selects existing OCI deployments and requires
`--soras_path`. An empty selection fails.

Host trust and authentication use normal OpenSSH configuration, with strict
host-key checks and noninteractive authentication. `--ssh_config`,
`--ssh_identity`, `--ssh_certificate`, `--ssh_known_hosts`, `--ssh_port`, and `--ssh_user` allow
explicit local overrides. Credentials are not written into manifests.
`--ssh_path` and `--rsync_path` select installed executables. The default
deployment deadline is 30 minutes; `--timeout` overrides it.

The public project defaults to the manifest project after removing a leading
`projects/`; `--project` overrides it. The tool validates each path component
and uses the remote layout
`<root>/projects/<project>/releases/<version>/<filename>`.
Normal publication rejects an existing artifact or extracted-site version.
It transfers the complete release into private staging, normalizes directories
to `0755` and files to `0644`, then promotes the directory under the project
lock. The collision check runs before transfer and again under the lock;
concurrent uploads cannot overwrite a completed version. A failed transfer
does not expose a partial release. Use the download component’s Ansible mirror for authoritative public replacement.

Remote shell arguments use POSIX single-quote escaping to preserve literal values.
SSH publication starts remote commands from `/`, so the content account need
not traverse the administrator's private home. Uploaded archive directories
are set to 0755 before promotion, independently of the remote umask. Repeating
`link` also repairs the selected archive directory's mode without changing bytes.

## Choose publication steps

`--steps` selects an ordered subset of `upload,extract,link`; it defaults to
`upload`. Each step can run alone. Unknown, repeated, empty, or out-of-order
steps fail before a remote change. A release manifest still supplies the
project and version for every step.

```sh
# Upload files only; archives remain ordinary downloadable files.
bazel_agent bazel run //infra/download:publish.yandex -- \
  --release_dir "$PWD/out/site-release" --steps upload

# Extract an already-published archive without changing current.
bazel_agent bazel run //infra/download:publish.yandex -- \
  --release_dir "$PWD/out/site-release" --steps extract \
  --site_archive website.tar.gz

# Select an existing extracted release; no local payload or archive is needed.
bazel_agent bazel run //infra/download:publish.yandex -- \
  --release_dir "$PWD/out/site-release" --steps link
```

Extraction accepts a published `.tar.gz` or `.zip` filename. It retrieves the
archive into a private local temporary directory, validates and extracts it,
then copies the site into private remote staging and promotes the completed
directory to `sites/<project>/releases/<version>`. Local free space must hold
both the archive and its extracted content. No archive helper is installed on
the host. Extracting an existing site version fails without changing it.

The extractor rejects traversal, absolute paths, links, device entries,
duplicate file names, and missing root `index.html`. Archives may contain at
most 100,000 entries. Archive ownership, executable bits, and extended
attributes are not preserved. Server-created content inherits the destination's
configured SELinux labels.

Linking atomically replaces `sites/<project>/current` after checking the
selected directory, its root index, and managed directory ancestry. It then
deletes other extracted versions of that same website, under the project lock.
All downloadable archives remain in `projects/<project>/releases/`. Cleanup
does not follow symlinks or cross nested filesystems. A cleanup failure reports
failure after selection; the new site remains selected and repeating `link`
can finish cleanup. Serialize complete publication workflows for one project:
linking also removes versions extracted by another unfinished workflow.

A failure before selection leaves the current site unchanged, although a
completed upload may precede failed extraction. Upload-only and extract-only
operations do not delete sites. Use a new version for changed website content.
To roll back, use the earlier version's manifest and `--steps extract,link`
to restore its retained archive; use only `link` if it is already extracted.
There is no separate rollback command.

`al_release_deployment` accepts `ssh_*` metadata alongside its existing OCI
variant. `ssh_public_url` supplies generated artifact URLs with the public
project/version layout. `al_release_binary` can package a fixed SSH destination
with `oras = None`, while taking `--release_dir` at runtime.

## Infrastructure synchronization

The download component owns local-to-public store synchronization through its
Ansible playbook and `//infra/download:sync.yandex` target. The release CLI has
no sync command. See the [download runbook](../../infra/download/README.md)
for check-mode preview, authoritative replacement, and failure boundaries.

The isolated `//infra/download/test/sync:sync_test` acceptance fixture
exercises this CLI's publication, duplicate-version refusal, archive-preserving
site cleanup, and the Ansible mirror against a real loopback OpenSSH server.
It never uses deployment hosts. Its repeatable receipt records content hashes,
permissions, links, and the checked failure cases.
