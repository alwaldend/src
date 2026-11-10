## Context

The approved comfy-cli 1.25.0 targets self-hosted servers over HTTP. The rental
exposes HTTPS with Basic authentication; changing the rental's exposure is unnecessary.

## Goals / Non-Goals

Use upstream commands for workflow submission, polling and output retrieval.
Keep credentials out of CLI arguments, configuration, subprocess environment and logs.
Infrastructure retains ownership of endpoint discovery and Vault references.

## Decisions

Proceed with a temporary Go reverse proxy using existing standard-library and
Cobra dependencies. A fork or transport monkey-patch would create more upstream
maintenance; plain public HTTP would violate the existing transport boundary.
The bridge binds loopback, checks the connecting TCP socket's Linux UID through
procfs, and rejects browser Origin requests. It accepts only the invoking user;
same-user processes already share this account's Vault authentication boundary.
Missing procfs ownership evidence fails closed. TLS verifies hostname and trust
before authenticated requests. Upstream redirects are refused.
The child CLI is forced to the bridge's local address and cloud routing is rejected.
The bridge and child terminate together; no persistent service is installed.

## Risks / Trade-offs

Live CLI acceptance exposed stale cached output URLs: a completed run retained
its temporary bridge port, and a later download reused that dead URL. Use a
private profile per authenticated invocation and delete it after the child exits.
Subsequent status/download commands query authoritative server history. A fixed
port would break concurrent calls; a persistent proxy would add an unnecessary
service. The raw packaged CLI keeps its regular task-owned profile.

- Linux procfs is required for the authenticated bridge; the raw packaged CLI remains portable.
- The CLI has upstream configuration and telemetry behavior. Its authenticated target uses
  a task-owned configuration/cache directory and disables upstream tracking.
- API command success does not establish model quality or reference fidelity.

## Verification before implementation

Exercise the packaged CLI as a subprocess against a local HTTPS fixture. Check
authentication, hostname/trust rejection, response redirect handling, workflow
submission/history/output transfer, child credential removal and cleanup. Check
that missing login, absent rental URL, cloud/alternate-host flags and browser
Origin requests fail without disclosing credentials. Preserve a sanitized result artifact.
