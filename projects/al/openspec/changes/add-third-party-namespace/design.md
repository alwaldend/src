## Context

See proposal.md. Existing storage and publishers hardcode projects. The requested snapshot is already available locally and contains 1,313,732 bytes.

## Goals / Non-Goals

Publish a versioned X schema without retaining the schema body in Git or requiring X authentication during builds. Preserve website selection and existing project publication.

## Decisions

Add an optional namespace with projects as the backwards-compatible default. Validate a closed set of namespaces, derive bundle names from matching prefixes, and use one owner for generated URLs and SSH paths. Reject website extract/link for third_party.

Schema version: 2.173; content SHA-256: 4377d1791c01632c9cf2860fc833a4cec76eff74f89d232fd7355e765e7bf0a9. Use an immutable schema version with the digest in its version label and pin exact integrity in Bazel.

## Risks / Trade-offs

Independent hosts require explicit publication to both. Verify anonymous HTTP bytes and checksums before removing the source snapshot. No live mutation is authorized beyond the two requested namespace configurations and this schema upload.

## Review corrections

RELEASE-SSH-AGENT-001: an explicit identity must continue to work when its private
key is held by an SSH agent. Make IdentityAgent=none and IdentitiesOnly=yes an
explicit --ssh_disable_agent option, disabled by default. A local SSH endpoint
and in-memory agent exercise public-key-file selection through the built CLI;
the same identity must fail when agent use is explicitly disabled.
