## ADDED Requirements

### Requirement: Graceful command cancellation

The launcher SHALL propagate SIGINT and SIGTERM to its active Terraform child,
wait for that child, and retain a bounded forced-termination fallback. Callers
SHALL be able to interrupt the launcher without leaving Terraform holding its
output pipe and preventing caller cleanup.

#### Scenario: Interrupt a nested invocation

- **WHEN** the caller is cancelled while Terraform is running
- **THEN** the active Terraform child receives a graceful interrupt and exits
- **AND** the launcher waits for it before returning control for caller cleanup
