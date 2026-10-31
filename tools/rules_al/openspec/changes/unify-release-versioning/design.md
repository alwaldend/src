## Context

The agreed release unit is `(project, version)`. Git tags identify repository source, while each destination independently records which project artifacts exist. Publishing a project does not build or publish other projects. CI is optional.

## Decisions

The final public owner is `tools/rules_al`, a standalone Bzlmod module shipping
the AL executable, reusable configuration/plugin APIs, release rules, and their
execution toolchain. Version, release, delivery, and Bazel commands are direct AL
subcommands. The command and plugin runtime remains scoped to `al run`; inspecting
versions or delivery state does not initialize credentials or plugins. Existing
typed delivery and release receipts retain their schema identities and guards.
The private executable cache remains removed. Bazel builds execution tools, and
the Bazel runner still replaces its process after validating the command.

Move component-owned lifecycle skills with the implementation and use the `al-`
prefix. Keep language, infrastructure, and upstream OpenSpec skills with their
owners. Repository integration tests that require the parent Vault injector stay
in the parent workspace, while the AL module's tests remain self-contained.

Workspace, rebase, complete checks, maintenance, and repository OpenSpec integration
remain skill procedures over their existing owners and targets. Move the existing
OpenSpec integration skill to `al-openspec`; upstream `openspec-*` projections remain
unchanged. Add `al-maintenance` for generated metadata and dependency procedures.
No extra command groups are needed for these workflows.

Follow the local ownership pattern used by `anki_as_code`: thin executable startup,
`internal/cli` command adapters, local domain packages, and public reusable APIs in
`pkg`. Keep the established release and delivery guards with their operations,
without introducing a facade that duplicates those responsibilities.

Keep version algorithms in their existing internal package and expose them through
`al version`. Remove the retired AL, release, delivery, and Bazel-agent projects
after migrating their consumers and maintained records. Retain delivery receipt
schema strings for compatibility; executable names and command adapters use AL.

First-party portable bundles require stable version, channel, source commit, and clean tree status. Development bundles may be inspected but cannot be published. Third-party versions remain owned by their dependency manifests. Preview `head` documentation remains separate from publishable bundles. Provide a portable bundle target for each project's nonempty declared payload set; the website aliases its existing bundle. Documentation-only declarations do not imply a releasable artifact.

Stage the selected payload set and a canonical publication manifest on the destination's content filesystem. Validate local bytes before transfer and remote SHA-256/size before promotion. Under the existing destination lock, compare an existing manifest and its files or atomically rename the complete staged directory. Publication manifests contain only identity and selected payload metadata, excluding environment and transport details. Never modify an existing release.

Extraction checks the published manifest and archive digest, and binds the retained site to that publication. Activation checks the binding under the same lock. Deletion is a separate command, serialized with upload/extraction/link; it refuses a currently selected website, removes both release trees, and retains no tombstone. Later publish follows the usual absent-release path. Shared source tags are never deleted by artifact deletion.

Source-ref publication retains reviewed plans/receipts, acquires an exclusive local Git-common-directory lock, and uses atomic Git push plus explicit expected values for creates. Patch plans publish only their immutable tag and verify the already-published release branch at the patch commit. Patch zero creates the branch and tag together. Remote Git provides cross-host ref serialization; the local lock coordinates linked worktrees.

## Risks and verification

This changes previous rsync overwrite semantics deliberately. Existing unmanaged directories fail closed and require explicit deletion before publication. All file transfers remain OpenSSH/rsync; no service or dependency is added. Nightly-to-weekly promotion rebuilds version-bearing artifacts under the weekly identity.

Exercise the built CLI in isolated Git repositories and a local SSH/rsync transport fixture. Cover identical retries, conflicts, corruption, interrupted transfer, concurrent promotion, independent projects, extraction/link, active deletion refusal, deletion and republishing, and diverged patch branches. Retain manifests and receipts as repeatable test artifacts. No live release operations are authorized by implementation.

The legacy OCI push cannot provide a conditional create under the OCI distribution
API used here. Fail closed rather than silently overwrite registry tags. Retain
OCI metadata generation; a future OCI publisher needs an explicit registry-side
immutability contract.

## Standalone rules module and review corrections

The public owner is `rules_al`, consumed through `@rules_al`.
Public rules, providers, and filenames omit the `al_`/`Al` prefixes; project skills
use `al-` names. The module declares its own existing pinned dependencies,
provides `defs.bzl`, builds its own CLI, and has no parent-repository labels.
Repository Git history bundles and optional OCI executables are explicit inputs.
The existing generic Git schema and parser are merged into this module and
shared by release changelogs and the git-info command. Used local Git toolchain
rules move with them; unused Git changelog wrappers and the separate Git CLI are
removed. Existing site data and tools consume the merged implementation.

Publication, extraction, and activation share the selected destination item set.
Local runfiles symlinks may resolve to regular payloads; copied staging bytes
remain size/hash verified. Website instructions require a nightly/release version.
A standalone Bazel fixture verifies actual bundle outputs and an isolated CLI
fixture covers filtered destinations and symlink payloads before implementation.

## Approved synchronization integration

The user approved reconciling fetched trunk `5cfe3ae3d8299448ace8e7190d7b4b86bc5f847a`
while preserving immutable publication in both namespaces and explicit deletion.
Keep the download owner's authoritative Ansible mirror as its separate operation.
Activation retains extracted versions for rollback. Managed destination ancestry
rejects symlinks and nested mounts; recursive deletion checks both operands and
their descendants before removing either tree. The configured content root may
itself be the mounted filesystem. Preserve SSH environment-reference identity
selection and start remote commands from `/` for content-account access.

## PR review refinements

Version commands use native Cobra subcommands and persistent flags. The public
checkout selector is `--repository`, with `--repo` retained as an alias. Command
cancellation reaches Git's owned process group and interrupts the stamped Bazel
client before a bounded forced stop.

An identical publication is verified under the destination lock before local
payload access or transfer. New identities keep their post-transfer locked
verification and promotion so concurrent publishers preserve the immutable
contract. Deletion preflight checks mounts on all descendant entry types.
The external consumer fixture supplies its own Bazel profile and explicit
first-party module overrides, independently checking public rule outputs.

## Terminal and operation evidence

Git subprocesses retain their own cancellation group. When AL owns the foreground
terminal, it lends that terminal to Git and restores the original foreground
group after the child exits. Noninteractive execution keeps the existing bounded
group cancellation. Terminal prompts and cancellation must both work.

An automated consumer check materializes the existing templates under an explicit
ignored output directory, builds through the supplied AL runner without parent
Bazel configuration, verifies manifest identity and payload bytes, and retains a
JSON result. It runs outside a Bazel action to avoid recursive server locks.

The optional artifact --receipt file records resolved project/version, canonical
publication digest, destination, and ordered step observations. Write a private
atomic checkpoint before remote effects and after each confirmed step. A started
step remains unknown until the complete operation confirms success; an error or
interruption cannot prove the destination was unchanged. Preserve earlier completed
steps in an incomplete receipt. Receipts are local observations, not authority or
proof of the destination's current state. They exclude credentials and raw errors.
