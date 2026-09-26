## Why

The collection packages reusable Ansible roles, but it has no repeatable
Molecule tests showing that `host`, `traefik`, and `forgejo` converge and
behave correctly on a fresh machine. Local QEMU tests will exercise those
deployment-used roles without requiring a production inventory or hypervisor.

## What Changes

- Consume the shared Molecule runner defined by the linked tools-owned
  change, reusing the packaged collection as its declared input.
- Add independently runnable scenarios for `host`, `traefik`, and `forgejo`
  under this collection's `extensions/molecule/`, plus an aggregate suite.
- Run each scenario on its own disposable local QEMU VM, using a pinned
  Fedora cloud image, fresh writable disks, and generated test inventory.
- Verify fresh convergence, an unchanged second run, relevant configuration
  updates, and observable host or service behavior. Preserve sanitized
  results and prove cleanup on success, failure, and handled cancellation.
- Fix only narrowly related role defects exposed by these scenarios, such
  as unconditional service reloads, while preserving deployment defaults.
- Keep full deployment-playbook tests, other independently tested roles,
  remote hypervisors, and live infrastructure operations outside this change.

## Capabilities

### New Capabilities

- `role-integration-testing`: Bazel-invoked QEMU scenarios and behavioral
  acceptance for the collection's `host`, `traefik`, and `forgejo` roles,
  including role fixtures, idempotence, configuration changes, and persistence.

### Modified Capabilities

None. Existing collection packaging remains the source of role runtime
inputs; this change adds a test consumer without changing its public layout.

## Impact

Implementation affects collection scenario and suite declarations and any
minimal corrections required in the three roles or their directly exercised
helpers. The shared runner and Ansible dependency integration are owned by
[add-qemu-molecule-runner](../../../../../../tools/molecule/openspec/changes/archive/2026-09-26-add-qemu-molecule-runner/proposal.md).
That change owns generic lifecycle, cleanup, execution, evidence, and
credential-isolation behavior; this consumer does not redefine it.

Acceptance requires all three real role scenarios to pass against the shared
runner and produce inspectable behavioral results. Planning validation alone
does not establish that acceptance. This proposal does not authorize
persistent host configuration or production deployment.
