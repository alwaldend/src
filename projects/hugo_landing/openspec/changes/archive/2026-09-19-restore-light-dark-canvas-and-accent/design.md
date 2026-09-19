## Context

Docsy compiles `_styles_project.scss` after Bootstrap, so the previous shared
canvas file set its custom properties there. That position is fine for a value
Bootstrap reads at runtime, but it cannot change anything Bootstrap derived at
build time. Two consequences produced the current defect:

- Selecting a color mode only swaps `data-bs-theme`; a rule that sets the same
  `--bs-body-bg` for `:root`, `[data-bs-theme="light"]`, and
  `[data-bs-theme="dark"]` overrides all three with one value. The observed
  stylesheet had `--bs-body-bg:#000` on all three selectors, which is why the
  toggle appeared inert.
- Bootstrap resolves `$primary` into literal values at compile time: the
  `--bs-primary*` custom properties, `.btn-*`, the `*-subtle` variants, focus
  ring colors, and the dark-mode link tint are all baked. Overriding
  `--bs-primary` afterwards leaves those rules on the original color.

## Decisions

### Split the palette per mode and set it after Bootstrap

Both modes need the same set of properties but different values, so the shared
file exposes `alwaldend-light-palette` and `alwaldend-dark-palette` mixins and
applies them to `[data-bs-theme="light"]` and `[data-bs-theme="dark"]`
respectively. `:root` shares the light selector so a page without the attribute
still gets a defined canvas.

The two mode selectors are equally specific and mutually exclusive, so the
later rule in source order wins only within its own mode; the explicit mode
selector also overrides the `:root` default. Values are restated from the
Bootstrap light and dark palettes that Docsy already compiled, and only the
page canvas is replaced, so text, surfaces, and borders stay legible.

`html` reads `var(--bs-body-bg)` instead of a captured `$black`, which keeps
the pre-stylesheet paint and overscroll surface in step with the active mode.

Docsy's footer extends `.td-box--dark`, whose literal `#212529` surface would
put white text on a white canvas in light mode. The shared file therefore pins
the footer to the active canvas and body color rather than to that literal
surface.

### Set the accent before Bootstrap

`$primary` must be assigned where Docsy loads `_variables_project.scss`, which
is before Bootstrap imports its variables. The shared accent lives in
`_shared_accent.scss`, and each site's `_variables_project.scss` imports it, so
one value still reaches every derived rule and no site repeats the hex.

The assignment is deliberately not `!default`: Docsy's `variables_forward`
assigns `$primary` before the project file is loaded, so a `!default` would be
inert. Compiling the accent through Bootstrap yields `--bs-primary:#7c3aed`,
`--bs-link-color:#7c3aed` in light mode, and the derived
`--bs-link-color-rgb:176, 137, 244` tint in dark mode.

### Keep one implementation across sites

The apex site builds from its own source archive, so it consumes both shared
files through `pkg_files` exports from the landing shell, matching the existing
`shared_theme` pattern. The landing sites pick the accent up through their own
`_variables_project.scss`. No site restates the accent value.

## Verification

Build `//projects/alwaldend.com:site` and `//projects:landing_sites`, then
inspect the compiled stylesheet and computed styles:

- The stylesheet's last `--bs-body-bg` for light mode is `#fff` and for dark
  mode is `#000`, and the override block follows Bootstrap's `_root.scss`.
- `--bs-primary` and the light-mode `--bs-link-color` are `#7c3aed`.
- In a browser, the computed body and footer background is `rgb(255,255,255)`
  in light mode and `rgb(0,0,0)` in dark mode, with links resolving to the
  accent in light mode and the lighter tint in dark mode.
- No site-local accent hex remains.

## Out of scope

Per-site presentation beyond the shared canvas and accent, alternate palettes,
and the choice of which semantic colors Bootstrap badges use.
