## Purpose

Look up canonical Danbooru tag names from the pinned vocabulary while preparing
image generation prompts in ComfyUI.

## ADDED Requirements

### Requirement: Local tag lookup

The Danbooru tag search workflow and App SHALL search the pinned CSV locally
and present a selectable text result with at most the requested number of tags.

#### Scenario: Enter a tag or part of a tag

- **WHEN** a user runs a query in the Danbooru Tag Search App
- **THEN** it returns the matching canonical tag names in prompt-ready text
- **AND** exact matches precede prefix matches and other substring matches
- **AND** matches in each group are ordered by post count from highest to lowest

#### Scenario: Filter deprecated tags

- **WHEN** a query matches deprecated tags
- **THEN** they appear only if the user enables the deprecated-tag option

#### Scenario: Find no matches

- **WHEN** a nonempty query matches no eligible tags
- **THEN** the result states that there are no matching tags
