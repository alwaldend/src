# 086 rear triangulation without facial graphic constraints

Root PROCEED from084cb6bff1fa0e654f1720236a5a90dab56346f41e36f999dfc08fe746ea26294cb.
Pinned5.2.1/build9e2066aef7ef, same linked worktree/background root writer,
factory-startup/disable-autoexec/four threads. One build plus one causal
API/data repair if checkpointed; no post-render variant. New hair_086_candidate.blend.

Independent diagnostic990fe584b1ddf20feee1ef415b6621b9d7c2af8ea87954d6c018be739d84b3bf
proves copied front-graphic CDT makes0.18–0.48mm rear chord depressions and
large normal errors, despite correct analytic vertex positions. Back faces are
all brown; named eyebrow objects are hidden. Normal-only correction cannot
remove the actual grooves. Linear subdivision would preserve them.

Keep all6596 head vertex positions and their ordering exactly. Preserve every
front and gusset polygon/material. Replace only back triangle connectivity
with native constrained Delaunay in XZ using the existing3298 back points and
only the outer boundary loop as constraint; no facial graphic domains. This
avoids introducing points or changing the underlying head shape. A whole-back
triangulation replaces the same erroneous data dependency consistently; local
four-patch stitching adds unnecessary boundary complexity. Every original
point and boundary edge must survive, with consistent outward winding.

Re-seat only head-pile strands whose nearest original receiver triangle is
entirely on the back. Preserve each strand length/radius and every front/gusset
strand exactly; use the new back at the same XZ,10um offset, rotate original
strand direction with receiver normal. Allow at most1mm anchor displacement.
No other object, material, modifier, camera, hair layer or body change.

Pre-save: finite closed positive head,0self pairs, exact vertex/front/gusset
data, exact nontarget fingerprints, only head/head-pile changed. Repeat the
existing36 strip-center sample locations, comparing actual back rays to the
unchanged analytic formula; maximum inward chord error must be<=60um (versus
483um baseline). Preserve source hashes. Record all sample changes and pile
counts. Save once, clean-open front/rear/3Q as highest-risk trio. Reject visible
marks remaining, new rear/facial distortion, pile streaks, or important layer
regression. Full packet and independent image review before bounded retention;
no whole clay stage, hair quality, rig or goal acceptance claim.
