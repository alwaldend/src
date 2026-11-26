---
title: Ansible
description: Bazel rules for ansible
languages:
  - bzl
tags:
  - bzl_rules
---

`al_ansible_package` from `//tools/ansible/rules:al_ansible_package.bzl`
declares mapped sources, their runtime package, an `<name>_bin` executable and
AL command wrappers. For example:

```starlark
al_ansible_package(
    name = "ansible",
    srcs = ["inventory.yaml"] + glob(["group_vars/**/*.yaml", "playbook_*.yaml"]),
    configs = ["//infra/example:al"],
    plugin_labels = {"ansible": "", "vault_env": "default_default"},
)
```

The default command uses `--inventory inventory.yaml playbook_deploy.yaml`.
Pass an `args` dictionary to declare custom command suffixes and arguments;
the macro supplies the executable. Sources retain package-relative paths.
`renames` changes individual destinations, while `packaged_srcs` supplies the
complete ordered list of mapped source groups (including `:srcs_files`).
`srcs_name` and `files_name` override those source target names, and
`binary_visibility` grants access to the packaged executable independently.
The binary includes the shared infrastructure Ansible configuration and
repository collection. Run `bazel test //path/to/package:<name>.syntax` to
check the default playbook or append `.syntax` to a named command to check its
playbook. Each test uses the command's arguments and packaged extra-vars files
without loading deployment plugins.

`ansible_lint` runs `ansible-lint` from Bazel-managed Python
dependencies. A Go test runner exposes the CLI binaries under their
hyphenated names, configures writable Ansible paths, and invokes the
linter against workspace-relative sources.

Python dependencies are declared in `pyproject.toml` and locked in `uv.lock`.
Update the lock with `//tools/ansible:uv_lock.update` and review its
diff. Build `//tools/ansible:uv_lock` to generate it without source writes.
These use the shared [uv locking workflow](../py/README.md), whose sandboxed resolution action
requires network access to PyPI.

The current lock uses ansible-core 2.22 for Python 3.15 targets. Its matching
Ansible community bundle is 15.0.0rc1, selected because the stable community
bundle still requires an older core. The [compatibility change](openspec/changes/support-python-315-targets/design.md)
records that release-candidate constraint and deployment validation.

Molecule is not integrated yet. The likely path is to add Molecule to
these Python requirements and invoke it through a Bazel-native Go runner,
but practical scenarios require choosing and provisioning a Molecule
driver before any scenarios can run.
