# A83 C6 implementation-independent pre-mortem

## Evidence and scope

This review uses only the fixed C1b three-quarter and rear renders and the
supplied reference images. It does not inspect Blender objects, topology, or
the proposed implementation.

- Exact-variant controls: canonical turn frames 06--10 for three-quarter and
  16--22 for rear.
- Construction cross-checks: `physical_side.png`, all four `turn.gif` frames,
  and the selected `sofa.gif` frames.
- Protected baseline: `C1b_packet/three_quarter.png` and
  `C3_baseline/rear.png`.

The desired result is not merely a thinner lock. It is one unmistakably
constructed, softly filled hanging fabric panel that replaces part of the
current spherical read while preserving the accepted outer bounds, face,
fringe, bow seat, and remaining hair.

## Strongest failure mode

The most likely failure is **helmet plus pasted card**. A thin reprofiled lock
can be technically separate yet leave the intact round cap visible behind it.
If the lock only acquires a sharp edge, constant stand-off, or local seam, the
viewer sees the same helmet with a rigid leaf attached to its outside. If the
root does not disappear continuously beneath the bow/crown receiver, the
stand-off becomes a bright or background-filled wedge and repeats the prior
floating-edge failure. If both the card edge and cap edge survive in
projection, it produces a double silhouette or armor-plate outline.

This risk is stronger than the risk of insufficient asymmetry: the canonical
panel can tolerate modest contour uncertainty, but it cannot read as a
detached card over a dome. The physical side and canonical rear sequence show
one seated root, a broad soft panel face, and one narrow overlap shadow. They do
not show a second exposed supporting shell around the panel.

## First discriminating pixel cue

Inspect the fixed **three-quarter** render first, before rear or implementation
details. From the root beneath the camera-facing bow tail through the middle
third of the lock, there must be a single continuous depth order:

1. the panel visibly covers a substantial patch of the former dome;
2. its middle reads as a broad, comparatively even/planar tone rather than the
   cap's uninterrupted rolling highlight; and
3. exactly one soft contact boundary separates panel from cap, with no bright
   wedge, background leak, parallel rim, or second outer contour.

At the 512 px gate, the planar read must remain obvious at normal size, not
only when zoomed. As a frozen minimum, it must span at least `0.12 Wh` in width
over at least `0.30 Wh` in root-to-tip height, and the physical overlap edge
must remain continuously legible for at least `0.20 Wh`. A thin line, isolated
tip, or highlight change smaller than this is not the target.

Only if that cue passes, inspect the fixed **rear** render. It must show one
unambiguous layer boundary: one lock lies over its neighbor and exposes a soft
rounded fabric edge. The boundary must terminate naturally toward the root or
tip; a complete outline around the lock is a pasted applique/card failure.

## Frozen binary decision

### KEEP only if every condition passes

- In three-quarter view, the panel creates the minimum broad planar region and
  minimum continuous overlap length above; the previous dome highlight is
  visibly interrupted across that region.
- The root is seated under the bow/crown with continuous contact. There is no
  connected background or bright gap wider than one review pixel, no parallel
  rim, and no detached root corner.
- The free edge has finite soft thickness but does not read as a hard board,
  tube, armor bevel, or paper-thin card. It produces one narrow occlusion
  shadow rather than a drawn dark line.
- In rear view, one edge has clear front/back ownership and a naturally
  terminating overlap. The panel is mildly asymmetric and softly filled; it
  is neither a perfect teardrop nor a straight-sided slab.
- The hair's outer silhouette and lowest bound differ from C1b by no more than
  `0.01 Wh` at any protected extremum. No double silhouette, crossing,
  accidental tangency, or z-fighting is visible.
- The bow seat and every unchanged neighboring lock remain visually stable.
  A later fixed front regression must also leave the face, fringe, head width,
  and overall height within the existing acceptance tolerance.

### UNDO if any condition occurs

- The cap still reads as the same complete round helmet and the edit is only a
  seam, narrow flap, isolated lower tip, or small highlight change.
- The panel is outlined on two or more sides against the cap, exposing the
  intact dome as a separate backing shell.
- Any root gap, floating bright wedge, background leak, parallel/double edge,
  hard constant stand-off, or visible intersection appears.
- The panel's middle is perfectly flat and rigid, its border is mechanically
  straight or uniformly beveled, or its entire perimeter is readable like a
  card/armor plate.
- The named three-quarter cue is not unambiguous at native 512 px, even if
  topology or dimensions changed.
- The rear overlap is ambiguous, looks like a drawn line, forms a full
  applique outline, or merely trades the helmet for a symmetrical petal fan.
- Any protected outer bound moves by more than `0.01 Wh`, or the bow, face,
  fringe, or neighboring locks visibly regress.

## Pre-mortem disposition

**Proceed only as this bounded keep/undo test.** The proposed owner is small
enough to be reversible and is supported by the references, but it succeeds
only if it removes the helmet read from a substantial projected patch. Adding
a thin object or edge while retaining the complete dome underneath is already
defined as failure and should be undone at the first three-quarter render.
