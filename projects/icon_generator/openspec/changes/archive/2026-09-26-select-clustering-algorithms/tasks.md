## 1. Specify and exercise algorithm behavior

- [x] 1.1 Write E2E cases before implementation for selector validation, compact default, count preservation, replay, exact paired births, parent retention, descendants, and density endpoints; observe missing-selector failures in `algorithms-red.log`.

## 2. Implement and document selection

- [x] 2.1 Add compact, strands, and replication with compact as default; pass the project E2E suite including the retained disk-write failure case.
- [x] 2.2 Document algorithm-specific density and factor semantics, eight-neighbor adjacency, bounded reproduction, and memory limits; verify help and example arguments agree with the README.

## 3. Verify outputs and integration

- [x] 3.1 Compare prior compact and strand PNGs with explicit algorithm output, replay the E2E example manifest, inspect replication previews, and measure maximum-grid replication; record observations in evidence.md.
- [x] 3.2 Validate the strict OpenSpec change and project/site integration, then format and inspect the diff; record results in evidence.md.

After implementation and verification, archive through OpenSpec and deliver through
repo-delivery. Delivery receipts own exact-candidate checks, commit, push, and PR state.
