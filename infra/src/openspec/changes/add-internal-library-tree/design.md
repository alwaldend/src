## Context

The user selected internal/x_sdk for the new SDK and rejected committing generated client source. Existing tool rules permit execution tools; their source-update documentation needs to distinguish build generation from checked-in source updates.

## Goals / Non-Goals

Define a shared runtime-library boundary without introducing a new language-specific tree or migrating unrelated projects.

## Decisions

Proceed with internal/ as a repository-only library owner. projects/ remains valid for independently maintained products and reusable project code. third_party/ retains upstream dependency definitions; it does not own this first-party integration. A root Go internal directory also restricts imports to the repository module tree; repository-internal Bazel visibility enforces the intended broader boundary.

Generated clients are build artifacts. Execute the generator in the execution configuration and compile its output in the consuming component; do not treat the generator executable as a runtime dependency.

## Risks / Trade-offs

A new tree needs documentation and transitive-visibility integration. Validate the generated library, its existing uploader consumer, OpenSpec workspaces and the documentation build.

The internal README now states visibility, publication and production-consumer
requirements in the same format as the other top-level trees. SDK reference
documentation describes the pinned mirror and delegates maintenance to its
dependency owner.
