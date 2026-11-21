## Why

A migration can pass its new entry point's tests while leaving an obsolete
execution path reachable through the same workflow.

## What Changes

Require the repository policy read by review-agent to trace superseded paths
through callers, configuration, and tests.

## Capabilities

- `project-agents` gains a migration review check for obsolete paths.
