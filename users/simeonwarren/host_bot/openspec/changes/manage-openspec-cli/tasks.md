## 1. Reconcile OpenSpec installation

- [x] 1.1 Add the pinned npm installation and inventory variables, share prerequisites, and document the CLI; inspect packaged tasks and resolved prefix/version offline.
- [x] 1.2 Validate the Ansible binary, OpenSpec artifacts, repository quality, and affected-package lint without deploying.

## Validation evidence

On 2026-10-16, the Ansible binary built, strict OpenSpec validation passed,
all 26 repository-quality tests passed, and lint builds passed for the owner,
Ansible, and OpenSpec packages. The packaged runfiles map to the edited playbook
and inventory. Offline YAML/Jinja rendering resolved package
`@fission-ai/openspec`, version `1.14.2`, prefix
`/var/home/simeonwarrenbot/.local`, `become: false`, and `state: present`.
Live installation and idempotency were not exercised; no deployment was requested.
Exact committed-candidate validation is recorded by the delivery receipt.
