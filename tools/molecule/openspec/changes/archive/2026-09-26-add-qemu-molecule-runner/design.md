## Context

See [proposal.md](proposal.md) for scope. The existing `tools/ansible`
executables and `PackageFilegroupInfo` mappings supply Ansible, renamed
service binaries, and collection inputs. The existing Fedora cloud-image
pin supplies the first supported guest. No shared Molecule runner exists.
This cross-package integration and VM lifecycle justify a design artifact.

## Goals / Non-Goals

**Goals:** One reusable lifecycle and test-environment contract with declared
inputs, isolated state, bounded cleanup, and retained evidence.

**Non-Goals:** Role-specific assertions, remote hypervisors, production
inventories, or persistent developer-host provisioning.

## Decisions

### 1. Shared runner with collection-owned scenarios

Add a repository-internal `tools/molecule` package containing a small Starlark
test rule, a Go entry point under `cmd/`, and lifecycle implementation under
`internal/`. Reuse the existing Ansible executables and dependency workflow;
do not maintain a second Ansible environment or install tools with host pip.
Use Molecule's Ansible-native scenario lifecycle rather than introducing an
unverified QEMU driver plugin dependency.

The alternative of copying lifecycle commands into consumers would duplicate
cleanup, environment, and artifact handling. A published standalone Bazel
module is unnecessary for the initial repository-only consumers.

### 2. Local QEMU and a pinned Fedora guest

Start with the existing Fedora 44 x86-64 cloud image. Resolve QEMU,
`qemu-img`, seed-image tooling, and Molecule through explicit pinned tool
inputs; verify compatibility before selecting exact new versions. Tool and
image dependency definitions belong with their existing owner or in
`third_party/`, while the test behavior belongs in `tools/molecule`.

Each run gets a fresh qcow2 overlay, cloud-init seed, and optional blank
virtual data disks. Use headless QEMU, run-specific control sockets, and
user-mode networking with forwarded ports bound to loopback. Do not use
host bridges, tap setup, physical disks, or a privileged shared daemon.
Enable KVM when accessible; otherwise use an explicit, reported TCG fallback
with its own bounded timeout. An explicitly requested unavailable
accelerator fails preflight.

CPU, memory, disk sizes, guest image, accelerator, and forwarded services are
declared inputs. Defaults should fit one small guest, then be calibrated by
the initial end-to-end run. A fresh state directory and port allocation with
bounded collision handling isolate concurrent tests.

QEMU documents [accelerators and user-mode networking](https://www.qemu.org/docs/master/system/invocation.html)
and [backing images](https://www.qemu.org/docs/master/tools/qemu-img.html).
Reusing running VMs would reduce startup cost but weaken independence, so
shared state across Bazel test targets is excluded.

### 3. One lifecycle owner and bounded cleanup

Checked-in shared create/destroy playbooks call the runner's lifecycle
commands using the declared VM fixture. The runner owns the resource manifest
and supervises Molecule; its finalizer uses that same manifest to clean up
after failed phases or signals. Do not introduce a second provisioning state
database or duplicate cleanup logic in individual scenarios.

The normal sequence is create, prepare, converge, idempotence, verify, and
destroy. Configuration-update and persistence checks extend verification
without sharing guests across tests. Capture diagnostics before destruction.
Destroy is repeatable and operates only on resources whose ownership and
process identity match the current run. A cancellation grace period precedes
forced termination of owned children. After SIGKILL or host failure, an
explicit cleanup command checks the manifest and process identity; it never
kills processes by a broad name match.

### 4. Materialize the packaged candidate

The rule consumes scenario files, the existing package mappings, external
collections, CLI tools, and image labels as declared inputs. Materialize the
mapped collection into a writable test project with its expected
`ansible_collections/alwaldend/main` layout, preserving executable modes.
The test uses those exact inputs rather than a globally installed collection.

Generate inventory only for the test guest and use an isolated Ansible
configuration. Construct the child environment explicitly, exposing only
required tool paths and test variables. Writable state lives under
`TEST_TMPDIR`; selected sanitized diagnostics go to
`TEST_UNDECLARED_OUTPUTS_DIR`. Never archive the entire temporary directory.

Consumers own their prerequisite setup and behavioral assertions. The runner
provides shared lifecycle fixtures without copying any role implementation.

### 5. Explicit Bazel execution and evidence policy

These tests are opt-in integration targets, initially tagged `manual`,
`external`, and `no-remote`. The aggregate suite lists them explicitly.
`external` forces actual execution even when a prior result exists; local
execution and remote-cache exclusions alone do not guarantee this. See
[Bazel tag semantics](https://bazel.build/reference/be/common-definitions).

Keep builds sandboxed and offline. Tests may need documented network access
for Fedora packages, and narrowly scoped runtime sandbox exceptions for KVM
or networking if the supported runner requires them. Prove the required
exception before enabling it; never change global Bazel flags or the host's
network configuration. Do not claim full hermeticity while package
repositories remain mutable. Retain installed-package versions and distinguish
package-service failures from role failures.

Each result records the scenario, source/input identity, image checksum,
tool versions, accelerator, package inventory, phase status, verification
summary, and cleanup outcome. The rerun command and declared labels are part
of the evidence. Serial logs and service diagnostics must be sanitized before
retention. Synthetic failure and handled-cancellation runs establish cleanup;
two simultaneous runs establish resource isolation.

## Risks / Trade-offs

- Runtime package repositories remain mutable -> document network needs and
  retain installed versions; do not claim image pinning makes them hermetic.
- Accelerator or sandbox support varies -> validate prerequisites and the
  narrow runtime exception before enabling it; keep build actions offline.
- Port collisions and interrupted processes -> use run-specific ownership,
  bounded allocation retries, and explicit abrupt-exit recovery.
- New consumers could diverge -> keep the generic contract in this owner;
  consumers reference it and define only their own fixtures and assertions.

## Migration Plan

Implement tooling and a minimal VM smoke scenario before enabling the linked
collection tests. Add this owner's README, Bazel declarations, and OpenSpec
validation registration during implementation; this planning revision adds
only OpenSpec metadata and artifacts. Run lifecycle, failure, cancellation,
and concurrency acceptance, then validate the three consumer scenarios in
their owning change. Rollback removes the tooling without migrating production
resources. Archive this owner when its acceptance passes and update consumer
links to the archived change or stable owner specification.

## Implementation evidence (2026-09-26)

The proposal merged in PR #110. Implementation began after fast-forwarding the
clean task branch to `32c0494893c8f3d6a1b66f9fded798bdac92bf88`; no replay was
needed because the source trees were identical.

Acceptance cases were written in `test/acceptance.md`, the smoke playbooks, and
`test/acceptance_test.go` before the runner implementation. The linked collection
mapped its role inputs and assertions in `extensions/molecule/acceptance.md`
before any role corrections.

Dependency decision: proceed with hermeticbuild's published static QEMU
11.0.0.1 distribution (QEMU 11.0.0), including its matching firmware. Unlike
using host QEMU, these declared inputs are pinned by content. Building QEMU
and its native dependency closure from source would add a separate build-system
project. The trade-off is trusting this third-party binary publisher; its
source recipe and release attestations are public. Downloaded system, image,
and firmware archives matched publisher checksum files and GitHub asset
digests. The Bazel-provided system binary reports QEMU 11.0.0.

Initial smoke failures exposed localhost Ansible temporary-directory isolation
and QEMU daemonization's changed working directory. The runner now selects a
private localhost temporary path and validates the open overlay for cleanup.
Fedora cloud-init reported only a hostname-setting warning; the fixture now
preserves the image hostname because changing it is not a test prerequisite.
At that stage, remaining acceptance tasks were unchecked. Raw
transient diagnostics live under ignored `out/molecule-role-tests`.

A sandboxed KVM smoke run completed create, prepare, converge, idempotence,
reboot verification, and destroy; source/image hashes, package inventory, and
cleanup success were inspected. The first executable acceptance suite passed
concurrency, repeated fresh execution, injected assertion failure, and handled
cancellation. Recovery and package-mapping acceptance were extended in the later run below. A controller-module failure identified that Ansible's plain Python
subprocess lacked Bazel's dependency imports; the runner now constructs its
Python path from the declared `PyInfo` rather than installing host libraries.

Further acceptance corrections use the CLI spelling `side-effect`, wrap the
executable mapping fixture in `pkg_filegroup`, and preserve separate numbered
phase logs for initial and updated applications. Recovery validates the PID
namespace as well as process start times, unique QEMU names, and the open owned
overlay. The harness additionally exercises TCG, bounded boot failure, abrupt
supervisor death, and repeatable scoped recovery. Native Bazel aliases do not
execute tests when selected with `bazel test`; collection convenience labels
therefore use single-test `test_suite` wrappers.

A downloaded guest-image cache entry differed from the existing Fedora pin
before role execution (observed SHA-256 `33be10d98d890758a8a98bc6a3cc9b467345d649ffd8d6410b071ed3bd4c7b1e`).
Its cause was not established. A repository-owned `bazel fetch --force
--repo=@org_fedora_cloud` restored the pinned bytes
(`28680fe5b371a5a82ebf43a31926e086a168e59949d03969c5093e7071f90b7f`).
Earlier role failures remain diagnostic observations, not final acceptance;
acceptance must use and retain the restored input identity.

Candidate `3ee3787109a858b26f6cbfe2dc2d7a49396ba6a1` passed the expanded
runner acceptance and explicit smoke test on the restored Fedora pin. The
acceptance includes KVM and forced TCG, concurrent and repeated runs, renamed
executable mappings, assertion failure, missing input and accelerator checks,
boot deadline, SIGTERM, and SIGKILL recovery followed by two scoped destroys.
Every completed run reported successful cleanup; recovery retained its own
structured receipt. The artifact scan found no isolation sentinel or private
key PEM markers. Root quality (27 tests including the two OpenSpec validations),
affected semantic lint, and candidate compilation also passed. The linked
role suite was incomplete at that candidate: Traefik passed, host failed a fixture trust-path
assertion after passing idempotence, and Forgejo still rewrote configuration.

## Accepted result

Published candidate `1bfe561c394e52d514b2bc756154ff24cc07a320` passed all three
collection scenarios, explicit smoke execution, affected semantic lint,
collection/runner compilation, and all 27 repository quality/OpenSpec tests.
The runner acceptance passed at `3ee3787109a858b26f6cbfe2dc2d7a49396ba6a1`;
the later candidate changed only two consumer fixtures and OpenSpec records.
Runner source, rule, smoke/acceptance fixtures, tool pins, image, and executable
configuration were unchanged. The two explicit smoke invocations retained equal
input hashes and tool versions but different run IDs, proving fresh execution.
Both used the pinned Fedora image digest recorded above.

Local evidence is under `out/molecule-role-tests/candidate-3ee3787` and
`out/molecule-role-tests/candidate-1bfe561`; delivery validation is bound by
`fixture-corrections-prepare.json.validation.json`. All completed scenarios
reported successful cleanup and their private runtime directories were absent.
The 55 retained consumer/smoke artifact files contained no isolation sentinel
or private-key PEM markers. These are local reproducible test artifacts, not
published credentials or production validation.

[Implementation PR #111](https://github.com/alwaldend/src/pull/111) was published
and verified against the accepted candidate. Archival and final documentation
change no runtime inputs; final publication reruns lint, OpenSpec validation,
and repository quality checks.

## Session review

The task-local `ergonomics.json` records bounded evidence and corrections for
`molecule-test-alias-selection`, `molecule-forgejo-config-normalization`,
`molecule-fedora-trust-lookup`, `molecule-large-result-display`,
`molecule-merged-proposal-branch`, `molecule-validation-plan-permissions`, and
`molecule-observer-stop-identity`. Use single-test suites, summarize result
fields, keep validation plans private, and use bounded observer lifetimes.
Mutable Fedora package setup and real KVM/TCG acceptance account for expected
verification latency. No shared skill or host configuration was changed.
