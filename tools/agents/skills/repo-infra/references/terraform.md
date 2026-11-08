# Terraform in this repository

## Choose the stage

- `tf_setup` provisions hosts and bootstrap objects. It usually needs a
  provider login plus the `tf=setup` state backend.
- `tf` configures services and API-side resources with the `tf=main` backend.
- Reusable behavior belongs in `projects/tf_modules/<module>`.

Not every component needs both stages. A component that keeps its state in
files or a database may legitimately have no `tf` package at all; document the
reason rather than adding an empty one.

Prefer grouping related Terraform resources into modules when they form a
cohesive unit with shared ownership and lifecycle. Keep component-specific
modules with their owning stage, and reuse existing modules before extracting
shared behavior. A module should express a useful boundary, not merely wrap
each resource individually.

## Package conventions

Terraform packages declare `.tf` files and configuration inputs in `data`,
select provider labels through `providers`, and pass one or more `al_config`
labels through the `configs` attribute of the repository's `al_terraform_*`
maps from `//tools/al:al_terraform.bzl`; pass selection labels in `plugin_labels`. The
[reusable Terraform rules](../../../../../tools/rules_terraform/README.md) own
underlying command and test execution. Plain offline checks can use its maps
directly. Reuse `vault_tf_backend` from `//tools/vault/tf_backend:defs.bzl`
when declaring backend calls; component paths and labels remain explicit.
Provider pins belong to
[third_party/terraform](../../../../../third_party/terraform/README.md) and are
shared across roots; source `.terraform.lock.hcl` files are not used.

Reuse the established `backend "http" {}` block; `tools/vault/tf_backend` supplies its
endpoint, lock URLs, and credentials through the injected `TF_HTTP_*`
environment. The runner's separate `AL_TF_BACKEND_CONFIG_*` support converts
values to literal backend arguments; it is not this plugin's output.

## Implement safely

- Keep provider constraints compatible with the shared extension pins. Add or
  upgrade providers through the owning Bazel declaration, following the rule
  documentation; do not acquire providers through a runtime registry fallback.
- Mark genuinely secret outputs and variables `sensitive`, while remembering
  this does not remove values from state.
- Avoid unnecessary resource renames. When a rename is unavoidable, add a
  `moved` block or clearly document the state migration; never accept
  destroy-and-create by default.
- Format changed files before validation.

## Complete one-time migrations

Use temporary `import` and `moved` blocks only while reconciling state. After
the authorized migration succeeds, verify the final state addresses, remove
those blocks, format and validate the configuration, and run the owning
no-change plan without them before final handoff. Do not retain temporary
blocks as migration history; keep a sanitized receipt in the owning change.

Do not remove a still-required migration or apply it merely to satisfy cleanup
when live mutation is outside the granted scope. Record the pending migration
and cleanup prerequisite, preserve its safe transition, and report reconciliation
as incomplete. A permanent compatibility migration needs an explicit maintained
requirement; a completed one-time reconciliation does not establish one.

## Validate without mutation

```sh
bazel query '//path/to/tf:*'
bazel test //path/to/tf:tf_tests.fmt_test
```

Target names may be mapped variants; use query output and the package
`BUILD.bazel` as the source of truth. Keep offline implementation validation
separate from a live plan review: `bazel run //path/to/tf:tf.plan`
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

## Review and apply interactively

Follow the root `AGENTS.md` Terraform output and artifact policy. For inspection,
run the owning `:tf.plan` target without `-out` or `-json`. For an authorized
deployment, run the owning `:tf.apply` target in an interactive terminal. Review
the plan displayed by that invocation and answer `yes` only when it matches the
authorized scope; otherwise answer `no`. A standalone plan is advisory: apply
generates a new plan, so review it again before answering.

Keep the normal text output in the agent transcript and retain only sanitized
summaries as durable evidence. Terraform's sensitive markings provide the
existing redaction; an extra filtering layer is not required. Text output is
not a guarantee that unmarked secrets cannot appear.

Inspect the owning command contract before execution. Some scoped targets still
require a saved plan. Report that incompatibility instead of saving or consuming
a plan, bypassing their guard, or switching to a broader target or raw Terraform.
A change to that contract is separate source work, not an execution workaround.

Keep the backend attributes supplied by `TF_HTTP_*` unset in source and backend
arguments. The plugin creates a fresh loopback endpoint and credentials for
each invocation. Preserve the owning `:tf.plan` and `:tf.apply` targets so each
establishes its own environment. See the
[HTTP backend configuration](https://developer.hashicorp.com/terraform/language/backend/http#configuration-variables).

Finish formatting and all input-file edits before starting apply. Keep inputs
unchanged while reviewing its prompt; if they change, cancel and start a new
interactive apply. A failed apply can still have completed independent resources
or advanced state metadata without changing resources. Account for those
changes and review the next interactive apply's plan for the remaining work
within the existing authorization before retrying.

When using `bazel run --script_path` launchers, regenerate every
owning launcher used, including both plan and apply, after changing declared
data files. `bazel build` can refresh the manifest while an existing
runfiles directory still omits added inputs or retains dangling links to deleted
ones. Verify each launcher's inputs before planning or applying.

## Diagnose provider startup failures

When provider schema loading fails before a plan, inspect the underlying
plugin startup diagnostic. A Unix-socket bind error can result from a long
`TMPDIR`: the plugin adds its own filename beneath that directory. If the
measured resulting path exceeds the platform's socket-path limit, shorten
`TMPDIR` within the same ignored task scratch and retain room for that suffix.
Re-run the owning target only after this causal change. A schema-loading error
alone does not justify changing provider versions or disabling validation.
