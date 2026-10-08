# 075 preserve free hair shapes while fitting their roots to074 face

Exact source074f3f5573dc6b39bd9ed586d72cc9e722968847ece96073c964b3c7642ae861474;
shape donor071129180df3be51e140d53f3b0bc4119327fb8993adcd9b079fad9807d396d0abc.
New face_075_candidate.blend only. Root sole model/goal writer.074 rejected
after38.897s fixed triple: broadened face but cage distorted free hair into
flared lobes. Preserve074 head and attached head pile only as provisional;
no stage/module acceptance inherited. All body/collar/graphics remain074 exact.

Root decision PROCEED with attachment-only correction. Restore four071 side
hair meshes and their four pile curves exactly. Restore donor front locks,
two front pile curves and four tie/trim meshes before refitting. Their free
parts retain their native shape: per side translateX outward5.5mm, noZchange.
Fit oneYtranslation to the frozen074 head: for evaluated donor lock vertices
atZ83–117.5mm whose translatedXZ lies over the head, cast front-to-back head
rays and take minimum(head_frontY-0.6mm-donorY). This is a single deterministic
clearance fit, not a parameter sweep. Require translation in[-30mm,0].

For vertices/curve points belowZ118mm apply only that rigid translation.
Above118mm blend to the074 root cage with smoothstep over3.5mm; by121.5mm
use the receiver-relative anchor. The anchor interpolates original/new head
front charts and preserves original depth offset via front/back fade. Ties
end below117.85mm and therefore move rigidly with the free locks. No receiver
field is applied to the free ends; no contour-amplitude tuning is allowed.

Pre-save: head/graphics/head pile and all other non-target fingerprints exact;
four side panels/piles exact donor; below118mm front locks/strands/ties differ
from donor only by the one translation. Check evaluated lock manifoldness,
signed volume and self-overlap against donor, report head intersection and
sampled free-vs-root penetration separately. No new head/collar issue because
head/collar remain exact. Mouth-level visible width should remain0.57–0.64Wh.
Stop on a new critical free-panel self-intersection or missed numeric fit.
One causal pre-save API/data repair only after a contemporaneous checkpoint;
no post-render variant. Pinned5.2.1/build9e2066aef7ef,4threads, task-local paths.

Save once, clean-open fixed512 front/side/3q. Reject flared/buckled hair,
exposed floating roots, face-mask ledge or remaining recessed lower panel.
If plausible, render matched close and get implementation-blind review;
full fixed packet is required before retention. Keep whole goal open until
absolute gates and remaining whole-asset criteria pass. Use compact goal plan.
