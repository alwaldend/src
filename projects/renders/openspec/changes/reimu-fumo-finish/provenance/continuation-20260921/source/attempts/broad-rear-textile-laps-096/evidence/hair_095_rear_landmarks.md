# Two broad rear laps: measured footprint evidence

Read-only observation, 2026-09-06 02:40 UTC. No Blender run, model edit, builder, or frozen plan. Source candidate094 SHA256 `af5a61921ae309a69cfe98b0da092d206b2d0ff29c9123ac2a4886de5a0f9add`.

## Bounded recommendation

The references support two long, broad side/rear laps over the existing brown rear underlay, not enlarged lower lobes. Replace the footprints of **both `Macro066 released side cloth ±1 overlap 1`** first. Preserve both `overlap 0` pieces and the separate `Sewn047` hanging front locks initially. Removing all four sidecloth pieces is not justified by these rear-oblique views. Even replacing only overlap1 needs front and mirrored-front silhouette checks: their low lateral tips contribute below overlap0.

The photographed visible free boundaries are long, sweeping, and asymmetric in projection. Both views establish a foreground lap over a central brown field; they do **not** establish anatomical left/right identity, hidden mutual lap order, exact seam count, or the concealed root topology. The current094 continuous brown cover can remain the underlay. These are open footprint constraints, not complete ready-to-mesh polygons.

## Image registration and visible traces

Sources: [reference16](side_062_frame_16.png), [reference21](side_062_frame_21.png), each 498×498; [094 rear](hair_094_all_review/rear.png) and [094 side](hair_094_all_review/side.png), each 512×512.

Pixel origin is upper-left. Register each reference by root-entry anchor R and hair/collar-center anchor C: reference16 R=(311,69), C=(290,330), L=261.84px; reference21 R=(266,66), C=(264,337), L=271.01px. Let ev=(C−R)/L, eu=(ev.y,−ev.x), u=dot(P−R,eu)/L, v=dot(P−R,ev)/L. Thus u is approximately screen-right and v runs root-to-collar. This is **2D landmark normalization**, not recovered camera calibration or a pixel-to-millimeter mapping.

| Landmark on near lap | Ref16 pixel → (u,v) | Ref21 pixel → (u,v) |
|---|---|---|
| Root exit, partly concealed | (311,80) → (.003,.042) | (264,79) → (−.007,.048) |
| Upper visible free edge | (310,120) → (.012,.194) | (270,126) → (.016,.221) |
| Mid free edge | (313,195) → (.046,.479) | (256,194) → (−.033,.473) |
| Lower free edge | (320,265) → (.094,.743) | (229,250) → (−.132,.680) |
| Tip approach | (337,311) → (.173,.913) | (187,291) → (−.285,.832) |
| Visible terminal tip | (368,322) → (.294,.946) | (145,316) → (−.440,.926) |

Manual localization uncertainty: visible boundaries/tips ±6–8px; bow-occluded root anchor ±12px; collar anchor ±8px. Including registration uncertainty, treat normalized landmarks conservatively as approximately ±.06–.08, not precision targets. Reference21's visible width is roughly 94px near y190 and 81px near y250 (about .35L and .30L), with partly ribbon-masked outer boundary; width uncertainty about ±.04L. Do not mirror these widths literally into reference16's differently foreshortened view.

Reference16's screen-right lap runs almost straight through its upper half before flaring to its outer tip; the central underlay continues behind it toward approximately (295,333). Reference21's screen-left lap sweeps progressively left and finishes around (145,316); the central field continues behind toward approximately (237,345). Root entries are under the bow in both. Outer edges near ribbons are occluded, so their hidden closure must remain an explicit design choice. Visible occlusion requires each photographed near lap to stand in front of the central underlay, but no evidence requires two new laps to cross each other at the center.

094 rear has approximate root/crown anchor (255,145) and collar center (255,338), L≈193px. It reads as a continuous round central cap with narrow lower lateral strips. The references instead place a free edge over almost the complete root-to-collar interval, with substantial mid-height cloth width. Side094 shows the same smooth round cap; merely extending its lower tips would not introduce this missing long overlap boundary.

## Existing pieces and rear-chart support

Actual bounds from `sleeve_095_actual_assembly.json`, observed 2026-09-06T02:23:15.147157+00:00, sourcing exact094. All values millimeters; Y is model depth, not image depth.

| Existing pair | −1 X / +1 X | Y | Z | Footprint decision |
|---|---|---|---|---|
| overlap0 | [−60.55,−41.26] / [44.72,59.06] | [−22.00,5.44] | [92,148] | Preserve frontward lateral framing initially |
| overlap1 | [−61.14,−31.51] / [34.44,59.66] | [1.03,27.02] | [77,148] | Superseded by the two broad rear laps |

The existing roots top out at Z148, far below the rear crown Z194.39. New broad laps need a higher root/upper footprint; rescaling the current lower strips is not equivalent. The assembly export lacks full arrays for the separate tied front locks, so it cannot by itself certify their occlusion or preservation.

From `hair_094_repair_arrays.json`, select outer cover triangles whose original barycentric source IDs are all ≥3298: 6,494 rear triangles, 140 projected boundary edges. Rear-chart bounds: X[−58.895,57.300], Y[−9.260,40.992], Z[80.000,194.387]mm. Horizontal boundary sections are:

| Z mm | Supported X interval mm |
|---|---|
| 90 | [−31.26,37.01] |
| 110 | [−49.86,52.86] |
| 130 | [−58.85,57.28] |
| 148 | [−56.49,55.00] |
| 160 | [−53.87,52.29] |
| 180 | [−43.20,40.07] |
| 190 | [−28.21,23.75] |

This known XZ graph supplies broad mid-height support but narrows sharply near crown and floor. Literal chart clipping cannot preserve an old overlap1 tip at Z77: the chart ends at Z80. A lower extension would be a separate explicit geometry change. Rear-chart Y is not everywhere positive; globally treating positive Y as safe outward separation would be incorrect near its lower turn.

## Minimum measurements before a plan

Choose two closed XZ footprints consistent with the open traces, placing only their concealed root closures beneath actual bow receiving regions. Measure each footprint's section widths and projected free-edge/tip landmarks in the two oblique reference orientations; do not use unlike camera views as metric matches. Record the intended layer relation at every shared footprint interval—lap over underlay, and either noncrossing center boundaries or one explicitly chosen lap-over-lap order.

Measure actual bow/root clearance and cover-to-lap gap, including finite wall thickness and rim faces; root attachment cannot be justified by a label or by proximity alone. Check full rear/side silhouette and depth against094, plus front/mirrored-front width and tip-height regression after removing overlap1. A constant-Y wall is plausible only on the selected rear graph: crown/root curvature and bow intersections remain unmeasured. No collision, containment, attachment, or candidate approval is claimed here.

Provenance: reference16 SHA256 `209fb44b435093dc1ce7e7339cd05503091d6d7815c471cda1d4e5b4843fc432`; reference21 `685de63c429493eda0047253f875e1b20bd1ae7847d9610ed57ebf5ef4a1a333`; 094 rear `88d95ba8ad6001e4e714c488933cd4bf06daa62d1e2d1df81de5fc7b7d54d9bd`; 094 side `8c6dae2cf78e1cce98b0385b11010c032a60ff762a00484d3b50a7a27d402535`. Rear arrays trace to unchanged087 SHA256 `c171fea554811e825e2cd9d7f1d5a069c6b305e85888c6058d2390e1a20738e1` and geometry-helper digest `e10bd0a85162f53de7b8adde7913136a838850092bd988ea2f09bb13ca00d975`.
