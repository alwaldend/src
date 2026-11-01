## 1. Provisioning

- [x] 1.1 Declare GPU, encrypted storage, KMS, network, and SSH restrictions;
      verify Terraform formatting and offline provider validation.
- [x] 1.2 Package the component AL configuration and prerequisite AppRole;
      verify Bazel analysis and inspect injection labels.

## 2. Host configuration

- [x] 2.1 Install verified runtime artifacts and an operator-supplied model;
      verify packaged Ansible syntax and artifact/checksum failure handling.
- [x] 2.2 Configure local Codex policy and inference service; render
      representative files and validate TOML and units. The T3 extension
      replaces the initially validated account-level egress restriction.

## 3. Delivery

- [x] 3.1 Document inputs, deployment boundaries, and future GPU/tool/network
      acceptance; verify documentation against packaged commands.
- [x] 3.2 Run offline validation and repository delivery checks; record exact
      results before delivery. Git and delivery receipts own publication state.

## 4. Authorized DNS execution

- [x] 4.1 Carry the shared injector/backend runfile fixes from PR #242 so the
      owner commands can run before that prerequisite merges.
- [x] 4.2 Review and apply only the `module.dns` saved plan; verify the public
      records and retain sanitized execution evidence. GPU provisioning remains
      outside this operation.

## 5. DNS extraction

- [x] 5.1 Remove Cloudflare resources, provider selection, credential calls and
      scoped DNS commands from provisioning; retain Yandex zone/address inputs.
- [x] 5.2 Verify standalone provisioning with backend-disabled Terraform
      validation and owner checks; prepare the cleanup alongside the DNS PR.

## 6. Terraform extraction

- [x] 6.1 Extract tf_setup into a separate master-based PR with stage-owned
      AL wiring, preserving resource definitions and the backend path.
- [x] 6.2 Validate host-only packaging and prepare PR #233 without Terraform
      resources or Yandex credentials. Delivery receipts own publication state.
