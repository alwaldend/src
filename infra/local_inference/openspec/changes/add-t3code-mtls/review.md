## Review target

Reviewed the complete PR #233 diff from master
`5cfe3ae3d8299448ace8e7190d7b4b86bc5f847a` to
`672253432a9e01b28d3bb98c40f18628db25d098`. This is the forge's actual PR
base, rather than the older local master tracking reference. The review used
`review-agent` and the owning infrastructure conventions without delegation.

## Findings

[P1] Permit T3's local MCP — ansible/files/codex_requirements.toml:22

The empty managed MCP table rejects the `t3-code` server that the pinned T3
adapter injects into Codex. The UI can start while its agent-facing T3 tools
remain unavailable. Permit its exact loopback identity, deriving the port
from the same value used to start T3. Keep its session token supplied by T3.

[P2] Check inference health after pending activation — ansible/tasks/configure.yaml:34

Model, runtime, and unit changes notify a deferred restart. The existing
health check runs before that handler, so an old healthy process can pass
validation before the updated process fails at the end of the play. Apply
pending handlers before the readiness check.

Both findings are fixed in the follow-up. Pinned Codex discovered 66 real
T3 tools and completed a read-only call; changed identities were denied.

## Simplification decision

Verdict: proceed with one ordinary coding account for T3 and direct Codex,
plus the separate inference account. No deployment exists to migrate. Separate
coding users and workspaces add ownership/configuration without a requested
boundary. Remove the duplicated CODEX_HOME environment because T3's provider
home setting and direct CLI HOME already select the same directory. Remove
managed-only hooks, login-shell denial, forced subagent disablement, and
automatic MCP-install restrictions that are unrelated to hosted inference.

Keep the shared standalone installer: it reuses the official verified bundle
and avoids another Node runtime pin. Keep canonical DNS, scoped Vault PKI,
encrypted storage, and normal T3 functions. Collapse the redundant shared-role
daemon-reload task into the service-start operation, retaining its restart
handler. Do not add a new security framework or dependency.

## Assessment and limits

No other concrete blocking regression was identified. The review does not
establish GPU model loading, actual inference/approval quality, live DNS/ACME,
or reboot acceptance. The earlier HTTP/pairing fixture never exercised MCP;
the follow-up fixture covers real T3 discovery, a read-only call through
pinned Codex, and rejection of unapproved identities. Git and delivery receipts own
final implementation and publication state.

## Follow-up review of 7e53c11

[P1] Preserve namespaced functions through inference

A captured Codex request contained all 66 T3 functions inside a `namespace`
tool, while pinned llama.cpp b11960 discarded every non-`function` tool.
Direct MCP discovery/calls bypassed that inference path and therefore did not
establish model access. The first-party loopback bridge now flattens function
schemas and replayed calls, then restores namespaces in JSON/SSE replies.
The command-level acceptance fixture also drives actual pinned Codex through
a local MCP call and verifies its result on the next inference request.

[P2] Repair incomplete runtime extraction

Binary-existence guards skipped incomplete runtime directories after an
interrupted extraction. Remove those guards from T3, Codex, llama.cpp, and
the bundled CUDA libraries. A packaged Ansible fixture repaired partial
files and confirmed unchanged reapplication reports no change.

Both follow-up findings are addressed. Actual model quality and GPU loading
remain outside these deterministic protocol checks. The bridge rejects
unsupported tool schemas explicitly rather than silently losing tools.
