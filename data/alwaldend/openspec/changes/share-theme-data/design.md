## Decisions

Proceed with the existing identity-data owner. Keep neutral color values,
light/dark semantic roles, and diagram font choices in one JSON asset. Website
Sass and Mermaid variables are generated projections through the existing
`rules_template` tool, not independently edited copies. Bootstrap/Docsy rules
remain in the website. Mermaid preset configuration and variable mappings are
identity assets under `data/alwaldend/assets/mermaid`; generic theme rules stay
in `tools/mermaid`. This keeps ordinary diagram inputs out of the tools tree.

A `MermaidThemeInfo` provider carries concrete configuration and palette JSON
files, the pinned font, and its license. The generic SVG renderer retains its
historical standalone configuration input; native SVG consumers
select the shared preset as one target. The spatial renderer remains in PR #211 with the existing Sass input API.
Retain `mermaid_palette` as a consumer mapping compiled from generated Sass,
and export its configuration/palette/font inputs through `SITE_MERMAID_THEME`. Chrome remains responsible for rendering and text
measurement, but is no longer used to recover colors from CSS.

Do not merge the historical hand-drawn preset into the shared neutral identity.
Do not move website component styling or Three.js materials into data. Reusing
the existing template tool avoids adding a new compiler or theme project.

## Verification

Reuse the existing website and native SVG/raster browser suites, plus the
application stylesheet consumers. Compare generated Sass with the original
source and inspect both diagram appearances. Delivery receipts own exact-candidate validation and publication.
