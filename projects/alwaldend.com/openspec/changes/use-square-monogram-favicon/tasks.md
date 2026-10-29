## 1. Site integration

- [x] 1.1 Map the data-owned SVG to the site's favicon resource and remove
      the old artwork; verify one matching SVG in the site source archive.
- [x] 1.2 Update the site and asset documentation; verify obsolete font
      attribution is absent from the current favicon descriptions.
- [x] 1.3 Build the site and inspect homepage, blog, documentation, and
      project HTML; verify their favicon URLs resolve to the canonical SVG.

## Verification

- `//projects/alwaldend.com:site_test` passed.
- The source archive contains exactly one `assets/icons/favicon.svg`,
  byte-identical to the data-owned SVG.
- Homepage, blog index, docs index, project index, Bazel agent landing,
  and 404 HTML reference the same emitted SVG. Its content matches the
  canonical source and its URL contains that source's SHA-256 fingerprint.
- The detailed output inspection is retained under ignored
  `out/s8-favicon/site-favicon-verification.json`.
