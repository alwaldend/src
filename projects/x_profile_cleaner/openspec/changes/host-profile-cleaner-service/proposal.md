## Why

The cleaner currently serves one operator on loopback with injected user credentials
and process-local state. A public service needs isolated accounts, temporary X
authorization, reusable discovery data, and enforceable spending and resource limits
before it can safely accept other users.

## What Changes

- Add a hosted Go service with server-rendered HTML and htmx interactions, behind
  Traefik on one VM, with SQLite persistence and systemd supervision.
- Allow open registration through X OAuth 2.0 authorization code flow with PKCE.
  Separate application login from read/write cleanup authorization. Keep X access
  and refresh tokens only in memory and discard authorization after 24 hours of
  inactivity, logout, disconnect, or restart.
- Inject application OAuth credentials through deployment environment variables.
  Never use the operator CLI's injected personal X credential for hosted users.
- Enforce a $5 USD monthly pool shared by all non-owner accounts, a $0.25 monthly
  account cap, and a $0.10 cleanup cap. Users may choose lower budgets. The configured
  owner is exempt from service monetary caps but remains subject to chosen budgets
  and resource limits.
- Persist account-isolated snapshots, policies, immutable plans, execution progress,
  and conservative spend reservations. Reevaluate cached snapshots without API calls.
- Add bounded X archive import as another snapshot source, avoiding timeline reads.
  Display source, age, coverage, and cost before confirmation.
- Execute untrusted policy evaluation and archive parsing in separately limited
  systemd workers. Bound inputs, discovery, output, storage, and job admission.
- Preserve the existing policy-selection and exact-plan execution contracts. Resume
  interrupted jobs explicitly after renewed authorization; never resume mutations
  automatically after restart.

## Capabilities

### New Capabilities

- `hosted-cleanup-service`: account-isolated SSR interface, durable jobs, exact-plan
  execution, single-VM service operations, and recovery.
- `x-account-authorization`: OAuth login, separate cleanup permission, memory-only
  credentials, refresh, and inactivity expiry.
- `cleanup-spending`: shared/account/run budgets, owner exemption, request admission,
  durable reservations, and user-visible cost accounting.
- `cleanup-snapshots`: reusable cached reads, archive import, coverage, retention,
  and invalidation after mutations.
- `policy-execution-limits`: untrusted-input limits, systemd worker isolation,
  bounded scheduling, and failure without a partial executable plan.

### Modified Capabilities

None. There are no published main specifications in this workspace. The completed
`add-x-profile-cleaner` change remains the maintained baseline for CLI behavior and
core cleanup safety; the new capabilities extend it for the hosted service.

## Impact

The future implementation affects the project-owned web, application, API, policy,
planner, and command boundaries and adds SQLite-backed service state. Reuse the
repository's existing OAuth 2.0 and SQLite dependencies; htmx is the specifically
requested frontend dependency and must be pinned and served locally. Any additional
external dependency requires separate approval.

Deployment implementation will need an owning infrastructure change, packaged
Ansible/systemd/Traefik configuration, component-owned DNS and VM provisioning as
needed, and Vault secret references. This planning PR changes only this OpenSpec
change's artifacts. It does not alter the CLI, provision infrastructure, configure
the live X app, or authorize live reads or deletions. VM location, public hostname,
and any live rollout remain separate deployment inputs.

Articles, Highlights, bring-your-own app keys, historical search, and automatic
cross-source inventory merging remain outside scope. Archive import is explicitly
in scope; it does not promise complete historical or graph coverage.
