# Full-cover feasibility after 084 and 086

Observed 2026-09-06T00:39:27Z; linked worktree branch
`t3code/continue-fumo-desktop-use`, Git HEAD `c7601f0fc80e0a94b585910459639ffccbdbdbd4`.
Read-only code/evidence review; no Blender process, geometry generation, model
save, or source/goal/Git change. Exact 087 bytes independently hashed:
`c171fea554811e825e2cd9d7f1d5a069c6b305e85888c6058d2390e1a20738e1`.
The coordinator owns the verdict and any subsequent trial authorization.

084/086 materially improve the basis for a redesigned full cover, but do not
establish that rerunning 077 would work. The strongest case is now a continuous
cover with explicit shared layer correspondence, using 086's independent rear
triangulation and 079's clipping identity method. A thin duplicate of the current
envelope still has no evidence of fixing the helmet read.

## Which causes changed

| Earlier cause | Relevant advance | Still unsupported |
| --- | --- | --- |
| 077's raw clipped vertex had degree four and disconnected face fans; its shell bridge became nonmanifold. Independent coordinate rounding lost source-edge identity. | 079 caches canonical source-edge crossings, reuses source endpoints, rejects disconnected fans and non-degree-two boundary loops. Its saved front patch passed topology checks. | 084/086 do not repair the clipper. The 079 helper explicitly accepts only the front XZ chart: its winding, boundary classification, and front-index assumptions cannot simply be applied to the whole head. A full-cover cut needs the same identity rules across all charts. |
| 077 had 18 real inner/outer crossings and 26 core/inner crossings near rear Z127 mm, separately from the Z105-mm clipping defect. | 086 removes the inherited graphic-constrained rear triangulation. Its lower-lid diagnostic samples include Z127.596 mm, spatially matching the old crossing region; maximum strip chord error fell from 483.078 to 18.701 micrometres across 36 samples. | Spatial correspondence supports the causal connection but is not a shell replay. 086 proved one surface's hygiene and sampled shape improvement, not inner/outer/core ordering. Independent radial displacements can still intersect in triangle interiors or at chart transitions. |
| 079 ended the separate front patch at the head-outline/crown boundary, producing a visible side seam; the smooth cap still read as a helmet. | A full cover can make that chart boundary internal. Neither 084 nor 086 created such continuity. | Closing each chart independently would recreate a ridge. Preserving the entire current curvature while adding a shallow fringe already failed visually in 079. |
| 079 intersected front locks at Z111–121 mm and overlap-0 side panels at Z107–144 mm. | 084 buries side-panel roots by at least 0.75 mm at Z144 mm and above; its deliberate emergence region is Z134–144 mm. | Front locks remained exact. 084 is not clearance evidence for a new cover, especially below Z134 mm. Existing contacts cannot be relabeled as concealed roots: the 084 exception is spatially bounded and concerns the existing receiving head. |

The new [088 reference diagnosis](head_088_reference_diagnosis.md) supplies a
separate visual reason to rebuild the interface: continuous cushion-to-underside
roll, overlap-defined face opening, and broad side fabric flowing into the cover.
It supplies neither calibrated new amplitudes nor proof of current topology.
A full cover alone cannot claim to resolve the lower-face roll or blade-like locks.

## Credible construction alternative

1. Define a single connected cover domain across front, gusset/crown, and rear.
   Keep the face graphics on a continuous cushion. Reuse 079's canonical vertex
   and edge correspondence, generalized to the actual 087 surface; keep front,
   side, and rear chart boundaries as internal shared edges. Close the cover only
   at deliberate exposed or concealed garment boundaries, never at the crown's
   front/back chart division.
2. Build outer, inner, and underlying cushion from a common subdivided surface
   correspondence. Front and rear can use ordered graph layers over XZ; the
   turning gusset needs its own nonfolding chart and shared transition vertices.
   Within a graph chart, common projected triangulation permits checking layer
   order across triangle interiors. That argument does not extend automatically
   across chart transitions: validate those and the complete shell separately.
   Use the 086 rear connectivity, not a copied front graphic triangulation.
3. Give the opening a continuous thin return and controlled overlap onto the
   cushion; reserve inward accommodation for covered cushion. Define the lower
   face/rim transition from the reference section and preserve graphic landmarks.
   This is a new interface construction, not another uniform lift of 079 or
   independent padded bang pieces.
4. Seat side-panel roots against the new cover/core interface through deliberate
   concealed overlaps. Rebuild the front-lock receiving transition if clearance
   requires it. Keep 084's useful lower-panel shape only where the new interface
   and reference flow support it. Reseat existing pile after receiver geometry
   passes; pile must not hide exposed cut edges.

This costs more than a front patch because it includes the cushion return and
hair attachments. It directly addresses the failures a front-only overlay left
untouched. Retaining 087 while obtaining the evidence below is a credible fallback;
another unchanged 077/079 setup is not.

## Minimum evidence before a native model trial

- **A complete domain and seam map on exact 087:** original vertex/edge IDs,
  chart membership, exposed fringe, concealed termination, and internal seams;
  degree-two boundary loops and connected vertex fans after the proposed cut.
  Include the former Z105-mm clipping witness and report sliver altitudes, not
  only triangle areas. Show that no exposed boundary crosses the crown.
- **An ordered layer construction at actual triangles:** front/fringe sections,
  both former rear Z127-mm witnesses, and every crown/gusset transition. Report
  signed ordering and thickness in declared directions, distinguish directional
  from normal clearance, and demonstrate that the proposed correspondence does
  not fold. A new scalar thickness without this evidence is insufficient.
- **The complete attachment interface check:** proposed cover versus both front
  locks and four side panels, with resolved intersection segments and declared
  concealed zones. Include Z111–121 mm front locks and side emergence below
  Z134 mm. Verify actual depth order and root coverage; triangle bounding boxes
  alone neither prove nor disprove the intended contact zone.
- **One measured reference section/edge target:** align the controlling front
  and oblique views and record opening overlap, fringe return, lower-face roll,
  and side-panel emergence with uncertainty. Name what changes beyond shell
  thickness and the expected silhouette effect. The 088 image-only diagnosis
  cannot supply numerical geometry by itself.

These can be established in a read-only geometric preflight without saving a
candidate. Passing them supports one frozen construction trial, not visual
acceptance. That trial still needs complete shell/core collision checks and
fixed isolated/assembled front, side, rear, and both three-quarter image review.

## Evidence identity and limits

Primary local sources: [077 plan](hair_077_plan.md), [077 result](hair_077_result.md),
[077 causal diagnosis](hair_079_diagnostic.md), [079 result](hair_079_result.md),
[084 result](hair_084_result.md), [084 contact localization](hair_084_zone_diagnostic.md),
[086 result](hair_086_result.md), [086 rear diagnosis](hair_086_rear_diagnostic.md),
and [087 result](chin_087_result.md). Inspected implementation includes
`hair_077_build.py`, `hair_079_front_patch.py`, `hair_079_build.py`, and
`hair_086_build.py`; no scripts were executed.

Verified SHA256: 077 builder `c088ef61980b2da9a92d56363b4a8d18a0190fb4ab0eb82f2124edf0fb1c7278`;
079 helper `01f8e7759e1a3626d1af876139580d0922296647b04acdacf0c40f0d29a1e5c7`;
088 diagnosis `7156b3c37556c5b9caeaa67a8c0ad98e53a264e0cab6e88a17f1d3d3be06d775`.
Geometric and visual outcomes above are inherited receipts, not new validation.
The detailed 079 diagnostic JSON read was truncated and was not used; its complete
Markdown diagnostic supplies the crossing evidence. No full-cover feasibility,
calibrated landmark pass, or present-candidate visual approval is claimed.
