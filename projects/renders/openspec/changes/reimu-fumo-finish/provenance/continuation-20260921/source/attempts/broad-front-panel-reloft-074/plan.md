# 074 replace pinched lower face with a broad shallow front panel

User: "its still receded, continue". Root sole Blender and canonical-goal
writer in verified t3code/continue-fumo-desktop-use linked worktree, HEAD
c7601f0fc80e0a94b585910459639ffccbdbdbd4. Preserve073 as comparison baseline.
Use exact071129180df3be51e140d53f3b0bc4119327fb8993adcd9b079fad9807d396d0abc
to remove072/073 chin-only deformations; all other objects match073. Save new
face_074_candidate.blend. No whole-asset acceptance inherited.

## Decision review and changed representation

The user's repeated rejection supersedes root's prior relative-improvement
description. Independent image-only face_074_reference_diagnosis.md SHA
d788888f84812ac71db79375bafc2ba469c52068b4f2fedc70ac6a6de636db3a finds exposed
mouth-level face width0.50–0.53Wh versus canonical0.57–0.62/physical0.59–0.65;
mouth-to-bottom height already0.10–0.11Wh. Primary physical-front for cushion,
canonical-front for graphic pattern; turn08/10 for cheek continuity. Oblique
visible slope bands overlap25–40deg; no exact true chin depth/tilt is established.
Root073 isolated side/3q diagnostic (no save) reveals a visible mouth-region
depression and pinched lower panel. Ray evidence separates limited left skin
width from right lock occlusion atX35mm. Source073 profile JSON
c6a7c40638e49eaa55643cabb664d16c8583a0b8cf39ddaceeed3cb235330264.

Root verdict REVISE the representation: replace the lower radial-dome panel
with a broad front plane and seam-distance padding, not another chin push.
Keeping mouth depth fixed prevented the previous edits from correcting the
whole eye-to-mouth-to-underside transition. Width-only would leave the visible
depression; isolated forward chin variants are retired. Whole-head rotation or
height change lacks reference support. No new tooling/desktop route is needed.

## Frozen construction

Keep the existing front/back triangulated graphic domains and upper head.
BelowZ104.5mm broaden the lower pattern inX using
dx=sign(X)*.23*P(abs(X)-4mm)*Wz, P(t)=0 for t<=0,t²/(8mm) for0<t<4mm,
otherwise t-2mm. Wz=smooth((104.5mm-Z)/9mm)*smooth((Z-60mm)/18mm).
This leaves eyeX/Z and mouthX/Z exact; predicts ~0.60Wh exposed width.

For each front-panel vertex, measure distance d in widenedXZ to its stitched
outer boundary. New lower frontY=-24mm-25mm*(1-exp(-d/4.5mm)); blend from
existing frontY with w=1-smooth((Z-103mm)/20mm). BackY unchanged. This replaces
the retreating radial cross-section with a shallow padded front and forward
lower rim; mouth/eye depths may change deliberately, never detach as graphics
are in the same mesh. Z/head height and upperZ>=123mm remain exact.

Apply the same surface-relative cage field to six existing hair pieces, four
tie/trim meshes and seven attached pile curves: interpolate front/back depth
and new-front delta from the original panel's XZ triangles, fade depth shift
from1 at/ahead of front to0 at back, and fade beyond panel footprint over20mm.
Use the sameX widening. Roots outside the support remain exact. Body, collar,
bow, sleeves, garment, lights, camera and shaders unchanged. This is one
head/panel/interface hypothesis, not a change to hair's own panel construction.

Pre-save: finite closed positive-volume head, no self-intersections; exact
front graphicXZ, upper head and non-target fingerprints; no new head/collar
surface intersection. Report head/hair baseline and after intersection sets,
do not label existing contacts clean from counts alone. Mouth-level front rays
must show exposed cream width0.57–0.64Wh, measured as geometry not calibrated
photo acceptance. Bound max changed head displacement30mm andXshift12mm.
One causal pre-save API/data repair allowed only after a contemporaneous
checkpoint; no post-render parameter variant. Stop before save on a critical
new geometric defect.

Pinned5.2.1/build9e2066aef7ef,4threads, factory-startup/disable-autoexec,
task-local tmp/config/cache. Save once, clean-open fixed512 front/side/3q,
then matched800px close and isolated profile only if plausible. Inspect for
mouth trough, jutting jaw, flared hair, collar crowding and graphic distortion.
Independent blind review before any retention; full fixed packet required
after a scoped fast pass. Whole asset stays open below absolute8 gates.
Use supported compact goal plan (portable registry full64).
