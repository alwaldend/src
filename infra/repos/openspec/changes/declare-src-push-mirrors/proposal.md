## Why

The primary Forgejo src repository needs ongoing synchronization to its
existing GitHub and GitLab destinations. Other repositories can use the
same setup.

## What Changes

Declare the desired push destinations in the repository catalog. The
Forgejo Terraform mirror root consumes those entries directly.

## Capabilities

### New Capabilities

- `forgejo-mirror-intent`: Shared source of push mirror configuration.
