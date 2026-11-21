## Why

The scanner's temporary Git checkout attempts to fetch a Git LFS payload from
the source repository when the host has a Git LFS filter installed.

## What Changes

- Keep Git LFS pointer files in the scanner's temporary checkout.

## Capabilities

### Modified Capabilities

- `repository-scan` gains support for Git histories with LFS pointer files.

## Impact

The scan wrapper sets the environment of its scanner subprocess. Git history
coverage and finding policy continue to apply.
