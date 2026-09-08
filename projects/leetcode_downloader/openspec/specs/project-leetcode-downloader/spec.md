# LeetCode Downloader Specification

## Purpose

Export LeetCode submissions and generate source and documentation from stored
submission data. This baseline records checked-in behavior at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on 2026-09-08. The project
README reports direct CLI downloads blocked by bot protection; this baseline
does not verify current external service access.

Sources: [project README](../../../README.md),
[submission generator](../../../main/go/generator.go),
[documentation rule](../../../main/bzl/al_leetcode_submissions.bzl),
and [browser userscript](../../../main/tampermonkey/leetcode_downloader_tampermonkey.js).

## Requirements

### Requirement: Accepted submission export

The generator SHALL skip submissions whose status is not `Accepted`, select
output directories and extensions from the configured language mapping, and
write each accepted submission's TOML metadata beneath
`questions/<question-id>/submissions/<submission-id>` within its configured
language directory.

#### Scenario: Generate an accepted submission

- **WHEN** an accepted submission has a matching language configuration
- **THEN** the generator SHALL write its TOML metadata and SHALL additionally
  write its code file when `WriteCode` is enabled.

### Requirement: Markdown submission documentation

The Bazel submission rule SHALL render Markdown from submission data with
timestamp metadata, a LeetCode submission link, and a fenced code block labeled
with the submission language.

#### Scenario: Render stored submission documentation

- **WHEN** a stored submission is processed by `al_leetcode_submissions`
- **THEN** its generated Markdown SHALL contain its submission URL and source
  code with the configured language label.

### Requirement: Browser download control

The userscript SHALL add a download button to matching LeetCode problem-set
pages, request submission pages in batches of 20, and save accumulated results
as `submissions.json` when pagination ends or the user stops downloading.

#### Scenario: Complete a successful browser response

- **WHEN** a successful submissions response has no next page
- **THEN** the userscript SHALL append that response's submissions and initiate
  a download of the accumulated `submissions.json` data.
