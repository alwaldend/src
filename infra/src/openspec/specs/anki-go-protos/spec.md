# Anki Go Protobuf Integration Specification

## Purpose

Make all protobuf schemas from the pinned upstream Anki source available to Go applications through reproducible repository builds.

Sources: [dependency ownership](../../../../../third_party/com_github_ankitects_anki/README.md),
[schema catalog](../../../../../third_party/com_github_ankitects_anki/schemas.bzl),
[binding generation](../../../../../third_party/com_github_ankitects_anki/anki.BUILD.bazel),
and [consumer aliases](../../../../../third_party/com_github_ankitects_anki/BUILD.bazel).

## Requirements

### Requirement: Complete pinned Go schemas

The dependency integration SHALL generate Go messages and declared gRPC clients for every upstream Anki proto, retaining upstream import paths and the existing immutable source pin.

#### Scenario: Consume a previously omitted service

- **WHEN** an application imports the pinned tags or i18n schema package
- **THEN** its generated messages, descriptors and declared clients compile and can be used against the native backend

### Requirement: Restricted schema consumption

Schema targets SHALL retain repository-internal consumer visibility.

#### Scenario: Intended application consumes all schemas

- **WHEN** anki_as_code links all pinned schema packages
- **THEN** its build succeeds through the dependency package's declared aliases
