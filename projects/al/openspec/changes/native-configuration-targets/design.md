## Context

The existing al_config runs AL over Lua/JSON/YAML sources and merged dependencies. See proposal.md for the requested breaking migration.

## Goals / Non-Goals

Expose explicit target dependencies and preserve labels, resource dependency semantics and deliberate overrides. No compatibility Lua API or new dependencies.

## Decisions

Proceed with typed configuration fragments and aggregate al_config targets. Fragment rules emit protobuf JSON, carry runfiles, and use label-based executable and certificate inputs. Plugin payloads accept ordinary Starlark dictionaries through macros and encode the existing Json wire representation. The existing config dump operation remains the validator and merger. Explicit file substitutions run before protobuf decoding using JSON-escaped replacement values.

## Risks / Trade-offs

Fragment targets add graph nodes but make ownership and input edges inspectable. Starlark JSON conversion bounds nesting at 64 levels and fails rather than silently truncating. Verify nested payloads, explicit override precedence, file substitution and plugin runfiles through real generated configurations and the existing consumer acceptance command.

Composition resolves explicit per-target extensions independently of aggregate reference order.

Shared registrations and resource graphs belong in reusable configuration
targets. The injector exports a public `:config` target. Aggregate configurations
and command wrappers inherit its executable and transitive runfiles; plugin
binaries are not repeated in consumer `data`. The command cleanup fixture uses
generated configurations, with a small inherited override for template failure.

Vault TLS runfile identifiers are resolved at runtime through the existing optional runfiles discovery API. This keeps generated configuration portable while supplying concrete paths to both Vault client initialization and the environment injector. The independent consumer verifies real local TLS requests from outside the runfiles tree in directory and manifest modes.

## Dictionary and plugin ownership follow-up

## Context

See proposal.md. The approved authoring shape retains generic args.res and structured env values. Named definitions belong in dictionaries; composition uses label references and resource dependencies retain explicit list order.

## Goals / Non-Goals

Reduce configuration boilerplate without changing selected credentials, labels, endpoints or stage boundaries. No shortcut env strings, Lua or legacy parameter aliases.

## Decisions

Proceed with an al_config dictionary compiler over the existing native fragments. Named sections are vault_auth, vault_conn, plugins and plugin_calls. Generic plugin payloads use args; executable arguments use argv. Injector res is a map keyed by resource name. Evaluate resource dependencies in their explicit list order and traverse top-level resources deterministically. Common infrastructure macros declare small reusable targets.

## Risks / Trade-offs

- Name collisions: dictionary keys own names; reject contradictory embedded names and ambiguous sibling overrides.
- Resource order: verify dependencies independent of dictionary insertion order with a real AL/plugin/child end-to-end fixture against a local simulated Vault endpoint.
- Credential drift: compare all generated configurations after normalizing the intentionally changed resource map representation.

## Failure cases before implementation

Wrong payload key or args/argv collision; lost labels or names; dependency order affected by dictionary order; missing resources; duplicate resources across selected calls; unexpected environment values; fetching an unselected call; plugin runfiles missing; empty maps/lists corrupted; live endpoint or inherited authentication accidentally used. Test the actual generated configuration with synthetic loopback-only data and retain a result artifact.

## Settled follow-up scope

The user selected named dictionaries for configuration and injector resources, retained structured env objects and the generic res payload, and renamed plugin parameters to args. Executable argument vectors use argv. The injector implementation moves to projects/al/plugins/injector and remains an external plugin. A built-in generic injection engine with resource-provider registration was explicitly deferred; no registration or dispatch RPC is added here.

The later authoring refinement prefers small native fragment targets composed
through configs. Shared target macros own repeated graphs. Short declarations
remain inline; local defs.bzl macros require significant reuse or a substantial
readability benefit for complex graphs, long templates or construction logic. Keep the dictionary API
for small inputs. Command wrappers accept plugin_labels dictionaries and own
their CLI encoding. The generic al_injector_call macro builds the existing resource graph from
flat named type maps without changing its structured settings or dependencies.

Failure cases before implementation: changed merge order or runtime names,
lost executable/TLS/substitution inputs, wrong selector values or presence-only
matching, environment resources missing their source dependency, and accidental
duplication of source graphs. Existing actual-process cleanup, dictionary and
Terraform fixtures plus complete output comparison verify these boundaries.

Explicit extends relationships select the most specific named declaration.
Aggregate configs union definition identities and reject unrelated competing
definitions rather than using input order. Extension ancestry supplies the
precedence for deliberately layered calls. The independent consumer places an
auth override before a parent aggregate and verifies its effective no_auth
definition, then rejects an unrelated sibling override. Let the formatter order
aggregate target references; no blanket order-preservation annotation is needed.

## Flat injector call authoring

Proceed with the user-selected named `res_<type>` maps on al_injector_call, including res_env, res_file,
res_kv, res_op, res_oidc, res_vault_env, res_vault_ssh and res_process. Each entry contains
structured type settings plus optional deps, vault_auth and vault_conn metadata.
Translate these into the existing args.res graph; do not change the runtime
schema or introduce implicit dependencies. Resource keys remain explicit; the
call identity defaults to its Bazel target name with an optional config_name
override, matching al_plugin_call. Generic calls also default their plugin name to the
Bazel target name; explicit plugin selects a different registration. Export the helper through the public AL API.
Remove al_injector_environment rather than retaining two overlapping authoring
shapes. Shared/local helpers keep their existing ownership and extraction scope.

The alternative of keeping generic plugin/args/res wrappers preserves unnecessary
boilerplate; the narrower source/export helper cannot represent general graphs.
No engine, provider registry, runtime protocol, dependency or deployment changes.

Failure cases before implementation: duplicate resource keys across type maps,
non-dictionary settings, lost dependency order or Vault metadata, typed process
or OIDC name mistaken for resource identity, wrong default/overridden call name,
missing file substitutions or runtime inputs, and accidental changes to existing
resource payloads. Extend the independent consumer and existing actual-process
fixtures before implementing the macro; retain artifacts and compare all outputs.

Resource-type groups use the res_ namespace and are extracted from kwargs by
prefix; other call attributes are forwarded. The macro contains no resource-type
allowlist or fixed parameters to update when the runtime adds a type.
