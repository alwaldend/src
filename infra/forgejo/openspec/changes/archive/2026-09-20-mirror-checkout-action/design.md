## Context

The secure workflow uses upstream checkout commit
`11d5960a326750d5838078e36cf38b85af677262`. The user explicitly requested an
organization mirror and use of that mirror. The live upstream default branch
is `main`; the public destination initially returned 404. The pinned
`svalabs/forgejo` 1.5.0 provider supports native pull mirror creation.

## Goals / Non-Goals

Create only `alwaldend/com_github_actions_checkout`, preserve the current action
pin, and retain the `ci` job and local CI action. Do not mutate unrelated
repositories, Vault policy, users, or host configuration.

## Decisions

Use the existing catalog and Forgejo Terraform resource, with `clone_from`
selecting the authoritative `upstream_url`. Keep the existing `github` clone
mode unchanged for first-party repositories. Set a 12-hour pull interval and
disable Actions on the mirror. No imperative repository API bootstrap or new
scheduled automation is needed.

## Risks / Trade-offs

Mirror creation depends on upstream reachability. Synchronization follows
upstream refs, while the workflow stays on its reviewed SHA. A targeted plan
will restrict the authorized apply to the new repository and unchanged
prerequisites; any unrelated change, deletion or replacement stops execution.
Plans are private task scratch, not durable evidence.

## Verification

Verified on 2026-09-20 after adding the mirror to candidate `5b201308`:

- Catalog contract and both catalog/service Terraform format tests passed.
- Pinned Forgejo runner workflow schema validation passed with the mirror URL.
- The saved Terraform plan created exactly one mirror with no changes or
  deletions. Applying that plan added one repository. A second plan for the
  same target reported no changes. No other live resource was targeted.
- The public API confirmed the original upstream, public visibility, `main`
  default branch, 12-hour mirror interval and disabled Actions.
- Git advertised `v4.4.0` at the existing action SHA. Both `action.yml` (5,144
  bytes) and `dist/index.js` (1,362,984 bytes) matched upstream byte for byte at
  that SHA. The workflow retains `github.sha` for the checked-out source.
- Source delivery is tracked by the repository delivery receipt. This mirror
  verification does not claim a successful full repository CI run.

The initial read-only authentication-source probe failed host verification.
Retrying with the checked-in public SSH server CA in task-local known-hosts
succeeded without host changes. Private plans and bounded verification receipts
remain under ignored `out/repo_ci/mirror/`.
