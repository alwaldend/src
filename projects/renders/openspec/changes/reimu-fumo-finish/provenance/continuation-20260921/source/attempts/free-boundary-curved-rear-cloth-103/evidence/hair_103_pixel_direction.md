# Hair 103 independent pixel direction review

Reviewer ID: `reimu-fumo-hair-094-pixel-review-agent-v1`.
Reviewer role: implementation-blind image reviewer, `/root/hair_094_pixel_review`.
Observed: `2026-09-06 06:04:40 UTC`.
Record: new immutable review of the four PNGs hashed below.

**Direction verdict: RESET.** The shorter rear center and more tapered side end improve locally on 102, but a major new cream-colored side opening breaks the hair coverage and root appearance. The rear still reads as broad fitted sections meeting a rounded lower mass, without the controlling references' independently ending cloth layers. The local gains do not justify keeping this shape direction.

## Image evidence

I inspected all four actual 103 PNGs: rear, side, front, and mirrored three-quarter. I then compared the same four views from retained `arm_101_all_review` and failed `hair_102_failed_review`, plus canonical front and turn extracts 12/16/21. No scripts, arrays, contacts, preflight, plans, implementation messages, or nonpixel reports were inspected. This reviewer retains earlier pixel-review context but makes the present judgment from these rendered images.

| Area | Finding |
| --- | --- |
| Cream opening / side root | **Major new regression.** Exact side shows a large pale strip beginning high on the side of the head with a stepped rectangular upper edge, then running down into the exposed lower cheek. Mirror shows the same conspicuous pale opening. This is absent in 101/102 and unsupported by the reference hair coverage. It makes the side hair look interrupted and its root visibly unresolved. This describes the rendered coverage; it does not diagnose the underlying geometry. |
| Broad cloth versus cap/card | **Insufficient improvement.** Broad rear divisions are readable, but the forms remain smooth fitted sections. The side tail is more curved and tapered than 102's straight curtain corner, yet it still appears thin and rigid rather than establishing a broad soft cloth face with a clear lap. |
| Rear independent endings | **Still missing.** The center is shorter than 102 and no longer extends as far over the upper neck, a useful local correction. Its bottom remains one rounded arc; the large adjacent sections do not create the staggered free cloth ends visible in turn16/21. |
| Lower jaw / chin framing | **Mixed.** The shorter rear center reduces 102's excessive lower continuation. Front still has a heavier dark surround under the lower cheek corners than retained 101. The much larger exposed pale side area makes the side head/neck transition less coherent. No physical chin intersection is inferred. |
| Front identity | **Preserved in the direct front.** Eyes, central fringe, mouth, and tied cheek-lock outline remain recognizable and visually close to 101/102. The direct front hides the serious side-coverage defect; it cannot outweigh the side/mirror evidence. |
| Bow and arms | **No obvious new change in these four views.** Bow outline and gathering, contained upper-cuff hand, sleeve silhouette, and arm placement remain visually consistent with retained 101. This is a four-view observation only. |

The strongest evidence for REFINE rather than RESET is the reduction of the long rear curtain and the more tapered side end. The stronger contrary evidence is the new large stepped cream opening, combined with the continuing lack of independently ending broad rear cloth. Reference frame12 preserves dark side/rear hair coverage around a much smaller lower cheek exposure; frames16/21 show real lapping faces and separate ends. The current result does not maintain those relationships.

## Limits

Only the four named views were reviewed. There is no reviewed non-mirrored three-quarter view or presentation render for 103, so this is not a complete regression set or model-acceptance review. No geometry validity, hidden attachment, or contact conclusion is made. The oblique references do not support exact orthographic pixel measurements here. This report issues a shape-direction verdict only.

## Exact PNG hashes

Paths are relative to `out/reimu_fumo_finish/desktop_astra/`. All four SHA-256 values were independently calculated.

| PNG | SHA-256 |
| --- | --- |
| `hair_103_failed_review/rear.png` | `83a7dc344190170c2713142f50b74cce3d7a41ec285ebc2e22fdab03fd5e869d` |
| `hair_103_failed_review/side.png` | `66694445cd3461cd579ad13e41c612fe4c2d4db23ce0437ccf3c71f34e80f6ac` |
| `hair_103_failed_review/front.png` | `62f3d8a99985f475f14159f58d078faadf9674a413dab2302e69756d8df9d3fd` |
| `hair_103_failed_review/three_quarter_mirror.png` | `12828111a6b556ee69f0cc33f7ca27b7198c1d3098252eefa55cd4668bdfe788` |

Controlling reference paths: `projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png` and `out/reimu_fumo_finish/desktop_astra/side_062_frame_12.png`, `side_062_frame_16.png`, `side_062_frame_21.png`.
