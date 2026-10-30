## Purpose

Maintain the third-party artifact storage root with the same ownership, read-only serving labels and deduplication policy as existing publication roots.

## ADDED Requirements

### Requirement: Third-party content storage

The download-host role SHALL create third_party beside projects, own it with the publication account, label its content for read-only HTTP serving and include it in bounded content deduplication. The publication tag MUST verify storage is mounted before changing publication roots.

#### Scenario: Configure an existing mounted download host

- **WHEN** the authorized publication-only role tasks run
- **THEN** third_party has the same publication ownership and serving labels as projects, and maintenance covers both namespaces without provisioning storage
