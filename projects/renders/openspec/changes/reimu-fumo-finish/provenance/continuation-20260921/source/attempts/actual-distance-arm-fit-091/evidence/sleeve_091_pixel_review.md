# Sleeve 091 pixel review

Observed 2026-09-06T01:13:53Z. Implementation-blind review of the three frozen
091 images below, using the canonical front and reference turn frames 10/12
first, then the corresponding 087 images as baseline. No scripts, meshes,
preflights, implementation numbers, or model state inspected. The reviewer
previously saw 087 during a head-only diagnosis; that limits baseline blindness.
The verified linked worktree branch is `t3code/continue-fumo-desktop-use`.

**Verdict: reject 091 for bounded sleeve retention; reset this sleeve result.**
It does not meet absolute sleeve likeness or constructed-cloth criteria.
Three views cannot establish whole-asset acceptance or rear/mirrored behavior.

1. **Cuff direction and frontal silhouette are major mismatches.** In
   [091 front](sleeve_091_fast_review/front.png), both sleeves extend almost
   horizontally as short rounded mitts, with nearly upright outer cuff marks.
   The
   [canonical front](../../../projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png)
   shows flared cloth falling outward and downward from compact shoulder
   attachments; the cuffs cross the silhouette diagonally beside the skirt.
   The [091 three-quarter](sleeve_091_fast_review/three_quarter.png) preserves
   the padded-mitt read, so this is not only frontal foreshortening.

2. **The opening no longer has the reference's hanging cloth shape.** The
   [091 side](sleeve_091_fast_review/side.png) shows a narrow pointed
   teardrop/slit with a tightly pinched top. In
   [turn 10](side_062_frame_10.png) and [turn 12](side_062_frame_12.png), the
   opening is a broad hanging oval, with a softly collapsing lower fabric
   pocket and the arm/hand seated high inside it. In 091 the dominant cue is
   the folded rim around a small aperture; a distinct rounded hand high in
   the opening is not clearly readable. Absence of a protruding hand in front
   is useful but does not establish correct hand occlusion in depth.

3. **Attachment is connected but the drape is compressed.** There is no
   obvious floating shoulder gap in these pixels. However, the short root
   swells directly into the frontal padded mass, while its side profile
   pinches into a steep upper fold. It lacks the reference's readable flow
   from shoulder attachment into the wider hanging cuff. The smooth inflated
   face and regular rim give a soft object read, but not the photographed
   sewn-sleeve construction.

Head-independent sleeve scores, qualitative and uncalibrated: reference
likeness 4/10; silhouette/cuff direction 4/10; sewn-cloth construction 5/10;
attachment and hand occlusion 5/10. Major identity-defining sleeve failures
remain; relative improvement cannot waive them.

Compared afterward with [087 front](chin_087_all_review/front.png),
[087 side](chin_087_all_review/side.png), and
[087 three-quarter](chin_087_all_review/three_quarter.png): 091 removes the
obvious frontal hand protrusion and reduces the large rounded-square tube
opening, but loses the useful downward/outward sleeve direction and replaces
the opening with a compressed slit. This does not justify retaining 091 as
the next sleeve baseline. 087 is also not an accepted construction.

Construction direction for the coordinator: recover the sleeve's downward
cloth flow and broad side-facing oval cuff, with the hand recessed into its
upper interior and the lower cloth hanging below it. Preserve a seated
shoulder connection while allowing the fabric panel to drape. This is a
visible target, not a numeric adjustment or a prescription for hidden
topology. The reference-fidelity skill's absolute gate drives the rejection.

## Image SHA-256

Paths are relative to this report except the canonical source. Reference and
087 hashes were recorded during the preceding image review; 091 hashes were
observed at the timestamp above. These are pixel identities, not a verified
binding to native model bytes.

| Image | SHA-256 |
| --- | --- |
| `../../../projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png` | `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c` |
| `side_062_frame_10.png` | `37c1e2866fdbe97ce79a0b5bdddf216a7f151f8e61d2a94c7286c22eaf41fd07` |
| `side_062_frame_12.png` | `4ebed9233acb351f2a8093496eb3e544771f1ee01e503cec84a1ce86f22b7365` |
| `sleeve_091_fast_review/front.png` | `4600ba947ac9baf3ebe28c506f07889460b950e7544f281fbf8273db34f334c9` |
| `sleeve_091_fast_review/side.png` | `e3db2e27f39881f6ebded0f53809bb8521eed2b7ed96ea88ed1cba217ec047b5` |
| `sleeve_091_fast_review/three_quarter.png` | `21b84b9b091caf28d77b32f9200365f341efe2ec24e4fa0f2310e3eef9d74fd5` |
| `chin_087_all_review/front.png` | `18610642926ca2be36a12ff0b843ae7b7d69fd2bb0a437df93634374babeeca9` |
| `chin_087_all_review/side.png` | `1e1e6d50be55d2d3528106bc39650932305d95e4bf31c321be71776552aff987` |
| `chin_087_all_review/three_quarter.png` | `dca74e014f7abf74918ec6d3329b0f4d479e4c1633bc3c39670a0dcac9dda197` |
