## Context

The apex site's `assets/` package and `hugo_landing`'s `assets/scss/` package
both end up in Hugo's `assets/scss/` lookup path for their respective builds, so
Docsy's `_styles_project.scss` hook is available to both. The two builds share no
source archive, which is why a shared style needs a declared export.

## Decisions

### Own the canvas in the reusable shell

The visual shell is what both sites have in common, so the canvas rules live in
`projects/hugo_landing/assets/scss/_shared_theme.scss`. Landing sites get them
through their own `_styles_project.scss` import. The apex site consumes the same
file through a `pkg_files` export and imports it from its `_styles_project.scss`.
The apex tree keeps only rules that are specific to the apex site.

### Express the canvas as CSS custom properties

Bootstrap components and Docsy surfaces read `--bs-*` custom properties at
runtime, so the shared file sets those once on `:root` and on both color-mode
selectors instead of restating component rules. The file reuses the Bootstrap
dark palette that Docsy already compiles, and replaces only the page background
so the canvas is pure black regardless of the selected color mode.

`_styles_project.scss` is imported after Bootstrap, so Docsy's dark palette
variables are in scope. Values are captured at build time; the emitted rules are
plain custom properties.

## Verification

- Build `//projects/alwaldend.com:site` and `//projects:landing_sites`.
- Confirm each compiled stylesheet sets a pure black page background and that
  the shared rule is the last `--bs-body-bg` declaration.
- Confirm in a browser that the computed page background is black with legible
  foreground text.

## Out of scope

Per-site presentation beyond the shared canvas, alternate palettes, and the
apex site's blog section remain with their own owners.
