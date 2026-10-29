## Context

The data package owns S8. The site already has a shared partial that reads
`icons/favicon.svg` and fingerprints it. See the proposal for the missing
connection between these owners.

## Goals / Non-Goals

Use one canonical SVG and preserve the existing browser-facing resource
contract. Deployment and a WebP build rule are outside this change.

## Decisions

Map the data package's SVG into `assets/icons/favicon.svg` using the existing
Bazel packaging rules, and remove the obsolete site-local SVG. This keeps
the template and Hugo mounts intact without duplicating artwork or using a
filesystem symlink outside the declared build inputs.

## Risks / Trade-offs

- Incorrect archive mapping could leave the old icon or a missing resource.
  Inspect the packaged source and the favicon referenced by built pages.
