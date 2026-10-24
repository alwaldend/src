## Context

See proposal.md for scope. The immutable source is
`b3c926655884ab070d628314b907950707306cc9`; this independent branch starts at
`a2c655b20432b78e4feb75e16da1a95b73b3d30c` and targets master.

## Goals / Non-Goals

Preserve one backend owner through startup, command dispatch and shutdown.
Keep collection policy with callers and omit browser implementation.

## Decisions

- Extract the shared direct-client lifecycle and run initialization directly in
  Start. Register startup as active work before unlocking, check context before
  and after the factory, and drain that registration after failure cleanup.
  Failed or stopped initialization closes its unpublished backend. Stop rejects
  new work and joins initialization, commands and destruction. This removes the
  startup goroutine, completion channel and cancellation select while retaining
  asynchronous command waiting and shutdown cleanup.
- Update the existing native E2E before production changes. Pause at the
  post-initialization context check, cancel startup and verify that Start stays
  blocked. Exercise concurrent Stop's deadline and keep Done open until startup
  releases its backend. Check both cancellation alone and cancellation with Stop.
- Include the shared Rust C ABI by source from the native adapter. Existing
  handle and buffer ownership remains unchanged.
- Parameterize existing timezone calculations by location; WASI continues to
  pass time.Local and retains the embedded timezone database.
- Keep existing E2E consumers as acceptance evidence. Record failure cases
  before extraction: premature Done, late connection publication, lost native
  errors, command replay, timezone transition drift and broken Rust linkage.

## Risks / Trade-offs

- Native initialization cannot be interrupted. Start waits for the factory and
  failure cleanup even after cancellation; a stalled initializer delays both.
  This accepted tradeoff preserves the original native startup contract.
  Command cancellation and Stop deadlines still bound their callers' waits.
- Moving shared source can break build inputs or visibility. Compile and lint
  the native bridge and exercise real generated clients and WASI workers.
- The directclient and backend BUILD files retain the source's narrow browser
  consumer visibility. These grants introduce no dependencies on browser code;
  validate native targets with that future package absent. Browser follow-ups
  own their implementation and OpenSpec record.
