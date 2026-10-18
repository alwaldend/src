## Context

third_party/com_github_ankitects_anki/anki.BUILD.bazel has an aggregate glob covering all 25 upstream schemas and a 17-entry Go schema map. The pinned commit is 29bb700b951e3f0c0cb69b77c0180fc1fe33e6ba. The audit source is the publisher's tree at that commit; raw downloads live under ignored task scratch.

## Goals / Non-Goals

Complete Go schema coverage and aliases using the existing rules_go/protobuf toolchains. Preserve the pin, native build and approved dependency set.

## Decisions

- Extend the existing schema map with the eight missing files and their exact upstream imports. Keep one package per upstream schema and standard Go/grpc generation.
- Expose missing wrapper aliases, including existing internal links/sync packages, through schema_consumers. The application's catalog explicitly links all packages for runtime reflection.
- Validate all bindings through their consumer, including native calls to previously unavailable services, rather than adding a unit test that mirrors the map.

## Risks / Trade-offs

- Future upstream additions require updating Go integration: compare the complete tree when changing pins.
- Complete bindings enlarge the CLI's linked metadata: needed for full API coverage, without another runtime descriptor artifact.

## Candidate and next action

Baseline 3aaf25ae9. Add the audited entries and aliases, then compile the runtime catalog and exercise native CLI output.
