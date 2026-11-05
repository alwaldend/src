# Plan

Return to immutable A157 at SHA-256
`433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
This supersedes the rejected single-extremum Grab setup with deterministic
mesh-space selection and a distributed native transform.

The smallest coherent subsystem is the two lower-tail assemblies. Each side
contains its padded red tail and separate white perimeter ruffle. The upper
loops and center knot remain frozen because their measured front span is
already within the canonical silhouette band.

1. Copy protected A157 to a new task-owned path and record all mesh and object
   baselines before editing.
2. For the left assembly, enter multi-object Edit Mode with only
   `A42 Left independent draped bow tail` and
   `A42 Left narrow gathered tail ruffle` active. Select the full-thickness
   outer 20 percent bands using source-space thresholds `x <= -0.0879192203`
   and `x <= -0.0974285379`, respectively.
3. Invoke Blender's native global-X proportional translation operator with
   connected topology, smooth falloff, and radius `0.045`. Enter numeric
   displacement `-0.0267` through timer-spaced native key events and confirm.
4. Repeat once on the right assembly with thresholds `x >= 0.0879192203` and
   `x >= 0.0974405333`, then numeric displacement `0.0267`.
5. Selection may be established through Blender's edit-mesh API, but no code
   may assign vertex coordinates. Preserve each ruffle's Solidify modifier.
6. Require 132 selected vertices on each red tail and 67 on each ruffle.
   Require selected extrema displacement `0.0267 +/- 0.001`, root-band
   displacement no more than `0.001 Wh`, unchanged depth/thickness within
   `0.005 Wh`, and no non-target geometry change.
7. Save once after both operators complete. Clean-reopen the exact bytes in
   repository-pinned Blender 5.2.1 and render fresh 512 by 512 front and side
   views. Only if they pass, render rear and both three-quarter views.
8. Target complete lower-tail bow span `2.038 Wh +/- 0.05 Wh`, centered in the
   frame. At the fixed 512 camera this is about `503 +/- 12 px`.

Reject if a modal operator fails, an outer tail crops, either root moves or
detaches, a panel becomes triangular, a ruffle separates or collapses,
left-right asymmetry exceeds `0.05 Wh`, or a loop/knot silhouette moves more
than `0.05 Wh`. Stop after the paired operation and review; this cycle can
retain a bow subsystem but cannot pass the whole-character clay stage.
