## 1. Package graph construction

- [x] 1.1 Record behavioral failure cases before implementation in design.md.
- [x] 1.2 Write E2E coverage first, then add `analyze target`, root/analyze help,
      and depth-controlled `analyze package` aggregation and rule wiring; verify
      JSON/Mermaid equivalence, membership, counts, directions, filters, depth,
      ordering, escaping, and safe failures.
- [x] 1.3 Document command usage and graph limitations, verifying that examples
      invoke the analyzer directly without task-local graph construction.

## 2. Validation and delivery

- [x] 2.1 Generate the overview from the recorded root capture using the new
      command, compare all aggregate relationships with the captured data, and
      render and inspect its dark WebP through the existing Mermaid renderer.
      Capture generation, semantic comparison, and visual inspection passed,
      as recorded in verification.md.
- [x] 2.2 Run affected E2E and OpenSpec checks, record the candidate and results,
      and hand maintained changes to the parent for integration review and delivery.

## 3. Explicit deps/srcs filtering

- [x] 3.1 Record failure cases and write command/native-query E2E coverage before implementation.
- [x] 3.2 Implement graph-owned attribute filtering and metadata validation, both CLI flags, rule options, and output metadata while preserving defaults.
- [x] 3.3 Document the structural scope, complete capture contract, producer exception, and collapsed-input tool-provenance limitation.
- [x] 3.4 Run pinned formatting, E2E and OpenSpec validation, build the standalone analyzer, retain repeatable artifacts, and hand stable sources to the parent.
- [x] 3.5 Reproduce native non-label selector None defaults before correcting their handling, then verify both modes and retain standalone/render artifacts for parent delivery.
