# C5 reference-construction verdict

## Scope

This is an implementation-blind reference review. I inspected:

- `canonical_front_25cm.png`;
- all 30 frames of `canonical_turn_180.gif`, with frames 14--23 used as
  the controlling rear sequence;
- `physical_front.png` and `physical_side.png`;
- every frame of the supporting `turn.gif`;
- representative frames from `sofa.gif`;
- the accepted C1b front and three-quarter renders; and
- the C1b-derived fixed rear render at `C3_baseline/rear.png`.

I did not inspect the C5 candidate, its implementation, or its object data.

## Hard verdict: center-back seam is unsupported

**UNSUPPORTED -- do not add or retain a center-back seam.**

The canonical rear turn does not show a continuous stitched, piped, or
pressed line running from the bow root down the center of the hair cap. In
frames 14--23, the readable construction is a stack of long, leaf-shaped
fabric hair panels. The dark lines are free panel edges and the soft occlusion
shadows where one panel overlaps another. They terminate at panel tips or
disappear beneath neighboring panels; they do not behave like a seam embedded
in one continuous cap.

The canonical front and physical photographs likewise show continuous brown
pile beneath the bow. Any faint central tonal change is consistent with pile
direction, compression, or lighting and is not stable enough across views to
justify seam geometry. The low-resolution supporting turn is too coarse to
override that evidence and also does not show a persistent center seam.

C1b already reads too much like one smooth spherical helmet. A thin center
line would decorate that incorrect reading instead of making the hair look
like constructed plush. It would also risk looking drawn, floating, or like a
hard molded groove.

## Reference-supported construction cues

The references support these cues on the brown hair:

- **Overlapping fabric panels:** the rear is several long padded/felt locks,
  not a single cap divided by decorative lines.
- **Soft edge thickness:** panel boundaries have a small rounded fabric edge,
  visible most clearly in the canonical rear sequence and physical side view.
- **Narrow contact occlusion:** a weak dark wedge appears only where one panel
  physically lies over another.
- **Root compression:** the bow presses the crown locally; this is broad and
  low frequency, not a sharp crease extending down the back.
- **Short directional pile:** physical front, side, and sofa images show fine,
  low-contrast fuzz. It softens edges and highlights but does not create random
  bumps or a dark drawn line.
- **Mild asymmetric stacking:** adjacent rear locks sit at slightly different
  depths and their tips do not form a perfectly radial, mirrored fan.

## Smallest safe alternative local detail

Add **one real overlap edge to the existing central rear lock**, rather than a
center seam:

1. Keep the lock root beneath the bow, the outer head silhouette, and the
   front view fixed.
2. Along only one lateral edge of the central rear lock, from roughly mid
   length to its tip, separate that edge from the neighboring rear panel by
   about `0.5--1.0%` of head width in depth.
3. Give the exposed edge a soft rounded thickness of about `0.3--0.6%` of head
   width and let the resulting geometry create a narrow occlusion shadow.
4. Do not add an independent line, groove, stitch strip, or full-length gap.

This is the smallest reference-supported change because it converts an
already-present lock boundary into a physical fabric overlap without changing
the major silhouette or regenerating the hair subsystem.

**Keep condition:** in the fixed rear and three-quarter views, the central
lock reads as a thin sewn/felt panel resting over its neighbor, while the front
view and outer silhouette remain unchanged.

**Undo condition:** the edge reads as a hard crease, dark drawn line, floating
shell, thick armor plate, or produces no visible construction improvement.
