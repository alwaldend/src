# Hair construction review for 068

Observed 2026-09-05 18:29 UTC. Advisory review for the coordinator; no Blender
process, model, source, goal, or Git state was changed. The only output is this
ignored report. The task worktree is `t3code-a13ca48d`, branch
`t3code/continue-fumo-desktop-use`, observed HEAD
`c7601f0fc80e0a94b585910459639ffccbdbdbd4`.

## Recommendation

**Revise the proposal's priority, then prototype it.** Make the tapered free
rear/side hair silhouette the primary visible change. A continuous separate
fringe with a thin free edge is a sound second change, but it is unlikely to
resolve the user's “flat” complaint by itself. Do not change whole-head scale
or add uniformly inflated locks. The coordinator owns the final verdict and
acceptance.

The intended outcome is likeness to the supplied Reimu plush, including
readable fabric construction and volume. The proposed method is subordinate
to that outcome. My review scope permits inspection and this report only.

## Evidence and limitations

I inspected all seven requested images and read both prior diagnoses in full.
The physical photos control construction and material observations; the
canonical front controls graphic landmarks and front silhouette; frame 12
controls the otherwise hidden canonical side hair continuation. Physical
photos and turntable have different drapes and unknown cameras. There is no
new calibrated overlay here, so approximate pixel comparisons below do not
establish a fidelity pass.

This is not an implementation-blind approval review: the coordinator supplied
the existing representation and proposed fix before inspection. No rear
candidate or neutral sculpt was supplied. No model code or live Blender data
was inspected; the .1162 m head width, .0910 m depth, four existing panels,
and shared face/fringe surface are coordinator-reported facts.

Exact inspected sources, SHA-256:

| Source | Digest |
| --- | --- |
| `macro_067_fast_review/front.png` | `b02fc2dc7c37ec65bbbdb200506d373c96ded02ad251bf9fd85d13aa16eb412a` |
| `macro_067_fast_review/side.png` | `edad764eb07305bb4198164674ec40597ba31dec417698fa5f38fdb2b887453a` |
| `macro_067_fast_review/three_quarter.png` | `fb121ec08b77294f3cfe800d889819ffc1255220637e20165228b66b8fc9b67b` |
| `canonical_front_25cm.png` | `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c` |
| `physical_front.png` | `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51` |
| `physical_side.png` | `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d` |
| `side_062_frame_12.png` | `4ebed9233acb351f2a8093496eb3e544771f1ee01e503cec84a1ce86f22b7365` |
| `head_063_diagnosis.md` | `3663531bdade0d4e59288a4b118e086108c03c5769aaad53686806f87b1577e8` |
| `side_062_reference_diagnosis.md` | `e10d0eb96b59b5c7f4ade99516162bbc65df43d1de93dcc667c55ab4a700a812` |

Candidate and diagnosis paths are under this report's directory; reference
photos are under `projects/renders/assets/reimu_fumo/references/`.

## The two useful geometric changes

### 1. Release the lower side/rear hair into tapered folded ends

This is the strongest supported change. In 067 side, the exposed lower brown
mass near x225–310, y305–348 finishes in broad rounded U shapes. Three-quarter
shows repeated neighboring rounded strips attached to the same convex head
envelope. They read as a continuous molded bob with shallow segmentation.
Being geometrically thin does not change that visible large-scale form.

Physical side shows thin angular ends and a change of direction below the
side cap, approximately x145–231, y282–338. Frame 12 shows the most obvious
disconfirming evidence against the existing shape: the far hanging brown
piece has a free pointed continuation, reaching approximately x444,y304,
while the close cap is around x350–356. The prior diagnosis estimates that
extra projection as roughly .39 crown-to-chin H. That endpoint is hanging
hair, not head stuffing. The photographic construction is a supported upper
sheet that releases and bends, not a new padded head lobe.

The proposed surface-supported roots plus free continuation can reproduce
this. The release must change the silhouette and surface direction over a
visible lower region; merely making the current rounded hem pointier while
it hugs the head would leave the dominant failure. Use a coherent broad bend
or folded shoulder into a tapered thin end, with mild nonplanarity between
the supporting root and free tip. The same plane should not extend straight
from the crown to the endpoint.

Exactly two replacement panels is an implementation choice, not a recovered
sewing pattern. The references do not establish a complete panel count.
Two can work if their visible overlap and taper reproduce the observed
shape. Do not force two perfectly symmetric, equally sized rear triangles.
The .39H projection is specific to the turntable pose and should not be
copied as an exact world-space extrusion without a camera/drape check.

Coordinator update received after this review: preserve the four existing
object identities, taper and bend two pieces per side, with the rear piece
releasing backward and the other ending inward/down. This is compatible
with the recommendation; object count is immaterial to the pictured
construction. A proposed 1.5 mm crossfold is only a test setting. It should
support the broad direction change, not become the sole visible difference.
The proposed absolute tip coordinates cannot be approved without their
evaluated relation to the cushion and candidate pixels.

Expected visibility: large in side, clear in three-quarter, modest in front.
Front should retain the existing narrow tied locks and their face framing.
The free rear ends must remain behind these rather than creating extra
cheek spikes. Rear inspection is required because two broad replacement
pieces can reveal the cushion, cross, or leave an artificial center slit.

Primary regression risks are planar-card hair, a hard fold, unsupported roots,
a new rear spike, and widened front silhouette. Sub-mm thickness does not
prevent these failures. The acceptance evidence is a thin bent sheet outline
with seated roots and clean overlap in the existing fixed views, under the
existing landmark tolerances. This report does not add a new dimensional
target.

### 2. Give the continuous fringe a small independent cloth edge

The front brown/cream transition in 067 reads as a color boundary on one
smooth object. The proposed extraction of the actual brown front domain is
a credible way to retain the graphic outline while adding an edge and local
overlap. Canonical and physical front show a continuous forehead hair mass,
not several padded individual bangs. A single continuous sheet is the
strongest case for the proposal.

The important qualification is visibility. The head is approximately
206 ± 5 pixels wide in the 512-pixel front render. With the reported .1162 m
width, .35–.70 mm corresponds to only about .6–1.3 projected pixels. A
sub-mm free-edge lift can produce a useful contact shadow and thin edge in
three-quarter, but it cannot provide a large depth cue across the forehead.
That is a reason to keep its intended effect modest, not to exaggerate it
until it becomes armor.

Keep the layer supported at the crown/temples and close to the face, with
small local edge variation rather than a constant dark moat. Deriving the
sheet from the brown triangles is acceptable only if the visible free edge
remains smooth and intentional; inherited faceting is not a fabric seam.
Do not add separate seams along every brown triangle or a thick raised rim
around the entire cap. Avoid coplanar duplicate surfaces at attachment
transitions.

Expected visibility: subtle local separation in front, more useful along
the oblique edge in three-quarter, small near the exposed face boundary in
side. Main risks are a continuous black trench, detached pointed fringe,
faceting, and changes to the eyelid/forehead openings. Preserve the existing
front critical landmarks unless measured reference evidence supports an
intentional correction.

## Is a bigger mechanism being missed?

There is no strong pixel case for a whole-head proportion correction in this
review. The 067 front already has broadly similar overall-height/head-width,
bow-span/head-width, and sleeve-span/head-width relationships to canonical
front. This is a coarse observation with perspective uncertainty, not a
claim that every macro landmark passes. Side visibly has substantial volume,
and the prior reference diagnosis found substantial physical head depth as
well. “Flat” should not be translated directly into “increase every depth.”

The smooth, nearly uninterrupted brown crown gradient and very clean edges
are a major additional contributor to the molded or cartoon impression.
Physical front has strongly visible pile and broad irregular fabric response;
canonical front has finer but still visible pile. The render lacks those
breaks at its current display scale. The red bow also has smooth satin-like
bands of light. This is a material/presentation inference from pixels, not
proof of a particular shader error. Lighting, roughness, normal structure,
and missing pile cannot be separated by these images alone.

Thus the hair sheet changes address an actual geometric failure but are
unlikely to finish the intended-medium problem. Preserve fixed light and
camera for their geometry comparisons; once the construction tier passes,
judge surface response in a controlled close and full render. Do not deform
the whole head to manufacture texture-like shading, and do not treat pile
as a repair for the current U-ended hair geometry.

The sleeve opening and cloth at the lower skirt also remain conspicuous in
side/three-quarter, but the coordinator is already addressing the sleeve.
They are not grounds for expanding this hair review into unrelated edits.

## Alternative comparison and bounded review disposition

| Choice | Likely outcome | Main risk/cost |
| --- | --- | --- |
| Only add the small fringe edge | Better layer separation; little silhouette change | High chance the same “flat” complaint survives |
| Prioritize free tapered rear sheets, then the small continuous fringe edge | Addresses the strongest visible construction mismatch while retaining front landmarks | Requires side/rear overlap review; moderate reversible modeling work |
| Inflate all locks or increase whole-head depth | Adds volume without establishing referenced cloth construction | Likely worsens molded/balloon appearance and proportions |
| Material/pile change first | Could improve fabric impression across a large area | Leaves the identity-defining lower hair silhouette wrong |
| No geometry change | Retains the front correspondence | Leaves the strongest reference contradiction unresolved |

Advisory verdict: **revise**, in favor of the second row. Conditions that
would change it: a controlled candidate where the new rear shapes become
cards or spikes; loss of front tied-lock framing; or a calibrated comparison
showing the intended free continuation belongs only to an incompatible
variant. In those cases revise the sheet drape, not whole-head scale.

For the supplied intermediate pixels, subject recognition is yes. Provisional
whole-image scores are likeness 6/10, macro silhouette/proportions 7/10,
construction 5/10, identity features 7/10, contact/occlusion 5/10,
intended-medium 4/10, presentation 7/10. The five most visible discrepancies
are lower hair termination, smooth hair/cap fabric read, sleeve opening,
weak fringe layering, and overly regular bow/cloth surface response. The
head/hair alone fails the construction and intended-medium gate. This is
an intermediate rejection, not an approval candidate; the context and
missing-view limits above prevent a formal complete blind gate.
