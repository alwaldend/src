---
name: al-versioning
description: >-
  Calculate and create this repository's global development, nightly, and
  weekly release versions and supply Bazel stamping metadata. Use for
  first-party repository or project releases; do not use to version third-party
  dependencies.
---

# Manage repository versions

Use the Go tool at `@rules_al//:al -- version`. It owns one global
version shared by the repository and all first-party projects. Leave external
dependency pins and upgrades to their owning dependency workflow.

## Preserve the version contract

- Ordinary development is always `0.0.0-dev`. Default unstamped builds must
  not derive a version from the clock or current commit.
- Nightlies are exact tags on clean trunk commits:
  `vYYYY.W.0-nightly.YYYYMMDD`.
- Weekly release branches are `releases/YYYY.W`. Release patch zero is the
  immutable branch-point tag `vYYYY.W.0`; each later first-parent branch
  commit increments the calculated patch.
- `YYYY` and `W` are the ISO week-year and week. Never zero-pad `W`: leading
  zeroes are invalid in SemVer numeric identifiers.
- The configured trunk defaults to `master`. Pass global `--trunk NAME`
  before the command only when the repository's actual trunk name differs.
- One commit may carry one nightly tag and one release tag when a nightly is
  promoted. Attached trunk and release branches select the channel
  automatically; detached checkouts require explicit context.

## Inspect before changing refs

Run:

```sh
al version show --format json
```

`show` and `bazel-status` are read-only. Tag and branch commands mutate local
Git refs; run them only when the user requested the corresponding nightly or
release action. They refuse a dirty tree, the wrong branch, an in-progress Git
operation, a moving source branch, or a conflicting version tag. Use
`--dry-run` first when validating an unfamiliar release state.
Dry runs perform the same branch, cleanliness, and ref-collision preflights as
real mutations, including an atomic prepare-and-abort of the ref transaction;
they omit only the committed ref update and branch switch.

## Create a nightly

From a clean trunk checkout:

```sh
al version nightly-tag --dry-run
al version nightly-tag
```

Use `--date YYYY-MM-DD` only for a deliberate UTC calendar override. The tool
creates the tag locally; pushing it is a separate delivery action and requires
the user's release scope and the repository delivery workflow.
If the current commit already has a different nightly tag, the tool refuses a
second nightly instead of creating an ambiguous commit. Wait for a new trunk
commit. It also refuses to add a nightly after a release tag already exists;
the supported co-tag order is nightly promotion into a release.

## Start or patch a weekly release

Start the current ISO week from clean trunk. The tool atomically creates the
branch and patch-zero tag, then switches to the release branch:

```sh
al version release-start --dry-run
al version release-start
```

On `releases/YYYY.W`, every first-parent commit changes the calculated patch.
Inspect it, then create the exact release tag:

```sh
al version show --format json
al version release-tag --dry-run
al version release-tag
```

Do not manually move release base tags or reuse a weekly branch name. A
calculated patch that disagrees with an exact tag is an integrity failure to
investigate, not a reason to force the tag.
If trunk has not moved since an earlier weekly branch point, `release-start`
refuses to put a second release tag on that commit. Create the next release
only from a new trunk commit.

## Supply detached-checkout context

Branch context is unavailable on detached HEAD. If the exact commit has both
the nightly and release tags, select which immutable tag is being rebuilt:

```sh
al version \
  --channel release show --format json
```

An untagged commit from a release branch needs its release line to calculate
the patch from first-parent history:

```sh
al version \
  --release 2026.35 show --format json
```

Use the same global option before `bazel` for detached CI stamping. Never
guess a release line: obtain it from the CI ref that selected the checkout.

## Stamp Bazel release artifacts

Use the bootstrap entry point for a one-command build. It asks Bazel for a
source-current launcher under ignored `out/release/`, then the tool runs
Bazel with that launcher as the workspace-status command:

```sh
tools/rules_al/cmd/al/al.sh version bazel -- \
  build --config=release //path/to:artifact
```

Direct `al bazel ... --config=release` fails closed because a workspace
status command cannot safely build itself and must not trust a stale
`bazel-bin` output. Maintained release and CI workflows must use the bootstrap.
The stable keys are `STABLE_VERSION` and `STABLE_VERSION_CHANNEL`. Nightly and
release builds also emit `STABLE_GIT_COMMIT` and
`STABLE_GIT_TREE_STATE`. Development emits those dynamic values as volatile
`GIT_COMMIT` and `GIT_TREE_STATE`, so changing a development commit does not
invalidate Bazel's stable-status cache inputs. Nightly and release stamping
refuses dirty trees.
Rules that consume stamping must opt into Bazel stamping and read these keys
through their language or packaging rule's supported mechanism.

The `bazel -- ARGS` subcommand delegates the supplied Bazel command and
preserves its exit code. It does not make that command read-only or authorize
its side effects. Apply the normal authorization rules for `run`, deployment,
publishing, infrastructure, and any target that changes local or external
state.

## Publish reviewed release refs

Generate a deterministic, reviewed release-ref plan from the resolved version
state. The plan is not an authorization and never moves an existing immutable
tag:

```sh
al version \
  release-plan --lease release-refs --plan out/delivery/release-plan.json
```

The plan binds the exact version/channel/commit/tree state, the target refs
(nightly tag only, or release branch plus patch tag), the expected remote
preconditions, and the atomicity requirement. Review it before consumption.

Publish the exact reviewed plan only with explicit release scope:

```sh
al version \
  release-publish --plan out/delivery/release-plan.json \
  --receipt out/delivery/release-ref-receipt.json
```

The guarded publisher fetches expected remote state, acquires the distinct
`release-refs` lease, publishes with atomic multi-ref support when required,
and verifies the remote before emitting a `ReleaseRefReceipt`. A remote that
cannot guarantee atomic publication, an occupied lease, or any observed
mismatch is an explicit refusal or unknown, never generic success. Keep
versioning, delivery, release, review, and OpenSpec changes as separate
authorities; each consumes typed references only.

## Publish independent projects

Build a single project with the shared bootstrap and its `releases:bundle`
target (the website uses `//projects/alwaldend.com:release`). Publish the bundle
with `al release publish` and an explicit SSH environment/destination.
Published identities are immutable: identical retries succeed, conflicts fail.
Use `delete` at the same destination for an inactive release, then ordinary
`publish` to republish. Deletion preserves shared Git tags and keeps no
restoration state. Read the release README for transport options.
