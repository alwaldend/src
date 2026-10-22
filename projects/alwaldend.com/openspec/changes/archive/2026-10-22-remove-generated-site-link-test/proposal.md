## Why

The generated-site link test has never passed since it was added in
`894494f4`. It fails on unmodified `HEAD` with 330 missing-file findings, and
its assertions conflate two different link namespaces: Markdown links written
against repository source paths, and the published URL space the site actually
serves. Maintaining a permanently failing test blocks repository health checks
without proving a site regression.

## What Changes

- **BREAKING** Remove `//projects/alwaldend.com/test/site:site_test` and its
  source.
- Remove the test from the documented local validation commands.
- Withdraw the baseline requirement that promised the test could report
  unresolved internal destinations.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-alwaldend-com`: the requirement that the generated-site test reports
  unresolved internal destinations becomes a resolution-behavior requirement,
  since the test is removed.

## Impact

`projects/alwaldend.com/test/` is deleted; the project README and the
`project-alwaldend-com` baseline drop their test references. The site build,
preview, and deployment targets are unchanged. Link resolution behavior in
`layouts/_partials/alwaldend/resolve-destination.html` is unchanged.
