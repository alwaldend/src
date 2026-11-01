## Purpose

Provide versioned third-party artifact publication with the same explicit, reproducible ownership and storage guarantees as project releases.

## ADDED Requirements

### Requirement: Third-party publication

The system SHALL support third_party/<name>/releases/<version>/<file> alongside existing project releases, preserve projects as the default namespace, and refuse unsafe or unsupported namespace paths.

#### Scenario: Third-party snapshot publication

- **WHEN** an operator explicitly publishes a named version to the third_party namespace
- **THEN** those bytes are available anonymously at the corresponding URL and existing project releases remain available
