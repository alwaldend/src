# Shared backend extraction validation

The source is `b3c926655884ab070d628314b907950707306cc9` from PR #191.
The independent base is `a2c655b20432b78e4feb75e16da1a95b73b3d30c`.
Browser implementation remains in follow-up changes; the shared core retains
its narrow future-consumer visibility grants without requiring that package.

## Behavior evidence

The following six Bazel tests pass with synchronous startup:

- `//projects/anki_grpc/pkg/ankiclient:ankiclient_test`
- `//projects/anki_grpc/internal/wasmhost:backend_test`
- `//projects/anki_grpc/test/wasm:worker_test`
- `//projects/anki_grpc/test/wasm:inprocess_test`
- `//projects/anki_as_code/internal/collection:store_test`
- `//projects/anki_as_code/pkg/ankiascode:ankiascode_test`

Native artifacts verify command cancellation, synchronous startup cancellation,
concurrent Stop deadlines, cleanup ordering, error trailers and independent
instances. The updated startup E2E failed before the synchronous change with
`Start returned before initialization completed`; the full native suite passes
after inlining initialization into Start. The timezone artifact verifies New York gaps/overlaps,
Lord Howe half-hour transitions and Apia's skipped date through Rust Chrono.
Managed collection and application checks retain database and package outputs.

## Independent review

The initial review identified a cancellation/publication race in asynchronous
startup. The user selected synchronous startup, which removes the competing
startup wait: Start owns initialization and failure cleanup before returning.
The revised native E2E checks cancellation alone and with concurrent Stop;
each emits a `direct-startup-*.json` report. Fresh independent review found no
actionable defects in the synchronous contract, including its accepted
stalled-initializer limitation. Partial-factory failures, slow or failing
destruction and terminal-command errors were reviewed statically; the native
E2E controls the post-factory context check.

## Delivery evidence

Task-local `out/browser-pr-split/` retains
the source manifest, shared-file patches, review result, failure and recovery
logs, artifact inventory and guarded preparation/validation/publication receipts.
Those receipts own final candidate identity and publication state. The selected
final plan requires root quality, all affected semantic lint, both Rust Clippy
targets and the six compatibility tests above. The synchronous-start revision's
evidence is under `out/browser-pr-split/synchronous-start/`; its publication is
held until the parent completes composed browser validation.
