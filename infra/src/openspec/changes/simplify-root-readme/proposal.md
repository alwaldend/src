## Why

The root README repeats agent policy and detailed build documentation, obscuring
its purpose as the repository entry point.

## What Changes

- Keep the directory map and setup links; put external links and license first.
- Remove the introductory paragraph and redundant Markdown title; retain the
  front-matter title.
- Replace agent workflow and build details with links to their owners.

## Capabilities

No new or modified capabilities. This is a documentation-only simplification;
`skip_specs: true` applies.

## Impact

Root README only, with its repository-owned change record. Acceptance requires
reviewing the resulting document and checking local link destinations.
