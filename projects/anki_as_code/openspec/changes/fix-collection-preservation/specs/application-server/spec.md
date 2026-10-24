## MODIFIED Requirements

### Requirement: Application-owned server operations

App SHALL expose ServerRun using an explicitly injected server runner. CLI SHALL
prepend serve for the server serve command before calling ServerRun; App SHALL
NOT expose a separate ServerServe forwarding method. App SHALL reject missing
runner configuration and cancelled startup before invoking the runner, preserve
runner error causes, and leave caller arguments unchanged. The executable SHALL
supply native hosting and process signal adaptation without assembling Cobra
server commands. CLI handlers SHALL forward command context, display name and
adapted native arguments without running the native server or invoking its runner.

#### Scenario: Serve through the application

- **WHEN** CLI receives server serve with native transport arguments
- **THEN** it passes serve followed by those arguments to App.ServerRun
- **AND** App invokes the configured runner with the display name and that argument vector

#### Scenario: Rejected startup

- **WHEN** startup is cancelled or no runner is configured
- **THEN** App returns an error before native execution
