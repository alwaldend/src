## Why

Server commands bypass App and are assembled in the executable, so the application interface no longer represents all user operations. Package comments need explicit ownership contracts checked against actual callers and side effects.

## What Changes

- Route server operations through App and inject native hosting from the executable.
- Keep Cobra commands in internal/cli and limit them to syntax, context and stream handling.
- Review all Anki as code Go packages and document responsibilities and exclusions in doc.go.
- Update the architecture diagram and preserve installed CLI, help and signal behavior.

## Capabilities

### New Capabilities

- `application-server`: Application-owned native server operations consumed by CLI and library callers.

### Modified Capabilities

None.

## Impact

Application API, command composition, package documentation and behavioral checks. No dependencies, native server changes, deployment or backward-compatibility requirements.
