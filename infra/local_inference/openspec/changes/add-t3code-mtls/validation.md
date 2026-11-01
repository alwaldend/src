## Source and disposable runtime acceptance

The original HTTP/pairing observations below were recorded for candidate
`672253432a9e01b28d3bb98c40f18628db25d098` on 2026-11-01.
The MCP follow-up supersedes its deny-all policy. No live infrastructure was
changed. Task scratch is `out/local_inference/t3/`; delivery receipts bind
mandatory gates and publication to the final committed candidate.

- The official T3 0.1.6 Linux x64 standalone archive matched the checked-in
  SHA-256. Its CLI accepted the service's `serve` flags and native `auth`
  commands. The full extracted bundle served the browser assets.
- Packaged Ansible syntax passed against a localhost-only inventory. Six
  rendered outputs (Codex config, llama unit, T3 unit, Traefik configs, and
  initial settings) were inspected; rendered TOML parsed successfully.
- Backend-disabled Terraform init/validate passed using Yandex 0.203.0 and
  Cloudflare 5.22.0 from the repository mirror. A scratch mirror of the root
  and its declared local modules kept backend and provider state out of source.
- DNS, Terraform formatting for both owners, and the three owner OpenSpec
  checks passed: nine tests. AL and Ansible packaging builds passed.
- Codex 0.166.0 parsed managed requirements inside a network-isolated
  namespace. Auto-review and Full permission settings were accepted. Attempts
  to enable apps/plugins left them disabled, and a supplied remote MCP server
  remained disabled by managed requirements.
- A disposable local T3 instance accepted the seeded local model and Auto
  settings, reported anonymous sessions as unauthenticated, accepted a pairing
  grant, and rejected reuse of that one-time grant.
- The packaged native Traefik 3.7.17 binary routed a trusted client to that actual
  T3 instance. Missing client certificates, a different client CA, and wrong
  SNI were rejected. HTTP redirected to HTTPS. The fixture used the rendered
  route and TLS options, substituting local ports, a disposable CA/server
  certificate, and an installed test certificate instead of live ACME.

Repeatable fixture code and non-secret results are `runtime-check.py` and
`runtime-results.json` under that scratch directory. Temporary pairing state,
startup logs, and test private keys are removed after inspection. The fixture
ran inside a network namespace, so it could not use hosted inference.

The pinned T3 runtime maps Auto to Codex `on-request`, `auto_review`, and
workspace-write permissions, and Full to `never` plus full access. Inspection
of pinned Codex reviewer setup shows it retains the selected provider and
falls back to the parent model when a separate review model is unavailable.
This establishes configuration compatibility; it does not demonstrate the
operator's future model making useful approval decisions.

## Remaining live acceptance

No Terraform plan/apply, inventory-host connection, Vault enrollment, or
production credential generation ran. GPU/model loading, actual inference
and tool calls, local model approval quality, deployed DNS/Vault enrollment,
service ownership, browser pairing, and reboot persistence require separately
authorized deployment. Full removes the command sandbox when selected;
accounts retain ordinary networking in both modes.

## Session ergonomics

- `T3-NATIVE-AUTH-COMPATIBILITY` (fixture-tested): 0.1.6 ignores the role's
  legacy `T3CODE_AUTH_TOKEN`. Standalone mode now uses native pairing rather
  than an ineffective secret injection. Check the pinned CLI contract first
  when upgrading this runtime.
- `LOCAL-FIXTURE-PACKAGING` (fixture-tested): an external fixture playbook
  needs the canonical DNS input beside it, matching the production package.
  The runtime check uses the native executable packaged in the role rather
  than a Bazel launcher without its runfiles environment.
- `LOCAL-FIXTURE-TLS-CA` (fixture-tested): the host's strict verifier requires
  CA key-usage extensions. The disposable CA declares them explicitly.

The user also requested the latest stable proxy. The shared dependency pin
was updated from 3.7.16 to 3.7.17 and the fixture repeated with a version
assertion. Upgrade evidence is maintained under
`infra/src/openspec/changes/update-traefik-3-7-17/`.

## MCP follow-up and review acceptance

The requested complete PR review identified the MCP denial and readiness
ordering defects in `review.md`. Both were corrected. The component now
shares T3's coding account, home, and workspace with direct Codex sessions;
its redundant CODEX_HOME environment and unrelated local restrictions were
removed. A legacy image-generation alias was removed after pinned Codex
reported that the canonical `image_generation` requirement already owns it.

Eleven assertions passed in a network-isolated fixture using the real T3
0.1.6 server and Codex 0.166.0 app server. Codex discovered 66 T3 MCP tools
and successfully called `t3_thread_list`. The matching name/URL was enabled,
while another name and an altered URL were disabled. Native authentication
remained enabled; temporary MCP credentials never entered source or evidence.
Auto and Full retained the local provider, hosted features remained disabled
despite enable attempts, and local subagents could be enabled. No model
inference ran. Fixture state, private logs, and credentials were removed.

The repeatable fixture is `out/local_inference/mcp_review/mcp_check.py`, with
non-secret results in `mcp-results.json`. Packaged Ansible rendering and a
localhost-only syntax check passed. Source inspection verifies that inference
handlers run before readiness and systemd reload occurs before T3's initial
start. Actual VM user creation, GPU startup, and reapplication remain live
acceptance gaps. Final mandatory gates run through guarded delivery.

## Responses namespace correction

The subsequent review of `7e53c11` found that direct MCP dispatch bypasses the
model request. Pinned Codex sends namespaces that llama.cpp b11960 and b12032
discard. The component now packages a first-party loopback bridge and the
checked-in `//infra/local_inference/test/protocol:protocol_test` target.

The fixture exercises the real bridge executable against a deterministic
function-only backend: namespace name collisions, forced tool selection,
replayed calls, JSON/SSE replies, concurrent mappings, error propagation,
malformed streams, cancellation, readiness, and route restrictions. Its
optional pinned Codex 0.166.0 case completed a local MCP tool call, submitted
the tool result on the following inference request, and received a final
response. Hosted search is disabled and only the fixture's harmless marker
tool is explicitly approved. This does not test approval-model quality.

Packaged Ansible repaired a partial synthetic bundle and reported unchanged
reapplication; all eight rendered outputs and a localhost-only packaged
syntax check passed. Inference readiness now traverses port 8080 through
the bridge to llama.cpp on port 8081. These checks did not deploy anything.
Task-local logs are under `out/local_inference/compatibility/`; the delivery
receipt owns final candidate checks and publication.

`MCP-INFERENCE-PATH-COVERAGE` (fixture-tested): discovery and direct tool
invocation are insufficient acceptance evidence for model-driven tools.
The checked-in fixture now includes actual Codex dispatch and result replay.
