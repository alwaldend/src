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

Extract the Vault configuration into an independent master-based PR. Include
the canonical local-inference DNS document and its minimal export so PKI has
one hostname owner. The user authorized applying infra/vault/tf; review and
apply an exact saved plan through that owner. GPU host deployment is excluded;
DNS application uses the separate local-inference root in PR #233.

## Runtime prerequisite

The standalone AL runner has no apparent-name mapping for the consuming root.
Use the canonical `_main` runfile paths for the shared Vault injector and
Terraform backend plugins, which are declared by this root workspace.

The user also authorized applying the local-inference DNS declarations through
their owning Terraform root after the Vault identity is available.
