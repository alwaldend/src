# Phase 1 shared-contract result

Implemented the first Phase 1 dependency slice and repaired the goal workflow
needed to carry it:

- Added the repository-internal `tools/agents/api/v1alpha1` package with
  shared references, atomic effects, authority and budget envelopes,
  independent path-policy and information axes, provenance/completeness/
  retention vocabulary, stable availability reasons, evidence applicability,
  operation contracts, and artifact envelopes.
- Added strict deterministic JSON codecs and validation fixtures for the
  Phase 1 rejection gates.
- Made checked-out attempt-free goals valid despite Git omitting empty
  directories, while recreating `attempts/` safely for first publication.
- Resolved checkpoint plan, result, evidence, review, and criteria inputs
  relative to the store workspace and rejected paths that escape it.

The bounded candidate passed owner-scoped tests, builds, lint, Gazelle review,
Buildifier, diff hygiene, and live goal validation. This attempt should close
with `refine`: the shared foundation is ready, but the registered-universe,
scratch-isolation, safe-operation migration, and resource-baseline criteria
remain for later attempts.
