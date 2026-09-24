## 1. Declare the component identity

- [ ] 1.1 Inspect the existing AppRole, DNS-access, SSH, and Yandex folder modules; document required and forbidden access before writing policies and verify each permission has a consumer.
- [ ] 1.2 Declare `src_infra_download` and its existing-group membership/outputs; verify source composition grants only its state namespace and required provider integration.
- [ ] 1.3 Declare DNS challenge secret reads and component SSH permissions through reusable modules; verify credentials remain references and no unrelated identity policy changes.

## 2. Validate and coordinate bootstrap

- [ ] 2.1 Run the owning Terraform formatting and packaging checks and review the policy diff; verify both allowed and denied paths using offline evidence without issuing credentials.
- [ ] 2.2 Link provider assignment and component AL prerequisites in the owning documentation; verify all identity names and output references agree with the download plan.
- [ ] 2.3 Under separately authorized live scope, apply the identity before provider assignments and verify scoped authentication; retain sanitized evidence without tokens or secret values.
