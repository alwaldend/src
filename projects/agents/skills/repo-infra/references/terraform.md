# Terraform in this repository

## Choose the stage

- `tf_setup` provisions hosts and bootstrap objects. It usually needs a
  provider login plus the `tf=setup` state backend.
- `tf` configures services and API-side resources with the `tf=main` backend.
- Reusable behavior belongs in `projects/tf_modules/<module>`.

Not every component needs both stages. A component that keeps its state in
files or a database may legitimately have no `tf` package at all; document the
reason rather than adding an empty one.

## Package conventions

Terraform packages declare `.tf` files and `.terraform.lock.hcl` in `data`,
pass one or more `al_config` labels, and expose `terraform_binary_map` commands
plus `terraform_test_map` tests. Add module or provider inputs to `data` when
the configuration reads them.

Provider constraints and the lockfile workflow are per package. Reuse the
established `backend "http" {}` block; `tools/vault/tf_backend` supplies its
endpoint, lock URLs, and credentials through the injected `TF_HTTP_*`
environment. The runner's separate `AL_TF_BACKEND_CONFIG_*` support converts
values to literal backend arguments; it is not this plugin's output.

## Implement safely

- Match the existing provider version constraints and lockfile workflow; do not
  regenerate a lockfile by hand.
- Mark genuinely secret outputs and variables `sensitive`, while remembering
  this does not remove values from state.
- Avoid unnecessary resource renames. When a rename is unavoidable, add a
  `moved` block or clearly document the state migration; never accept
  destroy-and-create by default.
- Format changed files before validation.

## Validate without mutation

```sh
bazel_agent bazel query '//path/to/tf:*'
bazel_agent bazel test //path/to/tf:tf_tests.fmt_test
```

Target names may be mapped variants; use query output and the package
`BUILD.bazel` as the source of truth. Keep offline implementation validation
separate from a live plan review: `bazel_agent bazel run //path/to/tf:tf.plan`
contacts the configured Vault, backend, and providers and can acquire a backend
lock. Run it when the requested scope includes that live review. If access is
unauthorized or unavailable, still run formatting and Bazel build/tests, then
report the precise limitation.

Treat `tf.apply`, `destroy`, imports, and state-changing operations as
mutating. A plan that proposes replacements records their impact; executing
those replacements is a separate operation. Execute mutations only within the
explicitly requested and reviewed scope. Summarize a plan for additions,
changes, destroys, replacements, and sensitive or security-impacting changes
without copying secret values.

## Review and apply a saved plan

Keep the backend attributes supplied by `TF_HTTP_*` unset in source and backend
arguments. The plugin creates a fresh loopback endpoint and credentials for
each invocation; the pinned Terraform HTTP backend resolves unset attributes
from the current environment when loading a saved plan. Preserve the owning
`:tf.plan` and `:tf.apply` targets so each establishes that environment. Literal
backend configuration saved in a plan takes precedence over the environment
and can retain an obsolete endpoint. See the
[HTTP backend configuration](https://developer.hashicorp.com/terraform/language/backend/http#configuration-variables).

For an authorized apply, save the scoped plan to an absolute path in private
task scratch, review that exact plan, then pass its path to the owning apply
target. The saved plan already records target selection; do not regenerate it
with different flags during apply. An interactive apply can instead hold at
its confirmation prompt while its proposed changes are reviewed. Both flows
must preserve the authorized scope and the secret handling of plan artifacts.

## Diagnose provider startup failures

When provider schema loading fails before a plan, inspect the underlying
plugin startup diagnostic. A Unix-socket bind error can result from a long
`TMPDIR`: the plugin adds its own filename beneath that directory. If the
measured resulting path exceeds the platform's socket-path limit, shorten
`TMPDIR` within the same ignored task scratch and retain room for that suffix.
Re-run the owning target only after this causal change. A schema-loading error
alone does not justify changing provider versions or disabling validation.
