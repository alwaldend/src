---
title: Rules AL
description: Repository lifecycle commands and reusable Bazel rules
languages:
  - bzl
  - go
tags:
  - releases
---

`rules_al` is the standalone Bzlmod module for AL, the repository lifecycle executable.
Its direct commands are `al config`, `al run`, `al version`, `al release`,
`al delivery`, `al bazel`, and `al doctor`.

The executable owns signals and startup. `internal/cli` assembles command adapters;
local release, delivery, versioning, Git, and Bazel packages retain their domain
operations and guards. Reusable configuration, plugin and resource-lifetime APIs
live in `pkg/`. CLI entry points use the same owning implementations as Bazel
execution tools. Credential plugins start through `al run`; ordinary lifecycle
inspection and artifact generation do not start them.

AL's public Bazel API is `@rules_al//:defs.bzl`, including `al_config`, AL run
wrappers, release packaging rules, and their providers. The module supplies its
own AL execution toolchain. Existing pinned dependencies are declared locally.

Build the current executable with `al bazel build @rules_al//:al`. Bootstrap a
checkout without AL using `bazel run --config=agent //tools/al:install`; installation
requires the operator's authority. Release stamping can bootstrap a source-current
launcher with `tools/rules_al/cmd/al/al.sh version bazel -- <command>`.

Operator procedures live in the module's `al-*` skills. Workspace preparation,
checks, maintenance, and OpenSpec workflows use existing targets through those
skills; they do not add speculative command groups.

AL also owns repository versions and independent project releases.
The release identity is `(project, version)`: projects share a source version
and process, while each publishes and activates independently. Operators can
run the entire process without CI.

## Bazel module API

Consumers depend on `rules_al` and load its unprefixed rules from
`@rules_al//:defs.bzl`: `release`, `release_files`, `release_bundle`,
`release_binary`, `release_deployment`, and `release_deps`. Providers are
`ReleaseInfo`, `ReleaseFilesInfo`, and `ReleaseDeploymentInfo`.

```starlark
load("@rules_al//:defs.bzl", "release_bundle")

release_bundle(
    name = "release",
    srcs = [":archive"],
    project = "projects/example",
)
```

`status_file` can supply an explicit stable-status file. Otherwise the bundle
consumes Bazel workspace status when stamping is enabled, and ordinary unstamped
builds emit the constant development version. The standard `stamp` attribute
controls this behavior. The module builds its own tool and declares
all software dependencies at repository pins. Run `al bazel test //...`
and `al bazel build //...` from `tools/rules_al` to check it independently.
The `tests/consumer:check` command builds the public API without parent Bazel
configuration imports and automatically verifies its output manifest and payload.
Its README describes launcher invocation and the retained `result.json`.

The former generic Git tool and schema are merged here. `git-info` generates
repository metadata using the same parser and Git schema as release history.
The exported `git` extension configures a consumer's local checkout and system
Git executable; register the resulting `git_toolchain` to use `@rules_al//:git`,
`git_binary`, or `git_info_file`. History generation deliberately reads that
checkout locally and is excluded from remote execution/cache upload. Ordinary
portable bundle generation uses only declared inputs.

```starlark
git = use_extension("@rules_al//:extensions.bzl", "git")
git.local_git(name = "repository_git")
use_repo(git, "repository_git")
register_toolchains("@repository_git//:git_toolchain")
```

The `release` documentation rule accepts optional `git_bundle` and `page_template`
inputs. This repository's website supplies its history bundle and shared Hugo
markup through `//projects/alwaldend.com/pkg/release_pages:defs.bzl`.
No published module target refers back to the parent repository.

## Prepare source

The calendar-version contract remains:

- development: `0.0.0-dev`;
- nightly trunk tag: `vYYYY.W.0-nightly.YYYYMMDD`;
- weekly branch: `releases/YYYY.W`;
- release tag: `vYYYY.W.PATCH`, with patch zero at the branch point and one
  increment for every subsequent first-parent commit.

The year and week use the ISO calendar; the week is not zero-padded. Inspect
source with `al version show --format json`.
Version operations use native Cobra subcommands: `al help version show` describes
their flags, and shared flags work before or after the subcommand. Select another
checkout with `--repository PATH`; `--repo` remains a compatibility alias.
SIGINT and SIGTERM cancel active Git and stamped Bazel subprocesses. Foreground
Git authentication can use the controlling terminal; AL restores terminal ownership
after the child exits.
Use `version nightly-tag`, `version release-start`, or `version release-tag`
for explicitly authorized local ref preparation. Each supports `--dry-run`.
The `al-versioning` skill documents detached-checkout and calendar details.

For a weekly release, start on clean trunk and run `version release-start`.
Review and publish the shared source refs:

```sh
al version release-plan \
  --lease weekly-release --plan out/release/source-plan.json
al version release-publish \
  --plan out/release/source-plan.json --receipt out/release/source-receipt.json
```

Source publication uses a Git-common-directory operation lock, guarded creates,
atomic push, and remote verification. Existing tags are reused only at the
same commit. For a patch, commit and tag on the weekly branch, deliver the
branch normally, then publish its tag; trunk may have diverged. Shared tags
are not project publication status. A stale `release-refs.lock` directory after
a killed process must be investigated before an operator removes it.

## Build one project

```sh
tools/rules_al/cmd/al/al.sh version bazel -- \
  build --config=release //projects/autoscroll/releases:bundle
```

Projects with declared first-party payloads have a `releases:bundle` target.
The website's bundle aliases `//projects/alwaldend.com:release`. The Android,
kustomization, and Terraform module release pages currently declare no payloads;
they provide documentation until their owners define release artifacts.
Existing `head` release pages are development documentation, not publishable
release identities. Nightly-to-weekly promotion rebuilds version-bearing
artifacts under the weekly identity.

`release_bundle` packages declared payloads as `release.json` and
`files/<filename>` in a Bazel directory output. It consumes `STABLE_VERSION`,
`STABLE_VERSION_CHANNEL`, `STABLE_GIT_COMMIT`, and `STABLE_GIT_TREE_STATE` from
the release bootstrap. First-party release bundles require a concrete version
and clean source identity. Development bundles can be inspected but cannot
be published. Third-party versions continue to come from their dependency owner.
The build does not contact deployment destinations.

The CLI equivalent is `al release generate --project projects/<name> --version_file
<status-file> --add_file <payload> --output_dir <empty-directory>`. Repeat
`--add_file` for multiple payloads. Missing or duplicate status values, unsafe
paths, missing payloads, colliding filenames, and empty bundles fail. The
output directory must be absent or empty. A completed manifest is written
only after payload copying finishes.

## Publish over SSH

SSH publication uses the operator's OpenSSH and rsync executables and requires
rsync and standard GNU filesystem utilities (including `flock` and `mountpoint`)
on the destination; it needs neither Soras nor Ansible. No receiver
service is installed. Select the environment and destination explicitly:

```sh
al bazel run @rules_al//:al -- release publish \
  --release_dir "$PWD/bazel-bin/projects/autoscroll/releases/bundle" \
  --environment local --ssh_host download-host --ssh_user admin \
  --publish_user download --ssh_root /srv/download
```

The download component provides preconfigured `//infra/download:publish.local`
and `//infra/download:publish.yandex` commands. These use the configured SSH
administrator and content account; provision the selected host first.

`publish` is an alias for `deploy`. `--ssh_deployment` accepts a deployment JSON
with an `ssh` object. Explicit host or deployment selection publishes every
manifest item; otherwise `--environment` selects matching SSH item metadata.
OCI manifest/page generation remains available, but OCI publication fails
closed: the former mutable-tag push cannot enforce immutable releases. Select
an SSH destination for publication. An empty selection fails.

Host trust and authentication use normal OpenSSH configuration with strict
host-key checks and noninteractive authentication. `--ssh_config`,
`--ssh_identity`, `--ssh_certificate`, `--ssh_disable_agent`, `--ssh_known_hosts`,
`--ssh_port`, and `--ssh_user` provide local overrides. Credentials never enter
manifests. Remote commands start from `/` so the content account does not need
to traverse the administrator's private home. `--ssh_path` and `--rsync_path` select installed executables; the
default deadline is 30 minutes (`--timeout`).

The public project defaults to the manifest project without `projects/`;
`--project` overrides it. `--namespace third_party` selects external artifacts.
The public layout is `<root>/<namespace>/<project>/releases/<version>/<filename>`.
Publication copies the complete selected set to private staging, verifies sizes
and SHA-256 hashes locally and remotely, then atomically renames the completed
release under a destination lock. A canonical public `release.json` binds the
identity and payload set without transport configuration. Directories are
`0755`, payloads `0644`. Shell arguments preserve literal values.

Identical existing releases are verified under the destination lock and succeed
without local payloads, staging, or another transfer. New publications retain a
final locked check after transfer to handle concurrent publishers. Conflicting manifests,
extra files, and corrupt payloads fail. Failed transfers expose no partial new
release. Directories from the former overwrite workflow must be explicitly
deleted before publication through this contract.

Use `--receipt /absolute/path/to/attempt.json` with publish/deploy or delete to
checkpoint local operation evidence. The private JSON receipt binds the selected
namespace, project, version, source commit, canonical publication digest and resolved
SSH environment/host/root. It records each requested step as `pending`, `unknown`
when started, and `completed` after confirmation. A failed operation preserves
completed steps and leaves an unconfirmed step unknown; even a transport error can
occur after a remote change. `status` is `running`, `completed`, or `incomplete`.
A killed process may leave `running` evidence with an unknown step. Use a separate
receipt path for each attempt to preserve history. Receipt persistence must succeed
before remote work starts. Receipts are observations, not authorization or proof
of current remote state; repeat the appropriate idempotent step to verify recovery.

## Extract and activate a website

`--steps` selects an ordered subset of `upload,extract,link`, defaulting to
`upload`. Every step can run independently. Unknown, repeated, empty, or
out-of-order steps fail. External artifacts support upload only.

```sh
al bazel run //infra/download:publish.yandex -- \
  --release_dir "$PWD/out/site-release" --steps upload,extract,link \
  --site_archive website.tar.gz
```

Extraction accepts a declared `.tar.gz` or `.zip` archive containing a root
`index.html`. It retrieves and verifies the published archive, extracts locally,
and promotes the site to `sites/<project>/releases/<version>`. Local free space
must hold both archive and extracted content. Traversal, absolute paths, links,
devices, duplicate names, and more than 100,000 entries fail. Archive ownership,
executable bits, and extended attributes are not retained.

Retained sites carry a binding to the exact publication. Linking checks this
binding and the public payloads under the destination lock, then atomically
selects `sites/<project>/current`. Managed ancestry must use real directories without nested mounts. Activation
also rejects mounts within the selected extracted site. Failed earlier steps leave activation
unchanged. Roll back by selecting another retained release with `--steps link`.

## Delete and republish

Use the same manifest and explicit SSH destination with `delete`:

```sh
al bazel run @rules_al//:al -- release delete \
  --release_dir "$PWD/out/site-release" --environment local \
  --ssh_host download-host --ssh_root /srv/download
```

Deletion removes the project release and its extracted site at that destination.
It refuses the currently selected website: first select another retained release
with `--steps link`. Deletion refuses symlinked ancestry and mounted target directories or descendants, including bind-mounted files,
before either tree is removed. Deleting an absent release succeeds. Shared source tags
remain. No tombstone or restoration manifest is retained; ordinary `publish`
works again after deletion, including with different artifact bytes.

## Infrastructure synchronization

The download component owns the explicit local-to-public Ansible mirror through
`//infra/download:sync.yandex`. AL publication preserves immutable identities; the
separate mirror replaces public content from the authoritative local store and
removes public-only content. It preserves locally retained archives and sites.
See the [download runbook](../../infra/download/README.md) for check-mode preview,
serialization, staging, and failure boundaries. AL has no synchronization command.

The manual `//infra/download/test/sync:sync_test` fixture uses real loopback SSH,
rsync, and an isolated mount namespace to exercise AL publication and the packaged
Ansible mirror. Its receipt records hashes, permissions, links, and failure cases.
