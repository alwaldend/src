# A88 S04 render-only measurement gate

## Scope and controlling evidence

This gate measures only the fixed 512 px renders. No Blend file, camera,
geometry, or tracked file was opened or changed.

- Candidate: `snapshots/A88_crown_s04.blend`, SHA-256
  `375156e4951a49bdb2e13877baade7059729495a56b3bd431f48440f6e57a802`.
- Candidate packet manifest: `packets/s04/manifest.json`, SHA-256
  `018f79685d50a495f9c6db2729d38cb6a3846d430a6a885a0786a9ebc4ba0ce8`.
- Exact parent: A87 S04 front-hair survivor.
- Local comparison: A88 S03.
- Primary references: `canonical_front_25cm.png` for the front hair
  silhouette and `canonical_turn_180.gif` for side/rear continuity.
- Supporting references: `clean_front.png`, `physical_front.png`,
  `physical_side.png`, `turn.gif`, and `sofa.gif`.

The candidate head width is `Wh = 371 px` (`x=72..442`) and is the
normalization unit below. Bright exposed receiver pixels were counted where
all beauty-render RGB channels exceed 245. Component preservation uses the
flat front ID pass, assigning antialiased pixels to the nearest declared ID
color. These tests are intentionally separate: the beauty pass measures the
visible pale defect, while the ID pass measures owner silhouettes.

All frontal references show a continuous brown crown beneath the bow. Their
center construction line is a dark seam/fold, never a pale receiver hole.
The canonical turn and physical side also show hair continuing downward and
rearward as overlapping fabric panels; they do not support a broad, bare,
horizontal receiver shelf under the crown edge.

## Measurements

### Front crown exposure and center seam

| State | Bright receiver pixels in crown ROI | Exposed bounds | Maximum horizontal run | Result |
| --- | ---: | --- | ---: | --- |
| A87 S04 parent | 14,939 | `x=89..425`, `y=62..149` | 294 px | Deliberately bald parent opening. |
| A88 S03 | 217 | `x=231..282`, `y=62..72` | 52 px (`0.140 Wh`) | Crown coverage works, but a large diamond-shaped pale center hole remains. |
| A88 S04 | 50 | `x=244..272`, `y=67..69` | 29 px (`0.078 Wh`) | 99.7% less exposed area than the parent and 77.0% less than S03, but still a clearly visible pale slit. |

S04's remaining center opening is 29 px wide by 3 px high
(`0.078 Wh x 0.008 Wh`). The 3 px height is small, but width is more than the
`0.05 Wh` major-gap tolerance and the cue is categorically wrong: the
references require a dark fabric seam or crease, not illuminated skin. S04
therefore does not pass the center-seam gate.

### Front left/right root gaps

Trimming the S03 side protrusions restores the exact exposed receiver gaps
already present in the incomplete A87 parent:

| Gap | A88 S03 | A88 S04 | Normalized S04 bounds |
| --- | ---: | ---: | --- |
| Image-left root | 0 bright pixels | 49 px, `x=72..77`, `y=179..192` | `0.016 Wh` wide by `0.038 Wh` high |
| Image-right root | 0 bright pixels | 90 px, `x=431..439`, `y=178..195` | `0.024 Wh` wide by `0.049 Wh` high |

The S04 bounds and counts are byte-for-pixel identical to the A87 parent in
these ROIs. Thus S04 does not regress the already-frozen A87 shapes, but it
also does not perform A88's integration job. The canonical and clean front
show continuous hair at both roots; neither supports these bright notches.

### Three-quarter root/rear-interface exposure

S04 removes S03's long side tabs, which improves the helmet diagnosis, but it
does so by cutting the crown at a nearly horizontal upper-side edge. This
exposes a much larger receiver shelf:

| Fixed view and inner ROI | A88 S03 | A88 S04 |
| --- | ---: | ---: |
| Three-quarter, `x=90..220`, `y=140..180` | 0 bright pixels | 3,878 px; begins at `y=148`; max continuous width 126 px (`0.340 Wh`) |
| Mirrored three-quarter, `x=280..422`, `y=140..180` | 0 bright pixels | 3,354 px; begins at `y=148`; max continuous width 107 px (`0.288 Wh`) |

Both exposed shelves persist for the full 33 px measured vertical interval
(`0.089 Wh`) and continue below it. Missing rear hair is a later owner, so
these pixels are not a demand to complete the hairstyle in A88. They are
nevertheless valid interface evidence: a future rear panel would have to hide
a 0.29--0.34 Wh horizontal shelf and a hard T-junction rather than meet a
plausible tapered or overlapping crown boundary. That fails A88 gate 5.

### Crown silhouette movement

Relative to the A87 receiver's outer front silhouette, S04 begins one pixel
higher. Across the first visible crown rows (`y=61..79`), each side moves
outward by 12.5 px on average and as much as 24 px (`0.065 Wh`). Across the
stable bulk of the crown (`y=80..147`), the mean side movement falls to
5.4 px (`0.015 Wh`) and the maximum is 9 px (`0.024 Wh`).

The bulk silhouette is within the usual 3--5% tolerance, but the 6.5% apex
excursion exceeds the 5% silhouette-extremum gate. This is not the dominant
failure, yet it is another reason not to promote S04 unchanged. S03 has
essentially the same top-crown expansion; S04's principal geometric change is
instead removal of S03's visible side extensions: the visible crown ID bounds
change from `x=60..451, y=61..219` to `x=78..433, y=60..147`.

### A87 fringe and side-lock regression gate

The protected lower shapes pass exactly in the front ID render:

- Fringe at and below `y=150`: zero XOR pixels; identical per-row extrema.
- Image-right-ID lock at and below `y=150`: zero XOR pixels; identical
  per-row extrema.
- Image-left-ID lock at and below `y=150`: zero XOR pixels; identical
  per-row extrema.

Whole visible-mask IoU is 0.9935 for the fringe, 0.9657 for one lock, and
0.9784 for the other. All whole-mask differences occur at their upper roots,
where the new crown panels change occlusion; none changes the protected lower
fringe edge, lock lengths, tip shapes, or lateral extents. Therefore A87's
provisional front-hair module is preserved, and the failure is localized to
the new A88 crown/seam/interface owner.

## Verdict

**Reject A88 S04 as the crown-module survivor.** It is a useful diagnostic,
not a promotable checkpoint.

S04 makes two measurable advances: it removes 99.7% of the parent's pale
front-crown area, and it preserves A87's lower fringe and lock silhouettes
exactly. Those wins do not override three absolute failures:

1. the center seam remains a 29 px (`0.078 Wh`) bright receiver hole rather
   than a dark fabric seam;
2. the front root gaps remain exposed on both sides; and
3. both three-quarter views expose 0.29--0.34 Wh horizontal receiver shelves
   beneath an implausibly hard rear/lateral crown boundary.

The evidence also explains the S03/S04 tradeoff: extending the existing crown
downward closes the gaps but becomes a helmet-like side wall; trimming it
removes the wall but exposes the receiver. The next bounded cycle should keep
the exact A87 lower silhouettes and the two-panel crown concept only if it can
change the *boundary construction*: close the two center surfaces as a
coupled seam, then taper or overlap the lateral/rear edge so a later rear-hair
panel has a sloped fabric interface rather than a horizontal shelf. A simple
translation or another uniform downward extension is already disconfirmed by
S03 and S04.
