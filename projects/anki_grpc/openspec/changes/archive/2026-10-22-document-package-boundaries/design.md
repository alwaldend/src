## Context

The three Go packages separate native hosting, client/collection lifetimes and transport verification. Their comments need explicit caller and application exclusions.

## Goals / Non-Goals

Document verified current boundaries without changing executable behavior or native code.

## Decisions

State managed versus borrowed lifetimes, explicit configuration, application policy exclusions and ownership of native runtime behavior in doc.go. Verify comments against the package implementations and existing native E2E cases. Keep protocol and storage semantics with Rust and Anki.

## Risks / Trade-offs

Documentation must describe actual guarantees; startup and shutdown contexts have distinct meanings, and cancellation does not establish rollback.
