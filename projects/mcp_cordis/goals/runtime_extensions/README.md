# MCP Cordis goal

This is a durable project goal. Future work should resume from the current
attempt and preserve accepted evidence and strategy changes here.

## Goal

Deliver a standalone, workspace-local MCP server at `projects/mcp_cordis`
that reuses Cordis for hot runtime packages, persists reusable packages in
the project, stores disposable packages under `out/mcp_cordis`, and ships
useful non-sensitive starter packages derived from recurring past-session
work.

## Status

Complete: the delivered MCP is a thin wrapper around official Cordis Loader,
Include, HMR, and Timer packages. Reusable modules use project-local
`cordis.yaml`; disposable modules use ignored workspace output. The final
hosted review corrections bind repository reads, searches, directory
inspection, and Git commands to verified file or directory handles, bound
permanent process-inspection failures, and remove a process-start timing
assumption from cleanup coverage.

## Current state

- Delivered candidate: `attempt-11/exact-consolidated-rebase`
- Rejected parent commit:
  `7cfef0719075ad372c3bb257ad216b35770356b2`
- Rejected parent tree:
  `34153eca0f582af5c641f81bf8c7209b0045ab9a`
- Current fetched and rebased base:
  `63e7b9f0be1e054373415914ff3d2ea2282aa3da`
- Published candidate: PR 32
- Stage: delivered and review-reconciled
- Last accepted checkpoint: exact aggregate published on PR 32 with the
  delivery receipt verified against the remote branch and PR
- Failing or unverified criteria: none
- Dominant issue: none
- Exact next action: merge PR 32 when desired

## Current plan

1. **Complete:** complete read-only preflight and choose the package/build
   boundary.
2. **Complete:** scaffold the project and pin dependencies reproducibly.
3. **Complete:** implement runtime lifecycle, stable MCP tools, and two-tier
   persistence.
4. **Complete:** make all starter packages pass executable tests.
5. **Complete:** make lifecycle evidence deterministic and publish PR 32.
6. **Complete:** integrate PR 24's `projects/agents` changes without
   regressing newer main-branch guidance.
7. **Complete:** correct all three valid review findings with regression tests.
8. **Complete:** independently review Attempt 4; verdict was refine.
9. **Complete:** correct the independent-review gaps in Attempt 5 and pass the
   focused, integrated, build, validation-aspect, and Buildifier gates.
10. **Complete:** obtain independent review; verdict was refine.
11. **Complete:** implement Attempt 6's lifecycle, compatibility, package,
    and skill-policy corrections, including the second-review strategy reset.
12. **Complete:** reject the custom package-manager architecture after the
    user's standard-solution review and freeze Attempt 7.
13. **Complete:** implement official Cordis Loader, Include, HMR, normal
    modules, and standard `cordis.yaml` storage.
14. **Complete:** rerun every invalidated focused and integrated validation
    gate and obtain fresh independent review. Review accepted `c05bd45a` with
    no actionable findings.
15. **Complete:** rebase, validate the exact aggregate commit, republish PR 32,
    resolve its review threads, and verify the remote candidate. The remote
    head matched the accepted local candidate after publication.
16. **Complete:** correct the hosted review's recursive fallback, byte-offset,
    and UTF-8 preview findings; pass focused regressions and the complete MCP
    test/build/Buildifier packet.
17. **Complete:** correct the follow-up review's Unicode case-fold, explicit
    file/glob, and partial-startup listing findings; rerun the same gates.
18. **Complete:** correct unavailable-scope and bounded-read endpoint findings;
    use fresh diff scrutiny to fix repeated-startup error loss and natural-body
    UTF-8 handling; update `repo-delivery` to require correctness revalidation
    after code changes.
19. **Complete:** replace the rejected transactional source-HMR protocol with
    atomic persistence plus Cordis HMR and reproduce the remaining native HMR
    overlap race independently.
20. **Complete:** serialize and drain HMR reload work in the reproducibly
    pinned dependency, retain the fallback corrections, and rerun every
    invalidated local delivery gate.
21. **Complete:** consolidate the exact owned range, reconcile the advanced
    base's skill-discovery and goal layout, validate the literal rebased
    candidate, publish it, and reconcile hosted review.

## Records

- [Acceptance criteria](acceptance.md)
- [Requirements and constraints](requirements.md)
- [Failure ledger](failure_ledger.md)
- [Evidence manifest](evidence.md)
- [Artifact log](artifacts.md)
- [Current attempt: Attempt 11](attempts/011.md)
- [Attempt history](attempts/)
