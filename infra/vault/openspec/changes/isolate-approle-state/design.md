## Decisions

Proceed with independent state, not nested modules in the central state. The
user explicitly requested moving AppRoles and prohibited restarting Vault.
Select NAS and development for the active deployment, plus AL for headroom.
Use the existing administrative Vault identity and separate backend paths under
src_infra_dc1_vault/tf_backend/approles/<role>, following the download project.
Preserve live object names, IDs, policy contents and bootstrap membership.

Move NAS and AL source out of the central module tree into their own roots.
The development root owns its AppRole, DNS policy, runner-read policy, SSH and
PKI roles. Core Vault retains shared mounts, global memberships and NAS groups;
look up the moved entity/group IDs by name rather than reading remote state.

## Migration

The immutable protected recovery state is serial 197 and includes the six
completed NAS access operations. Do not rerun apply using the older backend.
Copy that recovery state to private scratch and use the repository-pinned
Terraform state mv with local state/state-out files. Keep all object attributes
and IDs, preserve the core lineage, and verify disjoint complete ownership.

Move NAS module.approle, module.dns and vault_policy.nas_accounts from their
old module.src_infra_nas prefix to the same relative addresses in its root.
Move AL module.approle and module.dns out of module.src_projects_al.
Move development module.src_infra_dev_approle, module.src_infra_dev_ssh,
module.src_infra_dev_pki_server, module.dns_access["src_infra_dev"] and
vault_policy.src_infra_dev_runner to module.approle, module.ssh,
module.pki_server, module.dns and vault_policy.runner in its root.

Verify the new remote states are empty and guard the existing core backend
version against concurrent progress. Push through each owning tf.state target
without force. Retain default backend locks and lineage/serial checks. Stop on
any unexpected existing target state or concurrent write. Confirm each resource
belongs to exactly one final state, then review normal-text no-change resource
plans. Apply output/data-only state reconciliation interactively if required.
Do not restart or reload Vault. Resume the NAS/cache deployment after recovery.

## Failure cases before implementation

- Preserve the recovery original and all resource attributes while moving state;
  never print or commit state contents, credentials or private keys.
- Missing/duplicate migration selections or resource count/attribute mismatches
  must stop before any remote push.
- Root source/packages must include their dependencies and choose one distinct
  backend with the existing operator identity; no remote-state dependency.
- Preserve all policy strings, group/entity IDs, mount paths, DNS declarations,
  SSH domains and PKI permissions; plans must not recreate live objects.
- Reject nonempty destination states, newer core progress or no-force serial/
  lineage failures; never bypass locks or use force to complete migration.
- Clean protected temporary state only after successful remote verification.
