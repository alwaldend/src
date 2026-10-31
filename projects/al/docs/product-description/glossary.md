# Glossary

## Invocation and workspace

**AL.** The repository lifecycle executable described here.

**Invocation.** One execution of AL with fixed arguments and inherited environment.

**Operator.** The person or agent choosing a command and its authorized scope.

**Workspace.** A directory containing the build configuration for a Bazel invocation.
A Git repository can contain multiple workspaces.

**Worktree.** A checkout with its own working files and index. Linked worktrees
share Git objects and refs through a common Git directory.

**Candidate.** The exact commit and tree prepared for delivery, identified by OIDs.

**OID.** A Git object identifier used to bind an observation to exact content.

**Trunk.** The repository's main development branch, normally `master` here.

**TTY.** An interactive terminal; pipes and redirected files are different streams.

**Runfiles.** Runtime files located by a Bazel-launched executable. Copying only
the executable does not necessarily supply them.

**Plugin.** A configured helper started by `al run`, with resources kept alive
while its consumer runs. A label selects plugins or individual calls.

**Label.** A key or key/value selector for configured plugins and calls; also,
in a Bazel command, a target address such as `//path:target`. Context distinguishes them.

## Releases and durable state

**Source version.** The repository-wide identity calculated from source history;
its calendar rules belong to the module README, linked by the versioning document.

**Project release.** A publication identified by namespace, project, and version
at one destination. Projects choose when to publish a shared source version.

**Bundle.** A local directory with `release.json` and declared files under `files/`.

**Manifest.** The JSON description of a release, including identity and payloads.

**Payload.** A declared file whose name, byte size, and SHA-256 digest are recorded.

**Digest.** A cryptographic hash of content used to detect different bytes.

**Namespace.** The `projects` or `third_party` part of the public artifact path.

**Destination.** The selected SSH host, content root, and account context.

**Staging.** Temporary content prepared before a complete publication becomes visible.

**Publication.** A complete verified artifact set installed under its public identity.

**Publication plan.** A metadata-only preview of selected artifact identities,
destinations, the upload operation, and declared payloads. Website-operation dry
runs use the same report format for extraction or activation. Neither establishes
readiness or reserves destination state; source-ref publication uses a different plan.

**Release inspection.** Observations of local payload validity and remote
publication, retained-site, and activation state for a selected release.
The observations do not authorize mutation or make a later invocation atomic.

**Retained site.** An extracted website bound to a publication, available for activation.

**Activation.** Selecting a retained site through the destination's `current` link.

**Immutable.** Existing published identity and bytes cannot be overwritten by AL.
Explicit deletion permits ordinary publication again; it does not preserve a tombstone.

**Lease.** An expected ref identity checked before publication. It detects conflicting
progress; it is not permission to overwrite unrelated work.

**Receipt.** A local record of observations. Artifact-operation, source-ref, and
delivery receipts are different formats with different recovery behavior.

## Events that end or interrupt

**Completion.** The command reports success for the operation it actually performed.
It does not imply that every later release or delivery step has happened.

**Cancellation.** SIGINT (Ctrl+C) or SIGTERM requests cooperative shutdown. The
feature determines which subprocesses and resources respond.

**Interruption.** Any event preventing normal completion, including cancellation,
closed streams, network failure, changed state, or forced termination.

**SIGHUP.** A terminal-hangup signal; AL installs no special handler for it.

**Forced kill.** Termination that prevents application cleanup, such as SIGKILL.

**Unknown.** An artifact receipt step started but has not confirmed its outcome.
Unknown does not mean that the destination was unchanged.

**Refusal.** A guard rejects an operation. Earlier completed work may still exist;
read the operation's recovery boundary.

**Safe local observation.** A read-only invocation or an invocation confined to
synthetic task scratch, with no live publication or real credential use.

## Configuration and execution

**Configuration.** The merged explicit AL settings used by configuration inspection or command execution.

**Standard streams.** The child or AL process's stdin, stdout, and stderr channels.

**Exit status.** The numeric result returned to the invoking shell; zero denotes success for the selected operation.

**User Bazel configuration.** Per-user defaults in `~/.bazelrc`, including optional
agent output preferences. Bazel reads these alongside workspace configuration;
AL does not insert a configuration selection into build commands.

**Cleanup.** Shutdown and resource release after a child or plugin stops; it can fail and is not guaranteed after forced termination.

## Source and delivery vocabulary

**Channel.** Development, nightly, or release context used to calculate and validate a source version.

**Source ref.** A Git branch or tag naming source history, shared independently of project artifact status.

**Stamping.** Supplying source version and identity metadata to a build.

**Detached checkout.** A checkout whose HEAD names a commit directly rather than a local branch.

**First-parent history.** The chain obtained by following the first parent of each commit; it determines release patch position.

**Pull request (PR).** The forge's review object associated with a proposed branch change.

**Forge.** The service hosting repository metadata and review operations.

**Path scope.** The complete set of task-owned paths bound to a delivery candidate.

**Atomic update.** A particular operation makes its selected changes together; it does not imply that a multi-command workflow is one transaction.

**Review authority.** Short-lived receipt evidence binding an authorized review mutation to exact observed expectations; it is not permission by itself.
