## Why

The browser application needs a discoverable visitor entry point. The project landing should lead to an application page that uses the main site's established theme and navigation.

## What Changes

- Add a Try it now button on the Anki as code landing page.
- Add a themed workspace at `/projects/anki_as_code/app/` with an automatically
  loaded example, one TOML editor and Open JSON, Check and Download JSON actions.
- Add Open collection, Plan, Apply and Download collection actions through the
  existing browser App, plus explicit Save locally and restoration of the last
  saved workspace on startup.
- Package application runtime files through the existing project content registry and verify the rendered site in Chromium.

## Capabilities

### New Capabilities

- `browser-app-site`: Discoverable, themed browser configuration and collection
  workspace with local snapshots and portable downloads.

### Modified Capabilities

None.

## Impact

Project-owned site content, browser assets, package mappings and browser E2E tests change. The existing apex site's layout, Bootstrap components, color modes and registry are reused without new dependencies or a parallel site. Live deployment and VS Code integration are outside this request.

The site build also requires the missing `docs` target for the existing Go gRPC package README. Add that documentation mapping without changing the dependency or its patch.
