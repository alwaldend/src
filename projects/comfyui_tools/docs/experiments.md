# Buffered ComfyUI experiments

`//projects/comfyui_tools/cmd/experiments` runs arbitrary workflows against an
existing endpoint through the native ComfyUI CLI. It maintains one executing
job plus a configurable pending buffer after warm-up, independently observes
completion and downloads original media. It does not provision instances,
install models or manage credentials. The caller supplies an authenticated CLI
and a producer containing a finite set of ready candidates. The
[experiment skill](../skills/comfyui-experiments/SKILL.md) owns the direction,
evaluation, scrap/refine/fork procedure. Queue buffering schedules ready work;
it does not evaluate outputs or authorize endless successor generation.

Build the executable and any launchers before starting the sustained process:

```sh
bazel build //projects/comfyui_tools/cmd/experiments
bazel-bin/projects/comfyui_tools/cmd/experiments/experiments_/experiments \
  run --config /absolute/task/config.json --directory /absolute/task/run
```

The runner requires Linux process groups and advisory file locking. Child
commands receive literal argv without shell expansion and inherit the caller's
environment. Use absolute paths for executable wrappers and all inputs. Do not
start nested `bazel run` children beneath a long-lived Bazel invocation; build
first or prepare a `--script_path` launcher. Live runtime commands need network
access; builds and E2E fixtures do not.

## Configuration and producer contract

```json
{
  "cli": ["/absolute/authenticated-comfy"],
  "producer": [
    "/absolute/workflow-producer",
    "--matrix",
    "/absolute/matrix.json"
  ],
  "queued_buffer": 2,
  "max_experiments": 12,
  "timeout_seconds": 3600,
  "duplicate_attempts": 8
}
```

For a separately and explicitly authorized continuous stream, replace `max_experiments` with
`"continuous": true`. The two modes are mutually exclusive. Iterative direction development uses
finite runs and review before successor jobs; an earlier utilization request
does not select continuous mode. Buffer defaults
to two pending jobs, timeout to 3600 seconds per subprocess operation, and
duplicate attempts to eight. The timeout includes native queue waiting; allow
enough time for the configured lookahead. Native execution can continue after
a local timeout or interruption.

The producer receives one JSON object on stdin:

```json
{ "sequence": 0, "directory": "/absolute/task/run/000000" }
```

It writes exactly one JSON object on stdout, with diagnostics on stderr:

```json
{
  "workflow": "/absolute/task/run/000000/candidate.json",
  "description": "Compare raking light with the earlier diffuse-light composition",
  "identity_exclusions": ["/7/inputs/seed", "/10/inputs/filename_prefix"],
  "notes": {
    "positive_prompt": "A geometric still life in raking light",
    "negative_prompt": "blur"
  }
}
```

The workflow is a native UI or API JSON object, at most 64 MiB. The producer
owns model bindings, parameters, reference files, prompts and batch size; no
such fields are hard-coded into the runner. Return `{"done":true}` to end
the stream gracefully. Sequence numbers are monotonic requests, including
rejected duplicate candidates, rather than counts of accepted jobs.

Stop is rechecked after a producer returns and before submission, so a drain
requested during a slow producer does not add another job.

The runner snapshots exact workflow bytes before submission. Duplicate
identity is canonical JSON after removing explicitly declared JSON-pointer
fields. Object fields are removed; array positions become null so indexes do
not shift. Every pointer must exist; root exclusions are forbidden. Exclude
seeds, output prefixes and irrelevant UI bookkeeping consistently. Keep
meaningful parameters in the identity. This detects declared structural
duplicates; the producer and agent still assess useful semantic variation.

## Ownership, interruption and control

The private run directory contains `journal.json`, a single-controller lock,
per-experiment native receipts and workflow snapshots, original `media/`, and
`gallery.html` with collapsed notes, prompts, settings and hashes. Journal decoding supports the records its writer can produce, including
journals larger than 64 MiB. Opaque numeric notes retain their exact JSON
values through execution, restart and report generation. Directory
permissions are restricted to the current account. CLI captures are bounded
and raw diagnostics never enter the gallery. Inspect inputs and notes for
protected data before sharing any artifact.

```sh
experiments status --directory /absolute/task/run
experiments stop --directory /absolute/task/run
experiments stop --directory /absolute/task/run --cancel-owned
```

`stop` writes a control request for the running controller to drain. Explicit
`--cancel-owned` additionally targets only recorded outstanding prompt IDs
through native `jobs cancel`; an idle run consumes the request when resumed.
Cancellation is reconciled against the exact native job status: a rejected
interrupt retains ownership, and a job that completed before cancellation
retains artifact collection. A successful cancel command alone is not proof
of termination. No operation clears the global queue or changes an instance.

SIGINT/SIGTERM stop local subprocess groups and preserve remote jobs. Resume
with the same configuration and directory. Recorded prompt IDs are observed
again without resubmission. A durable successful submission receipt can recover
an ID after interruption before journal update. A missing acknowledgment is
recorded as `submission_unknown` and prevents automatic new submissions.
Inspect the private receipt and endpoint before any manual reconciliation.
A valid successful CLI envelope followed by a wrapper cleanup error retains
the acknowledged result and sets `cleanup_warning`; it is never blindly
resubmitted.

Retrieval transport failures get at most three attempts, with one- and
two-second delays and fresh artifact directories; generation is never
resubmitted by those retries. Invalid or escaping media is rejected immediately.

Producer, native job and exhausted artifact failures stop replenishment and retain
failure phases for inspection. Observer transport errors and timeouts retain
the owned ID in `observation_failed`, so restarting can observe it again or
apply an owned cancellation request. These errors do not establish a terminal
native outcome. Confirmed failure and cancellation summaries in unsuccessful
wait envelopes remain terminal across restart. A new configuration needs a new run after owned
work has been reconciled; do not use a new directory to evade an ambiguous
submission. Control requests persist, including across restarts.

`queued` in the local journal means submitted and awaiting a terminal native
observation. It does not distinguish executing from pending or establish GPU
utilization. Observe the endpoint for that distinction. Slow producers,
authentication, endpoint outages or insufficient lookahead can still exhaust
the buffer. Original outputs and execution success do not establish visual
quality, text accuracy or reference fidelity.

## Verification

```sh
bazel test //projects/comfyui_tools/test/experiments:experiments_test
```

Command-level fixtures verify buffering during slow downloads, seed-only
duplicate rejection, restart without resubmission, concurrent-controller
exclusion, lost acknowledgments, explicit producer exhaustion, cleanup errors,
owned cancellation, failed jobs and escaping media. They emit JSON receipts
under Bazel's undeclared test outputs. They use no live endpoint, credentials
or rentals. The [agent skill](../skills/comfyui-experiments/SKILL.md) covers
useful variations, feedback, queue observation and visual review.
