## Context

See proposal.md. The existing AL-specific layer only composes the generic AL
wrapper with the Terraform rule and expands named operation maps.

## Decisions

Decision review verdict: proceed with optional wrapper composition. Keeping
the maps in AL would retain the ownership defect; importing AL into the
standalone module would violate its reusable contract. A caller-supplied macro
keeps each implementation with its owner, at the cost of explicit wrapper
selection in consumer BUILD files.

The generic maps pass an inner executable, runtime data, and common target
attributes to the wrapper. Wrapper-specific configuration stays in a separate
dictionary, with reserved and duplicate arguments rejected. Wrapped tests
retain testonly on both targets and size on the test wrapper.

Shared provider pins and their real-provider fixture move to third_party/terraform.
The repository root command declarations move beside tools/al's shared config;
they contain no Terraform rule implementation. All provider versions remain
unchanged and extension resolution still selects one version per source.

## Risks / Trade-offs

- Wrapper composition could drop plugin configuration or test attributes;
  compare every migrated expression and test wrapped/unwrapped target contracts.
- Moving dependency includes can break module loading; retain the owning
  generator launcher and regenerate MODULE.bazel before analysis.
- Validation remains offline; no new live infrastructure operations are needed.

## Validation

The focused composition and provider run passed all 12 tests. The corrected
consumer/quality run passed all 86 tests, including the offline real-provider
fixture and command-map regressions. Standalone tests passed all 16 tests and
the standalone build passed. Semantic lint covered all migrated packages and
the new dependency package, and the documentation site built successfully.

Two formatting findings (a new dependency BUILD file and test-provider docs)
and a stale fixture label were repaired before the passing consumer run.
Independent review compared all 100 migrated invocations across 46 consumers
without an argument mismatch, and identified a wrapper testonly collision that
was corrected before the final checks. Provider declarations are byte-identical
to their prior location (SHA256 b1fbc0a302b4316162197e7c4864172fecdb1fa88a9ef1c20f7939ee8127e5df).

Rendered rule-module and dependency documentation and generated GitLab/root
invocations were inspected. Git and delivery receipts own final candidate and
publication identity; all verification in this correction is offline.
