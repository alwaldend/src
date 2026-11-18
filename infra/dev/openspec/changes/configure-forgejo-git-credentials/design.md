## Context

The installed fj 0.7.0 does not implement Git's credential helper protocol.
The operator explicitly authorized reuse of the native fj login.

## Decisions

Use the pinned client's host record for the verified Forgejo HTTPS repository.
The helper only handles Git credential retrieval. Let fj own refresh and
storage, and require its credential file to be private. Other hosts and
repositories receive no credential. Keep all credential material out of the
playbook and logs.

## Risks / Trade-offs

The helper depends on the installed fj release's credential file format.
Validate against an isolated fixture. Confirm compatibility when upgrading fj.
