## 1. Correct the capacity

- [x] 1.1 Expand the existing integration fixture before implementation; verify the failure at the old status limit.
- [x] 1.2 Increase the status bound to 4 MiB and verify complete inventories and overflow rejection.
- [x] 1.3 Refresh the delivery launcher and select this owner in the aggregate delivery scope and validation plan.

The integration inventory contains 1,201,152 bytes and verifies all 6,144
staged migration paths. It failed at the old ceiling and passes at 4 MiB. The
existing overflow case still rejects truncation without a partial inventory.

- [x] 1.4 Extend large-collection preparation/publication coverage before implementation; verify the 256 KiB receipt failure, introduce a separate 4 MiB preparation-receipt bound, and retain oversized-receipt rejection.

The large-receipt E2E case prepares and verifies publication of all 2,200
paths, then rejects a receipt larger than 4 MiB. It emits
large-collection-verification.txt as a repeatable Bazel test artifact.

- [x] 1.5 Correct the oversized preparation-receipt check to use the actual
      4 MiB preparation reader, rather than the 256 KiB review-receipt reader.
      Rerun the large-collection E2E artifact through the delivery validation plan.

## Immutable path inventory capacity

- [x] Expand migration and full publication E2E fixtures above 1 MiB before implementation.
- [x] Increase immutable changed-path output to 4 MiB while retaining overflow rejection.
- [x] Refresh the launcher, verify full inventories/publication artifacts, and select the owner validation gates.

Expanded fixtures verify an 8,192-path migration above 1 MiB and a 6,000-file
guarded publication. Both reproduced the old immutable-path ceiling; they pass
with complete inventories at 4 MiB, while overflow still fails closed.
