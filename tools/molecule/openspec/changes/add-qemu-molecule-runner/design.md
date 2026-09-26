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
