## 1. Common assets

- [x] 1.1 Relocate the shared assets and update Hugo, application, and Mermaid consumers.

## 2. Forgejo

- [x] 2.1 Generate and package the Forgejo theme with its existing mobile diff rules.
- [x] 2.2 Verify the affected builds and representative output.
- [x] 2.3 Package and verify the canonical monogram at both Forgejo SVG icon URLs.
- [x] 2.4 Match and verify site navigation, button, prose-link, focus, and footer styling.

The Hugo source archive includes the common Sass partials and monogram at their
public asset paths. The native Mermaid renderer, diagram snapshots, and site
tests pass. Forgejo 17.0.1 served the packaged header after the built-in theme.
Chromium confirmed the palette, typography, and radii in explicit Light and
Dark and in Auto with both preferences. A narrow viewport confirmed wrapping
and horizontal diff scrolling styles.

The packaged Ansible branding tasks installed the same SVG at both public icon
paths. A running loopback Forgejo served both files byte-for-byte and linked
them from the page without a restart. Chromium confirmed the navigation
weight and rounding, primary fill, keyboard focus outline, prose underline
offset, and footer border.
