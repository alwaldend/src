# Reimu Fumo attempt 3

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 3 — clean sewn-pattern macro forms

**Candidate:** planned `out/reimu_fumo_attempt_003_macro.blend`, SHA-256 pending,
black-silhouette macro stage, review packet
`attempt-003-macro-fixed-five-view`.

**Failure targeted:** The inherited head reads as an egg or helmet, the dress as
a cone, sleeves as tubes, bow as wings, torso as a human body, and feet as
tangent balls. Those failures survived at least two primitive-derived attempts.

**Hypothesis:** One shallow front/rear/gusset head cushion plus compact sewn
envelopes for the torso, skirt, sleeves, bow, and foot pods will produce the
correct nonhuman plush silhouette without identity graphics or materials.

**Plan written before implementation:**

1. Start from factory-empty Blender data and retain no baseline render mesh.
2. Use `Wh = 1` and the frozen world axes, dimensions, contacts, and cameras
   from `LANDMARKS.md`.
3. Construct a single shallow head cushion with a broad flat central face,
   rounded perimeter gusset, slightly narrower lower edge, and deeper rear.
4. Construct only a compact hidden torso, broad front/rear skirt envelope,
   detached bell-envelope sleeves, flattened gusseted foot pods, and asymmetric
   folded bow-loop and tail pockets. Seat every root inside an overlapping
   manufactured seam or garment volume.
5. Exclude the face, fringe, locks, ruffle detail, fabric texture, armature, and
   final materials so none can rescue a wrong macro form.
6. Save one candidate, render fixed front, side, rear, three-quarter, and
   perspective clay views, plus four flat black silhouettes from the exact same
   bytes.
7. Measure the applicable frozen landmarks and run an implementation-blind
   macro review. Reject before detail if any fixed regression failure remains.

**Work performed:** Built a factory-empty candidate with a unified seven-ring
head cushion, compact torso, lofted skirt envelope, two bell-envelope sleeves,
two gusseted foot pods, and five inflated bow parts. Created the four frozen
orthographic cameras, a perspective presentation camera, neutral studio
lighting, five clay renders, four black silhouettes, a landmark-aligned front
overlay, and a front/side reference contact sheet. No baseline render mesh was
retained and the tracked reusable model was not modified.

**Evidence:** Candidate SHA-256
`9aa12e0563d3ba71f1ecbb2181ebf8dc90c55c890fa333a82114b9a4cc40dbff`.
Measured evaluated bounds in `metrics.json` were: head
`1.018 × 1.049 × 0.740 Wh`, skirt `1.098 × 1.075 Wh`, feet approximately
`0.313 × 0.288 × 0.370 Wh`, sleeve width `0.334 Wh`, overall bow span
`1.710 Wh`, and overall height `2.120 Wh`. These are inside or close to their
applicable frozen numeric bands. The 42% aligned physical-front overlay still
showed round mouse-ear bow loops, oversized triangular tails, bean sleeves,
and an undifferentiated bare head. The side render showed a vertical slab bow,
box head, exposed neck column, wedge skirt, and foot projecting like a leg.

The implementation-blind reviewer answered no to unlabeled same-subject
recognition and described the intended constructed-plush read as weak. Scores
were overall macro likeness 2.5/10, silhouette and proportions 3/10,
construction 2/10, contact and occlusion 3/10, intended-medium read 2.5/10,
and presentation 6/10. Major failures were the missing identity-defining hair
mass, wrong mouse-ear bow, floating neck transition, primitive sleeves and
feet, and radial surface pinching.

**Criterion results:** Reference likeness fails. Measured component bounds pass
the applicable subset of criterion 2, but the full silhouette criterion fails
because the hair, bow, side garment, contact, and occlusion silhouettes are
wrong. Plush construction and presentation quality fail. Reusable structure,
animation readiness, and technical integrity remain unverified at this stage.
Repository delivery still applies only to the rejected migrated baseline.

**Decision:** Reject at the macro-silhouette gate. Do not advance this geometry
to identity detail, rigging, materials, or the tracked reusable asset.

### Progress and approach audit after attempt 3

- **Improved:** normalized component bounds, camera repeatability, and the
  distinction between a planar head face and rounded gusset all improved
  measurably over the inherited baseline.
- **Regressed or unchanged:** no user-visible acceptance criterion passed.
  Recognition, manufactured construction, hair identity, bow construction,
  side silhouette, and seated contact remain absolutely poor.
- **Absolute result:** 2–3/10 is still a failed plush model, not an approval
  candidate. Numeric landmark compliance did not translate into likeness.
- **Approach evidence:** the panel concept remains appropriate, but the chosen
  radial n-gon caps plus subdivision create visible pinching and foam blocks.
  Inflating every part with the same cushion generator erases the different
  manufacturing logic of hair, ribbon, sleeves, skirt, and feet. Omitting hair
  until a later identity tier also removed an identity-defining macro mass.
- **Repeated-defect diagnosis:** the same bow, garment, foot, hair, and side
  failures survived again. Their representations must be replaced; parameter
  tuning of Attempt 3 is prohibited by the reset rule.
- **Highest-leverage problem:** build the head, padded hair cap and locks, and
  thin folded bow as one seated sewn assembly, because that cluster controls
  recognition, head silhouette, attachment logic, and the intended-medium
  read in every view.
- **Next approach:** retain only the frozen measurements and cameras. Rebuild
  all visible forms from quad-dominant pattern surfaces with explicit front,
  back, gusset, cuff, hem, and attachment edges. Include macro hair masses in
  the next silhouette gate; use a low pooled skirt with tucked pods and no
  exposed neck.
