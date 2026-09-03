# Baseline coverage evidence

The normalized routing cases are recorded in
`tools/agents/catalogs/skill-cases.json`; the aggregate CLI emits
`tools/agents/catalogs/skill-coverage.json` bound to capability-catalog
digest `sha256:0775209e0054d7c3d68e905699c3a07b6b36b4f2a2a0e589c12ef1782202d7c5`.

The matrix contains three fixture-tested representative cases against 24
registered skills and reports `truncated: true`. It covers positive routing,
an inert payload, and a conflict signal, but adjacent negative, exclusion,
and composition cases remain explicit gaps, as do live and high-risk
writable coverage.

Remaining criteria are open by explicit status, not by inferred success.
