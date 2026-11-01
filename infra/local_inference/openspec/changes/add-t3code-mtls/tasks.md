## 1. Integration

- [x] 1.1 Add DNS, HTTPS ingress, native T3 authentication, and scoped PKI;
      verify provider validation, DNS checks, and packaged AL configuration.
- [x] 1.2 Deploy normal T3 Code through the shared role, seed local Codex
      and Auto mode, permit Full access, and remove account-level filtering;
      verify rendered settings, units, policy, and the packaged playbook.

## 2. Acceptance

- [x] 2.1 Exercise the pinned T3 CLI and mTLS configuration with disposable
      local fixtures; preserve repeatable artifacts without live deployment.
- [x] 2.2 Update documentation and define the mandatory quality/lint gates
      for guarded delivery. Git and delivery receipts own PR publication state.

## 3. MCP correction and requested review

- [x] 3.1 Permit the authenticated loopback T3 MCP by name and exact URL;
      reject other identities and retain the local model and hosted-tool policy.
- [x] 3.2 Share the ordinary T3 coding account with direct Codex sessions,
      remove duplicate settings and unrelated local-tool restrictions, and
      check inference health after applying pending restarts.
- [x] 3.3 Exercise the allowed MCP with pinned Codex and T3, check rejected
      identities and Auto/Full compatibility, and select guarded delivery gates.
      Git and delivery receipts own publication state.

## 4. Responses compatibility

- [x] 4.1 Adapt Codex function namespaces for llama.cpp with a local bridge;
      preserve tool identity in streaming replies, JSON replies, and history.
- [x] 4.2 Add a reproducible protocol acceptance target and exercise pinned
      Codex through the bridge, including a tool result and follow-up response.
- [x] 4.3 Package and configure the bridge, document its boundary, and select
      delivery gates without deploying; receipts own final check results.

## 5. Deployment packaging and role ownership

- [x] 5.1 Move the existing Codex, llama.cpp, CUDA, and T3 archive pins into Bazel repositories and package controller-side inputs; remove deployment-time archive downloads.
- [x] 5.2 Extract llama account, runtime, model verification, service, and readiness into a reusable role; preserve the V100 image checks in the deployment.
- [x] 5.3 Extract the existing Traefik upgrade and its change record into an independent PR against master.
- [x] 5.4 Verify packaged inputs, rendered units, archive repair/idempotency, and select protocol/publication gates without deployment.

## 6. Merged bridge integration

- [x] 6.1 Consume the merged projects/llama_responses_bridge target and remove
      duplicate component source; verify package inputs and protocol acceptance.
- [x] 6.2 Rebase onto current master, validate the exact candidate without
      deployment, and inspect that the aggregate retains only infrastructure work.
- [x] 6.3 Use the merged independent AL module's public rule API and current
      Bazel command entry point after the lifecycle-tooling migration.
