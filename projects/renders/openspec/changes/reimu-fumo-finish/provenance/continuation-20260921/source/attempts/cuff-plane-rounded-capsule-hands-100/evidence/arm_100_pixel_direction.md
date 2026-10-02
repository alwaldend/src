# Arm 100 image-only direction review

Reviewer: `reimu-fumo-hair-094-pixel-review-agent-v1`, `/root/hair_094_pixel_review`.
Observed: `2026-09-06 04:33:49 UTC`.

**Direction verdict: KEEP THIS SHAPE DIRECTION; no visual reset is needed for the requested hand-volume/containment balance.** It is worth resolving the known construction failure around this direction. This is a failed-build diagnostic, not a candidate acceptance or geometry-validity judgment.

I inspected the three ARM100 images first, then the corresponding ARM099 and retained BOW098 images, canonical front, physical side, and turn extracts 10/12. No model, code, geometry, preflight, or construction implementation was inspected.

| View | Visual result |
| --- | --- |
| Front | The hands remain contained, preserving 099's useful improvement over the low protrusions in 098. The sleeve silhouette stays clear and matches the canonical front's concealed-hand read. |
| Side | The tiny upper-rim crescent in 099 becomes a substantial rounded hand within the upper cuff. It now has readable volume while leaving cloth visible below, closer to physical side and turn10/12. |
| Mirrored three-quarter | A rounded hand is again clearly visible high inside the near cuff. This directly improves 099's mostly vacant opening without returning to 098's lower protruding placement. |

The strongest contrary evidence is that the hand still sits close to the upper cuff and the lower sleeve remains open. Within these three images, however, this no longer removes the hand's readable volume: the oval is visibly present in both controlling oblique/side views. The requested balance is therefore materially better than 099. I found no visible return of the low front protrusion and no shape regression in the three supplied views that warrants resetting this direction.

This report adds no new acceptance criteria and makes no inference about the known cloth-contact failure. The conclusion is limited to whether the visible shape direction is worth continuing.

## Diagnostic image hashes

Paths are relative to `out/reimu_fumo_finish/desktop_astra/`. SHA-256 values were independently calculated.

| Image | SHA-256 |
| --- | --- |
| `arm_100_failed_review/front.png` | `eae25b899e6857d647b80c25e04ca561c92768dbf18f590c50af3d10412bfe74` |
| `arm_100_failed_review/side.png` | `6c86af3d2ae9071ec00c6b5375663fe2e6bafa429cbd42f9aa52934595b2ca6f` |
| `arm_100_failed_review/three_quarter_mirror.png` | `0fe23a03ed8e70bc7d39913bb08b201850d30c55ead2a56d99fe4d4bdcd8237b` |
