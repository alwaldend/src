## Purpose

Provide reproducible pinned openapi generation through explicit Bazel targets with declared inputs and outputs.

## ADDED Requirements

### Requirement: Pinned OpenAPI generation

The tool SHALL expose a pinned upstream CLI and generate client source offline from declared schema and configuration inputs.

#### Scenario: Successful generation

- **WHEN** a developer generates from the declared schema input
- **THEN** the pinned generator produces a compilable Go client without network access or credentials
