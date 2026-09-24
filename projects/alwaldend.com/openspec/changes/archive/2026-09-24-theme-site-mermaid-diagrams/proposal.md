## Why

Mermaid images retain a white canvas when the site switches to dark mode.
Their presentation should follow the site's shared palette and theme selector.

## What Changes

- Derive adaptive SVG assets from documentation diagrams using the site preset.
- Use the site's existing light/dark tokens without a browser diagram runtime.
- Use clean, rounded nodes and containers, straight connectors, and a pinned sans-serif font.
- Preserve alternative text and leave every blog diagram unchanged.
- Verify actual browser rendering and theme selection with repeatable artifacts.

## Capabilities

### Modified Capabilities

- `project-alwaldend-com`: Mermaid images follow the site's selected color mode.

## Impact

The shared Hugo resource pipeline, site assets, and browser validation change.
The shared renderer in `tools/mermaid` gains a site preset; the default renderer
and blog publication assets remain unchanged. The Host Bot documentation selects
the site preset from its existing diagram source. Deployment is outside this change.
