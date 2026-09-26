---
title: Molecule VM tests
description: Bazel-packaged Molecule scenarios on disposable local QEMU guests
---

This tool owns the shared runner and VM lifecycle. Collection scenarios belong
to `projects/ansible_collection/extensions/molecule`. The supported contract is
recorded in the [runner specification](openspec/specs/molecule-qemu-runner/spec.md).

Linux x86-64 is the initial platform. Bazel supplies Molecule, Ansible, QEMU,
qemu-img, firmware, seed-image tooling, and the pinned Fedora cloud image.
The guest input is owned by `//third_party/org_fedora_cloud`, QEMU by
`//third_party/com_github_hermeticbuild_qemu`, and the Python tools by
`//tools/ansible`. Each result records their observed versions and input hashes.
The host must provide OpenSSH client tools. KVM is selected when accessible;
otherwise the runner reports TCG and uses its longer boot deadline. Requesting
KVM explicitly fails if it is unavailable. No host installation or network
configuration is performed.

The public `molecule_test` rule takes scenario files and a scenario name,
`PackageFilegroupInfo` collection inputs, a qcow2 image label, CPU/memory/root
disk limits, optional blank data disks, guest ports, accelerator selection,
boot deadlines, Ansible arguments, and diagnostic service names. A scenario
provides `prepare.yml`, `converge.yml`, and `verify.yml`; shared create/destroy
playbooks are supplied by the runner. An optional `update.yml` is followed by
another converge, idempotence, and verify cycle. Scenario files and package
destinations are staged separately, preserving executable input modes.

The Python import path comes from the declared Ansible target so controller
modules use the same pinned dependencies as the CLI. Guest modules use the
guest's Python. OpenSSH uses the run's private key and an isolated configuration;
the runner does not inherit inventory, proxy, Vault, cloud, or SSH-agent settings.

Integration tests are opt-in (`manual`), execute on every explicit invocation
(`external`), and remain local (`no-remote`). Guest package repositories need
network access (`requires-network`); their contents are mutable, so these tests
are not fully hermetic. Build actions remain sandboxed and offline.

```sh
bazel_agent bazel test //tools/molecule/test:smoke_test
bazel_agent bazel test //tools/molecule/test:acceptance_test
```

Each run stages its declared candidate under `TEST_TMPDIR`, generates temporary
SSH credentials, and creates only file-backed virtual disks with loopback port
forwarding. It does not inherit developer Ansible configuration or infrastructure
credentials. Selected diagnostics and `result.json` are retained under Bazel's
undeclared test outputs. Private runtime files are never archived wholesale.

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

Both KVM smoke execution and user-mode networking have run successfully under
the normal Linux test sandbox. No sandbox exception is enabled. TCG uses a
600-second boot deadline, compared with 240 seconds for KVM; the full test has
Bazel's `eternal` timeout. Fedora package installation is the main source of
runtime variation.
