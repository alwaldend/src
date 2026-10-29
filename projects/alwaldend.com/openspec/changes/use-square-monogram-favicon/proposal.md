## Why

The approved S8 SVG exists in the data tree, but the site still packages the
old Architects Daughter icon. The built site must use the selected artwork.

## What Changes

- Package the canonical data SVG at the site's existing favicon resource
  path and remove the obsolete local artwork.
- Retain the shared favicon partial and its content fingerprinting.
- Update the site documentation and verify the resource referenced by built
  homepage, blog, documentation, and project pages.

## Capabilities

### Modified Capabilities

- `project-alwaldend-com`: Use the canonical square monogram as the shared
  site favicon.

## Impact

Changes site asset packaging without adding a dependency or duplicating the
SVG. WebP remains a generated preview. Deployment is a separate operation.
