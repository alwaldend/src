## Acceptance evidence

The implemented entry point is `al bazel mcp serve`, as clarified by the user.
The six advertised tools use the official MCP Go SDK at the repository's existing
v1.8.0 pin. Existing AL affected-test selection remains unchanged.

- Standalone AL `bazel test //...`: 24 tests passed, including executable MCP
  discovery, all six operations, validation failures, child failure, truncation,
  timeout, request cancellation, and transport disconnect.
- Standalone AL `bazel build //...`: 124 targets built successfully.
- Root `bazel test //:repo_quality_test`: all 27 checks passed.
- Root `bazel build --config=lint @rules_al//cmd/al:al
@rules_al//cmd/al:lifecycle_e2e_test @rules_al//internal/bazelcmd:bazelcmd`:
  configured semantic aspects passed. The standalone workspace does not define
  the root's `lint` configuration, so that check belongs in the root workspace.
- `OPENSPEC_PROJECT=projects/al bazel run //tools/openspec -- validate
add-bazel-mcp --strict`: passed.
- An isolated real Bazel workspace was exercised over stdio MCP using the built
  AL executable. Query returned the exact `//:payload` label, cquery and aquery
  succeeded, build generated `payload.txt` containing `Bazel MCP fixture`, and
  info returned the exact workspace. Testing the non-test target returned Bazel
  exit 4 as an MCP tool error. Normal EOF exited the server with status zero.
  The fixture used its own output base and ignored host rc files, and its
  task-owned Bazel server was shut down afterward.

The executable E2E suite retains `bazel-mcp.json` in Bazel's undeclared test
outputs. Real protocol observations, generated payload, and command logs remain
under ignored `out/bazel-mcp/`. No host installation, client configuration, or
live infrastructure was changed.

## Failure dispositions

- `mcp-output-frame-limit`: Local review found that large escaped output could
  exceed the SDK client's default frame limit after structured/text duplication.
  The per-stream maximum is now derived from that SDK limit with encoding margin;
  an executable E2E case transfers two 512 KiB streams of escaped control bytes.
- `mcp-disconnect-lifetime`: Initial disconnect verification exposed the need to
  cancel session operations when stdin terminates. The input wrapper propagates
  EOF to the session context. The fixture now disconnects the transport directly,
  since the SDK client's graceful close deliberately drains active requests.
  Request cancellation and transport disconnect both pass.
- `nested-lint-configuration`: The nested module has no `lint` rc configuration.
  The required root-workspace invocation passes against AL's actual targets.
- `sandbox-cache-access`: Bazel and Git LFS use existing shared cache directories
  outside the sandbox. Authorized commands used those existing directories;
  no persistent host configuration or alternate cache was introduced.

## Session ergonomics

Discovery found AL's workspace resolver and the existing MCP SDK before any new
external dependency was considered. Three overly broad inspection outputs were
truncated; subsequent reads used bounded paths and fields. The default command
name was corrected to the user's exact `al bazel mcp serve` before final checks.
Future continuation should start with this change's tasks and the local delivery
receipt rather than rereading the repository-wide documentation.
