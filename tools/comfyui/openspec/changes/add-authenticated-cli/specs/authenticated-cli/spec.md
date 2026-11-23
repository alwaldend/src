## Purpose

Provide repeatable terminal access to an authenticated self-hosted ComfyUI without task-specific API clients.

## ADDED Requirements

### Requirement: Packaged upstream CLI

The repository SHALL provide the approved official ComfyUI CLI through Bazel
with locked dependencies and without installing host packages.

#### Scenario: Run a workflow

- **WHEN** a user supplies a compatible workflow and reachable ComfyUI
- **THEN** the packaged CLI SHALL provide upstream workflow and job commands

### Requirement: Authenticated HTTPS execution

The bridge SHALL accept an HTTPS endpoint, trusted CA and injected login, forward
requests only to that endpoint, and keep credentials out of child arguments,
environment and saved CLI configuration. It SHALL verify TLS and reject requests
with foreign browser origins, other local users and alternative CLI routing. It SHALL close its
listener when the CLI exits and propagate command failure.

#### Scenario: Valid login and trust

- **WHEN** the supplied CA and login match the endpoint
- **THEN** upstream commands SHALL reach the authenticated API through a temporary local connection

#### Scenario: Invalid trust or missing credentials

- **WHEN** authentication or certificate verification cannot be established
- **THEN** execution SHALL fail without disabling verification or printing credentials

### Requirement: Render and retrieve in one invocation

The authenticated CLI SHALL provide a command that runs a supplied workflow and
downloads its generated files through the same bridge. For native workflow
JSON it SHALL accept input overrides and local image references without
modifying the supplied workflow.

#### Scenario: Render a workflow with a reference image

- **WHEN** the user supplies a native workflow, an image slot and a local image
- **THEN** the command SHALL upload the image, set the slot for this run, wait
  for the job, and report the prompt ID and downloaded local file paths
