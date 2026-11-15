## Why

Development deployment fails before configuration on dev1: its Python 3.15 standard library removed dataclasses._is_type, which pinned ansible-core 2.20.7 still accesses. Ansible core 2.22 supports Python 3.15 targets.

## What Changes

- Update the existing Ansible dependencies to core 2.22.0 and its matching community collection bundle.
- Regenerate the owning uv lock without introducing new external packages.
- Verify packaged syntax, semantic lint and live fact gathering on both development hosts before continuing deployment.

## Capabilities

### New Capabilities

- `target-python-compatibility`: Supported deployment to Python 3.15 targets.

## Impact

tools/ansible CLI dependencies and their infrastructure consumers. This is necessary for the user's authorized infra/dev deployment. System Python and host roles remain unchanged.
