## Current source acceptance

AL exposes small authentication, connection, plugin and call targets plus
al_config composition and dictionary inputs for small declarations. Atomic
extends relationships select the most specific complete named definition;
aggregates and wrappers deduplicate shared ancestors and reject unrelated
competitors. Differently named call descendants follow their ancestors for
intentional resource/backend overrides. Aggregate reference order is irrelevant.

Plugin arguments use args, executable arguments use argv, and injector args.res
is a dictionary of named structured resources. Environment values remain objects.
Resource dependency lists retain their explicit Last semantics. The external
injector and schema live under projects/al/plugins/injector; a public config
supplies registration and transitive runtime files. al_injector_call translates
flat named resource-type maps into the existing args.res graph for inline
declarations and reusable or complex configuration macros.
Wrappers encode plugin_labels dictionaries, including presence-only selection.

All 22 standalone AL tests and the all-target build pass. The independent
consumer verifies nested arguments, substitutions, plugin/TLS runfiles and real
local TLS requests from outside the runfiles tree in directory and manifest
modes. It places an auth override before its parent, verifies one effective
auth definition in both an aggregate and a multi-config command wrapper, and
retains an expected analysis failure for unrelated sibling auth definitions.
Evidence: out/migrate-al-config/composition-consumer/result.json and ambiguity.log.

Real AL/injector/child fixtures pass dictionary selection, explicit resource
precedence, ordered Last dependencies and missing/cycle startup rejection.
Cleanup fixtures verify file modes, lifetime, success, startup failure and
cancellation and retain receipts in undeclared test outputs. Terraform fixtures
exercise exact-value and presence-only dictionary selectors, verify child
environment and argument boundaries, and retain JSON receipts. The parent
integration with the pinned Vault dev server also passes.

The repository-owned [migration evidence](../../../../../infra/src/openspec/changes/migrate-al-configuration-to-bazel/evidence.md)
records output equivalence and parent integration checks. Root quality and
affected semantic lint are mandatory prepared-candidate delivery gates.
No live authentication or deployment is asserted. The built-in generic injection
engine and resource-provider registration protocol remain explicitly deferred.

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
