# Reimu Fumo attempt 0

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 0 — inherited session baseline

**Candidate:** `out/reimu_fumo_sculpt_v31_working.blend`, SHA-256
`9d0ac214676d0ca61e971c26998da4698b1d7645f506cbdb05b946e9dba41839`,
untextured sculpt stage, review packet `v31-fixed-five-view`.

**Failure targeted:** Earlier work tried to replace the original poor Fumo with
a recognizable standalone plush.

**Hypothesis:** Repeated scripted adjustments to ellipsoids, flat panels, and
curves would converge on the references.

**Plan used:** Incrementally adjust head ratios, hair caps, bow loops and
tails, skirt, sleeves, feet, materials, and pose, rendering selected views
between changes.

**Work performed:** More than thirty local scripted variants were produced.
Later variants added fixed cameras, reference collections, broader sleeves,
lower skirt contact, thinner front hair, and a broad/shallow head experiment.

**Evidence:** The latest complete reviewed candidate was v31. Its blind review
scored overall quality 3.0/10, silhouette 3.5/10, sewn construction 3.0/10,
identity 3.0/10, contact 4.5/10, plush read 3.5/10, and presentation 6.0/10.
The reviewer could not identify it as the same variant without a label.

**Criterion results:** Reference likeness, measured silhouette, plush
construction, presentation quality, reusable structure, animation readiness,
and technical integrity all fail or remain unverified. Repository delivery is
temporarily satisfied only for the rejected baseline.

**Decision:** Reset the modeling strategy. Do not continue polishing v31.

### Progress and approach audit after attempt 0

- **Improved:** camera consistency, basic floor contact, and review discipline.
- **Regressed or stalled:** absolute likeness, head/hair construction, compact
  body proportions, bow softness, and plush material read.
- **Absolute result:** still clearly poor; relative improvement did not make it
  acceptable.
- **Invalid approach:** disconnected procedural primitives were refined before
  the macro silhouette and sewn-panel construction were correct. Side and rear
  failures persisted while front-view details consumed most iterations.
- **Highest-leverage problem:** rebuild the head/hair/bow and compact body as
  coherent soft fabric volumes using measured multi-view landmarks.
- **Next approach:** discard the current geometry as a modeling base. Use it
  only as documented negative evidence, establish measured blockout gates, and
  add detail only after the untextured form passes absolutely.
