## 1. Runner cache configuration

- [x] 1.1 Attach the cache disk in VM IaC and inspect provider semantics to preserve the boot disk.
- [x] 1.2 Publish disk device and capacity from Terraform to Vault and inject it automatically into Ansible; remove the manual output/extra-vars interface, retain guest safety checks, and verify provider validation and packaged syntax.
- [x] 1.3 Configure the action cache and 700-GiB GC; document deployment order.

## 2. Verification and delivery

- [x] 2.1 Build the packaged playbook and run Terraform, syntax, lint and repository quality checks without deployment.
- [x] 2.2 Record sanitized verification and deployment prerequisites.

- [x] 2.3 Keep namespace/cache declarations in independent infra/nas/approles/src_infra_dev state and targets; verify NAS operator/backend selection, VM-only provisioning and complete apply order.
