# A79 guarded candidate scaffold: independent static audit

Verdict: **REVISE before any future Blender build. Safe as the current
no-run scaffold.**

Audited after two matching hash reads:

- `build_candidate.py`:
  `6497e4922102cbb74e0b4e4c140160cd417e1e96b971d2d2999c0540b52d19e5`
- `validate_reopen.py`:
  `fcdfb017ecc57b9382f0b69810586a2b0df87939b97192b6d3573405852729e3`

No Blender process was run. There is no `authorized_inputs.json`, candidate
`.blend`, or build/failure report in `candidate/`. Variant A is a formal
`RESET`: its report has `build_authorized_by_pure_preflight: false` at line 2
and `verdict: REVISE` at line 879. The builder therefore fails closed before
`bpy` import. Protected hashes remain exact:

- rung 003: `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- tracked model: `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

## Findings

1. **Blocking for reuse: pre-save mechanical gates are not integrated.**
   `build_candidate.py` performs the static indexed-mesh checks at lines
   480-756, then saves at line 1121. It never imports or calls
   `presave_validate.validate_current_scene`. The rejected pure report still
   names signed root coverage and exact self/correspondence/candidate crossing
   tests as pending at lines 855-857. The resulting build report would itself
   record `mechanical_preflight_pass: false` at line 1163. A uniquely named
   quarantined save may precede semantic rendering and clean reopen, but root,
   crossing, and self-intersection checks must pass in memory before that save.
   Integrate the adapter, pin its hash, require its report `pass`, and include
   the report path/hash in the build report before enabling a later attempt.

2. **Blocking for reuse: the consumer schema is deliberately incompatible
   with the reset artifact and must not be silently relaxed.** Lines 600-709
   require positive pair IDs, skins `-1/+1`, triangular paired skins, one
   paired-boundary bridge quad, and inner-skin roots. The rejected frozen
   architecture used bridge-only `0/0` vertices, paired quads, multi-ring
   bridges, and two target-facing outer roots. The current hard rejection is
   correct. Any successor must first publish a versioned interface contract,
   then encode an exact ordered-annulus rule and per-root skin semantics in
   both producer and consumer; merely permitting `0` attributes or quads
   would under-validate the bridge.

3. **Non-blocking while no-run: clean reopen is structural, not mechanical.**
   `validate_reopen.py` verifies the builder hash before importing it (lines
   115-142), confines the report-supplied candidate to `candidate/` and
   rejects protected aliases (lines 130-140), revalidates hashes/source state,
   and compares exact frozen geometry. It correctly leaves
   `mechanical_preflight_pass` and `render_allowed` false at lines 281-284.
   Once finding 1 is fixed, reopen should also verify the pinned pre-save
   report rather than treating `builder_preflight_pass` as sufficient.

## Checks that are sound

- Workspace derivation is correct for `out/<task>/candidate` (lines 29-31).
- Authorization requires the exact role/path/hash set and independent stable
  hashes (lines 258-378); snapshots and live rehashes reduce input drift
  (lines 381-412).
- The static mesh guard checks finite millimetre coordinates, duplicate and
  degenerate geometry, closed two-face edge use, connectedness, Euler `2`,
  consistent winding, and positive signed volume (lines 480-596). Construction
  converts millimetres to metres exactly once at lines 915-918 and preserves
  explicit topology and attributes.
- Source preservation now normalizes only the two authorized visibility
  fields and checks the exact additive object/mesh/material/collection
  inventory (lines 827-899).
- Output naming is unique and non-overwriting, remains under `candidate/`,
  cannot alias either protected input, disables `.blend1`, and rehashes every
  frozen input around the single save (lines 972-1127).

This verdict evaluates the scaffold, not the rejected visual design. The
current correct action is no build and no render.
