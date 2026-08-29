# A83 cycles C1--C5 evidence

Observed through 2026-08-31T18:52:16+03:00.

## Immutable inputs

- Protected rung 003:
  `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/reimu_fumo_working_rung_003.blend`
  (`sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`).
- Tracked reusable model:
  `projects/renders/blender/fumo/reimu_fumo/reimu_fumo.blend`
  (`sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`).
- Neither input changed during these cycles.

## Reviewed local deltas

1. C1 moved only the upper cap. It exposed the receiver as a bald patch and
   was undone.
2. C1b applied the same narrow upper-cap correction to the cap and nested head
   receiver together. The fixed front and three-quarter pair showed a genuine,
   small reduction in the helmet-like upper width with no exposed gap. The
   blind reviewer selected **KEEP**. The isolated working checkpoint is
   `a83_C1b_coupled_cap_receiver_narrow.blend`
   (`sha256:d2357588b42b18285f31fcf780f2be5e76111a002a25b9ac25cd569be6cbf8d1`).
3. C2 changed coupled rear depth. Its fixed renders showed no material visual
   gain, so it was undone.
4. C3 bent one central rear lock while preserving its root. It produced only
   slight seam/tip drift and no material construction gain, so it was undone.
5. C4 tucked the coupled mid-side region. Its controlling pixels were
   effectively unchanged, so it was undone.
6. C5 added a proposed center-back seam. The three-quarter PNG was byte-
   identical to C1b, and only 257 rear pixels changed as isolated specks. A
   reference-only review also found no center seam in any canonical turn
   frame: the visible lines are overlapping free panel edges. C5 was undone.

The standing verdict files have SHA-256 digests:

- C1b: `e1e8970777faa8e7ccfa591aa252d231064156e15303af543a225944d02a3e6f`;
- C2: `389d85a8ccd1ba044134726861cb5857d066c33dcf7fdb7da88047a12edf0567`;
- C3: `ecf316f28b99dea01e4d7d108b222199a3163b6f73ca872a5a81df583326449d`;
- C4: `686c1c8aaedd5ff18759e53ac456fcb9a9cbea3da68c01bb8517dc98f9593872`;
- C5: `7d7db17cbdfc6b9481cfdd301ed9f27bde1615f81120437619199809137905a1`.

## Honest process and result verdict

C1b is the only retained local improvement; it is an internal working
checkpoint, not an approval candidate. C2--C4 establish that another broad cap
parameter change is not the next useful modality. C5 further shows that
decorative surface detail cannot substitute for physical panel construction.

The process now starts from C1b, changes one smallest connected owner, renders
the controlling and regression-risk views immediately, and keeps or undoes the
edit from pixels. Whole-head generation is not permitted without documented
structural-limit evidence and a discriminating test.

## Next discriminating test

Reprofile only one existing camera-facing rear lock into a visibly separate,
thinly stuffed fabric panel. Create a real soft overlap and contact shadow,
not a drawn line. Keep only if the fixed three-quarter view gains a substantial
planar fabric read and the rear gains a clean layered edge without changing
the accepted face, fringe, bow seat, or outer bounds. Otherwise undo.
