# Acceptance evidence

The command-level E2E fixture supplies unrelated neutral API workflows and a
simulated native CLI. It exercises a serial endpoint while downloads take
longer than rendering, verifies original PNG content and hashes, rejects
seed-only repetitions, excludes concurrent controllers, resumes acknowledged
IDs, refuses lost acknowledgments, drains or cancels owned jobs, handles
producer exhaustion and rejects failed jobs and escaping media. Its JSON
receipts are emitted through Bazel undeclared test outputs.

The initial restart fixture signalled after remote acceptance but before the
acknowledgment was durably recorded. That boundary belongs to the explicit
ambiguous-acknowledgment case. The corrected restart case waits for recorded
IDs. The subprocess fixture also needed atomic endpoint-state writes so killing
an observer could not truncate the simulated remote state. Failure diagnostics
are now retained in test output artifacts.

On the already-authorized existing endpoint, an external task-local producer
supplied distinct native UI workflows with two-image batches. The runner
acknowledged ownership before waiting, produced original 1536 by 1280 images,
recorded hashes and refilled replacements independently of retrieval. A native
observation recorded one owned job running and two pending. Restarting the
local controller adopted its existing IDs rather than submitting those studies
again. No instance lifecycle or persistent host configuration was changed.

Source contains no subject-specific producer, private endpoint, credentials,
rental state or task artwork. These remain in ignored task artifacts. Runtime
execution and observed queue state establish neither continuous uptime nor
artistic quality. Visual review found inadequate text accuracy and material
integration in some outputs; feedback changes the external producer without
altering completed workflow snapshots.

The skill's offline Promptfoo target validates staging and configuration only.
It does not establish agent reasoning or visual-review quality. No live model
calls or rentals are included in ordinary tests.

Review correction E2E cases first reproduced all four findings against the
previous candidate. The corrected command retains live ownership after an
observer disconnect for both resume and owned cancellation, resumes a valid
journal larger than 64 MiB, prevents a post-stop producer result from submitting,
and preserves the integer seed 18446744073709551615 through initial execution,
restart and gallery output. These are command-level fixtures with local native
CLI simulation; they introduce no live allocation or model calls.

Additional command E2E cases reproduced unsuccessful native wait envelopes being
mistaken for observation errors, rejected interrupts being marked cancelled,
and completion races discarding original outputs. The corrected fixture follows
the pinned CLI's error/details summaries, verifies terminal failure and external
cancellation persist across restart, and distinguishes a recoverable wait
timeout. Cancellation cases retain live ownership after a rejected interrupt,
including a misleading cached cancellation status, and collect media when the
native job completed before cancellation.

The agent procedure now uses finite ready directions with one unevaluated job
per direction and explicit scrap, refine or independent fork decisions. Retired
directions cannot recur automatically, and pauses persist. Offline skill cases
cover these decisions and the queue-only display preference; configuration
validation does not establish that an agent will make correct visual judgments.
The low-level runner retains continuous mode for separately authorized uses;
it does not perform direction evaluation or invent successor jobs.

A mixed-outcome command fixture additionally reproduces an external cancellation
followed by another observer's disconnect. Recovery preserves the cancelled
experiment's terminal failure and does not replenish work while reattaching to
the disconnected job. Requested cancellations retain their empty failure field.
