## Purpose

Keep repository browser rendering reproducible while retaining Chromium's
process sandbox and applying available stable browser security updates.

## ADDED Requirements

### Requirement: Reproducible patched browser

The shared browser dependency SHALL use a stable publisher release with archive
integrity and include the security fixes motivating this update.

#### Scenario: Resolve the browser

- **WHEN** a repository renderer resolves its declared browser dependency
- **THEN** it receives the exact pinned archive verified by its integrity value
- **AND** the browser version includes the October 27, 2026 security fixes

### Requirement: Sandboxed rendering

The maintained Draw.io and code screenshot launchers SHALL retain Chromium's
process sandbox and SHALL fail if it cannot initialize, without retrying with
sandboxing disabled. The maintained browser video example SHALL follow the
same launch policy.

#### Scenario: Render under Bazel

- **WHEN** a supported non-root host runs the rendering commands
- **THEN** Chromium's renderer confinement remains active
- **AND** the existing SVG and PNG output contracts are preserved

#### Scenario: Sandbox unavailable

- **WHEN** the host prevents Chromium from initializing its sandbox
- **THEN** the command fails instead of launching an unsandboxed browser

### Requirement: Preserve content restrictions

Draw.io SHALL continue rejecting non-local page resource requests, and code
screenshots SHALL continue treating supplied source as inert text, disabling
page JavaScript, and rejecting page resource requests.

#### Scenario: Capture source containing markup

- **WHEN** a source file contains HTML or JavaScript text
- **THEN** the screenshot displays that text without executing it
