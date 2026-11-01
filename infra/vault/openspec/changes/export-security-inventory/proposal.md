## Why

The shared security workflow must consume this owner's existing inventory.

## What Changes

Expose `//infra/vault/ansible:inventory` to
`//infra/security/ansible:__pkg__`, preserving existing access.
The target contains only the authoritative `ansible/inventory.yaml`.

## Capabilities

No runtime behavior changes; this record covers build visibility only.
Inventory values, group membership, and host configuration remain unchanged.

## Impact

One owner-local BUILD declaration. Consumer packaging and execution belong
to `infra/security`. No deployment is authorized by this record.
