## Context

`//projects/alwaldend.com/test/site:site_test` was introduced together with
`resolve-destination.html` in `894494f4`. It tokenizes every generated HTML
page and checks that each internal `href`/`src` resolves to a generated file.
It fails on unmodified `HEAD`: 330 distinct missing targets, of which 105 are
linked only from normal pages.

The failures split into two classes. Repository source links such as
`/docs/infra/ansible/BUILD.bazel` or `/docs/tools/openspec/README.md` point at
files the site deliberately does not package, so they can never resolve in the
generated output. Print-edition pages (`/_print/...`) additionally rewrite
document-relative links against a deeper output path. The test treats both as
site defects, so it reports pre-existing, intentional behavior as failure.

## Goals / Non-Goals

**Goals:**

- Remove a permanently failing test that no longer reflects intended output.
- Keep the documented validation commands accurate.

**Non-Goals:**

- Changing how Markdown links resolve, or which files the site packages.
- Adding a replacement link checker in this change.

## Decisions

The test is removed rather than repaired. Its failing assertions encode a
link-namespace assumption the site does not implement: it requires every
repository-relative Markdown destination to exist in the published HTML tree,
but the site intentionally serves only packaged documentation. Repairing it
would require deciding which source links are allowed to dangle, which is a
site policy question, not a test fix.

The `project-alwaldend-com` baseline requirement that named the generated-site
test is withdrawn with the test. Resolution behavior itself remains specified;
only the claim that a test reports unresolved destinations is removed.

## Risks / Trade-offs

- Removing the test drops the only automated check for missing internal files
  and fragments. That check was failing before removal, so no passing signal is
  lost; a future replacement is possible but out of scope here.
- `//projects/alwaldend.com:site_test` (the `build_test` of `:site`) remains
  and still verifies the site builds.
