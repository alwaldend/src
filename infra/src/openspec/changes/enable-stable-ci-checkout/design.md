## Decisions

The local action keeps its existing default working directory. When enabled on
the dedicated serial runner, it prepares a separate checkout under the runner
account's home. It fetches from the temporary checkout, reuses its LFS objects,
and clears leftovers before launching Bazel. Git never needs a second network
checkout or job credential. The reusable workflow retains the existing
`master` push trigger and declares `workflow_call` for trusted callers.
