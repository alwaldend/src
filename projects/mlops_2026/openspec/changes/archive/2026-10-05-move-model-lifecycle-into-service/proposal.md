## Why

Address the latest lifecycle ownership and private-state review feedback.

## What Changes

The service owns model loading, context-manager cleanup and an idempotent deletion fallback. Inject it into the application and router constructors. Keep implementation state private with explicit public app/router properties.

## Capabilities

Internal ownership and convention corrections without a new external contract.

## Impact

Recommender implementation, documentation and repository Python guidance.
