---
name: repo-ansible
description: Implement and validate Ansible inventories, playbooks, variables, files, and collection usage packaged and executed by Bazel in this repository. Use for infra/**/ansible changes, al_ansible_binary, playbook targets, SSH/Vault injection, or deployment automation.
---

# Work with Ansible

## Follow the packaged execution model

1. Read the component's `ansible/BUILD.bazel`, parent `al.lua`, inventory,
   `group_vars`, and the closest sibling deployment.
2. Keep runtime inputs in the `pkg_files`/`pkg_filegroup` source set. Update its
   globs or explicit entries when adding a new file type or location.
3. Use `al_ansible_binary` with `//infra:ansible_cfg` and the repository Ansible
   collection. Execute through `al_binary_run`/`al_binary_run_map`, not an
   ad-hoc local command that bypasses packaging and injection.
4. Keep Vault injector and SSH/auth plugin labels aligned between `al.lua`,
   Bazel `data`, and `run_args`.

The root `infra/ansible.cfg` establishes Python, roles, collections, and YAML
result behavior. Do not add component-local overrides without a demonstrated
need.

## Write reliable automation

- Use fully qualified collection names where surrounding code does.
- Prefer idempotent modules over `shell` or `command`. When those modules are
  necessary, define `changed_when`, `creates`, or another explicit idempotency
  condition and quote untrusted input.
- Use handlers for restarts and notify them only when configuration changes.
- Preserve inventory/group naming and put defaults at the narrowest correct
  scope.
- Never commit secret values. Receive them through the component's Vault-backed
  injection flow and mark tasks that could reveal them `no_log: true`.
- Be cautious with destructive tasks, quorum changes, host reboots, and rolling
  operations; encode serial/health behavior where appropriate.

## Validate without deploying

Inspect generated labels, then build the packaged binary and use any existing
non-mutating test or syntax target:

```sh
bazel query --config=agent '//path/to/ansible:*'
bazel build --config=agent //path/to/ansible:ansible_bin
```

Do not run the normal `:ansible` or named playbook targets as a syntax check:
they can connect to inventory hosts and mutate infrastructure. If no repository
syntax-check target exists, validate YAML and packaging locally, and report that
live execution was intentionally not performed. Never use `--check` as proof of
safety unless every involved role/module is known to support check mode.
