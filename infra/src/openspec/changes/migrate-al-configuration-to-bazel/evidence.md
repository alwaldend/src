## Current source acceptance

All repository AL Lua configurations and shared Lua helpers are removed,
including the runtime XO bootstrap configuration. Existing component aggregate
labels remain. Small native targets compose through configs; explicit extensions
select overrides independently of aggregate reference order. Component auth
extends //tools/al:auth. The formatter may sort aggregate references normally.

Short, single-use declarations call existing macros directly in BUILD.bazel.
Local defs.bzl extraction is reserved for significant reuse or a substantial
readability gain from complex graphs, long templates or bulky settings. Common infrastructure patterns and
Vault Terraform backend calls use reusable owner-provided target macros.
Repository Terraform maps own AL wrapper selection, accept configs and
plugin_labels directly, and preserve standalone rules_terraform boundaries.
The repo-bazel skill records namespace, locality and target-composition rules;
its offline evaluation contracts cover these choices.

The captured 33-configuration baseline compares equivalent by effective named
definition, including executable/TLS inputs and plugin payloads. A further 288
selector cases preserve effective resource maps and backend-call sequences.
The few intentional resource/backend precedence relationships use explicit
call extensions. Aggregate reference ordering has no effect on auth overrides.
Comparison evidence: out/migrate-al-config/composition-equivalence.json.

Source validation passes 204 selected root tests, including repository quality,
37 offline Terraform format tests, AL/command/injector tests, skill configuration
checks and OpenSpec validation. Affected semantic lint passes. The linked
[AL evidence](../../../../../projects/al/openspec/changes/native-configuration-targets/evidence.md)
owns standalone and independent-consumer acceptance.
Prepared-candidate quality, lint and local review remain required by al-delivery;
Git and trusted delivery receipts own commit, publication and review status.

The native XO bootstrap target declares only the shared Vault connection and
its plugin. Approved target and optional issuer AppRoles are selected and
validated in the plugin process at startup, preserving the previous boundary
without evaluating Lua or capturing invocation values in build outputs.
Existing validation covers allowed roles, invalid targets and invalid issuers.
No live authentication, plan, apply, deployment, host configuration change or
new external dependency was performed.

## Resolved defects and bounded procedure evidence

- TLS-RUNFILE-LOOKUP: the initial review exposed relative runfile identifiers
  being opened beneath the working directory. Runtime TLS initialization and
  injector exports now share resolution; directory and manifest consumer
  fixtures verify real local requests outside the runfiles tree.
- CONFIGURATION-ORDER: the initial implementation relied on merge input order.
  Explicit target ancestry replaces that workaround and rejects ambiguous
  sibling definitions. No order-preservation directive remains in config lists.
- COMMAND-FIXTURE-SCHEMA: cleanup and dictionary fixtures consume the current
  structured resource maps and inherit the injector executable through config.
- LINKED-WORKTREE-LFS-ACCESS: normal Git works with the required filesystem
  permission for shared worktree metadata; no task-local LFS storage workaround
  or persistent host change was introduced.
- FORMATTER-INPUT-PATH: one formatter attempt used paths relative to its runfiles
  directory. Absolute source paths fixed invocation; configured formatting and
  quality checks pass. Reuse the absolute-path procedure, not another workaround.

The built-in injection engine and resource-provider registration protocol were
explicitly deferred. The injector remains an external AL-owned plugin.

## Proportional extraction correction

The cleanup inlines 42 trivial local wrappers and removes 18 unused defs.bzl
modules, including infra/github/defs.bzl. Six local macros retain the larger
resource graphs, SSH templates and long certificate/OIDC settings. The existing
shared injector/infrastructure macros remain reusable and are called directly.
Guidance and the repo-bazel offline evaluation cases require a significant
readability or reuse benefit rather than extracting every resource dictionary.

All 39 captured production and fixture configurations compare identical after
JSON parsing; names, labels, payloads, declaration order and credential references
are unchanged. Evidence: out/migrate-al-config/extraction-equivalence.json.

## Flat injector groups and target-derived defaults

The public al_injector_call API accepts named `res_<type>` maps. Structured entries contain explicit dependency and Vault
metadata alongside type settings; the macro translates to the unchanged runtime
args.res schema. All 51 existing injector declarations use this shape, including
shared infrastructure helpers. The former narrow environment helper is removed.
Short declarations remain inline and complex local graphs keep their owners.

Call identity defaults to the target name. Generic plugin calls also default the
selected plugin to their target name; explicit config_name and plugin overrides
remain independent. Dictionary entries use their key, keeping generated private
target names out of runtime identities. Eight redundant explicit plugin selections
were removed from existing declarations.

The independent consumer verifies all eight resource kinds, per-resource metadata,
ordered dependencies, typed process/OIDC name fields, public-file substitutions,
default and overridden identities, and rejection of cross-kind duplicate names
and scalar settings. It retains result.json and the expected failure diagnostics
under out/migrate-al-config/flat-prefixed-consumer. Its earlier run failed on the
missing public API before implementation. All 39 production/fixture outputs
remain identical after JSON parsing. Actual AL/plugin/child cleanup, resource
selection, graph validation and Terraform fixtures retain their existing behavior.
No runtime protocol, generic injection engine, provider registration, live
infrastructure operation or external dependency is introduced.

Resource groups use res_ prefixes and are extracted from kwargs. The suffix is
forwarded as the resource type, so future runtime kinds do not require new macro
parameters or an additional type list. Other call attributes pass through.
The finalized prefixed consumer passes all API, naming, rejection and TLS checks.
