# A77 P0 independent code audit

Audited `build_p0.py` at
`sha256:51d30ca2136322205963fefe11f008cad03a6ff3208bb185ddb633a1324e05d6`
and its current `build_report.json` at
`sha256:96771d7a19e1c2e2c1702619ff1e7e16543e31b7302c0cb8fa63ee495b2c5dff`.

## Verdict: NO-GO

Do not render or promote the current three-object build. The exact base
topology and macro depth metrics are useful, but the authored forms still
recreate the closed failure families: a generated rounded support under a
near-concentric mantle, plus a ruled solidified rear card. The current report
also explicitly fails its mechanical preflight (`0.422 mm` anchor RMS and
`1.237 mm` maximum versus `0.25/0.45 mm`). This is not eligible for a bounded
beauty correction.

## Representation assessment

- **Head rest: better than A74's literal face card, but not the designed
  discriminator.** `X_ROWS`, `FRONT_ROWS`, `REAR_ROWS`, `TOP_DROP`, and
  `BOTTOM_RAISE` drive every control through linear row interpolation
  (`build_p0.py:32-71,182-193`). The named direct-edit sets are only vertex
  groups; no four coherent crown/rear/temple edits occur
  (`build_p0.py:526-576`). This contradicts the design's explicit ban on a
  per-height/generated final field (`sparse_form_design.md:116-119,157-184`)
  and retains rounded-box/mattress regularity. The reported depth
  (`.783 Wh`) passes, but evaluated lower-rear turn-in is only `9.11 mm`, short
  of the intended `12-15 mm` profile change.
- **Mantle: high helmet/canopy risk.** One connected swath owns crown, both
  temples, rear, and upper face; it is copied, uniformly offset `0.8 mm`, then
  uniformly Subsurf/Solidify thickened (`build_p0.py:578-610`). There is no
  authored exposed-boundary correction and no pair of subordinate lower-rear
  lobes. More importantly, independent Catmull-Clark boundary stencils mean
  copied controls do not remain close-fitting: the report's evaluated mantle
  depth is `118.64 mm` versus `103.41 mm` for the head, a `15.23 mm` envelope
  increase despite only `0.8 + 3.6 mm` nominal outward construction. This is
  the projecting-canopy mechanism A74 exposed.
- **Rear leaf: high card/cape risk.** Every free row has constant Y and Z
  across its five X controls, and the only volume is uniform `3.6 mm`
  Solidify (`build_p0.py:614-653`). That is a ruled sheet, not a transversely
  stuffed/tensioned panel. Its five base root controls are placed `0.5 mm`
  outside the mantle, but the final evaluated Subsurf/Solidify boundary is
  never seated or checked. It can therefore float or bridge while the metric
  passes.

These are direct recurrences of A74's smooth canopy/card/blade failure. The
closed head avoids that attempt's one-sided white face card, and its central
face plateau is less spherical than A75 C0, but the retained interfaces can
still recreate A75's cavity/island failure because they are not actually
validated.

## Mechanical checks are not sound enough to veto those failures

- `nearest_distances()` is unsigned and samples whole objects, while the gates
  use only each object's minimum and only an upper limit
  (`build_p0.py:421-434,674-724`). One tangent vertex passes even if the root
  crosses or the rest floats. The current report already has eye minima near
  `0.001-0.004 mm` (below the required `0.05 mm`) and mouth minimum
  `0.114 mm` (below `0.15 mm`), yet those checks do not fail. Root-object
  medians reach roughly `1.5-17.7 mm` while near-zero minima still pass.
- Leaf contact is measured on the five **control** vertices, not the evaluated
  leaf boundary (`build_p0.py:680-684`); it does not enforce `0.2-0.8 mm` over
  80% of samples or test crossings/bridge geometry.
- Topology reporting verifies raw counts, quads, connectivity, and edge
  incidence only (`build_p0.py:454-487`). It omits finite/degenerate and
  coincident-element checks, consistent outward normals, the base valence
  contract, exact open boundary loops, evaluated self-intersection, identity
  transforms, single-user data, and exact modifier order/settings.
- `mechanical_preflight_pass` is advisory: the script saves a blend even when
  false, and it does not clean-reopen the saved snapshot
  (`build_p0.py:732-748`). That violates the design's reject-before-render and
  saved-artifact gates.

## Exact next boundary

1. **Reset to a head-rest-only coupon.** Build the registration box, then apply
   four explicit keyed displacement sets for face anchors, crown compact,
   lower-rear tuck, and temple/high-rear shaping. Record before/after vertex
   coordinates per set; remove the row interpolation as final geometry.
2. Solve the nine anchor constraints and mouth support jointly on the
   evaluated surface. Hard-fail unless anchors meet `0.25/0.45 mm`, eyes are
   noncrossing at `0.05-0.35 mm`, and mouth is noncrossing at
   `0.15-0.55 mm`. Render only front and the worse profile. A rounded
   box/muzzle/helmet result ends this representation.
3. Only after that pixel veto passes, derive the mantle. Iteratively fit its
   **evaluated** surface under retained roots, directly author the exposed
   front edge and two unequal lower-rear lobes, and gate pointwise signed
   clearance plus triangle crossings. A raw normal offset is insufficient.
4. Only after the mantle passes, author per-column Y/Z camber, edge roll, and
   asymmetric tip tension in the rear leaf. Fit and test the evaluated root
   boundary (`0.2-0.8 mm` for at least 80%, zero bridge/crossing), not its base
   controls.
5. Make topology/normal/modifier/contact failures fatal before saving a render
   candidate, then clean-reopen the exact saved file under pinned Blender.

Do not add seams, fibers, materials, more leaves, or more topology to rescue
this build. The high-leverage correction is the staged representation and
evaluated contact boundary above.
