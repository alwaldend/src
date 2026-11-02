## ADDED Requirements

### Requirement: One standalone repository lifecycle module

The standalone rules_al module SHALL provide the AL executable and its Bazel APIs.
AL SHALL expose version, release, delivery, bazel test, and doctor directly, without a repo
namespace. Lifecycle commands SHALL preserve existing guards and receipts and
SHALL initialize credential plugins only when the selected operation requires them.
Artifact publication SHALL use the publish command, and repository selection SHALL
use --repository. Retired deploy, --repo, and --soras_path compatibility spellings
SHALL not be accepted.

The independent AL project SHALL live under projects/al while retaining its
rules_al Bazel module identity. Reusable runtime behavior SHALL obtain
consumer-repository-specific policy and integration values from configuration,
not require the alwaldend/src repository or its source layout.

#### Scenario: Configure AL for another repository

- **WHEN** another repository supplies AL's required configuration and declared dependencies
- **THEN** the reusable executable and rules can operate without an alwaldend/src checkout
- **AND** repository-specific publication destinations and validation policy remain explicit inputs

#### Scenario: Use repository lifecycle commands without plugin startup

- **WHEN** an operator inspects a version or delivery provider, or generates a bundle
- **THEN** the AL command executes the owning implementation without starting plugins
- **AND** publication guards and the independent project version contract remain enforced

#### Scenario: Build with ordinary Bazel

- **WHEN** an operator builds a project with bazel build --config=release
- **THEN** Bazel obtains the current checkout's version identity from the installed AL workspace-status command
- **AND** no nested Bazel invocation or AL build wrapper is required

#### Scenario: Create local release refs separately from checkout

- **WHEN** an operator invokes al release create from clean trunk
- **THEN** AL atomically creates the weekly release branch and initial version tag locally without changing HEAD or pushing refs
- **AND** the operator selects an existing release branch through ordinary Git checkout, without an AL checkout command
- **AND** the retired al version release-start spelling is rejected

#### Scenario: Use personal agent defaults

- **WHEN** an agent invokes Bazel with a configured user rc file
- **THEN** Bazel applies that user's output and test defaults directly
- **AND** AL does not inject a repository-specific agent profile or couple ordinary checks to delivery preparation

#### Scenario: Rebase before source publication

- **WHEN** an operator invokes direct delivery rebase for an eligible task-owned feature branch
- **THEN** AL replays it locally and reports the exact candidate and fetched remote lease
- **AND** it leaves the remote branch and pull request unchanged until separate prepared publication

#### Scenario: Run an installed standalone command

- **WHEN** an installed AL command runs without a selected runfiles-backed requirement
- **THEN** missing Bazel runfiles do not prevent child execution
- **AND** selected runfiles-backed inputs still require valid resolution

#### Scenario: Replace command output and reject conflicting exports

- **WHEN** an operator selects a named stdout or stderr output
- **THEN** AL creates or truncates that output instead of retaining trailing bytes
- **AND** conflicting selected plugin environment values fail before the child starts without exposing their values

#### Scenario: Cancel a version operation

- **WHEN** a version inspection or source-ref publication receives SIGINT or SIGTERM
- **THEN** its cancellation context reaches the active subprocess and the command fails promptly

#### Scenario: Use native version subcommands

- **WHEN** an operator supplies version flags before or after the selected subcommand
- **THEN** Cobra parses the command and its help consistently with AL
- **AND** --repository selects the checkout without compatibility aliases

### Requirement: Prepare source independently of checks

`al delivery prepare` SHALL perform Git preparation and record the exact
prepared candidate and publication lease without selecting or invoking Bazel.
It SHALL NOT accept Bazel arguments or emit hardcoded affected-target, validation
gap, or suggested-check reports. Build, test, and lint requirements SHALL remain
consumer procedures executed separately. The retired delivery validate,
continue, and deliver commands, JSON validation plans, and validation ledger
SHALL be removed.

#### Scenario: Prepare a source candidate

- **WHEN** an operator prepares a task-owned source change
- **THEN** AL records the prepared Git candidate, scope, and remote expectations
- **AND** no build or test process is started by preparation

#### Scenario: Run checks separately

- **WHEN** an operator runs the repository's required checks after preparation
- **THEN** those checks use their ordinary command interfaces independently of delivery receipts
- **AND** no delivery validation plan or result ledger is required

### Requirement: Publish the exact prepared source

Source publication SHALL consume the preparation receipt and enforce its prepared
HEAD, index, tree, scope, and remote lease guards without a separate validated-head
assertion. Existing repository, pull-request, publication verification, and review
guards SHALL remain. Publication SHALL preserve unrelated local edits outside
the prepared scope when no rebase is required; review addressing SHALL retain
its clean-candidate requirement.

#### Scenario: Publish a prepared candidate

- **WHEN** the current candidate matches a valid preparation receipt and publication guards hold
- **THEN** AL publishes that exact prepared commit through the recorded remote lease
- **AND** it does not claim that the receipt certifies test results

#### Scenario: Rebase after the base advances

- **WHEN** source publication discovers that the remote base requires rebasing the prepared candidate
- **THEN** AL updates the candidate and receipt and returns checks_required before pushing
- **AND** the operator can run checks and publish again with that receipt

#### Scenario: Recover a failed review observation after publication

- **WHEN** source push and pull-request update succeed but the subsequent review inventory read fails
- **THEN** publication evidence is retained with publication_verification_pending
- **AND** verification can recover the observation without repeating the push or changing the pull-request identity

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
its declared dependencies. AL-owned rule and rule-macro names SHALL begin with
al_, provider names SHALL begin with Al, and owned skill names SHALL begin with
al-. Repository-specific Git inputs SHALL be explicit.

#### Scenario: Consume release rules independently

- **WHEN** a consumer loads al_release_bundle from @rules_al//:defs.bzl
- **THEN** the module-owned tool generates a verified bundle without parent labels
- **AND** the same source can build and test as a standalone workspace

#### Scenario: Release without a full repository history export

- **WHEN** an operator builds a publishable project bundle
- **THEN** its source identity comes from version, channel, source commit, and tree status
- **AND** no full-history JSON export or website data projection is required

#### Scenario: Include an optional changelog

- **WHEN** a release declaration supplies a Git history bundle
- **THEN** changelog generation preserves commit and annotated or lightweight tag associations
- **AND** the optional history input does not become a requirement for ordinary artifact publication

#### Scenario: Authenticate Git through the controlling terminal

- **WHEN** a foreground version operation starts a Git SSH helper that requests terminal input
- **THEN** the helper can read that input and AL restores the original terminal foreground group afterward
- **AND** cancellation still stops the owned Git group promptly

#### Scenario: Run the independent consumer check

- **WHEN** an operator runs the supplied consumer validation command
- **THEN** it builds the public API with independent Bazel configuration and verifies the resulting manifest and payload
- **AND** a machine-readable result artifact records the verified bytes and identity

### Requirement: AL project landing page

The main site SHALL expose AL at /projects/al/ through the shared project registry.
The standalone AL module SHALL own the landing content without depending on parent
site layouts or configuration.

#### Scenario: Discover AL from the project index

- **WHEN** a visitor opens the project index and selects AL
- **THEN** the landing describes current repository lifecycle capabilities
- **AND** it links to AL's documentation under /docs/projects/al/
