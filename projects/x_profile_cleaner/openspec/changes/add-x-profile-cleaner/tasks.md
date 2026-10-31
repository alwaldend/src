## 1. Acceptance and policy engine

- [x] 1.1 Write offline command-level failure cases before production code;
      verify the fixture covers graph branches, foreign-parent replies, repost
      age, filtering, precedence, schema rejection, and confirmation.
- [x] 1.2 Implement the author-only inventory model, forward graph, Lua
      callbacks, query functions, and documented shims; verify E2E plan actions.

## 2. CLI and HTTP execution

- [x] 2.1 Implement timeline collection and typed post/repost actions with
      injected authentication; verify paginated fixture requests and actor scope.
- [x] 2.2 Implement plan, saved-plan apply, and inline apply with confirmation
      and yes flags; verify no mutations occur on planning, refusal, or invalid input.

## 3. Documentation and delivery

- [x] 3.1 Document API contracts, commands, example policy, credentials, and
      discovery limitations; verify examples against the built CLI and shims.
- [x] 3.2 Run E2E validation and emit a repeatable plan artifact, formatting,
      semantic lint, repository quality, and OpenSpec validation; record results.
- [x] 3.3 Commit and publish the task-owned candidate through repo-delivery;
      verify publication or report the exact blocking condition.

## 4. Review corrections

- [x] 4.1 Extend the CLI fixture for edit-history aliases, graph references to
      earlier versions, and callbacks on converging graph branches.
- [x] 4.2 Resolve owned edit-history aliases consistently and simplify operation
      dispatch, preservation traversal, and command construction.
- [x] 4.3 Verify review corrections against the CLI acceptance fixture and strict
      OpenSpec validation; use the existing delivery gates for PR publication.

## 5. Policy and command simplification

- [x] 5.1 Extend CLI acceptance cases for ID arrays, deletion precedence, and
      exact post/repost mutation paths before changing production code.
- [x] 5.2 Use named kinds and one deletion level, derive reply/quote flags,
      and accept a string or dense ID array in direct selection functions.
- [x] 5.3 Pass the shared application instance into command constructors, move
      plan-source selection into the application, and update examples and comments.
- [x] 5.4 Verify the CLI acceptance output, example policy, formatting, and
      strict OpenSpec validation; use repo-delivery for candidate gates and replies.

## 6. Mandatory safety contract

- [x] 6.1 Record non-mutating planning, exact policy selection, exact plan
      construction, and execution limited to the plan as mandatory requirements.
- [x] 6.2 Trace those boundaries and verify successful and rejected CLI cases
      against an isolated API fixture, including invalid later actions and
      stopping after the first execution failure.

## 7. Shared SDK, lazy loading, and rate limits

- [x] 7.1 Extend offline CLI cases for SDK fields, maximum-size page requests,
      lazy kind loading, shared cached queries, pacing, cancellation, and
      unconfirmed mutation responses before production changes.
- [x] 7.2 Replace handwritten HTTP calls with the shared generated X SDK,
      preserving bounded responses, private authentication, strict mutation
      confirmation, and failure handling.
- [x] 7.3 Load inventory kinds and reply/quote indexes only when the policy
      requests them; cache timeline pages and keep callback iteration stable.
- [x] 7.4 Add cancellable per-endpoint quota pacing without replaying failures
      or changing action order; document API batching limits and behavior.
- [x] 7.5 Verify representative output, strict OpenSpec validation, and all
      required delivery gates for the migrated candidate.
- [x] 7.6 Cover and implement reply-parent traversal when a quoting reply
      connects an earlier own thread to the preserved seed, retaining filters
      and excluding traversal to quoted ancestors.

## 8. Shared OAuth 1.0a credentials

- [x] 8.1 Extend the CLI fixture before implementation to verify signed reads
      and mutations, pagination parameters, fresh signatures after quota waits,
      missing credential rejection, and absence of credential values in output.
- [x] 8.2 Use the SDK's shared OAuth 1.0a signer and the Article uploader's
      four environment variable names; reuse its existing injection configuration
      without duplicating secret references or provisioning credentials.
- [x] 8.3 Update authentication requirements and operator instructions; verify
      the offline CLI output, runner build, OpenSpec validation, and delivery gates.

## 9. Discovery window documentation

- [x] 9.1 Document that discovery windows can omit preservation seeds or
      connecting posts, distinguish missing-seed failures from partial graphs,
      and distinguish discovery bounds from age predicates.

## 10. Lua discovery windows

- [x] 10.1 Add command-level cases before implementation for independently
      optional bounds, adjacent chunks, signed bounded pagination, scoped queries,
      exact mutations, missing seeds, rejected bounds, and late window changes.
- [x] 10.2 Implement set_time_window with optional RFC3339 bounds, frozen before
      inventory loading; constrain every SDK page and the cached owned inventory
      without fetching outside the window or falling back to a full timeline.
- [x] 10.3 Document examples, inclusive-start/exclusive-end boundaries, default
      unbounded discovery, API coverage limits, and incomplete thread context;
      verify the CLI artifact and required delivery gates.

## 11. Lua API review improvements

- [x] 11.1 Write command-level failure and acceptance cases for strict arguments,
      named decisions, relative windows, optional metadata, traversal context,
      overlapping preservation requests, and exact saved-plan execution.
- [x] 11.2 Reject unsupported arguments and options; add fixed-time days_ago,
      optional rule names, and informational decisions without changing execution.
- [x] 11.3 Expose language, attachments, and text entities with fresh Lua copies;
      provide preservation seed context and explicit edge direction.
- [x] 11.4 Share callback-free traversal within each declaration while retaining
      per-seed custom callbacks, cancellation, and deterministic explanations.
- [x] 11.5 Update shims and examples, inspect acceptance artifacts, and verify
      E2E behavior, repository quality, runner builds, and semantic lint for delivery.

- [x] 11.6 Separate the Go planner from the Lua adapter; register plain rule
      descriptions with opaque callback identifiers and pass plain callback data.

- [x] 11.7 Add concurrent endpoint deletion with separate quota pacing, same-target
      dependencies, failure cancellation, and precise completed-action reporting;
      revise the serial-order requirement and verify offline acceptance cases.

## 12. Cancellation and operator progress

- [x] 12.1 Add CLI acceptance cases for interrupting confirmation with stdin open,
      separate JSON and diagnostic streams, and visible discovery, quota, and apply progress.
- [x] 12.2 Make confirmation cancellable and report structured progress to stderr
      without exposing credentials or post contents or changing planned actions.
- [x] 12.3 Document diagnostics, inspect emitted evidence, and verify cancellation,
      stream separation, quota progress, and exact execution in the offline CLI suite.

## 13. Local browser workflow

- [x] 13.1 Write browser acceptance cases before implementation for CLI and picker
      policy loading, planning, exclusions, exact apply, cancellation, logs, stale
      or fabricated selections, access controls, and responsive shared-theme output.
- [x] 13.2 Expose in-memory application planning and add a loopback server with
      authenticated, serialized jobs, server-owned plans, subset apply, and log polling.
- [x] 13.3 Implement the embedded policy editor, removal review, confirmation,
      cancellation, and live logs using the exported site theme.
- [x] 13.4 Verify browser and CLI behavior, inspect light/dark/mobile artifacts,
      and document the serve command and its session boundaries.

## 14. Browser review controls

- [x] 14.1 Extend browser acceptance cases for filtered bulk selection, explicit
      exclusions, hidden-selection counts, and collapsing the policy editor.
- [x] 14.2 Improve selection visibility, text readability, and responsive layout;
      consume the site's shared filled buttons and checkbox styling.
- [x] 14.3 Verify light/dark/mobile output and exact subset execution; record the
      separately authorized single-post lifecycle.

## 15. Diagnostics and simpler review

- [x] 15.1 Cover partial-response diagnostics and credential redaction before
      implementation; retain failure without partial plans or mutations.
- [x] 15.2 Report bounded X problem details with HTTP status and pagination context;
      retry the real policy read-only and record the actual result.
- [x] 15.3 Remove redundant UI status/instruction rows and the editor toggle;
      place Load file in the policy header and truncate post previews.
- [x] 15.4 Verify CLI/browser behavior, light/dark/mobile output, and delivery gates.

## 16. Compact metadata and readable plan

- [x] 16.1 Place fixed-height scrolling logs below Policy in the left column,
      independently of list height, and put rule names beside post IDs.
- [x] 16.2 Convert the saved real-policy plan into a searchable standalone HTML
      report with exact action coverage, readable content, and no apply capability.
- [x] 16.3 Verify browser output and report integrity, then publish the UI changes.

## 17. Full-height review workspace

- [x] 17.1 Cover independent list scrolling and persistent review controls in the
      browser acceptance suite before changing the layout.
- [x] 17.2 Fill desktop height with Policy and Logs beside a scrollable Review;
      compact its header and keep text expansion inline.
- [x] 17.3 Export only planned post/reply deletions as text and links, verify
      report coverage and responsive output, and publish the UI changes.

## 18. Resume interrupted execution

- [x] 18.1 Add CLI and browser acceptance cases for missing targets, partial
      completion, lost responses, reload, cancellation, and exact resumption.
- [x] 18.2 Treat missing mutation targets as skipped; expose confirmed action
      completions without changing the selected plan or automatic retry policy.
- [x] 18.3 Retain execution progress in the server session and provide Resume
      for unfinished actions, including after a browser reload.
- [x] 18.4 Exercise the full saved batch with an offline API fixture, document
      session and error boundaries, and complete delivery checks.
- [x] 18.5 Review execution and recovery, then record a fresh real-policy plan
      in the dark browser interface without applying it.

## 19. Apply state documentation

- [x] 19.1 Maintain a readable vertical Mermaid state diagram covering confirmation,
      account validation, parallel endpoint workers, quota waits, completion,
      failure, cancellation, and resuming unfinished actions.
- [x] 19.2 Publish the diagram at the top of the project README with the site
      theme and generate a dark WebP from the same source; verify both renders.
