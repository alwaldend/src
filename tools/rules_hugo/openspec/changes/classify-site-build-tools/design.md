## Context

`postcss` already uses the execution configuration. The `tools` list currently
uses the target configuration, so `--notool_deps` still follows Sass.

## Decisions

Declare `cfg = "exec"` for the tool list, retaining the source archive's
target configuration. Use Bazel's dependency classification rather than
inferring tool roles from target names in the graph converter.

## Risks / Trade-offs

- Tool outputs move to the execution configuration; validate actual site output
  and the existing Hugo rule tests before local handoff.
