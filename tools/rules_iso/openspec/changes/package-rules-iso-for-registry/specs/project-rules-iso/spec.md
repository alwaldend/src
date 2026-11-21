## ADDED Requirements

### Requirement: Consume a released ISO ruleset from a registry

The ruleset SHALL provide a source archive and a matching module declaration
for Bzlmod registry consumers.

#### Scenario: A fresh workspace uses the release

- **WHEN** a consumer declares a dependency on the registered archive version
- **THEN** it can load `@rules_iso//pkg/bzl:iso_extension.bzl` and build the
  resulting image and flash targets outside the monorepo
