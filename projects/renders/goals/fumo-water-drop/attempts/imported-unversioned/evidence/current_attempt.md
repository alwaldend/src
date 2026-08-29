# Fumo water-drop current attempt

[Back to goal](README.md)

## Attempt 00 — interface and scale preflight

### Frozen scene/asset boundary

The reusable source asset must expose exactly one collection named `FUMO`.
That collection has one top-level empty named `Fumo_Root`, at its local origin
with unit scale, whose `Z = 0` plane is the lowest intended support plane. All
renderable asset objects are descendants of that root. Asset units are meters,
`Z` is up, and the character faces local `-Y`.

`Fumo_Rig` is the unique armature descendant of `Fumo_Root`. It exposes the
existing neutral action `Fumo_Seated`. A child collection named
`FUMO_COLLIDERS` contains a closed, outward-wound, low-detail deforming mesh
named `Fumo_CollisionProxy`. The collider follows the same armature or root
motion as the visible asset and has no liquid, camera, light, tank, or scene
control data.

The shot appends a local copy of `FUMO`, validates it, and parents only
`Fumo_Root` to the scene-owned empty `FUMO_DROP_CTRL`. Ballistic keys live on
that control. The shot does not edit the reusable action, rig hierarchy,
weights, materials, mesh, or source file.

Before integration, a separate collection named
`FUMO_PLACEHOLDER__NOT_ASSET` may contain only neutral proxy geometry, a
conservative wire envelope, and visible text reading
`NEUTRAL PLACEHOLDER — NOT FUMO`. It may not be named `FUMO`, contain
`Fumo_Rig`, or imitate character identity.

### Scale and layout

- Conservative Fumo envelope: `0.30 x 0.22 x 0.25 m` (`X/Y/Z`), root at the
  bottom-center of the box. Only the exact `0.25 m` height is source-measured;
  width and depth deliberately overbound the unfinished asset.
- Tank interior: `0.80 x 0.65 x 0.65 m`; base at `Z = 0`, centered on X/Y.
  A later render shell may add `0.02 m` wall thickness without changing the
  frozen interior.
- Rest water surface: `Z = 0.36 m`, leaving `0.36 m` depth and `0.29 m` of
  wall above the water.
- Liquid-domain bounds: X `[-0.40, 0.40]`, Y `[-0.325, 0.325]`, and Z
  `[0.01, 0.71]`. The extra `0.06 m` above the tank wall is diagnostic splash
  headroom, not a promise that the final splash fits.
- Preview maximum resolution: `48`, for nominal voxel size `0.0167 m` along
  the `0.80 m` domain maximum. This is a collision/cache test, not final water.
  No final resolution is frozen before the approved collider exists.

### Timing

The scene is `48 fps`, frames `1–120`. Frames `1–24` hold
`FUMO_DROP_CTRL.z = 0.62 m`. Release is frame `25`; first contact is frame
`36`, when the root reaches the `0.36 m` rest surface. The `0.26 m` fall takes
`sqrt(2 * 0.26 / 9.81) = 0.2302 s = 11.05 frames`, so the chosen contact beat
is physically consistent within half a frame.

Frames `37–84` are reserved for impact/immersion and `85–120` for recovery or
settle. Attempt 00 defines no post-contact trajectory because buoyancy,
absorption, deformation, and liquid coupling are unknown.

### Camera, world, and lighting

- Hero diagnostic camera: perspective, `45 mm`, `36 mm` sensor, 16:9, near
  clip `0.01 m`, at `(0.90, -2.70, 0.90)`, aimed at `(0, 0, 0.43)`.
- World: neutral gray, strength `0.15`; no environment image.
- Key: white `1.20 m` area light at `(-1.10, -1.00, 1.80)`, aimed at
  `(0, 0, 0.36)`, `650 W`.
- Fill: white `0.90 m` area light at `(1.20, -0.80, 1.10)`, aimed at
  `(0, 0, 0.36)`, `260 W`.
- Rim: white `0.60 m` area light at `(0.40, 0.90, 1.60)`, aimed at
  `(0, 0, 0.40)`, `350 W`.

Exact lamp power remains preview-dependent, but their roles, neutral color,
and ownership are frozen. The diagnostic renders must expose tank edges,
water line, label, and placeholder silhouette rather than conceal them.

### Small falsifiable plan

1. Validate the machine-readable interface and scalar calculations.
2. Build only tank guides, domain wire, water plane, camera, three lights,
   drop control, neutral placeholder, envelope, and label in a blank temporary
   Blender file.
3. Render frames `1`, `24`, and `36`; reject on label, framing, clearance,
   path, or ownership failure.
4. After the approved asset exists, append and validate `FUMO`; do not fix the
   scene by mutating the asset.
5. Only then run one resolution-`48` liquid/collider proxy. A costly bake is
   outside Attempt 00.

### Current verdict

`NO-GO`. The generated validator passes its `15/15` implemented checks, but
absolute review finds two omissions: the sheet draws a bottom-center-root
contact envelope below the water line, and the packet does not prove protected
tracked blend hashes unchanged. Its numeric clearance and timing evidence may
be retained, but Blender Stage B is blocked.

Blender MCP returned no connection on 2026-08-30, so no `.blend` scaffold or
pixel claim was authorized through another path. The next attempt corrects and
revalidates only the preflight packet.

### Progress and process audit

The preflight converted an uncertain future animation into one inspectable
interface contract, a visual witness, and a deterministic validation report
without spending time on a fluid bake. The failed absolute review also exposed
gaps that its validator did not encode. The retained evidence advances scale,
timing, ownership, and clearance while leaving Stage B unauthorized.
