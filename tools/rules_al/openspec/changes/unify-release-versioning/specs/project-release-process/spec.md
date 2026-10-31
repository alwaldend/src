## ADDED Requirements

### Requirement: One standalone repository lifecycle module

The standalone rules_al module SHALL provide the AL executable and its Bazel APIs.
AL SHALL expose version, release, delivery, and bazel directly, without a repo
namespace. Lifecycle commands SHALL preserve existing guards and receipts and
SHALL initialize credential plugins only when the selected operation requires them.

#### Scenario: Use repository lifecycle commands without plugin startup

- **WHEN** an operator inspects a version or delivery provider, or generates a bundle
- **THEN** the AL command executes the owning implementation without starting plugins
- **AND** publication guards and the independent project version contract remain enforced

#### Scenario: Delegate Bazel through AL

- **WHEN** an agent supplies a supported Bazel command to al bazel
- **THEN** AL inserts the agent configuration after that command and preserves remaining arguments
- **AND** AL replaces its process so signals and the final exit status reach the Bazel client directly

#### Scenario: Cancel a version operation

- **WHEN** a version inspection, source-ref publication, or stamped Bazel command receives SIGINT or SIGTERM
- **THEN** its cancellation context reaches the active subprocess and the command fails promptly

#### Scenario: Use native version subcommands

- **WHEN** an operator supplies version flags before or after the selected subcommand
- **THEN** Cobra parses the command and its help consistently with AL
- **AND** --repository selects the checkout, with --repo retained as a compatibility alias

### Requirement: Independent projects share one release process

The release tool SHALL own calendar version calculation, stamping, source-ref publication, bundle generation and artifact publication. Each project SHALL release independently using the shared repository version. The workflow SHALL be executable by an operator without CI.

#### Scenario: Release only one project

- **WHEN** an operator builds and publishes one declared project bundle
- **THEN** no other project is required to build or publish
- **AND** its manifest binds version, channel, clean source commit, filenames, sizes and SHA-256 hashes

### Requirement: Guard source refs without coupling patches to trunk

Source-ref publication SHALL preserve immutable tags and verify exact remote ref expectations. Concurrent local operations SHALL use a real exclusive lock. Patch-zero branch and tag creation SHALL be atomic. Patch publication SHALL verify the existing release branch independently of trunk.

#### Scenario: Publish a diverged patch

- **WHEN** a clean tagged patch commit is on an existing release branch whose remote points at that commit
- **THEN** its tag can be published even when trunk points elsewhere
- **AND** a conflicting tag or branch observation causes failure

### Requirement: Standalone release rules

The rules_al Bzlmod module SHALL build independently of this repository using
its declared dependencies. Public rule names SHALL omit al_ and owned skill names
SHALL begin with al-. Repository-specific Git inputs SHALL be explicit.

#### Scenario: Consume release rules independently

- **WHEN** a consumer loads release_bundle from @rules_al//:defs.bzl
- **THEN** the module-owned tool generates a verified bundle without parent labels
- **AND** the same source can build and test as a standalone workspace

#### Scenario: Generate metadata in a linked worktree

- **WHEN** an operator generates Git metadata from a linked worktree
- **THEN** the merged command reads objects and configuration from its common Git directory
- **AND** the output contains its source commit and configured remotes

#### Scenario: Authenticate Git through the controlling terminal

- **WHEN** a foreground version operation starts a Git SSH helper that requests terminal input
- **THEN** the helper can read that input and AL restores the original terminal foreground group afterward
- **AND** cancellation still stops the owned Git group promptly

#### Scenario: Run the independent consumer check

- **WHEN** an operator runs the supplied consumer validation command
- **THEN** it builds the public API with independent Bazel configuration and verifies the resulting manifest and payload
- **AND** a machine-readable result artifact records the verified bytes and identity
