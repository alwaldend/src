# Upper-front depth continuation proposal

Read-only design review; no model or helper edits, execution, or save.
Baseline: `head_027c_candidate.blend`, SHA256
`19ce0cb14c7d679750422702d6df97753480cf8a4db7cd73f1203b0f28bf7416`.
Inspected helper SHA256
`ac21e50760c7345fbd036a2777359101791079ba49e50cb766111199345077a6`.
Evidence: fixed side and three-quarter images in
`head_027c_eevee_review/`, canonical front, physical front and physical side
under `projects/renders/assets/reimu_fumo/references/`. No new research or
Blender query was necessary. Physical side supports soft curvature, not an
exact depth measurement: its view and individual plush differ.

## Diagnosis and decision

027c is less bulbous at the temple, but its side view has a long nearly
vertical front wall meeting the crown at a pronounced corner. The helper
explains this: `max(-.0525, -.03915-.36*(z-.09168))` becomes constant at
Z = 128.7639 mm; on the centerline, the radial edge transition starts near
Z = 177.1710 mm. Thus almost 48.4 mm of upper panel inherits one depth before
turning sharply into the crown. This is not evidence that the whole lower
face should become spherical. Retain the current lower facial construction
and finite rolled gusset; replace only that upper depth plateau with one
tangent-continuous rising profile shared by the receiver and hair.

## One bounded profile

Keep world X/Z, every point at Z <= 125 mm, the rear panel, and the existing
front/rear seam functions unchanged. This includes the89mm protected contact
patch and the lower-face sampling grid through 125 mm. Existing materials,
lights, body, collar/tie, bow, feet and hem remain outside scope.

Above Z0=.125m replace the `broad` term's hard maximum with a cubic Hermite
profile H. Parameters are a first coarse hypothesis, not measured acceptance:

- Z0=.125m, Y0=-.0511452m, dY/dZ at Z0=-.36, matching the inherited panel.
- Z1=.1932208091020584m, Y1=-.022m, dY/dZ at Z1=.85, meeting the unchanged
  upper front seam rather than shifting the crown or bow.
- L=Z1-Z0; t=(Z-Z0)/L; use the ordinary cubic Hermite basis:
  H=(2t^3-3t^2+1)Y0+(t^3-2t^2+t)L*(-.36)
  +(-2t^3+3t^2)Y1+(t^3-t^2)L*(.85).
- Preserve the existing transverse stuffing term and radial edge blend
  E=smoothstep((r-.72)/.28). Therefore the depth change is exactly
  DeltaY=(H-old_broad)*(1-E), zero below Z0 and zero at the perimeter.
  There is no new surface-hit/overhang branch or pointwise hard clamp.

These are analytic centerline witnesses from the current helper, not sampled
candidate geometry. Positive change means rearward, away from the camera:

| World Z, mm | 027c front Y, mm | Proposed front Y, mm | Change, mm |
| --- | --- | --- | --- |
| 125 | -51.145 | -51.145 | 0 |
| 140 | -52.500 | -53.011 | -0.511 |
| 150 | -52.500 | -50.818 | +1.682 |
| 160 | -52.500 | -46.423 | +6.077 |
| 170 | -52.500 | -40.296 | +12.204 |
| 180 | -50.243 | -32.262 | +17.981 |
| 193.221 | -22.000 | -22.000 | 0 |

The small 140 mm forward lobe preserves a stuffed forehead; the substantial
170–180 mm retraction gives this first test a visible causal difference.
It does not shrink total head depth by 18 mm: the deepest forehead remains
within about 0.5 mm of its current depth while the upper cap becomes shallower.

## Coupled application and falsifiable witnesses

Use this same upper profile in the core and 027c's continuous analytic fringe
receiver. Reconstruct the dependent hood from its corresponding core faces,
retaining existing cloth thickness/padding and rooted overlap. Apply the
surface-depth difference to embroidery while preserving its X/Z and depth
residuals; keep lower cheek locks untouched except any upper root that must
follow its receiver. Leave `_panel_point(..., False)` and the rolled seam
unchanged. This is one profile refinement, not independent edits to each
visible patch. Expected object scope is core, hood, fringe and necessary
upper face-detail/root depth fitting; all other controls remain identical.

For direct deformation of frozen 027c instead of regeneration, use one
continuous through-depth extension: normalize s from the current front
receiver to rear receiver and multiply DeltaY by
1-smoothstep(clamp(s,0,1)). Front cloth outside the receiver gets full motion;
rear cloth gets zero. Use corresponding analytic panel coordinates, not
nearest-hit branches. Check a positive depth Jacobian before saving; do not
move the whole rear head with the front field. Regeneration with explicit
front/rear topology is simpler if the coordinator's wrapper supports it.

Three witnesses decide whether this bounded test is useful:

1. Fixed side/three-quarter pixels: the upper front should rise progressively
   into the crown, removing the wall-to-roof corner. Reject a new sloping
   slab, shelf, pinched cap, or reappearing round helmet sweep. The table is
   a predicted displacement witness, not a visual passing score.
2. Front and lower-face controls: identical landmark X/Z, identical 63-point
   lower-plane grid through 125 mm, identical retained underside face hash
   and unchanged collar-contact witnesses. Refit embroidery without
   flattening its thickness; report any new hair/contact intersections.
3. Coupled cloth: no pale crown/temple gap, ridge at an old ray-hit boundary,
   severed fringe root, or newly floating bow. Keep the rear contour and
   4–11 mm varying gusset intact; record actual post-change depth bounds.

The strongest risk is excessive upper-front recession: the 18 mm change near
180 mm may make the cap look pinched or expose the bow attachment. That risk
is visible in one fixed side/three-quarter comparison and does not justify
reopening body or material work. Restoring the old whole-head radial bulge
would surrender the lower plane; a constant-depth extrusion would preserve
the wall. Neither addresses this localized cause. The coordinator owns the
attempt decision and acceptance; this note claims neither.
