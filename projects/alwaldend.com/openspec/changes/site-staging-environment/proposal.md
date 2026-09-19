## Why

The project landings moved from dedicated subdomains into the main site at
`/projects/<name>/`, but that replacement has never been served. The live apex
still returns 404 for `/projects/`, so the per-project hostnames, DNS records,
and Pages repositories cannot be retired: tearing them down first would break
every published project link with no replacement.

Verifying the replacement by deploying straight to production would publish an
unreviewed layout to the site visitors already use. A staging publication gives
the replacement its own hostname, its own Pages repository, and the same
rendered output, so the migration can be checked end to end before the apex
changes and before any teardown.

## What Changes

- Add a `www-staging.alwaldend.com` DNS declaration beside the existing apex
  records, and a catalog repository that serves it.
- Add a staging build of the apex site that renders the same source with the
  staging base URL, and a deployment command that publishes it.
- Publish the site to staging, verify the landing, documentation, and taxonomy
  routes there, then deploy the same accepted output to the apex.
- Retire the per-project hostnames, DNS records, and landing repositories once
  the apex replacement serves them.

## Capabilities

### New Capabilities

- `site-staging`: a staging hostname that serves the apex site's rendered
  output for pre-production verification.

### Modified Capabilities

<!-- Retirement of per-project publication is owned by
     consolidate-project-landings-into-apex; this change only supplies the
     verification that unblocks it. -->

## Impact

- `infra/dns/dnsconfig.json` gains one staging record; `infra/dns` keeps
  owning the apex and staging declarations.
- `infra/repos/orgs/alwaldend/repos/www-staging.json` adds a catalog
  repository, and the catalog gains a creation-only `auto_init` default.
- `infra/github/tf` gains the corresponding repository attribute; the staging
  repository's default branch is protected while site content publishes to a
  separate unprotected `pages` branch.
- `projects/alwaldend.com` gains a staging site target and deployment command.
- The staging hostname, its record, and its repository are retired with the
  rest of the per-project publication infrastructure.
