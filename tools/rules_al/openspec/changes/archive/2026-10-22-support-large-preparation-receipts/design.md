## Context

Allow guarded publication of large path inventories.

## Goals / Non-Goals

Preparation receipts support up to 4 MiB; validation plans, state, logs, and other records retain 256 KiB reads and writes. Preserve integrity, schema, scope, and lease checks.

## Decisions

Large preparation receipts previously exceeded the 256 KiB bound. The shared reader must retain its original limit for all non-preparation callers.

## Risks / Trade-offs

Large guarded publication, oversized preparation rejection, and oversized validation plan refusal are exercised end to end. Exact candidate validation remains required before publication.
