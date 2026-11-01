## Why

The shared security workflow must consume this owner's existing inventory.
The XCP-ng host also needs an explicit Ansible interpreter selection that
survives the security consumer's group variables.

## What Changes

Expose `//infra/xcp_ng/ansible:inventory` to
`//infra/security/ansible:__pkg__` and
`//infra/security/appliances:__pkg__`, preserving existing access.
The target contains only the authoritative `ansible/inventory.yaml`.
Expose the existing root AL configuration to the appliance package so it can
select the owner-controlled XCP-ng and XOA SSH credential injectors.
Set `ansible_python_interpreter` on the existing XCP-ng host entry to
`/opt/ansible-python/3.14.4/bin/python3.14`, without duplicating inventory or
host configuration. Document the user-authorized manual bootstrap exception
and its preservation requirements in the
[Ansible README](../../../ansible/README.md).

## Capabilities

No new capability specification. This record covers build visibility and
owner-local interpreter configuration; inventory membership is unchanged.

## Impact

Existing owner BUILD exports, inventory, and documentation. Consumer packaging
and execution belong to `infra/security`. The interpreter follows the existing
repository CPython pin; no new third-party pin is introduced. Runtime
installation and native module verification succeeded according to the parent;
shared-role audit validation remains pending. The parent implementation agent
must review the candidate before deployment. The documented manual exception
permits no reboot. Deployment authority remains with the user's explicit
instructions.
