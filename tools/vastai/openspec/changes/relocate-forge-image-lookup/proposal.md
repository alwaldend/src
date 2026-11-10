## Why

The user requested a direct image tag instead of automatic Forge image discovery.

## What Changes

Retire the first-party image lookup executable and its protocol fixture after the personal rental switches to a direct tag. Preserve the upstream Vast CLI.

## Impact

Removes source, fixture, and consumer dependency wiring; no new dependency or runtime is introduced.
