## 1. Implement

- [x] 1.1 Implement extensible XCP-ng provisioning and retire the operator-confirmed
      absent legacy VM; inspect state and pass Terraform formatting/DNS checks.
- [x] 1.2 Configure host/dev_vm/runner deployment and Vault-backed registration;
      pass packaged Ansible syntax checks and render representative configuration.
- [x] 1.3 Configure restricted CI JWT authentication; inspect exact bound claims
      and apply the scoped Vault plan with no unrelated changes or deletions.

## 2. Deploy and verify

- [x] 2.1 Apply the scoped XCP-ng resource-set plan and verify runner access.
- [x] 2.2 Provision secure and its DNS record; verify saved plans, VM resources
      and DNS resolution and retire the stale legacy DNS record.
- [x] 2.3 Run packaged registration and Ansible deployment; verify runner
      service health, repository registration and repeat-deployment behavior.
- [x] 2.4 Run the smoke workflow on the declared protected validation branch;
      verify capacity, Vault token restrictions and revocation for the exact commit.
      Runner 13.1.0 and the approved Vault clock repair are deployed; CI run 4 passes.

## 3. Deliver

- [x] 3.1 Run required repository quality and affected semantic checks; record
      candidate identity and observed deployment/CI evidence.
- [x] 3.2 Archive verified specifications and deliver the committed feature
      branch and pull request through repo-delivery.
