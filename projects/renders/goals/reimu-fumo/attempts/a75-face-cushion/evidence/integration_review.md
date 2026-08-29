# A75 one-object integration-boundary review

## Recommendation

**REVISE, then proceed.** Replacing only
`Head_Cushion_Manual_Target` is a valid small, reversible experiment **only as
a visible face-carrier / aperture cushion test**. It is not a valid test of the
complete stuffed head receiver, profile volume, rear construction, or final
hair integration.

The mechanical boundary is clean: the source receiver is single-user,
parentless, constraint-free, driver-free, and modifier-free. Hiding it in a
copied scene and adding one new single-user object does not mutate the protected
rung. The visual boundary is narrower. The retained cap, fringe, temple panels,
cheek locks, rear leaves, and bow own most of the head silhouette and hide most
of the receiver outside the face opening. They therefore provide useful whole-
plush context and regression witnesses, but they make hidden receiver quality
unobservable.

Do not run A75 with the current broad claims unchanged. In particular:

- front and three-quarter renders may judge the exposed beige carrier;
- side and rear renders may judge only leakage, depth order, and regressions;
- the retained hair prevents side/rear renders from proving a complete shallow
  cushion construction;
- whole-plush construction/contact scores cannot reasonably be required to
  rise to `6/10` when the known-rejected hair, bow roots, and lower stack are
  deliberately frozen. Require a strong affected-region score and no whole-
  plush regression instead.

This recommendation is evidence, not authorization to change the protected
asset. The canonical writer should make only a disposable candidate.

## Evidence inspected

- protected rung 003 blend:
  `sha256:c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`;
- A66 exact object inventory:
  `sha256:6d60833cfb1388a6203cef987b08a44b5a35ae13539946f25b1e405652aaaeea`;
- A69 exact contact/interface inventory:
  `sha256:ea4c83d4f2d9128fcebe8552bb79cc6c641991a2b1d4ba07bef86a793834af43`;
- A74 blind review:
  `sha256:b0d855cd1aa5d9f4d6793c99334852561e8312115013651efca907e813a354ae`;
- A74 self-review:
  `sha256:91d1c1973877d9243e8b3387d63c3e5e9ce625a1a585dfbb4f3d3a70b6a2d32e`;
- immutable A74 five-view renders and protected rung-003 five-view renders;
- the complete relevant reference packet: canonical front and turntable,
  clean/physical front, physical side, sofa, and older turn references.

A74's decisive failure was not merely a bad silhouette. Hiding all old hair
owners made the rear bald by construction, while its one-sided support became
a white square card in every non-front view. Retaining complete old hair in
A75 removes that guaranteed failure and makes a bounded face-carrier comparison
possible.

## Strongest case for the proposed boundary

The proposed boundary isolates the most concrete A74 failure: a face carrier
must be a softly rounded stuffed form, not a one-sided rectangle. The current
receiver is independently editable and its exposed beige region is large
enough in front and both three-quarter views to judge face-plane curvature,
cheek/chin rounding, applique seating, shading continuity, and plush read. The
complete old hair also prevents a deliberately partial experiment from being
misread as a finished bald model. This is much faster and more falsifiable than
another simultaneous receiver/hair/bow rebuild.

## Decisive limitation and disconfirming evidence

The old hair is not a neutral fixture around an arbitrary receiver. It is a
thin shell fitted almost exactly to the current receiver and several retained
objects already cross that receiver:

| Retained object | Baseline receiver relation |
| --- | --- |
| `A44 continuous hair cap with smooth opening` | no crossing; minimum sampled clearance about `0.529 mm` |
| left/right temple transition panels | crossings, `112 / 119` triangle pairs |
| left/right cheek locks | crossings, `60 / 73` triangle pairs |
| left/central/right rear leaves | crossing / `0.913 mm` gap / crossing |
| left/right bow loops | crossings, `68 / 59` triangle pairs |
| left/right loop ruffles | crossings, `61 / 72` triangle pairs |
| `Garment42 compact bodice` | crossing, `140` triangle pairs |

The face witnesses are also registered very tightly to the current surface:

| Frozen witness | Baseline minimum sampled distance |
| --- | ---: |
| left/right eye applique | `0.076 / 0.080 mm` |
| left/right half-lid stitch | `0.174 / 0.174 mm` |
| left/right upper-expression stitch | `0.159 / 0.159 mm` |
| mouth dash | `0.320 mm` |

Consequences:

1. A substantially different full receiver cannot be inferred to fit merely
   because the exposed front looks better.
2. Most profile and all rear receiver construction is occluded by unchanged
   brown geometry, so it is not visually testable in this module.
3. A global "no crossing" gate would be false before A75 starts. Contact gates
   must be baseline-relative and distinguish retained known defects from new
   defects.
4. The hair/fringe defines most of the face-opening contour. A receiver-only
   change cannot repair a wrong hairline, cheek-lock silhouette, crown, or rear
   mass. If the exposed beige pixels barely change, the module has no useful
   visual leverage and must stop.
5. Requiring the unchanged whole plush to score `>= 6/10` in construction and
   contact confounds the module result with known old-hair and lower-stack
   failures. Rung 003 itself was recorded around `4.5/10` for constructed-plush
   read and `4.3/10` for contact/occlusion.

## Credible alternatives

| Alternative | Outcome quality | Risk / reversibility | Cost | Verdict |
| --- | --- | --- | --- | --- |
| Revised one-object exposed face cushion | Directly tests face-plane volume and applique support in whole-plush context; cannot validate hidden head construction | Low risk; fully disposable and independently reversible | Low | **Recommended next experiment** |
| Replace all 15 audited head/hair objects atomically | Can validate complete silhouette, rear coverage, roots, and construction | High strategy risk; A74 already showed a bad full boundary can erase useful context | High | Defer until one bounded form is demonstrably good |
| Add a padded face patch over the old receiver | Keeps every old contact while changing visible face | Produces doubled surfaces and repeats A74's card/mask failure mode | Low | Reject |
| Deform only the exposed region of the old receiver in a copy | Preserves the hidden cap-fit envelope and can improve cheek/chin shading | High-density legacy mesh remains hard to reason about and weak as reusable construction | Medium | Fallback if a replacement cannot maintain witness seating |
| Replace receiver and move/remake seven face witnesses together | Allows unconstrained curvature and correct flush placement | Confounds cushion quality with eye/mouth placement and expands a well-defined module | Medium | Later coupon only after the cushion representation passes |
| Keep rung 003 unchanged | Perfectly reversible baseline | Does not advance the stated visual goal | None | Baseline, not a solution |

## Exact A75 P0 mutation boundary

### Hide, by exact name only

- `Head_Cushion_Manual_Target`

Set only `hide_viewport = true` and `hide_render = true` on the copied source
object. Do not delete, rename, relink, edit its datablock, or change its
materials.

### Add

- one object: `A75_Face_Cushion_P0`;
- one single-user mesh datablock under an `A75_` name;
- one single-user beige diagnostic material under an `A75_` name;
- one explicit `A75_` collection for the experiment.

Do **not** add a seam witness at P0. The controlling references do not expose a
face-cushion seam clearly enough inside the hair opening, and a decorative seam
would add geometry without testing the representation. A later source-backed
seam coupon can follow a passing cushion.

### Preserve visible, unchanged

Preserve all remaining 14 objects in the audited head/hair replacement
boundary:

- `A44 continuous hair cap with smooth opening`
- `A44 left temple fringe panel`
- `A44 left temple transition panel`
- `A44 off-center main bang panel`
- `A44 right swept fringe panel`
- `A44 right temple transition panel`
- `A45 left tapered flexible cheek lock`
- `A45 right tapered flexible cheek lock`
- `A42 Left asymmetric rear lock`
- `A42 Off-center main rear lock`
- `A42 Short right rear lock`
- `A42 Main lock left seated seam`
- `A42 Main lock right seated seam`
- `Subtle crown center seam`

Preserve the seven exact face witnesses:

- `A45 left flush composite eye applique`
- `A45 right flush composite eye applique`
- `A45 left drooped half-lid stitch`
- `A45 right drooped half-lid stitch`
- `A45 left fine upper expression stitch`
- `A45 right fine upper expression stitch`
- `A44 tiny neutral embroidered mouth dash`

Preserve the complete bow/root/trim group, body, clothing, sleeves, feet,
lights, world, fixed review cameras, scene, and view layer byte/transform/
visibility unchanged. Leave all 41 already-hidden legacy head/face/hair
objects hidden and untouched. Do not edit shared `Face fabric clay`: it has
outside consumers on both cuff pads.

## Revised acceptance claim

The only claim A75 P0 may establish is:

> In the exposed front/three-quarter aperture, a replacement beige face
> carrier reads more like softly stuffed fabric than the rung-003 carrier,
> keeps the frozen facial appliques seated, and introduces no visible
> regression in the complete protected context.

It may not claim that the complete head cushion, hidden seam topology, profile
depth, rear volume, hair construction, or bow roots are approved.

## Early-stop gates

### Before building

Stop before Blender mutation if any intended improvement requires moving a
hair object, changing the hair opening, moving an eye/mouth witness, or
changing the bow/body. That is a different module. Also stop if the design
cannot name a front and three-quarter pixel change that should be visible with
the old hair retained.

### Static integrity after building

1. Recheck the protected source hash exactly.
2. Visibility diff against the appended source must contain exactly one source
   object: `Head_Cushion_Manual_Target` becomes hidden. No other source object,
   datablock, transform, modifier, material assignment, or visibility bit may
   change.
3. The candidate contains exactly one new visible mesh object and a
   single-user diagnostic material. No helper cap, rear closure, face card,
   bridge, or hidden support is allowed.
4. The new cushion must remain inside the retained cap outside the visible
   face aperture. The clean cap/cushion relation must have no triangle
   crossing and should retain a positive support clearance on the order of the
   `0.529 mm` baseline, not a multi-millimetre hollow gap.
5. All seven facial witnesses must remain in front of and nearly flush to the
   evaluated cushion. Use the existing `0.076--0.320 mm` distances as the
   baseline and reject any crossing or visible float; a conservative P0 bound
   is positive `0.05--0.60 mm` at the nearest evaluated surfaces.
6. Known baseline crossings with temple panels, cheek locks, rear leaves, bow
   loops/ruffles, and bodice are not evidence of a new A75 defect by
   themselves. Stop if A/B masks or evaluated contacts show a new visible
   penetration, a new object-pair contact category, or a gap exposed by the
   candidate.

### Fast visual packet

Render front and both three-quarter views first, with the protected rung as a
fixed-camera A/B baseline and with canonical references attached.

Stop immediately when any of the following is visible:

- square/card/mask/egg/mattress/helmet read;
- beige pixels outside the protected hair envelope;
- facial appliques embedded, floating, or casting a disconnected shadow;
- a new hard vertical side wall, sharp rim, unsupported air gap, or accidental
  tangency;
- no material visible change from rung 003 in the exposed face region;
- improvement depends on explaining hidden topology rather than the pixels.

If the fast packet survives, render side and rear only as regression views.
They must show no new beige leak, silhouette change, bald patch, or depth-order
change. Do not score invisible receiver construction from those views.

### Blind review gate

- Judge the exposed cushion region for rounded stuffed-plane read,
  cheek/chin transition, applique seating, and continuity with the hair.
- Require the affected-region construction and contact scores to be at least
  `6/10` and require an unambiguous preference over rung 003 in front and both
  three-quarter views.
- Require no regression in whole-plush silhouette, identity, contact, or
  presentation. Report whole-plush scores, but do not require this one module
  to cure unchanged subsystems.
- Reject if the reviewer cannot identify the improved region without being
  told how it was modeled.

## Condition that changes this recommendation

If the front/three-quarter A/B packet shows that the hair opening completely
dominates the visible shape, or if maintaining the frozen applique seating
forces the replacement back to the old receiver surface, close A75 as an
interface failure. The next honest boundary is not a decorative patch. It is
either a controlled front-carrier-plus-seven-applique coupon or, once its
representation is designed from source measurements, the full audited
15-object head/hair boundary.
