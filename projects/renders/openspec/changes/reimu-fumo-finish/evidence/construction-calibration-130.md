# Cloth calibration and bow construction studies 130–133

Continue the unfinished macro/construction work from protected model 127,
SHA-256 `1834951f4efa931238a06e560153f1f4d185f47fe0aa2e89a30612e3150c5079`.
No model or stage is accepted. Root remains the sole native writer using
pinned Blender 5.2.1, build `9e2066aef7ef`, four threads and new ignored outputs.

## Decision before execution

**Proceed with construction tests, conditional on their own evidence.**
The 129 ribbon's matching cross-width contours cannot make a gathered,
flared hem. A planar annular sector can supply unequal material lengths at
its attachment and free edges. Native cloth is a credible route to relax
those lengths, but 118 never settled, so a small calibration must come first.
Static constrained sheet relaxation remains an alternative if native cloth
cannot retain a known seated drape; repeating an extruded border is rejected.

The [Blender shape documentation](https://docs.blender.org/manual/en/3.4/physics/cloth/settings/shape.html)
distinguishes pinning from the rest shape. The
[cloth introduction](https://docs.blender.org/manual/en/4.4/physics/cloth/introduction.html)
describes separate spring behavior and cautions against initial intersections.
These mechanisms support the experiment; they do not supply a calibrated
textile material for the pinned runtime.

1. **Calibration 130:** retain a 36 by 65 mm rectangular cloth's known seated
   path over a rounded step, with only its first edge pinned. The rest metric
   is a flat rectangle. Use a fixed numerical profile: angular bending,
   0.20 kg/m² areal mass, tension 15, compression 1.5, shear 5, bending 0.0001,
   damping 5/5/5/0.5, air damping 0.2, 24 quality steps, six collision steps,
   0.2 mm cloth/collider clearance and 60 fps at half time scale. Run at most
   400 sequential frames or 600 seconds. Require final 20-frame maximum
   motion below 0.05 mm/frame, total 20-frame drift below 0.15 mm, 95th-percentile
   absolute edge strain below 3%, maximum below 10%, and sampled penetration
   no worse than 0.15 mm. Inspect the saved profile and contacts as well.
   One setup/contact repair is allowed; no numerical material sweep.
2. **Conditional gathered ribbon:** only after calibration passes, test one
   9 mm annular sector with inner radius 27 mm, outer radius 36 mm and angle
   8/3 radians. Its attachment/free-edge material lengths are 72/96 mm over
   a proposed 48 mm assembled span. Discrete stitch anchors leave material
   between them free to pucker; do not pin the entire edge to a shorter line.
   Keep the planar rest metric and leave the free rim unconstrained for the
   final solve. Require a narrow attachment, wider open mouths and changing
   curvature across the width before any model transfer. Use the calibrated
   profile without a parameter sweep; physical stability does not pass likeness.
3. **Bow 131:** the side-reference mismatch is primarily orientation. Rebuild
   the upper lobes as swept, folded fans across the crown and the tails as
   thin, mildly twisted two-layer strips with broad side-facing surfaces.
   Canonical front high/low outer lobe landmarks are approximately
   `(±0.72, +0.21)` and `(±0.63, −0.39)` head widths relative to the crown;
   tail extremes are approximately `(±1.00, −0.49)` and `(±0.70, −0.82)`.
   Uncertainty is 0.03–0.05 head widths, not calibrated correspondence.
   Preserve the authoritative 2.038-head-width front span. White tail trim
   follows the two longitudinal edges; the rounded terminal edge is plain red.
   These layer/orientation changes replace the upright frontal pillow family.
   Allow one candidate and one diagnosed contact/setup repair, with no global
   repositioning or silent root movement over 4 mm.

## Review protocol

Freeze source bytes and save each result separately. Clean-reopen fixtures for
front, side, elevated and top views. Bow work starts with fixed front/side,
then all five unchanged contract views and an uncropped presentation if worth
retaining. Use implementation-blind image review before showing the user.
All model acceptance criteria remain open; no fibers, final graphics or rigging
are introduced. An uncalibrated reference overlay cannot grant a landmark pass.

## Result

The first calibration hit its 600-second solver bound and exited before saving
a fixture. At frame 140 its maximum vertex step was 5.383 mm and its 20-frame
drift was 6.854 mm. It did not pass; no visual verdict exists for this unsaved
run. Preserve `assembly_130/build.log` and `progress.json` as failure evidence.

### Reopened setup decision

**Revise the calibration's contact setup once.** Independent source review
found that the pinned runtime's recorded minimum collision distance is 1 mm.
A fresh native property round trip on build `9e2066aef7ef` confirmed that all
four requested 0.2 mm cloth, self and collider margins read back as
1.000000047 mm. The initial 0.5 mm support gap, including the pinned edge,
therefore contradicts the collision envelope. This is a demonstrated setup
defect, not evidence that the material profile is wrong. A separate sharp
floor transition falls inside a mesh edge and compresses it by 12.449% before
simulation.

The single permitted repair uses explicit 1 mm margins, a 2.25 mm initial
normal gap, and mesh vertices at the path's crease boundaries. Keep the
36 by 65 mm flat rest metric, mass, spring settings and numerical gates.
Read back and record actual collision settings before simulation. Save a
diagnostic on the wall-time bound without granting a numerical pass. The
annular ribbon remains conditional on numerical and visual calibration pass.

The corrected fixture also **fails**. It saved a diagnostic snapshot,
SHA-256 `cd0a73a4b949aff7cae8d8659676d28fe8acda3bfee7c3a44a9b8700cb36212d`,
after frame 113 and 606.45 seconds, with `wall_time_limit` as its stop reason.
Its final 20-frame maximum motion is 7.713 mm/frame and drift is 10.008 mm.
The 95th-percentile absolute edge strain is 81.48%, maximum 144.91%; the
pinned edge remains exact. Sampled geometric contact is outside the support
and floor, but that cannot override the failed stability and rest-metric
tests. The initial 0.2409 mm minimum collision reserve and 0.02150% maximum
edge strain passed their setup checks. Thus correcting initialization is
necessary but insufficient for this numerical profile. No mass, stiffness or
scale sweep was performed. The authored annular-ribbon draft 132 remains
unexecuted and must not be treated as a calibrated construction method.
All four fixture views were clean-reopened. Side and elevated pixels show
the cloth collapsed into a knotted bunch beside the pinned edge instead of
retaining the known seated path over the rounded support. Root rejects the
visual result as well. A future solver study needs a separately justified
numerical strategy; repeating this material profile is closed for this cycle.

### Reopened bow decision before native execution

**Revise the route before rendering.** Static review found actual tail/fan
crossings in the first draft. A bounded local layer repair was written, but
the draft's tail reached Y=112.3 mm while the rear hair ends near Y=80.5 mm.
Exact turn frames 12/25 put the hanging tail forward of that hair tip. This
directional contradiction rejects the initial route before native execution;
the depth correction is a new design, not a silent extension of the 4 mm
contact-repair budget.

Carry the concealed tail roots across the crown before descending outside
the head. Orient their broad width in depth with the inner/lower edge farther
back, while retaining the terminal frontal corners. Bring the fan's lower
rail forward coherently. Check the complete strip edges, not just its
centerline, against head and hair; keep tail depth forward of the rear tip.
Retain the local 4 mm correction cap, welded binding and evaluated tail/fan
triangle audit. Source checks cannot replace the fixed front/side pixel test.

### Bow 131 native and pixel rejection

The native shell audit found 50 tail/fan crossing triangle pairs at the outer
rim. The builder returned a plural `defects` field while the wrapper checked
only singular `defect`; it saved a rejected setup with the old bow visible and
the replacement hidden. Preserve that file as diagnostic evidence only:
`f98c447b7ed2b0a90c721cbff46a0606bb79715f5b98aa80f8413554d9e3e415`.
The wrapper now recognizes both fields. The first fast renders show the old
visible bow and are not evidence of the new geometry. Separate diagnostic
front/side renders explicitly reveal the new geometry without changing the
saved bytes; their receipts record that visibility difference.

Root and independent image-only reviewer `/root/bow131_pixels` reject the
replacement shape before any further contact repair. It has a triangular
side tail, steep upper triangle, exposed rigid center block, hard layer
junction and sharp corners. The presumed terminal width vector has a 40 mm
Z difference over 60 mm depth, forcing a 34-degree side slope. Likewise the
fan's straight outer rail forces a steep side hem. Fixing intersections would
not correct these manufactured-shape failures. No 131 geometry is retained.

### Reference correspondence correction and next hypothesis

**Revise the material-point correspondence.** Comparing the canonical front
photo with near-frontal turntable frames 03–05 shows different cloth poses:
visually estimated tail span is about 2.0 head widths in the photo versus
1.35 in frame 04, while upper-loop spans remain about 1.5 and 1.45. The photo
spreads the tails and presents upright fans; the turntable curls the fans and
hangs the tails closer to the head. This supports different posing with
moderate-to-high confidence, not a different character variant. Camera
uncertainty remains, and these estimates are not calibrated measurements.

Study 131 wrongly treated projected extrema across those sources as matching
material corners. The next bounded construction study, 133, uses the canonical
photo for the intended bow pose, while the turntable controls cloth layering,
attachment and folded construction. It does not silently relax the review
contract or grant any silhouette pass. Use a near-constant-width rectangular
ribbon, a short gathered root, rounded terminal hem and separate longitudinal
bending, rather than a strip widening throughout its length. Give the upper
lobes actual folded returns rather than one straight ruled edge. Preserve
127 as the source, use one coherent candidate, and inspect fixed front/side
pixels before additional repair or full-model review. Every final criterion
remains open.

Study 133's first static draft also stopped before native execution: a 74 mm
ribbon around an 8.4–8.5 mm centerline radius inverted its inner rail. Adjacent
surface-normal dot products reached −0.99595 outside the gathered root.
Preserve the draft, then use one disclosed concealed-construction correction:
each short gathered tail attaches beneath its upper lobe, while the central
lobes remain seated on the measured crown surface. The visible references
establish the overlapping layers but not a continuous broad strip routed from
the central knot around the crown. The corrected tail follows an almost
straight hanging run with a fixed width frame. Audit its attachment as well
as the exposed cloth; no broad root exclusion may hide crossings.

The mid-panel attachment still crossed the lobe below its seam, so the final
static attachment follows the actual lower boundary at fan UV `(0.75, 0)`.
Its native diagnostic is
`b9dc15946378419e88d54fa3cc10e36a15e03f291a8b8f624e332f4ad468f1b9`.
The crown ray hits source hair at `(0, 0.018, 0.2215416)` m. The full bow span
is 267.34 mm, but 1,568 evaluated tail/fan triangle pairs intersect. The
wrapper explicitly saved a rejected diagnostic with new geometry revealed;
the original 127 bytes remain unchanged. A first wrapper run failed while
changing visibility through a live collection iterator. The corrected wrapper
copies that collection before changing visibility; model geometry is unchanged.

Front/side pixels expose tangled tail roots and the existing rigid upper
profile. Root allows one final, explicit construction correction before the
full review: retain the fans and crown attachment, replace the unobserved
accordion tail roots with near-frontal rectangular sheets in the canonical
spread pose, and fit them behind the actual evaluated fan/fold rear envelope.
The canonical projected transverse edge is approximately 44 by 40 mm, so
use about 60 mm material width; the earlier 74 mm choice imported a width
inference from the differently posed turntable. Keep X/Z outline targets,
use a roughly 10 mm boundary transition and measure full shell clearance.
This is an orientation/attachment correction, not an extension of 131's old
4 mm local repair budget. A narrower side profile is expected in this pose;
it does not grant a pass against the unresolved cross-view correspondence.
Run one fixed construction, then stop for native and independent image review.

The first surface-fit snapshot,
`50857caa1a48dd66dc69173ad1c866cdd3b228409504fa0ae59e3c24c2b38e77`,
has zero evaluated tail/fan crossings and a 268.98 mm full span. Two concealed
inner-root samples still penetrate the source head by 10.36 and 6.92 mm.
Complete the same fit by including the actual head/hair rear envelope and
conservatively combining both sides; X/Z, pose and clearance settings stay
unchanged. Fixed side pixels also expose the center binding's empty lumen.
Represent its compressed cloth contents with a small closed root bundle
inside that intentional sewn-binding region. This fills the visible opening
without changing the outer lobe or tail outlines; the tail/fan audit remains
separate from the intentional root bundle.

### Completed bow study and retention verdict

The completed contact snapshot is
`dd7015d893f349a5778d1ff4e1bdbca57f1ab2cb9545e0d02fe6460337d27424`.
Its full bow span is 268.9785 mm. The evaluated tail/fan shell audit reports
zero crossing pairs; it excludes witnesses enclosed by the center binding,
although this run reports none there either. All 100 sampled head/hair
contacts have positive signed-normal distances. These bounded checks do not
prove universal clearance. The tail surface fit moves points by up to
34.085 mm in Y; the separate zero local-repair value does not mean the
construction has no displacement. Source 127 remains byte-identical.

Five fixed views, front/side clay views and an uncropped presentation were
clean-reopened from the saved file. Independent implementation-blind review
rejects the full model: likeness 6, silhouette 7, construction 5, identity 7,
contact/occlusion 7, medium 5 and presentation 9 out of 10. Comparing the
canonical and turntable poses caused the reviewer to retract the claim that
upright lobes and outward tails are inherently wrong. The canonical front
supports that broad arrangement; it cannot establish side or rear depth.
The supported failures are straight panel edges, an exposed regular binding,
limited broad cloth curvature and thin, weakly gathered trim. The rigid skirt
over the exposed red base, molded hair and hollow sleeve funnels remain the
leading whole-model failures.

Relative pixel review found a concrete trim regression: the left upper white
border disappeared. Its normal was computed after reflecting the surface,
which reversed the cross-product orientation and sent the offset inward.
The correction builds the ruffle on the original surface, then reflects the
completed geometry. Core red panels, fitted envelope, pose and cameras remain
unchanged. The final trim snapshot is
`35671fc6750bd0cc7a4e7c16d2051859be3848fb495d6bc39bb59632606c591b`,
module SHA-256
`8f489786c9d3ca37f3fe0ce9f1970724e25b2cce0639646e408cfd1062542330`.
It has a 268.9941 mm full span, unchanged core bounds, zero evaluated
tail/fan crossings and the same positive sampled contact results. Fresh
front/side pixels restore the white border without a visible core-shape
regression. The final file has its own five fixed views, two clay views and
768-pixel presentation. Independent relative review confirms the correction
in front, side and presentation images, while retaining the same construction
and retention verdict. This corrects an implementation defect, not the broader
cloth construction failures.

**Retain 127 as the continuation source; preserve 133 as a separate unfinished
study.** Its upper-lobe silhouette and center-converging folds are useful,
but its complete bow is not established as an improvement. No geometry from
130 or 131 is integrated; 132 remains unexecuted. The next material hypothesis
is a small static sheet-relaxation fixture with explicit rest-length and seam
constraints, followed by pixel review before any garment transfer. That
hypothesis is unexecuted; it does not authorize another sweep of the failed
cloth profile. All full-model criteria remain open.

## Process review

The independent read-only ergonomics review covered 14 distinct files and
one bounded output inventory. It identified three task-local controls:

- `FUMO130_COLLISION_ROUNDTRIP`: reject unsupported requested settings and
  invalid initial contact before simulation. The corrected setup reads back
  actual values and measures initial collision reserve; this does not prove
  that contact initialization is the only cause of solver instability.
- `FUMO131_REFERENCE_CORRESPONDENCE`: distinguish projected extrema from
  identified material points and record the pose associated with each source.
  Do not infer shared material corners from independent silhouette outlines.
- `FUMO131_DEFECT_SCHEMA_VISIBILITY`: a failed builder must not produce an
  ordinary candidate receipt or silently render the old visible component.
  Use an explicit builder success result and verify replacement visibility;
  keep deliberately revealed failure snapshots labeled as diagnostics.

The final relative pixel review also exposed the mirrored ruffle's reversed
normal. Reflecting a completed offset surface preserves the intended trim
side; source checks and collision counts alone had missed the visible border
regression. The correction was verified from the exact saved file.

The old delivery launcher also had missing runfiles. Its first refresh hit a
92-second deadline during cold analysis; a longer bounded refresh completed
in 185 seconds. The same-session generated Blender runfiles wrapper, built by
the successful diagnostic-render command, started the corrected calibration
while delivery analysis occupied the Bazel server. Its path/hash are recorded
in ignored scratch and are not durable launch instructions. No host profile,
shared tool, canonical skill, or acceptance contract changed.
