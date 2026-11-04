## 1. Header integration

- [x] 1.1 Enable the navigation logo, connect the canonical S8 resource, and
      inspect the built header on the homepage and project indexes.
- [x] 1.2 Capture and inspect light and dark screenshots with the existing
      site screenshot command, then show the requested pages.

## 2. Validation

- [x] 2.1 Validate the change record and run the site, repository quality,
      and semantic lint checks.

## Verification

- The site, DNS configuration, and repository quality suites passed all 32 tests.
- Strict OpenSpec validation and semantic lint of the site configuration and
  asset packaging targets passed.
- Both packaged SVG resources match the canonical data asset. Homepage,
  projects index, and project documentation index HTML include its header SVG.
- Six full-page screenshots cover these pages in light and dark themes;
  reports and PNGs are retained under ignored
  `out/sort-projects-by-creation-date/header-logo-*` paths.
