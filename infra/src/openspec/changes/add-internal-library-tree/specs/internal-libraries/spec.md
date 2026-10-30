## Purpose

Provide an explicit repository-only boundary for first-party runtime libraries shared by multiple repository components.

## ADDED Requirements

### Requirement: Shared internal libraries

Shared first-party repository-only runtime libraries SHALL be owned under internal/<name>, use repository-internal target visibility and keep generated client source in Bazel outputs. They SHALL NOT require separate product landing pages or releases.

#### Scenario: Application uses a shared library

- **WHEN** a repository application depends on an internal library
- **THEN** Bazel can generate and compile that library using declared execution tools without committing generated client source
