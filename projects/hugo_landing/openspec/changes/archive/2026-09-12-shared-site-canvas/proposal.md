## Why

The apex site and the project landing sites are separate Hugo builds that render
the same Docsy component set. Their page canvas is currently each build's own
concern, so a visual decision such as the page background has no owner and would
have to be restated in every site tree.

## What Changes

- Provide one shared Docsy canvas style implementation in the reusable site
  shell, covering the page background and the foreground palette both sites use.
- Package that shared style into a landing site's assets tree and expose it for
  the apex site to package into its own tree, so both consume one implementation.
- Keep the apex site's remaining presentation rules site-local.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-hugo-landing`: The reusable shell now also provides the shared canvas
  styles and exposes them for the apex site, and the rendered landing sites use
  the same canvas as the apex site.

## Impact

- `projects/hugo_landing/assets/scss/` gains the shared canvas style file and the
  export that the apex site consumes.
- `projects/hugo_landing/assets/scss/_styles_project.scss` imports it, so every
  rendered landing site picks it up.
- `projects/alwaldend.com/assets/` packages the shared file into its own
  `assets/scss/` tree; the apex site keeps only its own presentation rules.
- The apex site's blog section and its packaging are owned by
  `projects/alwaldend.com/openspec/changes/add-docsy-blog-section`.
