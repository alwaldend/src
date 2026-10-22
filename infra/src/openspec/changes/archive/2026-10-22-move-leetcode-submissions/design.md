## Context

The submission package contains 1,930 tracked files. The user explicitly requires
that generated submission pages remain on the main site after relocation.
This task-specific instruction takes precedence over the user-tree publication
restriction in `users/README.md`.

## Goals / Non-Goals

Preserve submission values, downloader generation, and generated site pages while
relocating ownership. Fix the delivery capacity blocker and create a PR.
Submission formats and live site deployment are outside this change.

## Decisions

Keep the internal layout at the requested path. Allow the site's existing
submission package to consume the moved target and update its dependency label.
Record shared layout coordination here; each component maintains its own linked change.

Preparation previously committed the candidate but could not serialize its
receipt because exact path inventories exceed 256 KiB. Verdict: proceed with a
separate 4 MiB preparation receipt limit, retaining bounded stable reads, schema
validation, exact scope, leases, and the 256 KiB cap for other typed records.
Compact JSON alone does not solve the repeated long-path inventories; removing
exact paths would weaken verification. A narrow capacity change preserves the
existing guards and is reversible. Test a real large Git change through prepare,
publish, and verify, plus rejection of a preparation receipt beyond the new cap.

## Risks / Trade-offs

Preparation receipts permit a larger bounded allocation; no authentication or
scope check is removed. Generated submission TOML files retain their original Git bytes and are excluded
from generic formatting through .gitattributes, as required by the user's correction. The downloader owns their serialized layout. Verify all moved
submission files as exact renames and keep site generation enabled.
