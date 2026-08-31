---
name: codex-migration
description: >-
  Safely migrate Codex models, providers, authentication, or shared
  configuration on a host with concurrent users or agents.
---

# Migrate Codex configuration

Treat the existing shared configuration and its owning deployment source as
one change. Preserve unrelated settings and identify active consumers before
editing either copy.

Stage provider and model changes as an isolated `CODEX_HOME` or opt-in profile.
Do not redirect the shared default during validation. Keep credentials in the
host's secret manager and retrieve them at runtime through provider command
authentication; never write or print plaintext credentials.

Validate the exact candidate with all of these checks:

1. Strictly parse the candidate configuration.
2. Start a fresh ephemeral session and complete a minimal response.
3. Start another fresh session, execute one harmless tool call, and complete
   the follow-up response.

Never resume or fork a session across providers. Reasoning items can contain
provider-bound encrypted state, so a valid key and successful one-turn request
do not prove that migrated conversation history is portable.

In a concurrent environment, prefer a tested opt-in profile so existing
sessions retain their provider. Change the shared default only when the user
explicitly requests that disruption and the exact provider, model, auth path,
and tool round trip have passed. Restore the previous configuration if any
validation fails.

Before handoff, mirror every intentional host edit into its owning checked-in
role, validate the packaged desired state without deploying it, and report the
new-session command plus any credential prerequisite. A model request can cost
money or disclose in-scope content to an external provider; obtain authority
immediately before testing unless the current request already grants it.
