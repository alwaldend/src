## Context

See proposal.md. Existing storage and publishers hardcode projects. The requested snapshot is already available locally and contains 1,313,732 bytes.

## Goals / Non-Goals

Publish a versioned X schema without retaining the schema body in Git or requiring X authentication during builds. Preserve website selection and existing project publication.

## Decisions

Use the existing authenticated SSH and Ansible flows. Apply the publication namespace to both local and Yandex as explicitly requested; no infrastructure provisioning or DNS changes. Public content remains read-only; staging and state remain private.

Schema version: 2.173; content SHA-256: 4377d1791c01632c9cf2860fc833a4cec76eff74f89d232fd7355e765e7bf0a9. Use an immutable schema version with the digest in its version label and pin exact integrity in Bazel.

## Risks / Trade-offs

Independent hosts require explicit publication to both. Verify anonymous HTTP bytes and checksums before removing the source snapshot. No live mutation is authorized beyond the two requested namespace configurations and this schema upload.

## Authentication findings

The first namespace-only Ansible attempt failed host certificate verification before any host changes. Reuse the existing SSH injector key/certificate/known-hosts mapping for Ansible instead of weakening verification. The signer currently overwrites its generated public-key path with the certificate; preserve the public key and use the standard -cert.pub filename. Failure boundaries: mismatched key/certificate, missing host CA and non-mounted storage must stop before publication changes.

Publication wrappers explicitly select --ssh_disable_agent for their injected
temporary identity. Generic release commands retain the operator's agent policy.
This correction changes wrapper packaging only and requires no host deployment.
