## Purpose

Provide versioned third-party artifact publication with the same explicit, reproducible ownership and storage guarantees as project releases.

## ADDED Requirements

### Requirement: Third-party publication

The system SHALL support third_party/<name>/releases/<version>/<file> alongside existing project releases, preserve projects as the default namespace, and refuse unsafe or unsupported namespace paths.

#### Scenario: Third-party snapshot publication

- **WHEN** an operator explicitly publishes a named version to the third_party namespace
- **THEN** those bytes are available anonymously at the corresponding URL and existing project releases remain available

### Requirement: Public namespace index

The download root SHALL return a JSON directory listing of projects and third_party without redirecting or exposing staging, state or extracted sites.

#### Scenario: Browse download root

- **WHEN** a client requests GET / on the download host
- **THEN** the response is HTTP 200 JSON containing only the two public namespaces
