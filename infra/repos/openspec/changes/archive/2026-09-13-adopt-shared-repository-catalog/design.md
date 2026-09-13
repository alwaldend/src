## Context

See [proposal.md](proposal.md) for the requested adoption. The three forge
consumers already own separate providers, state backends, and credentials.
GitHub has existing repositories and Pages sites; Forgejo has existing Vault
OIDC accounts, service grants, and Terraform addresses. Those identities and
settings must survive the move to shared configuration.

The initial inventory contained 39 GitHub repositories. All three forge
adoptions are verified as recorded below. The user's subsequent
[rule-landing retirement](../../../../../src/openspec/changes/archive/2026-09-13-retire-bazel-rule-landings/design.md)
removed eleven selected landing entries from the catalog before GitLab
adoption. The remaining cohort is 28 GitHub repositories, their one-time
GitLab imports, and the F-Droid metadata fork. A reviewed replacement plan for this reduced
catalog superseded the earlier GitLab phase-one plan.

## Goals / Non-Goals

**Goals:** Keep one owner for repository identity, naming, and named roles;
give each provider a consistent projection; adopt existing resources without
recreation; enforce the requested developer restrictions while retaining
existing service access.

**Non-Goals:** Replace credential flows, reconfigure Pages, change providers
outside the approved GitHub update, or implement ongoing synchronization.
The user will configure synchronization
later. Its absence is an accepted boundary of this change, not unfinished
implementation work.

## Decisions

### Shared JSON files with a provider-free Terraform projection

The root `config.json` owns the schema version and defaults. Organization
records live in `orgs/<organization>/org.json`; each repository has its own
`orgs/<organization>/repos/<repository>.json`. These files own named roles,
repository identity, original external upstreams, and forge-specific
settings. The module under `infra/repos/tf` combines those inputs, merges
defaults, and derives the maps read by each consumer. It creates no resources
and requires no credentials.

Separate per-forge lists would repeat the same facts. A generated catalog
would introduce a second update workflow for a small declarative inventory.
A provider-free module keeps derivation in one place while preserving the
existing provider and state ownership boundaries. The user selected separate
organization and repository files so each record can be maintained directly;
the projection still presents the same consumer contract.

### Naming follows ownership and original upstream identity

First-party means owned by us, regardless of forge. A copy of our `src`
repository remains `src` on every forge. External forks and mirrors belong
to a configured organization and derive their names from the original
upstream's reversed hostname and full owner/repository path, with punctuation
normalized to underscores. Copies of those forks retain that derived name.

Prefixing every cross-forge copy with its immediate source would rename our
repositories and lose the external upstream's identity. Derivation therefore
uses the catalog's original upstream, not an intermediate organization copy.
The [catalog contract](../../../../README.md) owns the rule and examples.

### Import existing GitHub resources and add default-branch protections

The GitHub consumer adopts the existing organization, memberships, and all
repositories. It preserves current repository settings, Pages sources, and
custom domains. Default branches remain unchanged except for the requested
landing-repository migration to `master`. Existing managed addresses remain stable
or receive explicit `moved` declarations that preserve remote IDs.

The requested access restriction adds administrator-only protections to the
39 default branches. The user selected `master` as the landing repositories'
default while retaining `pages` as their Pages source. Create a missing
`master` branch from the existing repository contents before changing the
default. Protect `master`; retain the developer's existing ability to deploy
to `pages`. Complete this GitHub migration before GitLab imports so the
copies inherit the selected source default branches.

Use `github_repository_ruleset` for the 39 new default-branch rulesets with
the user-approved GitHub provider 6.10.2. Each ruleset restricts creation,
updates, deletion, and force pushes and grants an organization-administrator
bypass. It depends on the managed default-branch switch, so landing
repositories select `master` before the default-branch ruleset applies.
Set `update_allows_fetch_and_merge = false` explicitly. The two existing
`src` rulesets are imported unchanged, preserving their rules and bypass
actors.

Provider 6.6 incorrectly serialized the fork-sync setting as true for update
restrictions. Classic branch protection was considered while retaining that
pin. The user selected the provider update after the corrected serialization
was confirmed in [6.10.2 source](https://github.com/integrations/terraform-provider-github/blob/v6.10.2/github/util_rules.go#L294-L301),
allowing the final configuration to use rulesets. The reviewed apply and
read-only verification outcomes are recorded below.

### Adopt GitLab protections after repository creation

Phase one imports the existing group and memberships, configures group
default-branch protection, and creates the missing copies and fork. The
group defaults deny Developer pushes and merges to default branches.
Phase two imports the observed automatically created protections, then
manages their permitted roles in place.

Recorded `gitlab.id` values select project imports and both the protection
resources and their imports. Protection import IDs use the known
`<project-id>:<default-branch>` pair. New records without an ID enter phase
one under group defaults; after verification, recording the observed ID
enables phase-two protection adoption. The project identity postcondition
uses a null-safe fallback for new records and requires recorded IDs to match.
The [owning workflow](../../../../../gitlab/tf/README.md#import-and-apply)
defines the supported sequence.

The coordinator's verdict is **proceed** with this selection after independent
review identified and corrected the nullable-ID postcondition mismatch using
`coalesce`. The full follow-up plan confirms this source improvement
preserves the completed adoption without remote changes.

Creating a standalone protection resource over an existing default rule can
cause the pinned provider to unprotect and recreate that rule. Two full
reviewed plans avoid that path and allow the actual imported branches and
effective access to be checked before protection adoption. A provider or
subscription limitation must be recorded and resolved without deleting the
existing rule or silently claiming management succeeded.

### Catalog roles coexist with Forgejo's Vault identity and service grants

The catalog owns named organization administrator/developer assignments.
Vault remains authoritative for Forgejo login eligibility, entity UUIDs,
OIDC source mapping, and existing service-specific access groups. Catalog
users must resolve into the validated Vault login population. The consumer
adds the configured bot to the Developer team while preserving current
service administration, package-writing, and automation-writer grants.

Moving those credentials or service identities into the catalog would mix
repository intent with authentication state. Replacing current grants with
the named-role lists would also remove authorized service access. Validation
must reject a catalog developer that simultaneously receives administrator
access through the retained Vault groups.

### Primary design verdict

The coordinator accepts the shared data module and additive
administrator-only default-branch protections. Preserve existing repository
settings, service grants, and credential flows. The user-approved GitHub
6.10.2 update supports the selected rulesets; other provider choices remain
unchanged. Do not introduce a synchronization workaround. The linked
rule-landing retirement records the user's subsequent explicit destruction
scope. Any deletion or replacement outside that scope requires separate
user approval for its concrete scope.

## GitHub Adoption Evidence

The coordinator reports the GitHub provider 6.10.2 apply completed with 91
imports, 70 additions, 31 updates, and no deletions or errors. A complete
follow-up Terraform plan reported no changes.

Read-only API verification observed at 2026-09-13T19:44:42.867198Z completed
263 of 263 requests successfully. It confirmed all 39 repository IDs and
catalog default branches, 39 developer write grants, 39 active default-branch
rulesets with only an organization-administrator bypass, and the two active
catalog memberships with their intended roles. All 31 landing defaults were
`master`; their original `pages` commit IDs remained unchanged. All 33 Pages
configurations were preserved on supported fields, and both original `src`
rulesets retained their IDs, conditions, rules, and bypass actors.

All five forks explicitly returned `update_allows_fetch_and_merge=false`.
GitHub omitted that parameter for the 34 non-forks, where the documented
fork-only setting does not apply. Pages source comparison covered legacy
builds; workflow builds do not support that configuration field. These were
read-only configuration and identity checks, without test pushes or merges.

Private API responses and plan artifacts remain in ignored task scratch.
This is the pre-retirement baseline. The linked
[retirement evidence](../../../../../src/openspec/changes/archive/2026-09-13-retire-bazel-rule-landings/design.md#current-evidence-and-next-action)
owns the completed eleven-site deletion, surviving-resource comparison,
and clean GitHub follow-up plan.

## GitLab Adoption Evidence

The reviewed replacement phase-one plan completed with two imports, 30
additions, one change, and no deletions. All 28 source imports and the
F-Droid fork completed. Receipt: `out/repos/private/gitlab_apply_v6.log`.

Read-only API verification observed at 2026-09-13T20:10:57.729581Z completed
149 of 149 requests successfully, with no failed checks or truncation. It
confirmed all 29 projects, compared all 28 imported source commit IDs, and
verified the F-Droid upstream identity, selected defaults, memberships,
effective roles, and complete default-branch protection lists. Feature-branch
and merge-request permissions were inferred from those effective roles,
enabled features, and branch rules; no test writes were performed. Receipt:
`out/repos/private/gitlab_post_apply_summary.json`.

The API omitted the group's `developer_can_initial_push` field. That field
remains recorded as unavailable rather than asserted from the response;
observed project protection and role checks passed.

Phase two imported all 29 observed default-branch protections, with no
remote additions, updates, or deletions. The complete post-adoption plan
reports no changes. Receipts:
`out/repos/private/gitlab_protections_apply.log` and
`out/repos/private/gitlab_post.log`. These are one-time copies; ongoing
synchronization remains an accepted deferral owned by the user.

All 29 observed project IDs are now recorded in their per-forge catalog
configuration. Project imports and the ID postcondition preserve those
identities; protections select only records with a known, non-null ID.
The final follow-up plan for this reusable source selection reports no
changes: `out/repos/private/gitlab_adoption_post.log`.

## Forgejo Adoption Evidence

The coordinator reports that the reviewed Forgejo plan was applied after
the user explicitly approved creation of the `src_infra_openhands` account.
The apply added that account and the catalog bot's Developer membership,
enabled the Developer team's access to all repositories, and preserved the
existing repository ID through its explicit address move. The bot's observed
`is_admin` value is false. Existing administrator-team membership and service
grants remained unchanged, and the complete post-apply plan reported no
changes.

Each forge's deployment outcome is recorded separately above. Private rollout
artifacts are temporary; the validated candidate is recorded below.

## Validation

The affected package and root consumer checks passed. The
linked [retirement evidence](../../../../../src/openspec/changes/archive/2026-09-13-retire-bazel-rule-landings/design.md#current-evidence-and-next-action)
owns their receipts, the completed nested-module validation, and rendered
documentation checks.

The implementation candidate passed the mandatory quality and lint checks
and was published before this archive, as recorded below. The archive merges
the completed capability into the owning main specification; its final
publication is verified through the repository delivery workflow.

## Risks / Trade-offs

- Incomplete remote inventory or defaulted provider fields could change an
  existing resource. Preserve observed settings, inspect full plans, and
  confirm remote IDs and representative outputs after adoption.
- Address changes could appear as replacements. Use imports or explicit
  `moved` declarations and reject any unapproved deletion or replacement.
- Overlapping access rules or inherited grants could elevate the developer.
  Review effective permissions, including Forgejo service grants and GitLab
  matching branch rules, rather than relying on a role label alone.
- Changing a landing repository's default could disrupt Pages if its source
  follows the default implicitly. Keep the explicit `pages` source and custom
  domain, then verify developer deployment access and protected `master`.
- GitLab imports and forks complete asynchronously. Wait for completion and
  verify source/default branches before adopting protections.
- One-time copies become stale as their source changes. The user accepted
  this behavior and will arrange synchronization separately.

## Migration Plan

1. Validate catalog projections and all affected packages offline. Review the
   current remote inventory and state through the owning authenticated flows.
2. Import existing GitHub resources and review a full plan that preserves
   existing settings except for the requested landing defaults and access
   controls. Create missing landing `master` branches, change their defaults,
   and verify that Pages remains sourced from `pages` before GitLab imports.
3. Reconcile the shared catalog with the separately authorized rule-landing
   retirement and discard the previous GitLab phase-one plan. Review the
   remaining source cohort before continuing GitLab adoption; verify remote
   retirement through its owning change.
4. Review and apply GitLab phase one, then verify completed copies, fork
   identity, default branches, membership, and effective protection.
5. Add GitLab protection imports, review the complete phase-two plan, and
   apply only supported imports or in-place updates.
6. Review Forgejo's address moves, catalog-derived roles, and retained service
   grants; apply the reviewed plan and verify its access behavior.
7. Record exact candidate, plans, outcomes, remaining limitations, and the
   next action in the task evidence; validate and deliver the source.

If adoption fails, retain the existing remote objects and their protection.
Correct the source or import mapping and review a new full plan. Removing
resources, discarding state, or replacing accounts is not an adoption rollback
strategy. The linked retirement follows its separate explicit user scope.

## Validated Implementation Candidate

Implementation candidate `69ae7d11b3ae43fbb5621fd0b521ea8199d2102a`
(tree `259877706d04dfca1a21c4efe86fe712cd0bc5c4`) passed all 27
selected delivery checks, including root repository quality, Buildifier,
root consumer semantic lint, representative builds, and explicit build/test
checks in all twelve standalone modules. The standalone wrapper has no lint
profile; its supported builds and tests supplement the root checks.

The candidate was pushed and its [pull request](https://github.com/alwaldend/src/pull/87)
was verified by the owning delivery workflow. Infrastructure work is complete.
The pinned OpenSpec CLI archived this change and merged its six requirements
into the main repository-catalog specification.
Private rollout artifacts were used for the observations above and are
removed before handoff; these sanitized outcomes and catalog identities are
the durable evidence. Final publication and review identity belong to the
repository delivery receipt and Git history.
