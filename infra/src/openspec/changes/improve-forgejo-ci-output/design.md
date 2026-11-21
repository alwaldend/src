## Context

The Go CI command discovers workspaces and runs their phases serially. The Node.js
action inherits the Forgejo environment and relays the command's exit status.

## Decisions

- Emit the Forgejo group markers and phase diagnostics on the same output stream
  to preserve ordering in the job log.
- Use the runner's step summary file to append a Markdown table. Keep the phase
  exit status authoritative if writing the summary fails.

## Risks and Trade-offs

The job summary is written after the command has collected the phase results.
The phase logs remain available while the command is running.
