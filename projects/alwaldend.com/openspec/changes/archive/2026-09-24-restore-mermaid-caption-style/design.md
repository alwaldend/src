## Decision

Restore the previous appearance through Mermaid's `themeCSS` and flowchart
configuration. Dagre retains placement and edge routing. The renderer resolves
the shared Sass palette once per color mode and exposes those same colors as
CSS custom properties in Mermaid's input. It writes the returned SVG unchanged.

The preset restores the former node/rank spacing and node padding. Caption
borders, corner radii, and padding are ordinary theme styles measured by Mermaid.
The existing `subGraphTitleMargin.top: -16` setting places title plates on their
container borders, matching the earlier presentation. This is explicit native
configuration; it does not move or reorder generated SVG elements.

This retains custom theme CSS and the negative native title margin. Removing
all custom styling would not reproduce the earlier appearance. It removes
post-processing, not the styling itself.

## Scope correction

An earlier scratch prototype introduced nested layout and obstacle routing.
The user rejected that expansion. No prototype code enters the implementation;
Mermaid's existing Dagre implementation remains unchanged.

## Verification

Before implementation, the browser acceptance check was extended to verify
that visible caption frames are rounded, bordered, match the selected canvas,
and fit their measured foreignObject. Existing checks cover native node label
bounds, theme switching, mobile/print, embedded fonts, and blog preservation.
Inspect actual page screenshots as well as these mechanical checks.

## Work state

Native preset implemented. The browser E2E passed on 2026-09-24, including
caption bounds in both modes, menu behavior, mobile/print, and historical blog
image identity. Page screenshots were inspected in light and dark mode.
Delivery gates and publication are tracked by the delivery receipts.

## Session review

`MERMAID-SCOPE-EXPANSION`: the scratch prototype added a layout engine layer and
router before establishing that the earlier look required either. The user's
correction narrowed acceptance back to appearance. Future changes should first
compare the previous visual preset against the native theme path, and reserve
layout changes for an explicit layout requirement. No shared policy change.
