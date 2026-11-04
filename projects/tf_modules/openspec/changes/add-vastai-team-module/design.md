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
