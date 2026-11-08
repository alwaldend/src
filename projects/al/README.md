---
title: AL
description: Repository lifecycle commands and reusable Bazel rules
languages:
  - bzl
  - go
tags:
  - releases
---

AL is an independent project providing repository lifecycle commands and reusable
Bazel rules. Its source lives in `projects/al`; its Bzlmod module name remains
`rules_al`, preserving consumers' `@rules_al` labels.
Its direct commands are `al config`, `al run`, `al version`, `al release`,
`al delivery`, `al bazel test`, and `al doctor`.

Consumers supply plugin configuration, repository and base selections, and release
destinations. They run their normal build, test, and lint commands before source
publication. Calendar versioning, aggregate-commit rules, and scratch conventions
are integration contracts that a consuming repository provides or adopts.
This repository's configuration and installation wrapper live
in [`tools/al`](../../tools/al/README.md). The examples and verification evidence
below cover this repository and its isolated fixtures.

The [product description](docs/product-description/README.md) follows the
operator-visible command lifecycle, interruption and recovery boundaries, with
verification checklists and source-linked defect triage.

The executable owns signals and startup. `internal/cli` assembles command adapters;
local release, delivery, versioning, Git, and Bazel packages retain their domain
operations and guards. Reusable configuration, plugin and resource-lifetime APIs
live in `pkg/`. CLI entry points use the same owning implementations as Bazel
execution tools. Credential plugins start through `al run`; ordinary lifecycle
inspection and artifact generation do not start them.

AL's public Bazel API is `@rules_al//:defs.bzl`, including composable configuration targets, `al_config` aggregation and small
dictionary inputs, AL run wrappers, release packaging rules, and their providers. The module supplies its
own AL execution toolchain. Existing pinned dependencies are declared locally.

From this repository's root, build the executable with
`bazel build @rules_al//:al`, or install it with
`bazel run //tools/al:install`; installation requires the operator's authority.
Use ordinary Bazel commands for builds and tests. For Git-aware test selection,
`al bazel test --base master` selects affected tests in the chosen workspace;
`--dry-run` previews that selection. See the
[affected-test reference](docs/bazel/README.md#test-affected-targets).
The
[repository setup guide](../alwaldend.com/content/docs/misc/repo.md) describes
per-user Bazel defaults, and [Bazel setup](docs/bazel/README.md) describes AL
installation and diagnostics. Release stamping uses the installed `al` on PATH;
update that installation before building a release from changed AL source.

Operator procedures live in the module's `al-*` skills. Workspace preparation,
checks, maintenance, and OpenSpec workflows use existing targets through those
skills; they do not add speculative command groups.

AL also owns repository versions and independent project releases.
The release identity is `(project, version)`: projects share a source version
and process, while each publishes and activates independently. Operators can
run the entire process without CI.

## Bazel module API

Consumers depend on `rules_al` and load its `al_`-prefixed rules from
`@rules_al//:defs.bzl`: `al_release`, `al_release_files`, `al_release_bundle`,
`al_release_binary`, `al_release_deployment`, and `al_release_deps`. Providers are
`AlReleaseInfo`, `AlReleaseFilesInfo`, and `AlReleaseDeploymentInfo`.

```starlark
load("@rules_al//:defs.bzl", "al_release_bundle")

al_release_bundle(
    name = "release",
    srcs = [":archive"],
    project = "projects/example",
)
```

`status_file` can supply an explicit stable-status file. Otherwise the bundle
consumes Bazel workspace status when stamping is enabled, and ordinary unstamped
builds emit the constant development version. The standard `stamp` attribute
controls this behavior. The module builds its own tool and declares
all software dependencies at repository pins. Run `bazel test //...`
and `bazel build //...` from `projects/al` to check it independently.
The `tests/consumer:check` command builds the public API without parent Bazel
configuration imports and automatically verifies its output manifest and payload.
Its README describes launcher invocation and the retained `result.json`.

Release manifests carry the minimal version, channel, source commit, and tree
state needed to identify a build. Optional GitBundle history supplies changelogs;
AL does not export a complete repository-history JSON document. The exported
`git` extension configures a consumer's local checkout and system Git executable;
register the resulting `git_toolchain` to use `@rules_al//:git` or `al_git_binary`.
History generation deliberately reads that checkout locally and is excluded from
remote execution/cache upload. Ordinary portable bundle generation uses only
declared inputs.

```starlark
git = use_extension("@rules_al//:extensions.bzl", "git")
git.local_git(name = "repository_git")
use_repo(git, "repository_git")
register_toolchains("@repository_git//:git_toolchain")
```

The `al_release` documentation rule accepts optional `git_bundle` and `page_template`
inputs. This repository's website supplies its history bundle and shared Hugo
markup through `//projects/alwaldend.com/pkg/release_pages:defs.bzl`.
Consumers supply their own history and page-template inputs to the module API.

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
checkout with `--repository PATH`.
SIGINT and SIGTERM cancel active Git subprocesses. Foreground
Git authentication can use the controlling terminal; AL restores terminal ownership
after the child exits.
Use `al version nightly-tag`, `al release create`, or `al version release-tag`
for explicitly authorized local ref preparation. Each supports `--dry-run`.
The `al-versioning` skill documents detached-checkout and calendar details.

For a weekly release, start on clean trunk and run `al release create`. It
atomically creates `releases/YYYY.W` and `vYYYY.W.0` at trunk HEAD, leaving the
current checkout unchanged. Select the created branch with ordinary Git before
planning source publication or building release artifacts. Use the branch name
printed by `create` in place of `releases/YYYY.W` below:

```sh
al release create --dry-run
al release create
git checkout releases/YYYY.W
```

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
bazel build --config=release //projects/autoscroll/releases:bundle
```

Projects with declared first-party payloads have a `releases:bundle` target.
The website's bundle aliases `//projects/alwaldend.com:release`. The Android,
kustomization, and Terraform module release pages currently declare no payloads;
they provide documentation until their owners define release artifacts.
Existing `head` release pages are development documentation, not publishable
release identities. Nightly-to-weekly promotion rebuilds version-bearing
artifacts under the weekly identity.

`al_release_bundle` packages declared payloads as `release.json` and
`files/<filename>` in a Bazel directory output. It consumes `STABLE_VERSION`,
`STABLE_VERSION_CHANNEL`, `STABLE_GIT_COMMIT`, and `STABLE_GIT_TREE_STATE` from
the workspace-status command. First-party publication requires a concrete nightly or
release version and clean source identity. Development bundle manifests and
payloads can be examined locally; `release plan`, `release inspect`, and
publication require a publishable identity. Third-party versions continue to
come from their dependency owner.
The build does not contact deployment destinations.

The root `release` configuration enables stamping and invokes
`al version --repository . bazel-status` through the workspace-status script.
Install AL before starting the build; the script does not bootstrap a nested
Bazel invocation. Source refs choose the version, and dirty nightly or weekly
source fails stamping. For an ambiguous detached checkout, override
`--workspace_status_command` with the same status command and explicit
`--channel` or `--release` selection.

The CLI equivalent is `al release generate --project projects/<name> --version_file
<status-file> --add_file <payload> --output_dir <empty-directory>`. Repeat
`--add_file` for multiple payloads. Missing or duplicate status values, unsafe
paths, missing payloads, colliding filenames, and empty bundles fail. The
output directory must be absent or empty. A completed manifest is written
only after payload copying finishes.

## Publish over SSH

SSH publication uses the operator's OpenSSH and rsync executables and requires
rsync and standard GNU filesystem utilities (including `flock` and `mountpoint`)
on the destination. No receiver service is installed. Select the environment and
destination explicitly:

```sh
bazel run @rules_al//:al -- release publish \
  --release_dir "$PWD/bazel-bin/projects/autoscroll/releases/bundle" \
  --environment local --ssh_host download-host --ssh_user admin \
  --publish_user download --ssh_root /srv/download
```

The download component provides preconfigured `//infra/download:publish.local`
and `//infra/download:publish.yandex` commands. These use the configured SSH
administrator and content account; provision the selected host first.

`--ssh_deployment` accepts a deployment JSON
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

Use `--receipt /absolute/path/to/attempt.json` with artifact publish/delete or
website extract/activate to checkpoint local operation evidence. The private JSON receipt binds the selected
namespace, project, version, source commit, canonical publication digest and resolved
SSH environment/host/root. It records each requested step as `pending`, `unknown`
when started, and `completed` after confirmation. A failed operation preserves
completed steps and leaves an unconfirmed step unknown; even a transport error can
occur after a remote change. `status` is `running`, `completed`, or `incomplete`.
A killed process may leave `running` evidence with an unknown step. Use a separate
receipt path for each attempt to preserve history. Receipt persistence must succeed
before remote work starts. Receipts are observations, not authorization or proof
of current remote state. Inspect the destination and compare its publication
digest and observed state with the receipt, then retry any unresolved intended
step. Inspection does not rewrite receipts or authorize recovery operations.

## Preview and inspect artifact publication

`release plan` previews artifact upload using the same bundle selection,
destination overrides, and source identity validation as `release publish`.
This is separate from `version release-plan`, which plans Git refs. Website
commands provide their own offline preview through `--dry-run`.

```sh
al release plan --release_dir "$PWD/out/site-release" \
  --environment local --ssh_host download-host --ssh_root /srv/download
al release inspect --release_dir "$PWD/out/site-release" \
  --environment local --ssh_host download-host --ssh_root /srv/download
```

The plan emits `release/artifact_plan/v1` JSON with `mode: offline_metadata`.
It reports project/version/source, exact selected payload names, sizes and
hashes, canonical publication digest, destination, public/site paths and steps.
`selection` identifies explicit host/deployment overrides, which select all
items, or environment selection. Matching destinations are processed in stable
order. The plan reads bundle metadata only: it does not check payload bytes,
contact SSH or rsync, initialize credentials, or claim remote readiness. Missing
local payloads therefore do not prevent previewing an identical publication retry.
The destination includes configured login/content accounts; `port` and
`ssh_user` are null when OpenSSH owns their resolution, and `publish_user` is
null when commands run directly as the SSH login account. Authentication file
paths and credential contents are omitted.

Inspection emits `release/artifact_inspection/v1` JSON with
`mode: read_only_snapshot` and a completion timestamp. It independently reports
each local payload as `verified`, `missing`, `corrupt`, or `unreadable` and reads
the destination through bounded, noninteractive SSH. Local failures do not
prevent remote inspection. The destination content root must already exist.

| Remote status | Observed meaning                                                                                                |
| ------------- | --------------------------------------------------------------------------------------------------------------- |
| `absent`      | No publication exists at the selected version path.                                                             |
| `matching`    | Canonical manifest, complete file set, sizes, and SHA-256 hashes match.                                         |
| `conflicting` | The regular public manifest differs from the expected canonical bytes.                                          |
| `corrupt`     | The path or file types are unsafe, or the expected publication has invalid payloads or extra/missing files.     |
| `unavailable` | Transport, timeout, read access, missing content root, or response validation prevented a reliable observation. |

For first-party sites, `site.status` is `absent`, `retained`, `active`,
`conflicting`, or `corrupt`; it is `not_applicable` for external artifacts and
`unavailable` when the read fails. A site check verifies the publication binding
and regular root `index.html`, not every extracted file. `site.active` independently
reports whether `current` selects this version, including a corrupt or missing
site, and is null when selection is unknown. `site.current` distinguishes
`absent`, `selected`, `other`, `dangling`, and `unsafe` links. A valid retained
site can therefore have unknown selection if `current` is unsafe.

Completed inspections exit zero even when they find absent, conflicting, or
corrupt state. Unavailable observations emit their report and exit nonzero;
invalid local metadata fails before an observation report. Cancelling local
verification stops hashing between chunks and prevents the destination read;
no completed snapshot is produced. Neither command
accepts `--receipt`. Inspection creates no remote locks, staging directories,
or repairs, and uses no rsync transfers. The configured timeout bounds each
destination read; each captured stream is capped at 16 KiB and overflow cancels
the command. This is an observation without a lock, so concurrent publication
or deletion can change the state during or after inspection.

## Extract and activate a website

`al release publish` uploads artifacts. Run `al release website extract` to create a
retained site from a published archive, then `al release website activate` to select it.
These are independent commands; neither uploads artifacts. Website commands
accept the same bundle and SSH destination selectors as publication and apply
only to first-party projects.
`--site_archive` belongs only to extraction; activation and inspection derive
site state from the selected publication and retained binding.

```sh
bazel run //infra/download:publish.yandex -- \
  --release_dir "$PWD/out/site-release"
bazel run //infra/download:website.extract.yandex -- \
  --release_dir "$PWD/out/site-release" \
  --site_archive website.tar.gz
bazel run //infra/download:website.activate.yandex -- \
  --release_dir "$PWD/out/site-release"
```

Use `.local` instead of `.yandex` for the local environment. Direct CLI commands
take explicit destination flags, for example:

```sh
al release website extract --release_dir "$PWD/out/site-release" \
  --environment local --ssh_host download-host --ssh_root /srv/download \
  --site_archive website.tar.gz --dry-run
al release website activate --release_dir "$PWD/out/site-release" \
  --environment local --ssh_host download-host --ssh_root /srv/download --dry-run
```

Each `--dry-run` returns the offline metadata plan for that one operation. It
does not inspect archive bytes or remote state, contact SSH, or write a receipt.
Remove `--dry-run` to perform the chosen operation. Mutating commands accept
`--receipt`; combining it with `--dry-run` fails. Use `al release inspect` for
read-only publication, retained-site, and current-selection observations.

Extraction accepts a declared `.tar.gz` or `.zip` archive containing a root
`index.html`. It retrieves and verifies the published archive, extracts locally,
and promotes the site to `sites/<project>/releases/<version>`. Local free space
must hold both archive and extracted content. Traversal, absolute paths, links,
devices, duplicate names, and more than 100,000 entries fail. Archive ownership,
executable bits, and extended attributes are not retained.

Retained sites carry a binding to the exact publication. Activation checks this
binding and the public payloads under the destination lock, then atomically
selects `sites/<project>/current`. Managed ancestry must use real directories without nested mounts. Activation
also rejects mounts within the selected extracted site. Failed earlier steps leave activation
unchanged. Roll back by selecting another retained release with `al release website activate`.

## Delete and republish

Use the same manifest and explicit SSH destination with `delete`:

```sh
bazel run @rules_al//:al -- release delete \
  --release_dir "$PWD/out/site-release" --environment local \
  --ssh_host download-host --ssh_root /srv/download
```

Deletion removes the project release and its extracted site at that destination.
It refuses the currently selected website: first select another retained release
with `al release website activate`. Deletion refuses symlinked ancestry and mounted target directories or descendants, including bind-mounted files,
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

## Plugins

The [built-in injector](pkg/injector/README.md) evaluates named resource graphs
inside AL and owns generic environment, file and process resources. Eligible
plugins advertise resource types in their startup response and supply
them through generic resource-fetch RPCs. Startup returns an extensible
`PluginStartResponse` with an optional generic result; fetches return `ResourceResult`
directly. The standalone [Vault provider](plugins/vault/README.md)
owns Vault-specific resources. See the [configuration API](rules/README.md).
