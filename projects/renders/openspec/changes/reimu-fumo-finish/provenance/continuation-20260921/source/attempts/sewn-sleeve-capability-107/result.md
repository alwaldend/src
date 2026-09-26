# Coupon107 closed reset: softness works, reference construction does not

Frozen corrected coupon5880a74d61d3f55026fc1bcf57a8b89ea332e56f8318e3c49aad8433ae762053.
Source71c01d74418f63cfabe0f7386b2886dab056c76841678342dce34173120ee702.
Native receipt7cec8c4a00e4bd0a72b19fa658571af4b5ff6e85dcdd753c779cf2edb38e5e53.
Pinned5.2.1 LTS build9e2066aef7ef clean-reopened the frozen mesh, with no
dependency on RAM cloth cache. Fast512px fixed front/side/three-quarter:
front1e2793cfa4e63fce32173c53232b5f3d765a05258c3161427b28f1fa722b858f;
side6d0b12146cc057e9967acc9e0c19bc788113a5ebe407be95b25f9a8b72f7ab24;
three-quarter93ae96f529768f0ed231f6fe4046dc24b55d714a71d57e65a9ec1b30de0acf06.
All images in cloth_107_corrected_coupon_review. No asset promotion.

Independent image-only reviewer macro105-blind-pixels-A (prior-pixel familiarity,
no implementation) rejects sleeve capability. Root agrees. Soft folds exist,
but the sleeve is a crumpled high cuff/rosette, too short and tightly gathered.
Its opening faces the front, whereas the reference front shows a broad outer
panel and the side reveals the opening. Exposed arm/root transition fails.
Simulation success and fabric-like deformation do not prove faithful form.
This was a subsystem capability coupon in rejected106 context; fullasset
criteria001-008 remain unverified for this subject, not eight implied passes.

Execution evidence and corrections:
- API preflight first failed because5.2 rest_shape_key takes a ShapeKey pointer,
  not an integer. Corrected to actual Cut_flat_pattern key; no solver ran before.
- First native simulation timed out240s after80frames without a settled save.
  Initialc8911cde486ce7db2ddd9ed239c056212bfa0c19ce8c1129c7e3272f7554e9fb remains.
- Audit3bae68dde7f4904a81594c502d65cd7333df90ef6eaead38a4f5ed1d8720f0b6
  proved18fixedroot vertices and many free vertices inside misaligned arm,
  min signed-9.257mm. Both native collision margins clamped1mm, totaling2mm.
- One setup correction aligned support, seated pins outside, fixed inward
  winding, and completed initial-free-sample clearance with <=0.811mm projection.
  Corrected preflight60f9593d84eb3db5a1df4a2da636aa8e4d34bcd96df6717733a611f68a7eda89:
  arm minimum2.381mm, pins2.753mm; torso5.478mm; combinedmargin2mm.
  No mass/bending/rest-pattern tuning; no pressure. A preliminary corrected
  clearance check failed at1.590mm and was preserved in cloth_107_corrected.log
  before the bounded free-sample projection completed that repair.
- Correctedinitial350a53390764e0911bc1f3c14c5e5ad76247a5488420eccca1d7f5c82b0a5a21.
  Cache advanced through64frames; immutable evaluated mesh snapshots16/32/48/64
  saved. Loop erroneously continued to96, but cache ended64:80/96 snapshots
  repeat64 state, not extra simulation or independent evidence. Source is
  preserved unchanged so this receipt discrepancy remains auditable.

Decision/process review covers105-107: REVISE / ASK for new source authority.
Two whole-character cage studies and one bounded cloth coupon retained no
reference-faithful macro base. Faster image feedback is demonstrated; sufficient
artistic construction capability is not. No more solver variants or inherited
helmet/primitive cage tuning. Direct pattern/sculpt capability also remains
unproven; do not promise another clean rebuild will solve this.
The owning asset README currently prohibits outside geometry/materials/etc.
Root will ask whether the user permits adapting a suitably licensed base mesh.
No outside asset has been acquired or used, and no restriction was changed.
If not authorized, explain the capability limit and agree a different authoring
input, rather than silently importing an asset or continuing rejected variants.

Ergonomics: native initial collision audit should have preceded96frames; that
cheap check would have avoided the4minute failed run. Snapshot freezes fixed
loss of useful solver evidence on timeout. Root runtime-read values beat
configured numbers (1mm clamps). One API type mismatch and loop bound oversight
remain avoidable. No new infrastructure, skills, hostconfiguration or tooling
framework added. Protected101 SHA40ec32f8d2832f4a42e36c4ee4f27bdf213d1940ba0114a4fd2efcc476b9408b unchanged.
