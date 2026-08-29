# A89 S05 rear/profile measurement and ownership gate

## Scope and evidence

This is a render-only gate on A89 S05 against its exact A88 S07 parent and
the supplied references. No Blend or tracked file was changed.

- A89 S05 source SHA-256:
  `0ecef762838665ed8da7274279defdb5fa1973b58f1ada02ebeebbca09382bbb`.
- A89 S05 gate manifest SHA-256:
  `bdd3c0eedd8da42df9eb6b038f510747859a152513c0abcdddc09c4901323814`.
- A88 S07 source SHA-256:
  `0a30e2af3142081648bb3137ad75d6d1cc73de55e9f830f85bad1f85e92c8788`.
- A88 S07 six-view manifest SHA-256:
  `e500a58389af2223c0ad471f3936081b4bbbb4cf3bf87e6de79530bd2fa38be1`.
- Fixed-render normalization: head width `Wh = 371 px`, inherited from the
  accepted A88 gate.
- Rear ownership masks use the declared component colors in `rear_ids.png`.
  The three A89 panels are the green classes, the A88 crown is blue, and the
  receiver head is tan. Profile deltas are beauty-render differences above a
  2% grayscale threshold.

The canonical turn controls side/rear silhouette, coverage, and layer order.
The physical side supports panel thickness and overlap. The canonical and
physical front images control protected front identity, but S05 has no front
or three-quarter render, so those protections cannot be certified by this
packet.

## Rear coverage and connectivity

The union of the three A89 ID classes is one connected raster component with
a `406 x 345 px` bounding box (`x=53..458`, `y=54..398`) and 117,851 pixels.
That proves literal closure but not correct construction.

| Measurement | A88 S07 parent | A89 S05 | Gate |
| --- | ---: | ---: | --- |
| Root/bald holes between crown and lower panels | Open lower owner by design | **0 px** | pass closure only |
| Inter-panel background/head holes | n/a | **0 px** | pass closure only |
| Separate visible panel components | n/a | **1 union component** | fail leaf read |
| Visible center crown rim | about 10 px (`0.027 Wh`) | **0 px** | fail depth order |
| Exposed blue crown components | continuous arch in A88 | two lateral remnants, `936` and `950 px` | fail |

S05 removes gaps by placing its roots in front of the accepted A88 crown rim.
The center rim that was the explicit tuck-under receiver is completely
occluded. Only two blue side arcs remain (`x=56..105` and `x=406..455`). Thus
the render does not demonstrate the required under-rim overlap; it shows the
opposite depth order and turns crown plus lower panels into one shield.

## Lower rear tip band

At representative rear columns, the A89 hair endpoint and the first visible
receiver-head pixel are adjacent, so there is no narrow accidental gap.
However, a large tan receiver band remains exposed below the hair:

| Rear x | Hair ends at y | Visible head band | Normalized |
| ---: | ---: | ---: | ---: |
| 140 | 392 | 36 px | `0.097 Wh` |
| 180 | 364 | 69 px | `0.186 Wh` |
| 220 | 367 | 68 px | `0.183 Wh` |
| 256 | 388 | 48 px | `0.129 Wh` |
| 300 | 359 | 74 px | `0.199 Wh` |
| 340 | 377 | 55 px | `0.148 Wh` |
| 380 | 391 | 35 px | `0.094 Wh` |

The median exposed band is **55 px (`0.148 Wh`)**, with a
`0.094..0.199 Wh` range. In the rear and rear-three-quarter canonical-turn
frames, the brown leaves continue to the shoulder/collar line and no pale
back-of-head band is visible. The physical side independently shows the same
coverage relationship. S05 therefore fails the lower-tip band even though
its three points make a closed geometric edge.

The center and two off-center points are also too close in visual language:
they form a single shallow scalloped curtain rather than the reference's
overlapping leaves with independent free edges.

## Profile projection and side-fin widths

The A89-versus-A88 changed silhouette width in the owned mid/lower band is:

| Profile row | Left profile | Right profile |
| ---: | ---: | ---: |
| y=260 | 62 px (`0.167 Wh`) | 62 px (`0.167 Wh`) |
| y=300 | 59 px (`0.159 Wh`) | 60 px (`0.162 Wh`) |
| y=340 | 36 px (`0.097 Wh`) | 34 px (`0.092 Wh`) |

The left and right widths agree within `0.005 Wh`, so the nominally
asymmetric leaves are effectively mirrored in projection. More importantly,
they are seen largely edge-on as thin fins beside a broad exposed side plane.
The canonical turn and physical-side references show broad, overlapping
fabric faces carrying the rear half of the side silhouette, not narrow fins
attached to a bare cushion.

At the roots, S05 changes the upper profile silhouette by about 25--26 px
(`0.067..0.070 Wh`) beyond A88 and reaches `y=55`, two pixels above the A88
crown. Those protruding corners are another consequence of roots sitting in
front of, rather than beneath, the crown.

## Protected A87/A88 regression checks

- The profile ROIs containing the opposite A87 front cheek-lock free edges
  change by 193 of 24,500 pixels (`0.79%`) on the left profile and 50 pixels
  (`0.20%`) on the right at the 2% threshold. No outer free-edge displacement
  is visible; these small changes are compatible with new cast shadow.
- The central upper-crown ROI changes by 4,301 of 37,400 pixels (`11.5%`) on
  the left and 5,109 pixels (`13.7%`) on the right. This is not a safe
  protected-crown result: A89 roots visibly occlude and protrude beyond A88.
- A89 supplied no front or three-quarter S05 renders. Exact preservation of
  the A87 front contour, A88 front closure, and crown seam is therefore
  **unverified**, not passed.

## Verdict

**Reject A89 S05. It is not a safe provisional lower-rear module.** Preserve
A87 S04 and A88 S07; reject only the three A89 panels.

S05 closes literal holes, but its ownership solution fails three decisive
reference gates: the roots reverse the required crown depth order, the panels
read as one broad shield with edge-on side fins, and a median `0.148 Wh` bald
lower-rear band remains where the references carry hair to the collar.

The next bounded correction should keep A87/A88 frozen and replace or locally
re-seat only the A89 pieces. It must put roots behind the visible A88 rim,
extend staggered independent free edges by roughly the measured
`0.10..0.20 Wh` where the receiver remains exposed, and give the side/rear
leaves enough surface-relative curvature to present fabric faces in profile.
Render rear plus both profiles first. Before any provisional promotion, add
front and both three-quarter fixed views to prove the protected A87/A88 pixels
remain intact.
