## 1. Shared catalog and consumer source

- [x] 1.1 Complete the organization and per-repository JSON inventory,
      preserving observed settings and original external upstreams; verify that
      every selected GitHub repository and the F-Droid fork are represented once.
- [x] 1.2 Complete the provider-free projections and consumer dependencies;
      verify catalog tests cover first-party names, external upstream names,
      cross-forge copies, and per-forge settings without duplicate inventories.
- [x] 1.3 Complete GitHub adoption declarations, landing `master` creation
      and defaults, and administrator-only default restrictions; verify pinned
      provider support and preserve existing Pages configuration and `src` rules.
- [x] 1.4 Complete GitLab group, membership, one-time import, and fork
      declarations; verify the source has safe imports, Developer restrictions,
      no sync resources, and no unprotect-and-recreate adoption path.
- [x] 1.5 Complete Forgejo catalog consumption and explicit state moves;
      verify named roles resolve through Vault while service grants, identity
      mapping, and membership-validation guarantees remain intact.
- [x] 1.6 Reconcile the catalog with the subsequently authorized
      [rule-landing retirement](../../../../../src/openspec/changes/archive/2026-09-13-retire-bazel-rule-landings/tasks.md);
      verify the remaining cohort of 28 GitHub repositories and
      29 GitLab projects including F-Droid after the selected eleven retire.

## 2. Offline validation and plan review

- [x] 2.1 Run configured formatting, affected package builds and tests, and
      repository-required checks; verify the formatted source and representative
      catalog projections against the tested candidate.
- [x] 2.2 Review the complete GitHub plan and remote/state inventory; verify
      imports preserve IDs, only intended default/access changes occur, and no
      deletion or replacement is proposed without separate user approval.
- [x] 2.3 Review a fresh complete GitLab phase-one plan after reconciling the
      catalog with rule-landing retirement; the original 39-repository plan is
      stale. Verify existing group and membership adoption, catalog-selected
      copies and fork, group access defaults, and no unapproved deletion or
      replacement.
- [x] 2.4 Review the complete Forgejo plan with the verified OAuth source ID;
      verify state moves retain IDs and the developer assignment preserves
      existing service access without elevated developer grants.

## 3. Authorized adoption and observable access

- [x] 3.1 Apply the reviewed GitHub plan within the authorized scope; verify
      repository IDs, landing `master` defaults, Pages still using `pages`, and
      administrator/developer roles and active default-branch rules before
      GitLab imports. The coordinator's apply and read-only verification
      evidence are recorded in [design.md](design.md#github-adoption-evidence).
- [x] 3.2 Apply the reviewed GitLab phase-one plan and wait for Git imports
      and the fork to finish; verify contents, selected default branches,
      upstream fork identity, membership, and effective default protection.
- [x] 3.3 Import verified GitLab default-branch protections, review the full
      phase-two plan, and apply supported in-place changes; verify existing rules
      remain protected and record any provider limitation without claiming success.
- [x] 3.4 Apply the reviewed Forgejo plan within the authorized scope; verify
      preserved repository/account IDs, named developer access, retained service
      grants, and effective default-branch restrictions.
      The coordinator's verified outcome and explicit OpenHands-account
      approval are recorded in [design.md](design.md#forgejo-adoption-evidence).

## 4. Evidence and delivery

- [x] 4.1 Record exact candidate identity, plan summaries, verification
      outcomes, and any remaining limitations in durable evidence; confirm that
      the accepted deferral of synchronization is documented as scope.
- [x] 4.2 Validate the OpenSpec change and affected owner contracts, reconcile
      them with verified behavior, and archive only after all acceptance is met;
      verify the resulting specs and archive.
- [x] 4.3 Complete required delivery checks, commit, push, and offer the pull
      request; verify the publication receipt identifies the tested candidate
      and resolve applicable review findings.
