## ADDED Requirements

### Requirement: Reject unsupported special files without blocking lifecycle work

Repository read handlers SHALL reject unsupported file types without waiting
for an external reader or writer to appear. They SHALL finish their handler
and release acquired descriptors and runtime leases on rejection, preserving
existing workspace containment and supported regular-file behavior.

#### Scenario: Read a FIFO that has no writer

- **WHEN** a repository read selects a FIFO inside the workspace with no writer
- **THEN** it returns an unsupported-file error within the request's configured
  timeout without requiring a writer to unblock it
- **AND** the handler releases all acquired resources

#### Scenario: Stop or reload after rejected special-file access

- **WHEN** a special-file request has been rejected and no other work is active
- **THEN** an ordinary file read succeeds
- **AND** package stop or replacement completes within the lifecycle test deadline

#### Scenario: A target changes between path resolution and open

- **WHEN** the selected target is replaced by an unsupported type or escapes
  the workspace before the read is bound
- **THEN** the request rejects that target without consuming its content
- **AND** containment and descriptor cleanup remain effective
