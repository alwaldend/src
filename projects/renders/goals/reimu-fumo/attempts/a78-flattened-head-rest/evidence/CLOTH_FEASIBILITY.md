# A78 cloth feasibility

## Verdict: NO-GO

Do not use deterministic Blender Cloth pressure relaxation to generate A78's
head-rest form. Close the cloth branch before geometry, cache, render, or
integration work.

This verdict rejects pressure relaxation as the **form generator** for A78. It
does not reject the broader A78 construction hypothesis: a manually authored,
static front panel, rear closure, and gusset may still be viable. That would
need differential seam lengths, darts, and locally authored face relief before
any optional weak settle. It is a materially different method and should not
be smuggled into this cloth branch as parameter tuning.

## Decisive evidence

Attempt 56B already executed the proposed experiment more directly than a new
A78 probe would. It used independently authored front and rear cuts, a welded
six-strip tapered gusset, a closed manifold mesh, zero gravity, and Blender
closed-cloth pressure relaxation at 10x working scale.

Configuration B was technically successful:

- `17 x 17` quad front and rear cuts with a `64`-edge perimeter;
- `898` vertices, `896` polygons, no boundary or zero-area faces;
- maximum edge stretch `1.0170x`;
- maximum final-eight-frame RMS step `0.065 mm`;
- finite, watertight output in `0.79 s`.

It nevertheless produced exactly A78's forbidden result: a uniformly taut
rounded foam block/mattress, with nearly constant corner radii, straight slab
sides, decorative-looking seams, and an underside sag. Its source-mask
overshoot also failed in front, side, and rear. Configuration A supplied the
other endpoint: radial panels under pressure produced a button-tufted front,
drawstring rear crater, and thick bolster. These are representation failures,
not convergence bugs.

The earlier and smaller tests agree:

- Attempt 14's connected front/rear/gusset pressure cage remained in
  large-amplitude motion at frame 90. Frame 80-to-90 RMS motion was
  `0.069870 Wh`, depth changed `0.054189 Wh`, width normalization would have
  required `14.987%`, and the form inflated to almost equal width and depth.
- Attempt 42 bracketed weak and strong relaxation on a closed sewn pocket.
  The weak solve made no readable fabric change; the underfilled solve shrank
  width `31.4%`, grew depth `161.2%`, and rolled into a concave cup.
- Attempt 56B's explicit conclusion was that uniform pressure plus one
  continuous gusset cannot independently supply the shallow face, compact
  rear, crown allocation, and lower receiving transition.

The strongest disconfirming fact is therefore not that Cloth can be unstable.
It is that the closest prior solve was stable, finite, reproducible, and still
became the exact mattress representation prohibited by A78.

## Frozen-interface incompatibility

A78 also has a harder boundary than the isolated Attempt 56B coupon. Frozen
eyes, mouth, fringe, hair aperture, and roots require positive signed
clearance across declared bands, at least 80% band coverage, and no new
triangle crossing or beige island.

Uniform pressure has no term for those attachment-band objectives. Leaving
the face free permits global drift and depth inflation. Pinning the witness
neighborhoods preserves their initial locations but localizes strain at pin
falloffs, making puckers, a cavity ring, or a taut plate likely. Collision can
prevent some penetration, but it does not establish flush seating, positive
clearance throughout a band, or the required coverage. Adding custom spatial
constraints or differential pattern forces would create another solver
family, not validate this one.

The method therefore cannot simultaneously promise:

1. shallow bounded depth and a broad nonplanar cheek field;
2. compact rear support instead of a mattress or balloon;
3. preserved frozen hair and facial interfaces; and
4. readable differential sewn construction.

## Exact reproduced settings

No new A78 simulation should run. The following settings identify the stable
prior control that already falsifies the branch and should be cited rather
than recreated:

| Setting | Attempt 56B configuration B |
| --- | ---: |
| Working scale | `10.0x` |
| Frames | `1..56` |
| Quality | `10` |
| Time scale | `0.62` |
| Mass | `0.12` |
| Air damping | `7.0` |
| Bending model | `ANGULAR` |
| Tension / compression | `30.0 / 25.0` |
| Shear / bending | `21.0 / 1.05` |
| Tension / compression damping | `8.0 / 8.0` |
| Shear / bending damping | `7.0 / 1.5` |
| Pressure | enabled, `0.78` |
| Pressure factor | `1.0` |
| Gravity | `0.0` |
| Object collision / self-collision | disabled / disabled |

The nonconvergent Attempt 14 endpoint used `quality=12`,
`time_scale=0.65`, `mass=0.20`, `air_damping=3`, stiffness
`35/25/20/0.18`, damping `10/10/10/2`, pressure `1.5`, pressure factor `8`,
target volume `1.035x`, self-collision, and frames `1..90`. Moving toward that
endpoint adds inflation and instability rather than the missing construction.

## Probe and cache/apply decision

The A78 probe plan is deliberately empty:

1. Do not create a Cloth modifier or candidate mesh.
2. Do not bake or retain a point cache.
3. Do not alter the frozen rung-003 source or any goal record.
4. Preserve A56B and A14 as negative evidence and close this method.

If a future, separately authorized attempt materially changes the pattern
construction and needs only a weak post-form settle, use the established
reproducible freeze protocol:

1. Start pinned Blender `5.2.1` in background factory mode with auto-exec
   disabled; work only in a disposable copy under its own `out/` directory.
2. Set every Cloth and point-cache value explicitly, set cache frames before
   evaluation, clear stale cache state, and advance frames sequentially from
   frame 1. Record per-frame coordinate hashes and the final eight RMS steps.
3. Repeat the solve in a clean process and require identical topology and
   coordinate hashes (or a predeclared sub-micron numeric tolerance) before
   judging pixels.
4. At the accepted frame, freeze the evaluated dependency-graph mesh with
   `bpy.data.meshes.new_from_object(...)`. Do not make a live Cloth modifier
   or cache part of the candidate contract.
5. Save a modifier-free checkpoint, reopen it with the same pinned Blender,
   and verify coordinate hash, bounds, manifoldness, winding, signed volume,
   and every frozen-interface clearance before rendering.

That protocol can make a chosen evaluated state reproducible; it cannot make
uniform pressure discover the missing manufactured shape.

## Likely failure modes

| Failure | Expected symptom | Existing evidence |
| --- | --- | --- |
| Uniform-pressure mattress | Taut face, constant corner radius, slab side | A56B B |
| Pole or pin pucker | Tuft, crater, cavity ring | A56B A |
| Global inflation and drift | Depth approaches width; no convergence | A14 |
| Weak-settle null result | Only subpixel shading changes | A42 v1 |
| Underfilled collapse | Width loss, cupping, rolled edge | A42 v2 |
| Interface strain localization | Puckers around pinned eyes, mouth, roots | Mechanism plus A56B/A14 endpoints |
| False contact pass | No penetration but poor seating or band coverage | Collision cannot encode A78 gates |
| Cache-dependent artifact | Different frame/history or stale result | Avoided only by clean sequential replay and freeze |

## Time estimate

Closing the branch requires no geometry time. A replay of the already-failed
A56B control would take under a minute of solver time plus roughly one hour to
package and verify, while an integrated interface-aware probe would take
several hours and would still repeat a falsified representation. That time is
not justified. Spend the next attempt on an art-directed static panel/dart
construction, with pressure allowed only after the neutral form already passes
the A78 mechanical and visual gates.

## Evidence consulted

- `out/reimu_fumo_attempt_014_cloth_panel_sack/simulation_probe.json`
- `out/reimu_fumo_goal_migration_exact/reimu-fumo/attempt_14.md`
- `out/reimu_fumo_attempt_042_cloth_relaxation_study/README.md`
- `out/reimu_fumo_attempt_056b_sewn_head_cushion/REPORT.md`
- `out/reimu_fumo_attempt_056b_sewn_head_cushion/config_b/`
  `pattern_solver_manifest.json`
- `out/reimu_fumo_attempt_056b_sewn_head_cushion/build_solve.py`
- `out/reimu_fumo_attempt_067_head_topology/TERMINAL_REPORT.md`
- `out/reimu_fumo_attempt_078_head_rest/decision_review.md`
- `out/reimu_fumo_attempt_078_head_rest/goal_checkpoint/plan.md`
