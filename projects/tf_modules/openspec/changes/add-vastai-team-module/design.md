## Decisions

Use one module instance per team, with a map of named project keys and explicit
JSON permission documents. The caller owns Vault provider selection, secret paths,
and state. This separates reusable Vast.ai lifecycle from the consuming identity.
The infra/vastai root iterates a checked-in team catalog; adding another team or
project key changes that catalog rather than duplicating resources.

## Decision Review

Proceed with a team module and caller-owned secret delivery. Copying root resources
would leave every addition with independently maintained lifecycle declarations.
Putting Vault authentication into the reusable module would tie every consumer to
one role. The selected design permits explicit provider identities at the caller.

No live apply is authorized by these declarations. All creation requests transfer
zero credit, and neither module nor bootstrap declares compute resources.

## Imported teams

Keep one singleton team resource per module. Import an existing team's account
ID into that resource, or migrate its current Terraform address into the module
when already owned. Never duplicate the same team in multiple states. Imported
team state supplies identity and has null management credentials. In that case,
the API-key resource uses the caller-selected provider and the imported team ID.
New teams retain their generated management-key override for same-apply bootstrap.
There is no `existing_team_id` variable, count transition, or module moved block.
Existing managed resources retain their exact addresses.

Defer permission checks while a computed team binding is unknown. Once resolved,
imported-team owner binding still requires explicit default full access; a scoped
key must never be widened by selecting owner context. Terraform CLI fixtures
cover both import and managed bootstrap. Imported team destruction remains
refused; key-only cleanup preserves the team. Temporary migration blocks in
fixtures are removed after reconciliation and followed by a no-change plan.

The live HTTP 400 team-binding failure is unresolved. This refactor does not
justify repeating that unchanged request or recreating the funded team.
