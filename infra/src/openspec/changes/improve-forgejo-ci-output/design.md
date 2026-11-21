## Context

The Go CI command discovers workspaces and runs their phases serially. The Node.js
action inherits the Forgejo environment and relays the command's exit status.

## Decisions

- Emit the Forgejo group markers and phase diagnostics on the same output stream
  to preserve ordering in the job log.
- Use the runner's step summary file to append a Markdown table. Keep the phase
  exit status authoritative if writing the summary fails.
- Fold routine Bazel console records with their continuation lines. Keep error
  records and test result lines directly inside each phase. Keep each failed
  test's output in its own group.
- Fold the launcher diagnostics inside its build group and report its failure
  status when that build group is expanded.
- Use the pinned LLVM sysroot repository rule with the existing WASI SDK
  archive, so LLVM can consume the directory directly during module map generation.

## Risks and Trade-offs

The job summary is written after the command has collected the phase results.
The phase logs remain available while the command is running.
