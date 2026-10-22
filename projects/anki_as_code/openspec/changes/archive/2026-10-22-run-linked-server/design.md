## Context

Master currently constructs a zero-value App and discovers the standalone worker through runfiles/PATH. The API CLI change is separate and unmerged. This change targets master independently and preserves collection and publication behavior while requiring explicit native worker configuration for library callers.

## Goals / Non-Goals

Produce one callable application executable and reuse the managed child-process backend. No changes to protobuf commands, archive policy or reconciliation.

## Decisions

- Register a help-only server command group. Its run command forwards native arguments unchanged to the linked Rust CLI, which owns subcommands, help and validation. Its serve shortcut delegates to run with serve prepended; it carries no native parsing or server implementation. Avoid importing native server code into reusable CLI scaffolding or the application workflow library.
- Construct the product App with explicit current-executable and server serve arguments. Pass these settings through existing collection scopes. The zero-value App supports configuration-only operations; native operations require an explicit executable and argument vector.
- Restore default interrupt handling when entering server run (including calls through server serve) so it can stop a blocking native server call. Offline parent commands retain their normal cancellation behavior.
- Use self-launching rather than an in-process server: the existing backend can close stdin, bound shutdown, terminate an unresponsive worker and reap it. The extra Go runtime is an accepted cost. No executable embedding, compression, extraction or cache is needed.
- Remove production runfiles/PATH discovery and implicit server arguments. Missing configuration fails through the process runner without opening a collection. Explicit process settings retain caller-owned argument vectors.
- Remove the client library's standalone-worker runtime input. Native tests declare their executable fixtures explicitly; application workflows use the linked CLI fixture. Backward compatibility with implicit library startup is not required.
- Keep this PR based on master without the API command tree. Verify the existing native collection workflows and new installed/server cases together.

## Risks / Trade-offs

- Explicit worker constructor additions overlap a small part of the API PR; preserve both owners when later rebasing that PR.
- Argument handling must not intercept native server flags or help. Pass them unchanged to the existing Rust parser.
- Executable selection uses os.Executable; packaging must install the one linked application, with no companion worker requirement.
