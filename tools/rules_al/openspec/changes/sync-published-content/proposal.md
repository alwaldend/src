> Publication and retention behavior below records the incoming implementation.
> The approved `unify-release-versioning` change supersedes its duplicate-version,
> third-party replacement, and automatic cleanup policies with immutable retries
> in every namespace, retained sites, and explicit deletion. Infrastructure
> synchronization remains owned by `infra/download`.

## Why

Local publication must reject duplicated versions and retain downloadable
archives while removing older extracted websites after successful selection.
The user moved infrastructure synchronization to infra/download's Ansible
workflow, so the release tool must not retain a second sync implementation.

## What Changes

- Stage complete uploads and reject duplicate artifact or extracted versions.
- Select a complete website before removing obsolete extracted versions.
- Preserve archives for rollback by extraction and selection.
- Remove the superseded Go sync command and move integration acceptance to infra/download.

## Capabilities

### Modified Capabilities

- published-content-sync: retained release publication and extracted-site lifecycle; infrastructure mirroring moves to the download owner.

## Impact

Release CLI and SSH transport, the download integration fixture, and their
owner-local documentation and validation. No new dependency or credential model.
