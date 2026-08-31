# Bazel scratch-isolation result

Removed the mandatory runner's shared `out/tmp` behavior and the agent Bazel
profile's propagation of ambient temporary paths into repository rules,
actions, host actions, and tests. Bazel-managed temporary contracts now remain
authoritative, and a durable Bazelrc regression prevents reintroduction.

This attempt closes with `refine`. It removes one concrete cross-task collision
and cache-portability defect, but Phase 1 scratch isolation is not complete
until task/run manifests exist and Cordis scratch is namespaced by task and
run.
