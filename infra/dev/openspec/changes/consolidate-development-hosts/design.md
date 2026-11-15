## Context

Secure1 is 192.168.10.101; the obsolete runner1 declaration used 192.168.10.100. Dev1 is the existing workstation at 192.168.1.210.

## Decisions

Use one standard packaged Ansible invocation and the existing host-specific plays. A dedicated src_infra_dev AppRole supplies deployment credentials for both hosts and DNS. Its service policy reads only the existing OpenHands key; its SSH role is limited to the two host names and its PKI role to dev1's two service names. DNS uses scoped provider credential reads and the role's own backend namespace. It receives no access to legacy Terraform state.

Decision review: proceed. A shared deployment identity removes the nested command wrapper and its argument protocol without mixing workstation storage or services into the runner play. Keeping separate authentication inside that command requires the extra orchestration the user asked to remove.

Use one dns root consuming the complete canonical DNS document and the backend at alwaldend.com/vault1/approles/src_infra_dev/tf_backend/dns. Remove tf_setup and its duplicate provider/authentication wrappers. Retain the legacy Vault identities and backend contents until separately authorized state reconciliation is complete; the current root does not consume them. Retain migrated historical records unchanged.

Decision review: proceed. Keeping the two roots perpetuates the split the user explicitly asked to merge. Reusing a legacy backend would keep an obsolete identity in the active configuration. One owner-named backend removes both dependencies, with a documented adoption prerequisite because moving source does not move state. Live state transfer is outside this source change. The Vault-owned policy work is recorded in infra/vault/openspec/changes/consolidate-development-identity.

The inactive Actions worker needs no registration token, instance URL or disk input. Remove its unused injector and variables. Use secure1 and dev1 for inventory groups, playbooks and logical DNS keys. Preserve the simeonwarrenbot account, shared dev_vm role names, Codex host-bot profile and existing ACME account contact because they identify working interfaces or accounts.

## Rollout

Source changes do not deploy infrastructure. Under separate authorization, provision the shared Vault identity and DNS policy, freeze old DNS writers, reconcile both old states and import existing records into the new backend before applying desired DNS changes. Follow the owning dns README for the handoff and key mapping. Then deploy host and ingress configuration. Rollback requires freezing writers and reconciling state ownership against the prior source; reverting source alone is insufficient after adoption. No VM transfer is required.

## Legacy runner retirement

The user clarified that the Proxmox runner does not exist and requested IaC removal only. Remove VM 1200/runner1, its DNS record, the stale host1 inventory entry and Proxmox wiring. No live destruction or state modification is authorized or required by the revised request. Secure1 and dev1 remain unchanged.

Preliminary read-only investigation did not complete: Proxmox was unreachable, and the existing backend referenced a Xen Orchestra provider absent from the old Proxmox package. This pre-existing live-state discrepancy must be reconciled before a live DNS handoff. Preserve the backend contents and legacy access until then.
