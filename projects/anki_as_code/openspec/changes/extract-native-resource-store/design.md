## Context

Export and reconciliation already read the same generated Anki services. Their derived models and state responsibilities differ. See proposal.md for the extraction scope.

## Goals / Non-Goals

Share native client construction and resource enumeration. Keep configuration paths, state bindings, projection, planning, mutation ordering and backend lifetime outside the package. No new dependencies or hand-written native CRUD wrappers.

## Decisions

Add internal/resource.Store with documented exported generated client fields, New(connection), and Snapshot(ctx). Callers use those native clients for resource operations without forwarding methods. Snapshot contains upstream protobuf values and preset ID/name references. Each call reads afresh and orders decks by title and other lists by identity. Consumers own their returned values and project their own models. The store neither caches nor mutates snapshots.

A shared raw snapshot removes duplicate enumeration without sharing reconciliation models. Exposing the existing generated clients avoids another method layer. A generic resource abstraction would duplicate Anki's schema and error handling without adding behavior, so retain native types.

## Risks / Trade-offs

Sequential native reads are not transactional; callers provide exclusive offline access. Full snapshots still enumerate per entity through Anki's API and do not imply batching. Native clients permit writes; consumer packages retain managed-scope policy and mutation ordering.

## Migration Plan

Publish this library independently against master. Keep required identical prerequisite code in #128 and #160 so their integrations remain buildable; remove overlap by rebasing after this PR merges.
