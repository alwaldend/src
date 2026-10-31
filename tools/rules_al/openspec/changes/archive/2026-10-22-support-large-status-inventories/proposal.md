## Why

The requested Anki export adds more than 7,700 text files with hierarchical
paths. Its complete porcelain-v2 status occupies 1,879,106 bytes, exceeding the
delivery adapter's 1 MiB bound and preventing normal preparation.

## What Changes

Increase the operation-specific status, immutable changed-path, and preparation-receipt ceilings
to 4 MiB. The receipt records the exact aggregate paths and exceeded its
previous 256 KiB ceiling. Preserve bounded
command execution and rejection of truncated status. Expand the existing Git
integration fixture beyond 1 MiB and verify every staged path. Exercise full
preparation, publication, and verification of a large receipt.

## Capabilities

No specification delta; this corrects the supported inventory capacity.

## Impact

tools/repo_delivery status inspection, preparation-receipt encoding/reading,
and integration tests.
The generic diagnostic, review-receipt and index-inspection
limits are unchanged.

The preparation-receipt capacity has since landed on master in its owning
change. This PR retains the status and immutable changed-path inventory capacity adjustments and
large-collection integration coverage needed to publish the Anki export.

The Anki marker/suffix migration also exceeds the 1 MiB immutable changed-path
ceiling. Extend that operation to the same bounded 4 MiB capacity; retain
complete inventory checks and explicit overflow rejection.
