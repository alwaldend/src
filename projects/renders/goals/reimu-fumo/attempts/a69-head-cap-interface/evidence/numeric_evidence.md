# A69 C1 numeric evidence

Pinned Blender 5.2.1 clean reopen passed diagnostic verification for
`reimu_fumo_a69_c1_live_v4.blend`.

- Candidate SHA-256:
  `12f2a09d376d754c838fc2dddb9c528e59709e60119882bbf5563572903bcfe9`
- Intended curve-seam displacement: 4.000008 mm maximum, with no material X
  or Z motion.
- Cap/head separation: 0.512 mm minimum, 0.789 mm median, 1.098 mm maximum;
  zero samples below 0.25 mm.
- Rear-lock within-1-mm contact retention:
  - left asymmetric: 82.58% (441 of 534 baseline samples)
  - off-center main: 72.44% (326 of 450 baseline samples)
  - short right: 78.12% (332 of 425 baseline samples)
- Diagnostic contact gate: pass.
- Promotion contact gate: fail.
- Frozen-context fingerprint: exact after harmless cross-version metadata
  normalization.

The machine-readable source is
`out/reimu_fumo_attempt_069_head_cap_interface/c1_4mm_analytic_shape_keys/reopen_report.json`.
