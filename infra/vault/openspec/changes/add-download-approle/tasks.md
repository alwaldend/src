## 1. Declare the component identity

- [x] 1.1 Inspect the existing AppRole, DNS-access, SSH, and Yandex folder modules; document required and forbidden access before writing policies and verify each permission has a consumer.
- [x] 1.2 Declare `src_infra_download` and its existing-group membership/outputs; verify source composition grants only its state namespace and required provider integration.
- [x] 1.3 Declare DNS provisioning secret reads, component SSH permissions, and scoped internal ACME/EAB permissions through reusable modules; verify credentials remain references and no unrelated identity policy changes.

## 2. Validate and coordinate bootstrap

- [x] 2.1 Run the owning Terraform formatting and packaging checks and review the policy diff; verify both allowed and denied paths using offline evidence without issuing credentials.
- [x] 2.2 Link provider assignment and component AL prerequisites in the owning documentation; verify all identity names and output references agree with the download plan.
- [ ] 2.3 Under separately authorized live scope, apply the identity before provider assignments and verify scoped authentication; retain sanitized evidence without tokens or secret values.

Source and permission review evidence lives with the
[download implementation](https://github.com/alwaldend/src/blob/39b2dd9cea5fea84705a36cc6003fb5f46ba0547/infra/download/openspec/changes/add-static-hosting/evidence.md).
Live authentication and bootstrap remain unperformed and require separate scope.

PR review moved the identity, DNS policy, SSH, and PKI composition to the
standalone `infra/vault/approles/src_infra_download` root with its own backend.
Shared identity IDs are looked up by name; the central stage still owns its
membership lists and runs after the standalone root. The new identity has not
been deployed; no live state migration was performed.
