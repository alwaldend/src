## Context

The shared HTTP client currently constructs APIError after reading the body and
retains no response headers. The draft CLI already prints wrapped errors, so
extending the shared diagnostic exposes the evidence without adding a command
or another request.

## Goals / Non-Goals

Keep rejected-response diagnostics useful and bounded. Successful responses,
automatic retries, live quota probes, authentication, and publication are outside
this change.

## Decisions

- Clone the seven allowed header names, retaining all values. Copying all
  response headers risks including authorization or cookies.
- Keep original values and parse timing separately. A malformed or repeated
  counter must not become zero by default.
- Evaluate future resets against the response receipt time supplied by the
  client's existing clock. Use the latest exhausted-window reset and extend it
  with valid later Retry-After advice; keep standalone advice separate.
- Treat the result as a not-before boundary, not upload eligibility. A service
  outage or another unknown constraint can outlast a quota reset.
- Retain APIError alongside body read or close failures so callers retain both
  the HTTP evidence and the underlying I/O cause.

## Risks / Trade-offs

- Headers may be incomplete or ambiguous. Preserve them but withhold the
  combined boundary when an exhausted window cannot be timed reliably.
- Retry-After seconds can overflow a duration. Reject invalid or overflowing
  parsed timing while retaining the original diagnostic value.
- Cached errors reflect their observation time. Diagnostics identify UTC
  boundaries and never trigger a wait or retry themselves.

## Validation

Write HTTP E2E cases before implementation and observe the original header loss.
Use only a local recorded-response service with synthetic credentials, and
preserve safe error/header/request evidence under `out/x-article-rate-headers/`.

On 2026-10-02, 27 preimplementation HTTP cases exposed the discarded headers and
missing readable timing; a truncated response also lost APIError entirely. The
red run is recorded in `out/x-article-rate-headers/red.log` (Bazel invocation
`f4567f60-e756-4f12-8147-fc848d326130`).

The implementation on base `89ec9927` passed the full uploader E2E target and the
conversion command test (Bazel invocation
`6d96f75d-1118-40e3-92bc-2f24c53cf215`). Strict change validation passed. Inspected
evidence under `out/x-article-rate-headers/evidence/` confirms both-window timing,
Retry-After extension, unknown timing for 503 with one remaining slot, and
preserved headers plus the `unexpected EOF` cause for a truncated rejection.
Each case made exactly one request to the local fixture service; no real X
request was made. Independent review approved the implementation and inspected
the emitted error evidence with no blockers. Final candidate formatting, gates,
and delivery remain parent-owned.
