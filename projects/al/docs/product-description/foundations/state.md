# State, identity, and completion

## Summary

AL operates on several independent kinds of state. Choosing a source version
does not publish every project. Publishing artifacts does not activate a website.
Planning or inspecting a publication does not create or reserve its remote state.
Preparing or validating a delivery candidate does not itself update the PR.
Direct delivery rebase updates only the local candidate; publication remains a
separate guarded operation after preparation and validation.

## State boundaries

| State                        | What the operator can observe                         | Owning feature                                   |
| ---------------------------- | ----------------------------------------------------- | ------------------------------------------------ |
| Configuration                | Explicit files and their merged representation        | [Configuration](../commands/configuration.md)    |
| Plugin resources             | Resources available while a selected child runs       | [Run](../commands/run.md)                        |
| Source refs                  | Git branches and tags shared by project builds        | [Versioning](../releases/versioning.md)          |
| Bundle                       | Local manifest and copied payloads                    | [Bundle](../releases/bundle.md)                  |
| Publication                  | Complete verified payloads at one destination         | [Publication](../releases/publication.md)        |
| Retained site and activation | Extracted content and the selected `current` link     | [Website operations](../releases/publication.md) |
| Candidate and PR             | Exact commit, checks, remote branch, and review state | [Delivery](../delivery/workflow.md)              |

Projects share repository versioning and the operator workflow, but each release
is selected separately. Immutability applies while a publication exists. Explicit
deletion permits ordinary republishing; there is no restoration-only state.
Shared source tags survive artifact deletion.

## Receipts

An artifact-operation receipt records requested steps as `pending`, then `unknown`
before remote work, then `completed` after confirmation. Its overall status is
`running`, `completed`, or `incomplete`. A killed process can leave `running` and
an unknown step. A lost connection can occur after the remote change, so unknown
never proves absence. Receipt creation must succeed before the associated step.
Use a different receipt path per attempt to preserve history. Artifact upload,
website extraction, and website activation use separate invocations and receipts;
the receipt retains their backend names `upload`, `extract`, and `link`.

A source-ref receipt reports verified Git publication. It is not the checkpointed
artifact-step journal. A delivery receipt binds the repository, refs, candidate,
path scope, and publication expectations. The operator runs and assesses checks
separately; the receipt does not record their results. Its lock prevents
cooperating AL operations from changing that receipt
at the same time. None of these files grants permission, proves current remote
state forever, or substitutes for fresh observations.

## Recovery and concurrency

Whole workflows are not one transaction. An upload can succeed before a separate
website extraction fails; a branch push can succeed before PR metadata verification fails. Preserve
completed work and inspect the exact boundary before retrying. Do not edit trusted
receipts to manufacture matching state.

Release publication and source-ref operations use different locks. Artifact
publication serializes destination changes for a project. Source-ref operations
use the Git common directory, shared by linked worktrees. Forced termination can
leave a stale source-ref lock that requires operator investigation. There is no
single global lock coordinating all commands or unrelated external writers.

## Open questions and verification

- Source: `internal/release/receipt.go`, `internal/release/ssh_deployer.go`,
  `internal/versioning/releasepublish.go`, and `internal/delivery/receipt.go`.
- Inspect the code and tests when choosing recovery for an interrupted operation;
  receipt formats are deliberately not interchangeable.
- Checklist: [state](../verification/foundations.md#state).

Verified against alwaldend/src commit `7cd4714139984bf4f9b2ae1de9189470b336e12e`.
