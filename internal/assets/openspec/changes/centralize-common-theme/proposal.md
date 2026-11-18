## Why

The website and diagrams share a neutral Alwaldend palette. Forgejo uses a
different palette, and the shared theme sources are spread across owners.

## What Changes

- Consolidate the common theme sources and identity artwork in `internal/assets`.
- Apply the common light and dark palettes to Forgejo through its custom header.
- Match the site’s navigation, primary controls, prose links, and footer treatment.
- Use the shared Alwaldend monogram for Forgejo’s SVG favicon and navigation logo.

## Capabilities

### New Capabilities

- `forgejo-theme` defines how the common palette applies to Forgejo.

### Modified Capabilities

None.

## Impact

The Hugo site, application stylesheets, Mermaid renderers, and Forgejo's Ansible
package consume the relocated assets.
