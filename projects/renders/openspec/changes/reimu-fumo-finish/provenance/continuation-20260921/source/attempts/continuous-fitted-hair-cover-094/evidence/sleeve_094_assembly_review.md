# Sleeve 094 assembly decision review

Observed 2026-09-06 at approximately 01:50 UTC. Recommended verdict: **revise
the arm, attachment, and hanging sleeve as one bounded assembly**. The
coordinator owns the final decision. This review writes only this scratch
report; it does not implement or approve a candidate.

## Decision and its limits

Keeping both the current evaluated arm placement and its low/outboard,
fully clear sleeve root immutable is a practical modeling dead end for the
referenced sleeve. Another cuff-only placement or depth adjustment does not
address the dominant construction mismatch. This is not a proof that an
ellipsoid-shaped stuffed arm is inherently unsuitable: the arm may remain a
simple continuous stuffed pod after its extent and direction change.

The strongest case for preservation is real. An inexpensive continuous arm is
plausible for a plush, and the corrected cuff can admit its actual section
without a rim intersection. Exact shoulder geometry is occluded, and the
front-to-model mapping has not been registered. Those facts argue against
committing to precise replacement coordinates. They do not establish that
the current exposed hand and cloth relationship is reference-faithful.

The decisive evidence is the combination of visible construction and the
094 scalar result: the reference has a hand/arm bulb high within an opening
with substantial cloth below it, while the present construction exposes the
distal arm outside the cuff. Under the provisional corrected cuff, clearing
the rim still leaves the section centroid at 38.927% of cuff height and the
actual arm extending 11.226184 mm beyond its plane. These are different
problems from a cuff that merely intersects the arm.

## Evidence: trace, construction, and uncertainty

I inspected the six images below before reading the three diagnosis reports.
The assignment already supplied arm and measurement context, so this is an
independent decision critique, not an implementation-blind approval review.
No source mesh, Blender process, builder, or new parameter sweep was used.

| Evidence class | What it supports | What it does not establish |
| --- | --- | --- |
| Canonical front pixels and 093 trace | Outer white cuff ends near (712,649) and (647,738); long diagonal white outer edge; sleeve emerges beneath hair and joins the torso region. The front does not show the pronounced hand protrusions present in087. | A complete cuff plane in depth, a hidden throat ring, or exact world coordinates. The red dashed trim is inset and is not the outer cuff edge. |
| Reference side frames10 and12 | A broad oval aperture, visible stuffed arm toward its upper interior, and an extended lower cloth pocket. The attachment is high, under the hair/shoulder region; the sleeve reads as suspended fabric. | An exact orthographic side profile, a measured hand-height fraction, or the full hidden seam. The two frames have different azimuths and are not calibrated against087. |
| Current087 front/side/three-quarter pixels | Hand protrusions are conspicuous in front; the side aperture is relatively short and rounded rectangular, with a thick uniform rim. The opening/hand relationship differs from the reference's hanging sleeve. | A numerical error under matched cameras, or acceptance of unrelated body/hair/bow components. |
| 092 and094 evaluated-arm analysis | A rearward cuff-depth assumption caused actual planar clipping. Recentring the corrected094 cuff clears its rim at the arm section but leaves the arm low and distally extended. | Full 3D sleeve-wall, liner, skirt, trim, seam, or clearance acceptance. These are analyses of a provisional analytic cuff, not a completed sleeve. |

The supplied arm runs along P0=(±16,0,69) to P1=(±63,-12,43) mm.
The reported current throat center is |X|=38 mm, Z=56.829788 mm. The
provisional reference attachment witnesses fall around |X|=24–30 mm,
Z=66–81 mm. This discrepancy is a warning about the assembly, not permission
to average those witnesses into a new root. The upper witness is an obscured
loop/attachment region; the lower witness is where that loop merges into the
sleeve. Neither is a measured physical throat endpoint.

## Can camera mismatch explain it?

Camera mismatch is a credible contributor to exact cuff dimensions, shoulder
height, apparent opening width, and hand occlusion. The mapping
0.1165/368 m per pixel, center485, ground845 is explicitly provisional.
The source front photograph and current front render are not demonstrated to
share projection, elevation, or framing. The side reference views are oblique,
so a fraction derived from them must not be imposed on a world-Z cuff.

Camera mismatch is not presently a sufficient explanation for the whole
failure. Uniform image scale or crop cannot change the hand's relative
placement within a sleeve. View angle can change that appearance, but the
canonical frontal symmetry and the two supporting side views jointly favor
a shorter arm inside a longer hanging sleeve. More importantly, a camera
change cannot eliminate the reported 11.226184 mm geometric extension
through the chosen plane. A different, correctly registered cuff plane might
change that result; the current analysis does not prove that every possible
reference-consistent plane fails.

Before choosing dimensions, register the front using stable head/face
landmarks plus the existing body/ground relationship, then judge sleeve
contours without independently rescaling them. Match side azimuth and
elevation sufficiently to compare aperture height, hand occlusion, and the
shoulder connection. Keep the registered cameras fixed for subsequent
candidates. Do not use camera adjustment to conceal geometric crossings.

## Smallest credible assembly change

Treat this as one coherent construction correction, bounded to the two
arms and their sleeves/attachments:

1. Keep a simple continuous stuffed arm, but release the inherited distal
   position and downward axis. Shorten/retract and lift the distal bulb until
   it sits inside the upper cuff region and leaves visibly more empty cloth
   below. Preserve the buried proximal body attachment where feasible; do not
   preserve a long protruding tip merely because the current sleeve can be
   enlarged around it. A slight bend or varying radius is an option if a
   straight pod cannot seat cleanly, not a mandate to build an anatomical elbow.
2. Seat a shallow cloth root around the proximal arm at the shoulder beneath
   the hair. Determine its hidden extent from continuous attachment and
   side-view evidence, while retaining the visible loop as a separate witness.
   A root forced wholly outboard to achieve universal clearance is the wrong
   construction if it prevents the sleeve from attaching where the reference
   requires. Intentional sewn contact must be explicit and clean; this does
   not excuse accidental arm, garment, or hair crossings or waive existing
   clearance requirements on non-attaching surfaces.
3. Hang the sleeve below that high attachment as thin, softly padded
   front/back cloth panels, with a gently curved lower seam and a broad
   oblique oval cuff. Let its lower area extend beyond the short stuffed arm.
   The front panel/lip should provide the reference's frontal hand occlusion;
   the side opening should expose the upper arm and empty lower pocket.
   Avoid a thick, uniformly inflated tube or a rim enlarged solely to admit
   the old arm. Place decorative dashes only after these surfaces pass.

This preserves the simple plush vocabulary and limits the change to one
subsystem. It has more work and attachment/rig regression risk than moving
the cuff alone, but directly changes both failed relationships. The existing
body, skirt, and hair provide boundary conditions; do not move them to make
the sleeve fit. Coordinate the shoulder's hidden contact with the separate
hair-cover work rather than letting either subsystem assume the other will
hide a gap.

## Alternatives and reopening condition

| Option | Assessment |
| --- | --- |
| Keep the assembly unchanged | Preserves known pixels and existing attachment data, but preserves the visible rejection. Not a route to approval. |
| Preserve arm/root and tune only cuff depth or span | Cheapest reversible change;094 shows why clearing the planar section is insufficient. Likely repeats the low/protruding-hand failure or over-enlarges the cloth. |
| Camera-only correction | Necessary diagnostic work for exact measurements. Not a supported complete remedy for attachment and physical protrusion. |
| Coordinated short arm, high attached root, hanging padded sleeve | Recommended bounded revision. It follows the observed cloth/arm relationship while retaining simple stuffed forms. |

Reopen preservation only if a registered multi-view comparison demonstrates
that the current arm can remain high inside a reference-consistent cuff,
the front hand is correctly occluded, and a clean attached root is achievable
without distorting the traced panel silhouette. A zero planar intersection
alone does not meet that condition. No exact hand fraction, new clearance
exception, or hidden-root coordinate is asserted by this review.

## Input identity

Git observation: linked worktree `t3code-a13ca48d`, branch
`t3code/continue-fumo-desktop-use`, HEAD
`c7601f0fc80e0a94b585910459639ffccbdbdbd4`. The report directory is Git-ignored.
Image and report SHA-256 values were observed read-only; input files were
read in full, with no content truncation in the selected reports.

| Input | SHA-256 |
| --- | --- |
| `projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png` | `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c` |
| `side_062_frame_10.png` | `37c1e2866fdbe97ce79a0b5bdddf216a7f151f8e61d2a94c7286c22eaf41fd07` |
| `side_062_frame_12.png` | `4ebed9233acb351f2a8093496eb3e544771f1ee01e503cec84a1ce86f22b7365` |
| `chin_087_all_review/front.png` | `18610642926ca2be36a12ff0b843ae7b7d69fd2bb0a437df93634374babeeca9` |
| `chin_087_all_review/side.png` | `1e1e6d50be55d2d3528106bc39650932305d95e4bf31c321be71776552aff987` |
| `chin_087_all_review/three_quarter.png` | `dca74e014f7abf74918ec6d3329b0f4d479e4c1633bc3c39670a0dcac9dda197` |
| `sleeve_093_reference_landmarks.md` | `4eea0441bad88253de71343b2dc2f5f2f4d755c06c127c46bad15744a6f10f88` |
| `sleeve_092_cuff_diagnosis.md` | `633c0fbeb2ac40b862c3deaf45de07eb714b946e9fb4050b3796802f84bc500f` |
| `sleeve_094_cuff_feasibility.md` | `45c68c849c1a2090f620c6ce7d571d4f2d3daebb25fcc49bbc12e5d98a6f03c0` |

Paths other than the canonical front are relative to this report. The
decision-review skill structured the alternative test; the reference-fidelity
skill prevented treating provisional registration, planar clearance, or
relative improvement as a visual pass.
