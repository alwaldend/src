## Context

Sibling infra projects compose `al_config`, `vault_binary_map`, and the shared
`//infra:al` configuration. See [proposal](proposal.md) for scope.

## Goals / Non-Goals

Establish the same component boundary. The project has no Terraform root.
The authenticated CLI uses the existing upstream package and Vault injector.

## Decisions

Use the sibling GitHub project's AL authentication and backend registration
shape, with only Vast.ai's identity and namespace substituted. Expose standard
Vault wrappers so future credential setup uses the component identity. Keep the
existing upstream CLI in `tools/vastai` separate from infrastructure wiring.

Select API-key injection with `vastai=1`. Read `api_key` from the component's
`vastai` KV entry in the `secrets` mount and export `VAST_API_KEY`, which the
pinned CLI reads directly. The component CLI wrapper packages both the CLI
and injector and selects that label. Standard Vault commands select only their
Vault environment label, so they do not fetch or require the Vast.ai API key.

## Risks / Trade-offs

The role must be deployed before authenticated commands work. Document that
prerequisite; validate packaging offline without creating credentials.
Store the key through the component Vault target before running the CLI.
Do not pass credentials as arguments or persist them in a CLI key file.
