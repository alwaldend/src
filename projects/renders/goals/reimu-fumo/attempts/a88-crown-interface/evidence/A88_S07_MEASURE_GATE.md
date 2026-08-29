# A88 S07 narrow delta and six-view interface gate

## Scope

This is a render-only gate on S07 as a **provisional crown module**. The
explicitly separate lower-rear hair owner is not required here. I compared
S07 front/front-ID pixels with the exact A87 S04 parent and A88 S04, then
measured only the crown coverage and receiving boundary in S07's profiles and
rear. No Blend or tracked file was changed.

- S07 source SHA-256:
  `0a30e2af3142081648bb3137ad75d6d1cc73de55e9f830f85bad1f85e92c8788`.
- S07 six-view manifest SHA-256:
  `e500a58389af2223c0ad471f3936081b4bbbb4cf3bf87e6de79530bd2fa38be1`.
- Normalization: fixed-render head width `Wh = 371 px`.
- Bright exposed head pixels: all beauty RGB channels greater than 245.
- Component deltas: nearest declared color in the flat front ID pass.

## Front closure

| Measurement | A87 S04 parent | A88 S04 | A88 S07 |
| --- | ---: | ---: | ---: |
| Bright pixels in front crown ROI | 14,939 | 50 | **0** |
| Bright center-seam pixels | 4,360 | 50; 29 px maximum run | **0** |
| Bright image-left root pixels | 49 | 49 | **0** |
| Bright image-right root pixels | 90 | 90 | **0** |

S07 fully closes the top slit and both front root notches. Its center seam is
a color/value transition between the two crown panels, not exposed receiver.
This matches the controlling front references' continuous brown crown and
passes the front-coverage and seam-closure gates.

## Exact protected A87 lower-form delta

The lower free-edge zone begins at `y=220`, below all crown/root occlusion.
Against exact A87 S04, S07 has:

- fringe: **0 XOR pixels**, IoU 1.0;
- left-ID cheek lock: **0 XOR pixels**, IoU 1.0; and
- right-ID cheek lock: **0 XOR pixels**, IoU 1.0.

The per-row extrema are identical in all three masks. Differences above
`y=220` are confined to intended upper-root occlusion by the new crown: whole
mask IoU is 0.9900 for the fringe, 0.9578 for one lock, and 0.9695 for the
other. Therefore S07 does not change the protected A87 lower fringe contour,
lock lengths, tips, or lateral placement.

## Profile and rear receiving interface

Hair was segmented from the beauty renders by its brown/red chroma
(`R-G > 25`, `R-B > 25`) within the upper 260 px. The purpose is not to demand
the missing lower hair, but to test whether S07 supplies continuous coverage
and an overlap boundary that the next owner can use.

| View | Measured S07 crown/interface |
| --- | --- |
| Left profile | Crown is continuous across the central `x=145..350` span (`0.555 Wh`). Its lower envelope is `y=198..224`, a 26 px (`0.070 Wh`) gradual variation with median `y=223`; no pale hole crosses the crown. |
| Right profile | Same `206 px` central span. Lower envelope is `y=208..220`, a 12 px (`0.032 Wh`) variation with median `y=218`; no pale hole crosses the crown. |
| Rear | One continuous crown-rim component spans `x=58..453` (`1.067 Wh`). At the center it supplies about 10 px (`0.027 Wh`) of visible overlap depth; the curved side returns deepen progressively, reaching about 34--59 px near the lateral extremes. |

The profiles still look plain because the lower owner is absent, but their
boundary is long, continuous, and gently curved rather than the disconnected
receiver shelf in S04. The rear view is intentionally open below the crown;
its continuous arch-shaped rim provides a measurable tuck-under edge for a
separate rear sheet. The next owner must overlap beneath this rim rather than
butt against it, and its integration render must verify that the narrow
`0.027 Wh` center allowance remains covered under deformation.

The canonical turn and physical side support this ownership: the complete
rear silhouette is carried by overlapping lower fabric leaves, while the
crown remains continuous above them. They do not require S07 to contain those
leaves now.

## Verdict

**Pass S07 as the safe provisional crown-module survivor.** This is not a
standalone-sculpt approval and does not pass the complete hair subsystem.

S07 satisfies the bounded gates that S04 failed: zero exposed front crown,
zero pale top seam, zero front root gaps, exact preservation of every A87
lower free edge, continuous profile coverage, and a continuous rear receiving
rim. The only narrow contract is the rear-center overlap allowance
(`0.027 Wh`): freeze S07 and require the next lower-rear panel to tuck under it
with a rendered no-gap regression test from both profiles and rear. Do not
reshape the accepted front or lower A87 components merely to make this partial
module look complete.
