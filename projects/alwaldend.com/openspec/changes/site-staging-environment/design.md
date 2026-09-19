## Context

See `proposal.md` - Why.

The main site already carries every project landing and the deploy command that
publishes it, and the reusable infrastructure to add a repository and a DNS
record exists in `infra/repos` and `infra/dns`. The staging environment is
therefore a configuration addition plus a second build of the same source, not
new machinery.

Related change:
`projects/alwaldend.com/openspec/changes/consolidate-project-landings-into-apex/`
owns the landing migration and the retirement of per-project publication.
Its task 7.9 is blocked pending this verification.

## Goals / Non-Goals

**Goals:**

- Serve the exact replacement output on a hostname that is not the public apex.
- Verify landing, documentation, and taxonomy routes before any teardown.
- Reuse the existing publication command and the existing repository and DNS
  catalogs.

**Non-Goals:**

- A second Hugo build pipeline, a staging-specific layout, or a new deployment
  mechanism.
- A promotion or rollback workflow. Production publication stays the existing
  apex deploy command.
- Managing a staging certificate or a redirect for the former hostnames.

## Decisions

**Build staging from the same packaged source, differing only in base URL.**
The candidate that is verified must be the candidate that is published, so
staging is a second `al_hugo_worker` over the same `site_site` input with the
staging `--baseURL`. Rendering a separate content tree would verify something
other than what production serves.

**Reuse the deployment command instead of adding a staging-specific one.** The
publisher took `REPOSITORY` and `BRANCH` arguments already, so staging is an
argument change. A second script would be a second owner for the same
publication behavior.

**Give a new repository a default branch with `auto_init`.** GitHub rejects a
Pages source branch that does not exist, so a repository created empty cannot
be configured to serve Pages. The retired landing repositories solved this with
a dedicated branch resource sourced from the publication branch plus a
two-phase bootstrap apply. `auto_init` is creation-only, gives the repository a
default branch immediately, and lets the managed ruleset protect that branch
while content goes to a separate unprotected `pages` branch. Imports resolve
the attribute to `false`, so the existing repositories show no drift.

**Stage before deploying production, and tear down last.** The live check
showed the apex replacement returning 404 while twenty project hostnames
returned 200, so the order is: publish staging, verify it, publish the apex,
then retire the per-project records and repositories.

## Risks / Trade-offs

- [A staging hostname is publicly resolvable] -> It serves the same public site
  content and declares no secret; it is retired with the rest of the
  per-project publication infrastructure.
- [The staging Pages certificate must be issued before HTTPS is verified] ->
  Verify the served content first and treat certificate issuance as
  asynchronous evidence.
- [Teardown could remove a hostname still referenced by a published link] ->
  Retire only after the apex routes serve the replacements.

## Migration

Add the declarations and targets, publish and verify staging, publish the apex,
then retire the per-project DNS records and landing repositories under the
authorization already granted for that teardown. Vault is not touched.
