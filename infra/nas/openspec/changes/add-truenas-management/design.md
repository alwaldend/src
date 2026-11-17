## Context

The existing TrueNAS 25.10.6 appliance has one RAIDZ2 pool, seven datasets, and
Nextcloud. All nine resources are imported into the owner's independent
`tf_backend/truenas` state. No shares or periodic snapshot tasks are configured.

## Decisions

Use literal Terraform declarations for this appliance. A small inline dataset
table keeps the existing indexed resource addresses; the pool and app retain
their one-element indexes for the same reason. This removes the separate JSON
inventory, generic variables, validation rules, and unused resource types
without moving state.

App datasets use shared resource-level settings. The dataset table contains
only their names and the pool/encryption-root overrides; explicit null values
preserve properties inherited from a parent.

Record only locally configured dataset properties and encryption inheritance.
Protect pool, datasets, and app with `prevent_destroy`; keep recursive storage
deletion disabled. Nextcloud configuration and encryption secrets remain on the
appliance. The API credential comes from the existing Vault connection entry.

The provider pin and the existing DNS root are unchanged. The supported source
surface is this appliance's adopted resources, not a reusable NAS module.

## Validation

Run Terraform formatting and real provider-schema validation offline. Compare
the simplified source against the imported resources with a normal text plan;
accept only a no-change result. No new tests are needed for the removed generic
interface. Repository quality and semantic lint remain delivery gates.
