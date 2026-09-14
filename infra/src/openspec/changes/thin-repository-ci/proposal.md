## Why

The Forgejo smoke workflow embeds resource checks and Vault authentication in
YAML. The user requested a discoverable `repo-ci` skill with a hard thin-wrapper
rule, and migration of the current workflow to that rule.

## What Changes

- Document repository CI ownership and require executable logic in Bazel run
  targets or committed scripts.
- Extract the existing smoke implementation into its owning CI package with
  offline authentication regressions and a Bazel entry point.
- Check out the exact event revision with an immutable action pin, and declare
  its Node.js prerequisite in the runner's existing Ansible role.

## Capabilities

### Modified Capabilities

- `repository`: thin CI wrappers and discoverable authoring guidance.

## Impact

Touches skill discovery, the Forgejo workflow, and the runner CI package and
package prerequisites. Existing [runner contracts](../../../../forgejo_runner/openspec/specs/secure-runner-ci/spec.md)
and Vault restrictions remain owned by the component. Acceptance requires
offline checks and an authorized live run of the extracted command.
