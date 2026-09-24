## Why

Visitors need a themed way to browse public release files and see downloads
on release pages. The main website will also become the first archived static
site deployed through the release tool to the download hosts.

## What Changes

- Add `/downloads/` and a Downloads link in the existing site header.
- Add one reusable JavaScript listing component and Hugo partial/shortcode,
  consuming Nginx JSON listings from `https://download.alwaldend.com`.
- Reuse the site's existing presentation and embed the same listing component
  in release pages, scoped to the corresponding project and version.
- Package the built website as a release archive for SSH deployment and
  selected-release serving at `https://alwaldend.com`.
- **BREAKING**: Change the production deployment contract from GitHub Pages
  publication to explicit release-tool SSH deployment, with a separately
  authorized DNS cutover. Preserve local preview and unrelated staging.

## Capabilities

### New Capabilities

- `download-browser`: The themed downloads page and reusable live listing
  component shared with release pages.

### Modified Capabilities

- `project-alwaldend-com`: Replace the production publication requirement with
  release-tool SSH deployment while retaining preview and the shared theme.

## Impact

Affected areas include Hugo navigation, content, assets, shared release
rendering, release packaging, and production deployment wiring. Nginx and
storage remain in the [hosting plan](../../../../../infra/download/openspec/changes/add-static-hosting/proposal.md);
the SSH protocol and lifecycle remain in the [release-tool plan](../../../../../tools/release/openspec/changes/add-ssh-deployment/proposal.md).

The existing apex DNS records stay owned by `infra/dns`. Website builds must
remain independent of the live download server and infrastructure targets.
Acceptance covers browser rendering, navigation, direct file links, and a
published fixture site. This plan does not deploy the main website.
