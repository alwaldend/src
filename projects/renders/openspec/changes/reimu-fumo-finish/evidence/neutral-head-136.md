# Coupled neutral head construction 136

## Decision before authoring

**Proceed with one coordinated head, hair and bow construction test.** The
previous image reviews identify a bulbous face, helmet-like cap, stacked rear
hair and thin bow panels. These parts share a support and attachment problem.
Study 127 preserved the anterior cap while replacing posterior hair; this
cycle replaces the support and its fitted parts together. Keep its rejected
lower body unchanged as context. Do not continue the closed garment family.

Independent reference review disproves the initial assumption that the whole
head needs substantial depth reduction. Retained 127 has approximately
132.35 mm outer head/hair width and 99.00 mm cushion depth: 0.748 head widths,
already within the recorded 0.66–0.82 band. Using the smaller bare cushion
width as the denominator would falsely suggest excessive depth. The stronger
hypothesis is a broad low-curvature front panel with rounding concentrated
near its perimeter, a fitted fabric cap and surface-relative attachments.

Rebuild a closed front/back cushion with a gusset and preserve its projected
outline within 1 mm. Retain its overall depth rather than scaling all Y
coordinates. Refit existing diagnostic face graphics in depth only, preserving
their canonical X/Z layout and layer offsets. Replace the hair with a fitted
continuous crown/fringe/underlay, a broad swept rear flap and rooted cheek
locks; fit the red/white bands to their new receiving surfaces. Construct
the bow as folded fabric with a real return through depth, tied against the
new crown. It must not become another inflated plate or a thin card with a
token rolled edge.

Canonical front controls variant, scale, graphic placement and spread bow
pose. Turn frames 12, 19 and 25 control depth, layer order and rear silhouette.
Physical photographs clarify construction only where the variant agrees.
The turntable's folded bow pose does not override the canonical spread pose.

| Feature          | Bounded target and uncertainty                                                                                                                 |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| Main crown depth | Approximately 0.75–0.82 outer head widths; observed uncertainty about 0.06 widths                                                              |
| Face panel       | Central bulge at most about 0.02 widths, edge roll concentrated in a 0.08–0.12 width band; construction hypothesis, not a measured depth field |
| Cap edge         | Apparent exposed thickness about 0.01–0.025 widths; avoid an inflated perimeter                                                                |
| Rear reach       | Underlay approximately 0.08–0.16 widths below chin; swept flap tip 1.02–1.13 widths below crown, with at least 0.06 widths uncertainty         |
| Bow root         | Visible root gap below 0.01 widths; reject a gap or accidental tangency above 0.02 widths                                                      |
| Bow front        | Total span 2.038 widths ±0.05; top about 0.22 widths above crown; center binding largely crown-occluded                                        |

The strongest risks are replacing a balloon with a rigid puck, preserving a
helmet despite flattening the face, making the rear flap an isolated blade,
and burying bow fabric to hide attachment failures. A new face surface must
remain softly filled in side and three-quarter views. Broad bow folds must
not break canonical extrema or create intersections. No materials, fibers,
microdetail, camera changes or body-gap objective belongs to this cycle.

## Execution and review boundary

Root is the sole native writer using pinned Blender 5.2.1 build `9e2066aef7ef`,
background factory startup, disabled automatic scripts and four threads.
Three isolated source authors provide head, hair and bow modules; root
executes them sequentially so every part uses the actual receiving geometry.
Protected input is 127, SHA-256
`1834951f4efa931238a06e560153f1f4d185f47fe0aa2e89a30612e3150c5079`.
All model files, scripts, logs and renders stay under ignored
`out/reimu_fumo_finish/assembly_136/`.

Require closed head topology, unchanged lower-body geometry and cameras,
depth-only graphic refitting, bounded contact/self-crossing checks, and
clean-reopened front/side renders before judging the complete packet.
Use all five fixed views, neutral clay and an uncropped presentation for
implementation-blind review. Overlay alignment remains diagnostic because
the reference camera elevation is uncalibrated. Allow one diagnosed setup
repair; do not search shape coefficients after a valid visual rejection.
No candidate or final criterion is accepted by this plan.

## First native result and setup repair

Snapshot `d76a038de1ce4e59355c879f2621c80eca76b58c6fe5271253e108e270a6c611`
clean-reopens, but the first front/side packet visibly fails. All 38 protected
rendered meshes, cameras and face-graphic X/Z layouts remain unchanged.
The new head is closed and connected, with Euler characteristic 2, positive
volume, zero boundary/nonmanifold edges and zero audited self-crossings.
Projected outline error is 0.103 mm and its extrema are unchanged. Its depth
remains 99.003 mm; sampled graphic offset error is below 0.6 microns.

Hair fitting fails in specific ways. The swept flap intended to stop near
Z 80 mm instead reaches 46.655 mm after a 33.93 mm support correction. A signed
nearest-normal test against a thin closed cap misclassifies unsupported free
points beyond its hem. A front lock ray can also pass through the cap opening
and hit the rear sheet. The tiny crown fan and narrow lock tips fold under
normal thickness; at a lock tip the authored curvature radius is about equal
to half the shell thickness. Band ray-miss fallbacks move X/Z and reverse
surface order. Independent triangle checks find 81 cap, 177/139 lock, six flap
and 87/138 upper-band self-crossing pairs; these are not just quad-audit noise.

The bow's hard nearest-projection boundary reverses one root material-width
rail by 1.190 mm. Even-offset thickness then amplifies the folded geometry
into long spikes: evaluated total bow span is 462.82 mm, although tails and
their trim still span 268.30 mm. Disabling that amplification alone would
leave the invalid folded midsurface. The front and side renders expose the
spikes and warped attachments before any presentation or acceptance claim.
Read-only native inspection confirms that the raw loop still lies between
X 5–88 mm per side and Z 175–248 mm. With a requested 0.72 mm shell,
evaluated loop vertices move up to 165.71/261.96 mm from their nearest raw
vertices; corresponding tail displacements reach 41.31/41.48 mm. The protected
snapshot remains unchanged. These extreme movements establish a thickness
amplification fault rather than a reference or camera discrepancy.

**Proceed with the one diagnosed setup repair.** Preserve the first modules
and native snapshot. Hair must use the intended receiver branch and attachment
domain, never project unsupported free edges through a nearest-normal sign
assumption, and keep pole/tip closure compatible with fixed fabric thickness.
Band fitting must preserve surface order and reject uncovered samples. Bow
roots need a shared continuous crown chart with monotone material width and
bounded normal thickness. Preserve outer landmarks, full return loops and
nominal thickness; no global shape search follows. Rebuild from protected 127
with the unchanged new head, then repeat exact native and pixel checks.

## Repaired result

**Reject; keep 127 as the continuation source.** The repaired native file is
`fd6e960501fe1b5ee65031fd413433e7079ad568235581652c4be64d9cd0e43e`.
Head module `2e735a40dbae06faef5ac0a0fa5b566dfbda95431409643da63a3bb679970853`
is unchanged. Repaired hair is
`680e1f9b84f5be24467c2553d0a582dcb4f09caf1e4e365a21385ffc892b5270`;
repaired bow is
`49a9efc324c0e1186e7dc3b7896f067c46fde79869c4518c641a2d58a10c3f82`.
All 38 protected meshes, cameras and graphic X/Z layouts remain unchanged.
Source 127's bytes are unchanged. The saved result has five fixed views,
clay front/side, diagnostic comparison/overlay and an uncropped presentation.
All five fixed image hashes and the presentation hash were independently
checked against their exact-file receipts.

The bow's extreme spikes are gone. Its evaluated span is 268.267 mm, and
explicit shell offsets are bounded to 0.36 mm. The full returning loops remain.
However, independent checks still find 20 self-crossing pairs in the right
loop and 36/74 in the two tails. The repaired hair cap, locks and swept flap
have zero audited self-crossings, but the four white band edges have 20/16/20/16.
The cap also has 230 triangle crossings with the head. Bow/host crossings
remain; the raw audit includes intended binding overlap and other contacts
that have not passed classification. None of these results establishes a
complete contact pass. No independent audit pair limit was reached.

The underlay and swept flap now end at Z 80.37/80.20 mm rather than being
projected below their intended hems. Cap width is 129.175 mm and crown height
218.698 mm. Maximum sampled canonical fringe X/Z error is 2.778 mm; these are
bounded source-landmark comparisons, not calibrated-photo acceptance.
The locks' receiver gaps still reach 3.76/3.82 mm. White edges have uncovered
receiver samples and thickness-angle failures, with roughly 24 mm gaps.

Independent image-only reviewer `/root/model136_absolute` rejects the full
packet: likeness 6.5, silhouette 6.5, construction 5.5, identity 8, contact 8,
medium 5.5 and presentation 7. The largest visible failures are white spikes
from the hair bands, a thick rectangular bob with a separate blade-like rear
panel, the inherited exposed red seat and jagged hem, hollow rigid sleeves,
and broad flat bow tails. The contact score is a pixel judgment and does not
override the native audit. Missing texture alone does not explain the medium
failure.

After completing that absolute review, the reviewer compared 127. Removing
its visible cap seam, simplifying the rear hair arrangement and gathering the
bow center are useful directions. The face has no substantial likeness gain,
the side head is more boxlike, the white spikes are new, and the bow loses
ruffled edging while its tails become flatter. The detached rear panel already
existed in 127; it is not a newly introduced defect. No complete changed
assembly earns retention. Root agrees and does not promote the technically
valid head in isolation, since its receiving parts would require a separate
validated integration.

## Next construction prerequisite and process review

**Prove one cheek-lock-and-band component before another full assembly.**
Use the frozen 136 head as a support fixture, not as an accepted model. Author
the lock's front/back fabric surfaces from one shared boundary and attachment
seam. Derive the red and white band regions from the same surface coordinates
instead of casting rays onto successively thickened receivers. Require a
padded-fabric appearance, canonical outline, clean root contact, zero uncovered
samples and zero audited crossings in front, side and three-quarter views.
Only a passing component is eligible for integration. This closes the combined
136 substitution; it does not reopen its shape parameters or the lower body.

`REIMU-136-RECEIVER-DOMAIN`: source review identifies the remaining band fault.
Each white strip has two uncovered authored samples and 22 thickness-angle
failures. `fitted_band()` retains missed samples at its diagnostic ray origin
Y = −80 mm, directly explaining the whisker-like protrusions. More sampling
or a global thickness adjustment cannot repair this attachment representation.
The next component must share the receiving surface's coordinates and reject
uncovered samples before it can propagate into an assembly.

`REIMU-136-FAILURE-PROPAGATION`: this cycle made two native candidates, two
front/side packets, two independent native audits, one read-only setup
inspection and one full presentation packet. The latter deliberately preserves
whole-figure rejection evidence, but explicit module failures were already
known. For the next isolated component, require its native gate plus three
small renders before any full-model packet. Independent static reviews did
not establish valid projection domains or safe shell curvature; the native
checks supplied that evidence. No shared tool, skill or acceptance contract
is changed by this task-local correction.
