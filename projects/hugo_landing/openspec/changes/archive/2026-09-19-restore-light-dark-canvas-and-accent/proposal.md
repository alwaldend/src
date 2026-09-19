## Why

The shared canvas rules force one palette on every color mode: `:root`,
`[data-bs-theme="light"]`, and `[data-bs-theme="dark"]` all receive the same
Bootstrap dark palette and the same `$black` background. Selecting light mode
therefore changes nothing, so the `showLightDarkModeMenu` toggle is inert. The
site also renders Bootstrap's stock `#0d6efd` blue, which is not the accent
the site intends.

## What Changes

- Give light mode and dark mode distinct canvases: pure white (`#fff`) in light
  mode and pure black (`#000`) in dark mode, each with the matching Bootstrap
  palette for text, surfaces, borders, and links.
- Make the accent `#7c3aed` by assigning `$primary` before Bootstrap compiles,
  so links, buttons, badges, focus rings, and the dark-mode link tint derive
  from it instead of restating each rule.
- Replace the apex footer's site-local blue link color with the shared canvas
  link color, so no site-local accent copy remains.
- Keep one implementation: the shared shell owns both the canvas and the
  accent, and the apex site consumes them through the existing declared
  exports.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-hugo-landing`: The shared canvas gains distinct light and dark
  palettes and a shared accent, replacing the single dark palette that both
  modes received.

## Impact

- `projects/hugo_landing/assets/scss/_shared_theme.scss` splits the single
  palette into light and dark mixins and selects the canvas per mode.
- `projects/hugo_landing/assets/scss/_shared_accent.scss` is new and is
  imported before Bootstrap by each site's `_variables_project.scss`.
- `projects/hugo_landing/assets/scss/_variables_project.scss` is new and
  imports the shared accent for the landing sites.
- `projects/alwaldend.com/assets/scss/_variables_project.scss` and
  `_styles_project.scss` import the shared accent and drop the hardcoded
  footer blue.
- Both asset `BUILD.bazel` packages declare the new shared accent export.
