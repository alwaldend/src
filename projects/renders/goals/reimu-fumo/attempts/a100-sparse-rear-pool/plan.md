# A100: sparse rear-pool correction

## Objective

Test one causally distinct, reversible correction on the exact complete
rung003 high-water. Remove the unsupported horizontal rear skirt rail without
changing the rest of the plush.

## Frozen input

- Blend: `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/reimu_fumo_working_rung_003.blend`
- SHA-256: `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Canonical tracked model remains untouched.

## Exact owner and allowed change

- Editable surface: `Garment42 rear pooled dress panel`.
- Attached trim may follow only where necessary:
  `Garment42 rear pooled ruffle` and the nine
  `Garment42 hem stitch rear 00` through `08` objects.
- Add one reversible shape key to the panel.
- Freeze explicit vertex indices for the free rear boundary and move only
  those vertices approximately 9--10 mm in world `-Y`, from a maximum near
  `0.0492253 m` to the unchanged seat support ending at `0.040 m`.
- Preserve Z on the first move. Pin waist and side seams.

No coordinate field, proportional formula, remesh, generator, topology
change, new geometry, material edit, seat edit, front/side panel edit, foot or
leg edit, or upper-body edit is allowed. All other object identities and
datablocks are frozen.

## Evidence and gate

1. Verify the source hash and make an isolated task-local copy.
2. Open the copy in the persistent repository-pinned Blender 5.2.1 host.
3. Inspect the exact panel topology and record the explicit free-boundary
   indices and before coordinates before authoring.
4. Render the untouched source through the same fixed front, exact-side,
   rear, and worst-three-quarter views used for the candidate.
5. Apply the single sparse move and render immediately.
6. Keep only if the exact side loses the horizontal trailing rail, the worst
   three-quarter view reads as a compact seat-supported pool rather than a
   cape/ramp, the front does not regress, and the targeted lower
   silhouette/construction/contact is at least 6/10 with a clear A/B
   preference.

Stop on the first coherent render if the cape/rail category remains, a kink
or intersection appears, or improvement requires edits outside the owner.
Allow at most one adjacent-loop correction, and only after the first move
already passes categorically. On failure, close this route honestly instead
of starting another generator, rebuild, replacement, sculpt harness, or
unreviewed expansion.

## Process controls

- Target time to first decision-bearing A/B sheet: one focused authoring
  cycle after the attempt opens.
- There is one shared mutable `.blend` and one immediate pixel gate, so the
  coordinator remains the sole author. The already completed fresh reviewer
  supplies the independent correction verdict; no other workstream can
  affect the first pixels without creating shared-state or review overhead.
