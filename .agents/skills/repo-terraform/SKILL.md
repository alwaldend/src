---
name: repo-terraform
description: Implement, review, format, plan, and validate Terraform managed through this repository's Bazel macros and Vault-backed al.lua configuration. Use for .tf files, Terraform modules, providers, state backends, locks, tf_setup/main stages, or tf.plan/tf.apply targets.
---

# Work with Terraform

## Understand the component

1. Read the component `BUILD.bazel`, parent `al.lua`, sibling `.tf` files, and
   the closest similar component or module.
2. Determine whether the change belongs in `tf_setup` (bootstrap/VM setup), `tf`
   (service configuration), or a reusable module under `projects/tf_modules`.
3. Inspect the package targets with `bazel query '//path:*'`; generated names
   and label maps vary by package.

Terraform packages normally declare `.tf` files and `.terraform.lock.hcl` in
`data`, pass one or more `al_config` labels, and expose commands with
`terraform_binary_map` plus tests with `terraform_test_map`. Authentication and
backend configuration are injected by plugins named in `run_args`; never replace
that flow with committed credentials.

## Implement safely

- Match the existing provider constraints and lockfile workflow.
- Put reusable behavior in the established module tree; keep environment values
  in component configuration.
- Mark genuinely secret outputs and variables `sensitive`, while remembering
  this does not remove values from state.
- Avoid unnecessary resource renames. When unavoidable, add an appropriate
  `moved` block or clearly document the state migration.
- Format changed Terraform files before validation.
- Never hand-edit Terraform state or commit plan/state files.

## Validate without mutation

Prefer repository-provided targets over invoking a host Terraform binary:

```sh
bazel query '//path/to/tf:*'
bazel test //path/to/tf:tf_tests.fmt_test
bazel run //path/to/tf:tf.plan
```

Target names may be mapped variants; use query output and the package
`BUILD.bazel` as the source of truth. A plan can require Vault, backend, provider,
network, or cloud access. If unavailable, still run formatting and Bazel
build/tests, then report the precise limitation.

Treat `tf.apply`, `destroy`, imports, state operations, and replacement plans as
mutating operations. Do not run them unless explicitly requested, authorized,
and reviewed. Summarize a plan for additions, changes, destroys, replacements,
and sensitive/security-impacting changes without copying secret values.
