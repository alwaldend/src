## 1. Declare the component identity

- [x] 1.1 Inspect the existing AppRole, DNS-access, SSH, and Yandex folder modules; document required and forbidden access before writing policies and verify each permission has a consumer.
- [x] 1.2 Declare `src_infra_download` and its existing-group membership/outputs; verify source composition grants only its state namespace and required provider integration.
- [x] 1.3 Declare DNS provisioning secret reads, component SSH permissions, and scoped internal ACME/EAB permissions through reusable modules; verify credentials remain references and no unrelated identity policy changes.

## 2. Validate and coordinate bootstrap

- [x] 2.1 Run the owning Terraform formatting and packaging checks and review the policy diff; verify both allowed and denied paths using offline evidence without issuing credentials.
- [x] 2.2 Link provider assignment and component AL prerequisites in the owning documentation; verify all identity names and output references agree with the download plan.
- [x] 2.3 Under separately authorized live scope, apply the identity before provider assignments and verify scoped authentication; retain sanitized evidence without tokens or secret values.

Source and permission review evidence lives with the
[download implementation](https://github.com/alwaldend/src/blob/39b2dd9cea5fea84705a36cc6003fb5f46ba0547/infra/download/openspec/changes/add-static-hosting/evidence.md).
Live bootstrap and authentication evidence is recorded below.

PR review moved the identity, DNS policy, SSH, and PKI composition to the
standalone `infra/vault/approles/src_infra_download` root with its own backend.
Shared identity IDs are looked up by name; the central stage still owns its
membership lists and runs after the standalone root. The first deployment
required no state migration.

## Deployment evidence: 2026-09-25

The user authorized this operation with "Apply vault". The standalone root
was applied from revision `5ee276f1bdaa6435fb6b905664f545ee98faa4ee` using a
reviewed saved plan: 17 resources added, none changed or destroyed. The core
stage then applied a saved plan targeted to `vault_identity_group.approles`
and `vault_identity_group.ansible`: two updates, each adding only the download
membership, with no existing members removed.

Live planning found that the source omitted the Forgejo runner's existing
Ansible membership and would add a Yandex folder policy to that runner.
Preserving the membership and explicitly disabling that optional policy in
source reduced the targeted plan to exactly the two download additions.
The initial full core plan also proposed unrelated policy/DNS changes and
three Forgejo CI deletions. That full plan was not applied; this deployment
does not establish that the entire core root has converged.

Post-apply plans returned exit code zero with no changes for both the
standalone identity and the two targeted groups. Authentication through
`//infra/download:vault` succeeded. A `sys/capabilities-self` query verified:

- Read/write access to download's own state namespace.
- Read-only access to the two declared DNS credential paths.
- SSH host/client signing and role-scoped ACME EAB access.
- Denied access to the core Vault state and an unrelated EAB role.

The checks queried capabilities without reading credentials or issuing
certificates. Sanitized plan/action summaries are retained under
`out/vault-apply/`; raw plans, responses, and temporary runtime files are
removed after use. Yandex/XCP-ng assignments and download host deployment
remain with their owning changes and were not applied here.
