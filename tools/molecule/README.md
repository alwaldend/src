---
title: Molecule VM tests
description: Bazel-packaged Molecule scenarios on disposable local QEMU guests
---

Run a Molecule scenario on a disposable local VM:

```sh
bazel_agent bazel test //tools/molecule/test:smoke_test
```

The Linux x86-64 host needs distribution-provided QEMU
(`/usr/bin/qemu-system-x86_64` and `/usr/bin/qemu-img`, with their firmware and
modules), OpenSSH client tools, and network access for guest packages. Missing
executables fail preflight. No host installation or network configuration is
performed. KVM is selected when accessible; otherwise the runner uses TCG with
a longer boot deadline. Explicitly requesting unavailable KVM fails.

Bazel supplies Molecule, Ansible, seed-image tooling, the pinned Fedora cloud
image, and the collection under test. There is no separate Molecule environment
to set up. This shared adapter stages Bazel's collection package mappings,
creates and cleans up the VM, and retains test evidence. Collection scenarios
belong to `projects/ansible_collection/extensions/molecule`; the supported
contract is in the [runner specification](openspec/specs/molecule-qemu-runner/spec.md).

The public `molecule_test` rule takes scenario files and a scenario name,
`PackageFilegroupInfo` collection inputs, a qcow2 image label, CPU/memory/root
disk limits, optional blank data disks, guest ports, accelerator selection,
boot deadlines, restricted Ansible arguments, and diagnostic service names. A scenario
provides `prepare.yml`, `converge.yml`, and `verify.yml`; shared create/destroy
playbooks are supplied by the runner. An optional `update.yml` is followed by
another converge, idempotence, and verify cycle. Scenario files and package
destinations are staged separately, preserving executable input modes.
`ansible_args` accepts only `--diff` and `--skip-tags=tag[,tag]`, with tag names
containing letters, digits, underscores, periods, colons, or hyphens. All other
arguments fail preflight, including connection, inventory, authentication, and
extra-variable overrides. Ansible's reserved selectors (`all`, `tagged`,
`untagged`, `always`, and `never`) are also rejected, including within a tag
list, so they cannot suppress lifecycle execution. Declare scenario variables
in its playbooks.

The Python import path comes from the declared Ansible target so controller
modules use the same pinned dependencies as the CLI. Guest modules use the
guest's Python. OpenSSH uses the run's private key and an isolated configuration;
the runner does not inherit inventory, proxy, Vault, cloud, or SSH-agent settings.

Integration tests are opt-in (`manual`), execute on every explicit invocation
(`external`), and remain local (`no-remote`). Guest package repositories need
network access (`requires-network`). Host QEMU, its libraries, firmware and
modules, and mutable guest package repositories make execution non-hermetic.
Results record executable versions and hashes, including the host QEMU binaries.
Build actions remain sandboxed and offline.

Each run stages its declared candidate under `TEST_TMPDIR`, generates temporary
SSH credentials, and creates only file-backed virtual disks with loopback port
forwarding. It does not inherit developer Ansible configuration or infrastructure
credentials. Selected diagnostics and `result.json` are retained under Bazel's
undeclared test outputs. Private runtime files are never archived wholesale.
Task events retain generated IDs, action names, and outcomes. Rendered play,
task, and handler names are omitted because even `no_log` tasks can expose
values in those headings.

After an abrupt supervisor death, retain the private run directory and use the
runner's `destroy <absolute-run-directory>` command. Recovery checks the manifest
and exact process identity before terminating owned processes. Never use a
process-name-wide kill. Handled SIGINT/SIGTERM runs perform bounded cleanup
automatically; cleanup errors fail the test.

For an interrupted direct invocation, use the private path printed at startup:

```sh
bazel_agent bazel run //tools/molecule/cmd/molecule_runner -- destroy /absolute/path/to/molecule-run
```

Recovery must run in the original PID namespace. It refuses a namespace
mismatch rather than interpreting guest-test PIDs as host PIDs. Bazel normally
reaps the entire test namespace when destroying its sandbox. For a direct runner
invocation outside a sandbox, use the same runner with `destroy` and the exact
private directory reported by the run. There is no persistent QEMU daemon.

The runner's full acceptance harness exercises concurrent guests, failure,
interruption, recovery, and TCG:

```sh
bazel_agent bazel test //tools/molecule/test:acceptance_test
```

No sandbox exception is enabled. TCG uses a 600-second boot deadline, compared
with 240 seconds for KVM; the full test has Bazel's `eternal` timeout. Fedora
package installation is the main source of runtime variation.
