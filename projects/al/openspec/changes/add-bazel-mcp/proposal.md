## Why

Agents need discoverable Bazel operations with explicit inputs, cancellation,
and bounded results instead of ad hoc shell commands. AL already owns Bazel
workspace selection and is the appropriate host for this interface.

## What Changes

- Add `al bazel mcp serve`, a stdio MCP server bound to one Bazel workspace.
- Expose query, configured query, action query, build, test, and info tools.
- Return bounded stdout/stderr, exit status, and timeout/truncation indicators.
- Reuse the repository-pinned official MCP Go SDK and normal Bazel executable.
- Verify the executable through an MCP client and an isolated Bazel workspace.

## Capabilities

### New Capabilities

- `bazel-mcp`: Discoverable Bazel operations over MCP with bounded process lifetimes.

### Modified Capabilities

None.

## Impact

AL's Bazel command package, command E2E tests, documentation, and standalone
module dependency declarations. Existing `al bazel test` behavior is preserved.
No host installation, persistent client configuration, or deployment is included.
