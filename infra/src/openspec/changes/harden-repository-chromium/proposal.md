## Why

The shared browser is pinned to Chromium 150.0.7871.24, and two renderers
explicitly disable Chromium's process sandbox. Updating the shared dependency
and retaining its sandbox reduces exposure when rendering repository content.

## What Changes

- Upgrade the existing Linux Chrome Headless Shell pin to the publisher's
  current stable release, with archive integrity.
- Remove sandbox-disabling flags from Draw.io, code screenshots, and the
  browser video example. Preserve their existing input and network controls.
- Verify real rendering and renderer sandbox state under the existing Bazel
  execution policy; fail rather than retry without Chromium sandboxing.
- Correct current documentation describing unsandboxed launches.

## Capabilities

### New Capabilities

- `repository-browser-security`: Shared browser release and sandbox guarantees
  across repository rendering consumers.

### Modified Capabilities

None.

## Impact

The shared pin remains owned by `third_party/com_google_chrome_headless_shell`;
launchers remain owned by `tools/drawio` and `tools/code_screenshot`; the video
example remains owned by `tools/agents`. This change records their shared
security integration. T3 updates and persistent host configuration are excluded.
Acceptance requires command E2Es, representative SVG/PNG inspection, runtime
sandbox evidence, and the repository delivery checks.
