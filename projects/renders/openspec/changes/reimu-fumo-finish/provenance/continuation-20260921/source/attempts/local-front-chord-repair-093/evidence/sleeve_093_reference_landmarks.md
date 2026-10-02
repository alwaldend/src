# Sleeve 093 reference landmarks

Observed 2026-09-06T01:36:55Z. Source:
[canonical_front_25cm.png](../../../projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png),
1000 × 1000 pixels. Coordinates are `(x, y)` from the upper-left corner, with
right/down positive. These are manual pixel traces, checked using a read-only
4× nearest-neighbor detail display of source rectangle `(530,560)`–`(729,759)`.
No image file was altered or created. No candidate, script, mesh, preflight,
or model inspected. This is not approved camera registration.

The outer white cuff edge and inset red dashed trim are different curves.
Use the former for the visible physical cuff span; use the latter only for the
decorative trim. Neither front trace reveals the entire cuff opening in depth.

| Landmark | Pixel estimate | Uncertainty | Visible definition |
| --- | --- | --- | --- |
| Outer cuff, upper end | `(712, 649)` | ±4 px x/y | Far-right rounded corner where the sleeve's upper boundary turns into the long outer white cuff edge. Excludes the cast shadow. |
| Outer cuff, lower end | `(647, 738)` | ±8 px x/y | Lower end of the white cuff band just before/under the skirt frill. The frill and blur conceal the exact terminal seam, so this is a bounded corner estimate. |
| Red trim, upper end | `(688, 630)` | ±5 px x/y | Uppermost visible red dash where the inset decorative line meets the upper sleeve boundary. |
| Red trim, lower end | `(625, 722)` | ±5 px x/y | End of the lowest visible red dash, immediately above the folded lower seam and skirt frill. The inferred continuation meets the lower fold around `(626, 729)`, but that is not another visible dash. |

The outer cuff is gently curved. An approximate intermediate point on its
white fabric boundary is `(684, 694)` (±7 px), provided only to distinguish the
curve from a straight line joining the endpoints. In this image, `(714,684)`
is outboard of the main sleeve/cuff, while `(670,736)` is beyond the visible
lower cuff tip in the frill/shadow neighborhood. They are not its endpoint pair.

The **actual shoulder/throat seam endpoints are not both exposed**. Hair
covers the outboard attachment, and the white/red loop obscures the inboard
connection. The following are observable markers and panel-contour points;
they must not be silently substituted for a fully observed root cross-section.

| Root-region witness | Pixel estimate | Uncertainty | Meaning and limitation |
| --- | --- | --- | --- |
| Upper attachment marker | `(579, 588)` | ±12 px x, ±10 px y | Upper visible white/red attachment-loop region beside the neck and beneath the hair. Supports the shoulder's location, not an exposed geometric throat endpoint. |
| Lower attachment marker | `(561, 637)` | ±8 px x, ±8 px y | Lower end of that loop where it merges into the sleeve at the torso. This is the clearest visible lower attachment witness. |
| Upper main-panel emergence | `(636, 604)` | ±10 px x, ±8 px y | Sleeve upper boundary first clearly emerges outboard from behind the hanging hair, then runs toward `(712,649)`. This is an occlusion boundary, not the true root. |
| Inner lower-panel turn | `(563, 664)` | ±8 px x, ±10 px y | Main white sleeve contour turns away from the red torso and continues diagonally toward the lower cuff. This is already on the sleeve flank, not the lower attachment endpoint. |

For a root estimate, preserve the attachment-marker distinction and the stated
uncertainty. The canonical front does not justify a more precise physical
top/bottom throat pair. In particular, using the panel-emergence and lower-turn
points as the root would confuse hair occlusion and sleeve flank with the
attachment, repeating the cuff-definition problem at the other end.

[Turn 10](side_062_frame_10.png) and [turn 12](side_062_frame_12.png) were
inspected only to confirm the cuff's broad oval opening and the upper interior
arm/attachment relationship. Their oblique projections supplied no front
pixel coordinates. The reference-fidelity skill's uncertainty and controlling-
view guidance informed these distinctions.

Source SHA-256:

- Canonical front: `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`
- Turn 10: `37c1e2866fdbe97ce79a0b5bdddf216a7f151f8e61d2a94c7286c22eaf41fd07`
- Turn 12: `4ebed9233acb351f2a8093496eb3e544771f1ee01e503cec84a1ce86f22b7365`
