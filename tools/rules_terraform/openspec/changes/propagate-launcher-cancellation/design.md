## Context

The launcher currently creates Terraform with exec.Command and has no signal
context. The rental caller's default cancellation kills only this launcher.

## Decisions

Use a launcher-owned signal context for every Terraform subprocess. On cancellation,
interrupt the child and retain a bounded wait/kill fallback. The rental caller
must likewise interrupt its launcher rather than immediately killing it, with a
longer outer fallback so the inner child can exit first.

## Failure cases before implementation

An end-to-end fixture uses the real launcher between the rental command and a
fake Terraform child that waits for a signal. Cancelling the rental command
must produce the child's interruption checkpoint, reach the release operation,
and write the receipt. The pre-fix fixture must reproduce the held output pipe.
