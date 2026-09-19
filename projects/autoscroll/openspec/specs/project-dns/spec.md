# project-dns Specification

## Purpose

Record the retirement of this project's landing DNS infrastructure while
preserving its repository documentation and builds.

## Requirements

### Requirement: Keep landing infrastructure retired

The project SHALL have no dedicated landing DNS declaration, Terraform root, or
operational source export. Its landing page SHALL be published by the main site
under `/projects/autoscroll/` instead of a dedicated hostname.

#### Scenario: Inspect the project after landing retirement

- **WHEN** the project tree is consumed
- **THEN** it contains no landing DNS declaration, Terraform DNS stage, or
  landing build target, and the main site owns its landing page.
