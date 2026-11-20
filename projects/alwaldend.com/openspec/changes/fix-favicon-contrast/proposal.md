## Why

A browser can resolve the favicon’s dark-mode media query while drawing a
light tab. The transparent favicon then places light artwork on a light tab.

## What Changes

- Give the browser tab icon its own high-contrast background.
- Keep the header logo transparent and adaptive.
- Fingerprint the rendered favicon so browsers receive its new URL.

## Capabilities

- `project-alwaldend-com` covers the shared favicon across site pages.

## Impact

Changes the generated favicon and its URL. The header keeps the shared S8
geometry and its current presentation.
