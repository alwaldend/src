## Context

See `proposal.md` - Why.

`github_repository.project_landing` owns a `dynamic "pages"` block whose
`for_each` is empty whenever the project's registry key appears in
`bootstrap_projects`. Removing the Pages block tells GitHub to drop the custom
domain, and restoring it re-adds the `cname` and makes GitHub schedule a fresh
certificate.

The existing variable name and its documentation describe first creation: a
project whose `pages` branch does not yet exist. Applying that same variable to
an already-served site would now be the intended recovery action, so the two
intents must be distinguishable in review.

## Goals / Non-Goals

**Goals:**

- Provide a reviewable, checked-in way to force a GitHub Pages certificate
  reissue for one named landing site.
- Keep first-creation bootstrap and certificate recovery as separate,
  self-describing variables.
- Keep the recovery scoped to the affected project and leave other sites
  untouched.

**Non-Goals:**

- Managing certificates directly (GitHub issues them; Terraform cannot).
- Automating the two-phase apply or the content republish.
- Diagnosing why a specific certificate failed; the procedure applies to any
  site that needs reissuing.

## Decisions

**Add `certificate_reprovision_projects` rather than reusing
`bootstrap_projects`.** Reusing the existing variable would silently widen its
documented contract ("new projects whose pages branch does not exist yet") and
make an ordinary bootstrap apply capable of tearing down a live site's custom
domain. A separate variable makes the destructive phase explicit in the plan
and its name. Alternative rejected: a boolean `force_certificate_reissue`;
per-project membership matches the existing `for_each` shape and keeps the
two-phase apply as the idiom.

**Drive the same `dynamic "pages"` `for_each` from the union of both sets.**
The mechanism is identical, so both sets feed one expression: the Pages block
is omitted when the project is in either set. This avoids a second code path
that could diverge.

**Document a two-phase, filtered apply.** Phase one omits the Pages block for
the affected project only; phase two applies with both sets empty to restore
it. Content is republished only if phase one removed it; normally the `pages`
branch is unchanged and no republish is needed.

## Risks / Trade-offs

- [A Phase-one apply widens scope to other projects] -> Filter with `-var` to
  the single project and inspect that the plan touches only its Pages block.
- [The site serves HTTP while the custom domain is detached] -> The window is
  one apply plus certificate issuance; announce the maintenance window and
  verify HTTPS afterward.
- [Phase two is forgotten, leaving the site without a custom domain] -> Treat
  the change as one procedure; do not close it until phase two and HTTPS
  verification are complete.
- [GitHub does not reissue because the domain configuration was cached] ->
  Verify the issued certificate and, if it does not recover, re-run the
  documented procedure rather than editing provider state by hand.

## Migration Plan

1. Plan with only the affected project in `certificate_reprovision_projects`
   and inspect that the diff removes that project's Pages block and nothing
   else. Apply when authorized.
2. Plan and apply again with both variable sets empty, restoring the Pages
   block and custom domain.
3. Verify the pushed revision, Pages build state, DNS, and the HTTPS
   certificate served for the custom domain.

Rollback: if phase two fails, re-apply phase two with the Pages block restored;
no state surgery is needed because each phase is an ordinary Terraform apply.
