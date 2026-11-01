## Why

The local inference infrastructure needs its own deployment identity under
the existing AppRole and Yandex folder conventions.

## What Changes

- Add `src_infra_local_inference` using the existing AppRole module.
- Register it in the AppRole discovery and Ansible access groups.
- Add its global DNS access and exact-hostname server PKI enrollment role
  for the requested T3 Code mTLS endpoint.

## Capabilities

No capability contracts change. This extends the existing component inventory;
`skip_specs: true` records that distinction. Host behavior is owned by
`infra/local_inference/openspec/changes/add-yandex-local-inference`.

## Impact

Only checked-in Vault configuration changes. No Vault write or apply runs.
