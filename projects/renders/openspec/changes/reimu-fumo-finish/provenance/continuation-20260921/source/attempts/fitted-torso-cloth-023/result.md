# Fitted torso retained as an intermediate improvement

Candidate body_023_candidate.blend SHA-256
`61c5efe89e833f8c79b5327a439cb2e3688113a927c8dc094a98ed4441a718f1`.
Source022b SHA-256
`96e6deea298308573174a35699ea4cf7b99e827260b2c108de43f8f0c1266014`.
Retain this candidate for further modeling, not as a completed stage.
All full-goal acceptance remains open.

Torso width is now56.580009mm, approximately0.486Wh, matching the canonical
mid-bodice target. Height is unchanged. Sleeves use a separate0.90X field
with6.5mm inward root shift. New thin collar flaps and gathered tie follow
the evaluated chest and head underside instead of floating in front.

The first unsaved width probe pinched sleeves; separating their affine field
fixed that mechanism. Collar direct rays then missed the rounded shoulder;
a downward fallback also missed its lateral edge. The final representation
bridges roots to supported lower cloth and permits only2mm lateral inset
for the free outer edge. An exact boundary probe found one corner needed
2.0653mm, so its lower-X coefficient changed .88 to .875, moving it inward
.14145mm instead of weakening the2mm guard. All failures occurred before
saving; earlier candidate files were preserved.

Root inspected all five fixed renders. Implementation-blind reviewer
body_023_blind inspected references first and candidate views before baseline:
likeness5.5, silhouette5, construction4, identity6, contact5, medium3.5,
presentation7 out of10. Completed-model and visual-review criteria fail.
After baseline comparison, both reviewer and root retain the narrower torso
and fitted collar/tie as an improvement with no major new visible regression.
Persistent failures: helmet head/hair, cone-like skirt, taut bow, rigid funnel
sleeves, and black ellipsoid feet without cream leg sections.

Independent head_019_technical clean-reopened023 and022b with pinned Blender
5.2.1 LTS build9e2066aef7ef. All32 declared geometry controls and rig pose
match. Maximum Y/Z drift across width targets is6.52nanometers. New cloth
meshes are finite, closed, consistently wound, nondegenerate and have unit
Body weights. Evaluated root gaps: collars0.376–0.490mm, tie0.297–0.810mm.
The tie has no supported intersections. Each collar has seven supported
edge/chest intersections at its back thickness, maximum penetration0.079mm.
At these witnesses the visible front stays1.030–1.234mm clear. This is a
minor hidden-back caveat, not evidence for a globally collision-free claim.
No animation, whole-scene technical or final-material acceptance follows.

Fixed review contract SHA-256
`4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4`.
Front PNG `5ed3c1814df5c6266f00d0a6ee89f1bb34b08a93a8372193b624b76daa4746b6`.
Three-quarter PNG
`d511d0948df0a21a6e1b0f57af68213cacb09a2acb333410a2c2f515f9043cd0`.
Local artifacts are under out/reimu_fumo_finish/desktop_astra/:
body_023_review, body_023_writer_receipt.json, body_023_technical_audit.md,
body_023_technical_audit.json, and body_023_contact_detail.json.

Next decision: replace the flat white hem and black pods together with a
gathered strip over cream stuffed legs with black sewn toe panels. Keep red
skirt, torso, head, bow and established global height. Root owns canonical
writes; independent helper drafting and frozen-candidate review run in
parallel. Native background save/reopen/five-view feedback remains fast;
there is no reason to reopen desktop or screenshot troubleshooting.
