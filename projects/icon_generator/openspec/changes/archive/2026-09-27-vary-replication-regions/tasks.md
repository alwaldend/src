## 1. Specify and implement regional variation

- [x] 1.1 Write E2E cases before implementation for validation, compatibility, average density, spatial contrast, parent retention, paired births, replay, and small/partial grids; observe missing-option failures.
- [x] 1.2 Implement the shared smooth field, normalized seeding, and local reproduction probabilities; pass the E2E suite and compare variant 10 PNG bytes with all regional controls off.
- [x] 1.3 Document option units, density normalization, zero-reproduction seeding, and probability clamping; replay the new E2E example from its manifest.

## 2. Compare and verify output

- [x] 2.1 Generate separate 1024-square white-pixel comparison PNGs using variant 10's controls; inspect representative outputs and retain arguments, checksums, and final coverage.
- [x] 2.2 Measure maximum-grid execution and verify strict OpenSpec and project/site integration; record evidence in evidence.md.

Archive the completed change, format and inspect the aggregate candidate, and
use repo-delivery for exact-candidate gates, publication, and review verification.
