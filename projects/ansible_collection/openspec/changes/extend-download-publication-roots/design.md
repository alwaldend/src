## Context

See proposal.md. Existing storage and publishers hardcode projects. The requested snapshot is already available locally and contains 1,313,732 bytes.

## Goals / Non-Goals

Publish a versioned X schema without retaining the schema body in Git or requiring X authentication during builds. Preserve website selection and existing project publication.

## Decisions

Extend the existing directory, SELinux and deduplication definitions to include third_party. Tag the publication-related tasks so the requested namespace can be deployed without unrelated host provisioning.

Schema version: 2.173; content SHA-256: 4377d1791c01632c9cf2860fc833a4cec76eff74f89d232fd7355e765e7bf0a9. Use an immutable schema version with the digest in its version label and pin exact integrity in Bazel.

## Risks / Trade-offs

Independent hosts require explicit publication to both. Verify anonymous HTTP bytes and checksums before removing the source snapshot. No live mutation is authorized beyond the two requested namespace configurations and this schema upload.
