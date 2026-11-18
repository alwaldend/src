## Purpose

Define the settings and narrow-screen behavior for reviewing pull requests on
the managed Forgejo instance.

## ADDED Requirements

### Requirement: Managed branch cleanup and diff page size

Terraform SHALL configure the default deletion of eligible pull request source
branches after web merges for its managed Forgejo repositories. Ansible SHALL
configure the instance to render up to 100 files per diff page.

#### Scenario: Review a large pull request and merge it in the web UI

- **WHEN** a managed repository contains a pull request with more than 100 changed files
- **THEN** Files Changed shows up to 100 files per page and offers more files
- **AND** the merge form preselects deletion if the source branch is eligible

### Requirement: Narrow-screen diff layout

The instance SHALL keep the Files Changed view within the mobile viewport,
with reachable navigation, review controls, and file actions. Wide diff
content SHALL be scrollable within its own file panel as needed.

#### Scenario: Review a diff containing a long filename or line

- **WHEN** a reviewer opens Files Changed on a narrow viewport
- **THEN** the file heading and controls fit or wrap within the available width
- **AND** the diff does not cause page-wide horizontal scrolling
