## Why

Application workflows should pass addresses to a backend-owned constructor while remaining usable without optional native dependencies. Engine maps and concrete adapter construction do not belong in the application.

## What Changes

- Use the Client lifecycle contract and inject one address-only library constructor.
- Keep application and collection-store libraries pure Go; the executable explicitly imports optional compiled-engine composition.
- Preserve existing subprocess, TCP, Unix and archive workflows and delegate hosting to the backend library.
- Keep direct in-process execution and its application workflows in PR #184, independently based on master.

## Capabilities

### New Capabilities

- `client-connection-delegation`: Backend-owned construction with optional dependencies excluded from application libraries.

### Modified Capabilities

None.

## Impact

Application options, collection-store construction, executable wiring and affected Go callers. The backend contract is maintained by isolate-client-connection-platforms in projects/anki_grpc.
