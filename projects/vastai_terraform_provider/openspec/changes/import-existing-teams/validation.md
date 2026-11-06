## Acceptance evidence

The import scenario was written before implementation and first failed with
the provider's unsupported-import diagnostic. The completed full Terraform CLI
acceptance suite passes, including existing team creation and cleanup behavior.

Seven `team-import-*.json` files in the acceptance target's undeclared outputs
record valid import, absent identity, noncanonical ID, malformed metadata,
unavailable collection, external deletion, and refresh failure. Every recorded
operation is GET. Valid import produces a no-change plan with null credential
fields; attempted destruction fails and retains state. Confirmed external
absence removes imported state, while refresh errors retain it.

The live infrastructure reconciliation imported three existing owned team
identities without generating credentials. Their source records and observed
IDs belong to infra/vastai's reconciliation record, not this provider workspace.
