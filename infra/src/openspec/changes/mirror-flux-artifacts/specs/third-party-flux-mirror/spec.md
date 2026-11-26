## Purpose

Mirror pinned Flux and Flux Operator artifacts into Forgejo using the third-party publication workflow.

## ADDED Requirements

### Requirement: Flux artifact mirroring

The Flux and Flux Operator packages SHALL each provide Bazel targets to pull
their pinned charts and images and publish them to Forgejo using the
third-party AppRole.

#### Scenario: Publish an existing local mirror

- **WHEN** an operator runs the push target with that package’s staged artifacts
- **THEN** chart archives are published to the Helm registry
- **AND** image manifests retain their upstream digests and platform entries
