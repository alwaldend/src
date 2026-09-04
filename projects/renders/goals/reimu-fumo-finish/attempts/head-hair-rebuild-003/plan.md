# Head-and-hair representation reset

## Baseline and dominant failure

Start from protected A157, SHA-256
`433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
The body, bow, garment, sleeves, feet, cameras, and lighting are frozen.

The dominant failure is `head-hair-helmet-read`: the current continuous cap
and face opening read as a rigid helmet over a rounded doll head. Repeated
surface tuning cannot produce the photographed sewn construction.

## Hypothesis

Replace the visible head support with a shallow gusseted cushion and replace
the continuous cap with separately rooted rear, top, side, crown, and fringe
panels. Retain the existing flat eye graphics, mouth, and cheek-lock/ribbon
panels. The resulting neutral clay should show a readable front plane and thin
padded hair layers in front and side views without changing the frozen
whole-character layout.

## Gate

Render only front and side first from the exact saved candidate in pinned
Blender 5.2.1. Reset without an in-attempt correction if any of these holds:

- head or hair still reads primarily as a helmet, sphere, cube, armor, or
  unsupported cards;
- base-head depth leaves the `0.66-0.82 Wh` band or hair depth leaves the
  `0.71-0.87 Wh` band;
- visible beige face exposure or head silhouette misses the frozen canonical
  landmarks by more than `0.05 Wh`;
- hair panels float, clip, or create a tangency larger than `0.02 Wh`; or
- any frozen bow/body silhouette landmark regresses by more than `0.05 Wh`.

Only if the fast pair passes may the full five-view regression packet and an
implementation-blind review run. A retained subsystem is not a whole-clay or
goal pass.

## Parallelism

The model is a monolithic artifact, so one MCP writer performs the edit.
After the candidate is frozen, measurement and blind pixel review may run in
parallel on disjoint evidence outputs. The coordinator alone may retain or
reset the candidate.
