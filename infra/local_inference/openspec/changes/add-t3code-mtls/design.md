## Context

The existing host has loopback llama.cpp, a managed local-only Codex provider,
encrypted storage, and initially proposed per-account egress filtering. The shared
T3 Code role uses npm; the current official release also supplies a standalone
Linux executable. Traefik already owns reusable ACME and client-CA packaging.

## Goals / Non-Goals

**Goals:** Browser access using the standard Traefik/Vault/client-CA flow.
**Non-Goals:** Deploy, enable hosted models or tools, add ComfyUI, or change
the existing workstation configuration.

## Decisions

Verdict: proceed with the shared roles and direct VM mTLS enforcement.
Use an opt-in standalone installer in the T3 role, preserving its existing
npm mode. Pin the official release archive and checksum so service startup
uses an installed runtime without a package-manager download.

Run T3 Code as its ordinary `t3code` account, binding 127.0.0.1:3773. Reuse
the same account for direct Codex sessions and the managed executable and
system policy. Seed the local model and
Auto mode, with Full access selectable. The pinned Codex reviewer setup
retains the selected provider and falls back to the parent model when no
separate reviewer is available. Do not add custom T3 provider restrictions.
The user explicitly chose Codex's own permissions/tool controls rather than
account-level egress filtering, so remove that extra firewall and its units.

Expose `t3code.local-inference.alwaldend.com` directly through the VM's
Traefik. The TLS default and route require a trusted client certificate.
HTTP serves ACME validation and HTTPS redirection only. Public application
ports remain closed. A Yandex DNS zone owns the VM's dynamic address;
the component DNS document delegates that zone and declares the public CNAME.

Use T3 0.1.6's native pairing/session authentication. The older shared role's
`T3CODE_AUTH_TOKEN` variable is unsupported by this release, so retain it
only in legacy npm mode. Do not invent a substitute persistent credential.
Use a server-only Vault PKI role scoped to the exact T3 hostname. Ansible's
existing Traefik role obtains and protects EAB material only during separately
authorized deployment.

The user clarified that OpenAI-hosted capabilities should be disabled while
T3's MCP remains usable. Permit the `t3-code` server only at the loopback URL
derived from `t3code_port`; T3 supplies its per-session authentication. Keep
the provider and hosted-feature requirements. Do not constrain local hooks,
login shells, or local subagents for this unrelated purpose.

The requested review found an inference health check before pending restart
handlers. Flush those handlers before checking readiness so the new model and
unit are validated. Direct Codex uses T3's existing home/work directories,
eliminating a separate coding account and duplicate `CODEX_HOME` environment.

## Risks / Trade-offs

- Default UI settings are operator-editable; managed Codex requirements keep
  its inference provider local and disable hosted tools. Full access removes
  the command sandbox for that account when explicitly selected.
- The T3 server can run local tools with the T3 service account's permissions.
  mTLS and native application sessions protect access to that account.
- The remote browser necessarily sends workspace interactions to this VM.
  Inference stays local to the VM, not local to the browser's computer.
- DNS delegation, Vault reachability, live enrollment, and actual GPU-backed
  sessions require separate deployment acceptance.

## Failure Cases and Validation

The follow-up review captured Codex 0.166.0 advertising T3 and subagent tools
as Responses namespaces. Both b11960 and current b12032 skip these schemas.
Verdict: add a first-party loopback Responses bridge rather than fork either
runtime or silently remove tools. It flattens namespaced functions to stable
collision-checked names, rewrites replayed calls, and restores namespace/name
on JSON and streaming replies. Reject unsupported tool types explicitly.
The existing provider URL remains port 8080; llama.cpp moves to port 8081.
The bridge owns no credentials, inference, or tool execution. Remove it when
an upstream release passes the same acceptance cases natively.

Before implementation, exercise duplicate function names across namespaces,
plain-name collisions, forced tool choice, replayed calls, malformed/unsupported
schemas, incomplete upstream streams, JSON errors, concurrent requests, and
client cancellation. A command-level fixture must prove a function-only backend
receives all namespace tools and the client receives the original identity.
An opt-in pinned Codex check must complete a real local MCP call and submit its
result on the next inference request. Backend responses are deterministic;
model quality and GPU loading remain separate live acceptance gaps.

Before implementation, cover missing/invalid client certificates, wrong
hostname, backend port exposure, native authentication, runtime downloads,
wrong Codex binary/home/model, newly enabled cloud tools, and incompatible
Auto/Full policy. Render templates with synthetic credentials, verify the
pinned standalone CLI and loopback startup, and test the rendered mTLS proxy
with a disposable local CA. Keep artifacts under ignored task scratch.

MCP failure cases: T3's matching identity is rejected; a different URL or name
is accepted; session authentication is lost; local-model selection or hosted
feature requirements can be overridden; Auto or Full stops loading. Inspect
the pinned T3 endpoint construction and exercise actual protocol discovery
with its local server and pinned Codex. No GPU or production deployment is
needed for that check. Keep live model quality as a separate acceptance gap.
