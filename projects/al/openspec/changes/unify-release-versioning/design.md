## Context

The agreed release unit is `(project, version)`. Git tags identify repository source, while each destination independently records which project artifacts exist. Publishing a project does not build or publish other projects. CI is optional.

## Decisions

The final public owner is `projects/al`, a standalone Bzlmod module shipping
the AL executable, reusable configuration/plugin APIs, release rules, and their
execution toolchain. Version, release, delivery, and diagnostic commands are direct AL
subcommands. The command and plugin runtime remains scoped to `al run`; inspecting
versions or delivery state does not initialize credentials or plugins. Existing
typed delivery and release receipts retain their schema identities and guards.
The private executable cache remains removed. Bazel builds execution tools;
operators run ordinary Bazel checks independently of source publication.
`al bazel test` supplies dependency-based test selection without restoring generic
Bazel forwarding or a delivery validation ledger.

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
Public rules and their filenames use the `al_` prefix, and providers use `Al`;
project skills use `al-` names. The module declares its own existing pinned dependencies,
provides `defs.bzl`, builds its own CLI, and has no parent-repository labels.
Repository Git history bundles are explicit inputs. The existing generic Git
schema and parser support optional release changelogs. Used local Git toolchain
rules move with them; unused Git changelog wrappers and the separate Git CLI are
removed. The unused full-history JSON exporter, command, Bazel rule and target,
and website data hook are removed. Publishable bundles require only their stamped
version, channel, source commit, and tree status, not a repository history export.

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
checkout selector is `--repository`. Command
cancellation reaches Git's owned process group. Bazel owns its own client process.

An identical publication is verified under the destination lock before local
payload access or transfer. New identities keep their post-transfer locked
verification and promotion so concurrent publishers preserve the immutable
contract. Deletion preflight checks mounts on all descendant entry types.
The external consumer fixture supplies its own Bazel configuration and explicit
first-party module overrides, independently checking public rule outputs.

## Terminal and operation evidence

Git subprocesses retain their own cancellation group. When AL owns the foreground
terminal, it lends that terminal to Git and restores the original foreground
group after the child exits. Noninteractive execution keeps the existing bounded
group cancellation. Terminal prompts and cancellation must both work.

An automated consumer check materializes the existing templates under an explicit
ignored output directory, builds through Bazel without parent
Bazel configuration, verifies manifest identity and payload bytes, and retains a
JSON result. It runs outside a Bazel action to avoid recursive server locks.

The optional artifact --receipt file records resolved project/version, canonical
publication digest, destination, and ordered step observations. Write a private
atomic checkpoint before remote effects and after each confirmed step. A started
step remains unknown until the complete operation confirms success; an error or
interruption cannot prove the destination was unchanged. Preserve earlier completed
steps in an incomplete receipt. Receipts are local observations, not authority or
proof of the destination's current state. They exclude credentials and raw errors.

## Product description

Describe the Linux CLI from the operator's point of view in owner-local
`docs/product-description`. Pin source observations to the post-rebase executable
snapshot and distinguish source-derived claims, scripted observations, fixture
checks, and manual verification. Keep bundle generation, source refs, artifact
publication/activation, and guarded code delivery as separate state owners.
The initial description recorded defects before the separately authorized fixes.
The description is a source-identifying projection, not a second normative contract.

## Product-description corrections

The user's subsequent ownership decision makes AL an independent project under
`projects/al`, retaining `rules_al` as its reusable Bazel module name. Repository
tooling is a use of this project, not its ownership boundary. Consumer repository
names, paths, destinations, and validation policy belong to configuration; the
reusable implementation must not require the alwaldend/src repository. This is
an architectural portability requirement, with the existing independent-consumer
check retained rather than an exhaustive matrix of every CLI workflow elsewhere.

Standalone execution must not require Bazel runfiles when no selected consumer
needs runfile resolution. Named stdout/stderr outputs use replacement semantics
and create missing files privately. Conflicting selected plugin exports fail
deterministically instead of depending on map traversal; inherited environment
overrides retain their existing purpose. Direct delivery rebase changes only the
local candidate, preserving the fetched remote lease for preparation, separate
operator-run checks, and publication. Write behavior fixtures before these fixes.

Assess replacement of Lua configuration with typed Bazel providers using current
consumers. Keep analysis-time declaration separate from runtime credentials and
resource lifetime. The assessment does not authorize a repository-wide migration
or speculative lifecycle commands. The subsequent user request authorizes the
release inspection and publication preview described below.

### Configuration assessment

The inspected repository has 25 `al.lua` files and two helper libraries. They
construct static plugin/call/auth declarations; the two public SSH CA reads in
`infra/download/al.lua` and `infra/security/al.lua` are declared build inputs.
Credential acquisition and resource lifetime already belong to runtime plugins.

Decision: proceed with typed Bazel provider composition as the recommended
replacement for repository Lua authoring, in a separately requested migration.
The user explicitly deferred that migration; current Lua configuration remains.
Extend the existing `AlConfigInfo` boundary rather than add a competing config
system. Rules can collect typed plugin/call/auth data and declared executables,
validate names and references during analysis, then emit the existing JSON
manifest with `ctx.actions.write` and schema-aware `json.encode` input. Plugin
data uses the existing protobuf `Json` representation; encoding raw Lua-shaped
tables would not preserve that wire format. Preserve ordered dependency/local
composition and intentional last-match Vault auth/connection overrides. Public file contents require
a declared generation action or a runtime file reference; Starlark analysis does
not read arbitrary file contents. Vault login, SSH signing, temporary files,
backend listeners, environment export, and cleanup remain runtime operations.

Retaining Lua has the lowest migration cost but preserves two composition layers,
literal executable paths, and an execution step for static data. Handwritten
JSON/YAML removes Lua execution but loses typed label dependencies and keeps
manual composition. Provider-backed generation best matches current consumers;
preserve dependency ordering, resource templates, binary runfiles and non-Bazel
JSON consumers while migrating. Remove Lua loading only after those consumers
and standalone fixtures have been migrated and verified.

### Demonstrated lifecycle gaps

Prioritize read-only release inspection and payload verification: recovery
currently resolves unknown artifact receipts by repeating an operation, while
the existing publisher already owns the manifest, hash, site-binding, and active
selection checks. Reuse those checks to report absent, matching, conflicting,
retained and active state without publication effects.

Expose an artifact publication plan/preflight showing resolved destination,
selected payloads, steps and intended changes. Source-ref planning already exists,
but an explicit SSH host selects all manifest items while environment selection
filters them; an artifact preview would make that distinction visible.

Implement both through the release package rather than duplicating selection,
canonical manifest, digest, and site-binding algorithms in the CLI. Preview is
offline: validate the selected manifest and intended steps, show resolved SSH
destinations and payload metadata, and report the effect of explicit destination
overrides. It does not claim remote readiness or authorize later publication.
Inspection observes current remote state and verifies its expected publication
and site binding without creating directories, taking write locks, staging files,
repairing content, or changing activation. Observations are snapshots, not leases;
publication retains its own locked rechecks. Local payload verification must be
reported independently so an operator can inspect a remote release after local
payloads have been removed. Transport failures must not masquerade as absence.
Reports omit credentials and preserve distinct unknown/invalid/conflicting states.

Ordinary Bazel accepts startup options directly, including the bounded cleanup
invocation in `tools/bazel_session_cleanup`. Consider abandoned staging/lock
inspection before any cleanup command; age alone cannot establish safe ownership.
Keep formatting, full checks, maintenance, OpenSpec, installation, worktree
coordination and infrastructure synchronization with their existing owners;
do not create duplicate implementations merely to expand the command list.

### Inspection acceptance evidence

Built-CLI fixtures were written before implementation and preserve failing-before
logs for missing commands, newline-terminated current links, and cancellation while
hashing local files. The final isolated release suite retains 25 plan/inspection
reports, an environment-filtered preview, and a cancellation artifact. Tree
snapshots compare names, permissions, modification times, file hashes, and links
before and after observations. The cancellation fixture waits for local hashing
to begin in a sparse 128 GiB payload, then requires prompt nonzero exit after
SIGTERM without starting SSH. All destinations and payloads are local fixtures;
no live release was contacted. Inspection is still a non-atomic observation and
does not verify every extracted website file.

The four earlier fixes retain 13 installed-command observations and local-rebase
cases for present, absent, and already-based remote refs. The selected final gates
cover standalone AL tests/builds, parent consumers, generated skill discovery,
OpenSpec validation, representative independent-consumer output, repository quality,
and aggregate affected-package lint. Delivery receipts own exact candidate and
publication results; these fixture observations do not substitute for those gates.

### Independent project integration

The owner is projects/al; the exported Bzlmod name remains rules_al. The Go module
identity follows the source location, with parent imports and the existing local
replacement updated together before this feature's initial delivery. Root module
registration, Bazel profiles, configured plugin wrappers, release destinations,
and required check procedures remain consumer integration. The existing lint profile
name, quality target, calendar versions, aggregate-commit model, and ignored
scratch layout are AL integration contracts another repository can provide or
adopt; this change does not introduce a second configuration framework.

Inspection now honors only its explicitly selected or pull-request-selected base,
with the existing fresh snapshot and ancestry guards. An unrelated master tracking
ref is no longer an additional authority. Doctor accepts configured pin and source
binary paths plus declared rc-file inventory. Relative inputs resolve against its
workspace root, and rc declarations do not claim to discover or validate imports.
Its generic default lists .bazelrc rather than this repository's private layout.

The relocation uses owner-generated project/OpenSpec inventories and skill
projections; repeated generation is stable. Existing external-consumer coverage is
adapted to the source path and dependency overrides, without a new exhaustive
cross-repository workflow matrix or a claim that a registry distribution exists.

### Public naming and project presentation

All AL-owned rule and rule-macro symbols use the al_ prefix; provider symbols use
Al, following the existing AlConfigInfo convention. Rule/provider source filenames
follow that naming, while instance target names, the rules_al module identity, and
the public Git extension identity remain stable. Consumers and generated repository
templates migrate together; the pre-release unprefixed aliases are removed.

AL's landing content lives in projects/al/site/content and is exported by the
standalone module. The existing project registry maps al to rules_al so the shared
site packages it without parent labels inside AL. The landing describes current
capabilities and links to the owning documentation; no new site, DNS, or deployment
is introduced. Generated project inventories follow their existing owner workflow.

The requested dark Mermaid artifact contains horizontal Development, Source
publication, Release, and optional Website publishing subgraphs in one state
machine. Development covers editing, secret injection through al run, building,
and testing. Source publication follows a regular state chain for preparation,
checks, publication, review, and merge, without forced rows or nested layout
subgraphs. Release covers shared versioning,
per-project builds, publication, and deletion. Website publishing separately
covers extraction and activation through AL's actual publish steps.
Each command labels its own transition between meaningful states, in bold on one
line; immutable publication and deletion are explicit.
Inspect the raster at a practical display width; larger pixel dimensions alone do
not demonstrate readability. Keep its editable source and final WebP in task output.

The consolidation retains no compatibility-only owners, redirect stubs, or obsolete
entry points. Remove retired surfaces and migrate current callers to AL's canonical
API. Active repository configuration and installation wrappers remain consumer
integration, with their actual behavior and dependency graph determining ownership.

Cleanup removes tools/versioning, the unused root AL/Terraform aliases, the old
version --repo spelling, the release deploy spelling, and obsolete Soras/ORAS
executor arguments. Artifact publication uses release publish; repository selection
uses --repository. The separate OCI metadata rules remain functional.

## Ordinary Bazel and user configuration

Proceed with ordinary Bazel for builds, tests and lint. Remove generic `al bazel`
forwarding and `al version bazel`; they add no lifecycle operation. The focused
`al bazel test` command below adds affected-test selection.
Keep `al version bazel-status` as the identity producer. The root release
configuration invokes an adapter that executes installed AL with an explicit
current repository; it never bootstraps another Bazel process while a build
holds the server lock. Missing AL produces an installation instruction. Explicit
workspace-status overrides select detached release or nightly context when needed.

Move the agent output/test defaults from the repository's named profile into the
bot account's Ansible-managed user rc file. Ordinary options work independently
of the current repository; a global `common --config=agent` would require every
repository to define that name. Retaining the wrapper adds a command and excludes
normal startup-option placement without adding behavior. Keep the setup guidance
in the existing al-bazel skill, linking the owning config instead of duplicating
its settings. Documentation and checked-in desired state do not deploy the host.

Preserve doctor source/pin/scratch diagnostics and its declared rc inventory;
remove its hardcoded agent-profile claim. Verify clean stamped bundle identity,
dirty-source refusal, and missing installation, plus the existing delivery and
independent-consumer fixtures. Lua configuration migration remains deferred.

Keep build and artifact publication separate, with the verified release bundle
as their handoff. A later CI build job can retain that bundle; a publication job
consumes the same bytes without rebuilding or recalculating the version.

Expose initial weekly source-ref creation as `al release create`. It creates the
local release branch and initial tag atomically and leaves HEAD on trunk.
Switching to an existing branch belongs to ordinary `git checkout`; do not add
an AL checkout command or keep the former version release-start alias. Preserve
clean-trunk, existing-ref, and dry-run guards. The diagram distinguishes local
ref creation, explicit checkout, bundle building, and artifact publication.

## Git preparation and independent checks

`al delivery prepare` owns local Git preparation and records the exact prepared
commit, scope, rewrite authority, and fetched publication lease. It does not
select or invoke Bazel targets, accept Bazel arguments, or report hardcoded
affected labels, validation gaps, or suggested validation commands. Remove
`al delivery validate`, `continue`, and `deliver`, their JSON check-plan format,
and the validation result, environment, and log-hash ledger.

Operators run the repository's required builds, tests, and lint through their
normal commands before publication. These checks remain procedural requirements;
AL does not certify them through a second ledger. `al delivery publish` consumes
the preparation receipt and checks the current candidate against its recorded
HEAD, index, tree, and task scope. Remove the redundant `--validated-head`
assertion. Preserve explicit remote leases, repository and pull-request guards,
publication receipts, verification, and review handling. Unrelated local edits
outside the prepared scope remain permitted when the publication does not rebase;
review addressing retains its existing clean-candidate requirement.

If a freshly fetched base requires rebasing, checkpoint the updated candidate and
receipt, return `checks_required`, and stop before pushing. The operator runs
checks against that candidate and publishes again using the same receipt.
After a successful push and pull-request update, a failed review inventory read
retains publication evidence as `publication_verification_pending`; verification
can recover the observation without repeating publication.

## Affected tests through Bazel

The purpose is to keep unrelated targets out of the test run, so unrelated cache
misses cannot trigger unnecessary builds or test execution. The dependency query
still loads its graph; no performance improvement is claimed without measurement.

`al bazel test` is independent of delivery state. It compares Git changes against
an explicitly selected base, includes local working-tree changes, and asks Bazel
for the reverse dependencies and tests within the selected workspace and target
scope. It invokes ordinary Bazel test execution for the resulting labels. Bazel
owns dependency semantics, configuration, caching, and whether those tests need
to execute; AL does not infer dependencies from directory proximity or maintain
a competing build graph. Selecting one workspace does not claim coverage of
other standalone workspaces in the repository.

Only successful graph information can establish that no test is affected.
Deletions, build-configuration changes, or source paths that cannot be mapped
reliably require the complete selected test scope rather than a falsely narrow
list. A query failure or partial result cannot become a successful empty test
selection. Report why selection widened, and propagate failed test execution.
There is no arbitrary Bazel argument forwarding, validation receipt, delivery
plan, or publication side effect. Operators can still run `bazel test` directly
over their chosen scope and rely on its cache.

## Explicit website operations and minimal source identity

Proceed with dedicated `al release website extract` and `al release website activate`
commands. Website publication is optional and remains within the release workflow;
an additional top-level command family would separate operations with the same
artifact identity and destination contract. Artifact publish and delete commands
select their fixed operations; remove the public --steps selector rather than
keeping two ways to invoke website mutations. Each website command supports an
offline --dry-run preview; it performs no transport, starts no plugins, and writes
no receipt. Preview with --receipt is rejected. Destination resolution and backend
publication guards remain shared. Repository wrappers supply complete command
argument lists and retain the existing credential injection boundary.

Remove the full-history Git JSON export and its site data projection: the checked-in
site has no reader, while portable release bundles use only stamped identity.
Keep the separate optional changelog bundle and its parser, with sequential tag
reads and packed-object verification. This removes an unused feature without
changing the shared version scheme or project release independence.

Naming and cleanup acceptance uses the existing independent consumer, whose
verified manifest and payload now build through the prefixed API without parent
configuration. The built CLI records equal canonical --repository invocations
and rejects the three retired invocations; source delivery and publication
behavior tests retain their original guards. The shared site and DNS checks pass,
and inspected HTML includes AL in the project index and release taxonomy, its
restored landing, and the prefixed API documentation. Delivery receipts bind the
aggregate candidate to Git preparation and remote publication. Check results
remain separate evidence from the ordinary commands run against that candidate.
