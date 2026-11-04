## Why

Alphabetical project lists hide newly created projects. Visitors should see
newer projects first, using explicit creation dates maintained in landing
front matter.

## What Changes

- Set a creation `date` manually in every registered project's landing page.
- Sort the homepage project links and `/projects/` index by descending date.
- Require this front matter field in the repository Hugo skill.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-landing-pages`: Require explicit creation dates and newest-first
  ordering in the project index.
- `project-alwaldend-com`: Order homepage project links by creation date.

## Impact

The shared site owns ordering. Project-owned landing content supplies the dates;
`tools/agents/skills/repo-hugo` documents the authoring contract and its existing
evaluation covers it. No dependencies, hosting, or deployment change.
