# Concrete A02 independent acceptance review

## Verdict

`ACCEPT` candidate
`a9e7d470ae3de3236c1b808b087ee36191dcc357fcdda821340eeeb25fa84ccf`
against criteria revision 2. Acceptance is limited to the neutral rigid
concrete-drop scaffold and replacement interface. It does not approve Fumo
integration, Reimu likeness, plush construction, deformation, plush behavior,
or a final animation.

## Criterion findings

- `criterion-001`: pass. The exact candidate clean-opens hash-stably; `FUMO`
  owns the named interface, placeholder root, proxy, front marker, and required
  custom properties.
- `criterion-002`: pass. All eight native frames show a featureless neutral
  rigid proxy and the readable warning `NEUTRAL PLACEHOLDER / RIGID DROP - NO
  PLUSH`; inventory properties limit the claim consistently.
- `criterion-003`: pass. Metric scale is 1.0, collider height is 0.25000003 m,
  visible proxy height is 0.24500003 m, and the floor contains the full sampled
  trajectory plus collider extent.
- `criterion-004`: pass. The scene uses 24 fps and frames 1--72; frames 1--12
  are held and kinematic, release occurs at 13, contact at 22, and the proxy is
  settled through frame 72.
- `criterion-005`: pass. Every native 640 by 360 frame was inspected. The
  label, complete proxy, floor, 25 cm witness, and identifiable impact region
  remain clear. Initial proxy margins are 21.238 px top and 261.564 px bottom;
  the label has approximately 24 px top/right margins and no proxy/impact
  overlap.
- `criterion-006`: pass. Sampled descent is 0.845343 m, maximum penetration is
  0.000687 m, and center-motion span over frames 60--72 is 0 m.
- `criterion-007`: pass. Candidate, parent, contact sheet, protected assets,
  and every manifest-listed artifact rehashed correctly. A02's pre-state is
  A01's final semantic snapshot and the broad snapshot diff contains only the
  warning label's X and Y location components. A01 and A02 mechanics canonical
  JSON hashes are identical.

## Scope and caveats

The mutation assigns `Placeholder Warning Label.location` exactly once,
changing `[0, -0.1150000021, -0.7200000286]` to
`[0.2170832306, 0.1609133482, -0.7200000286]`. The camera, physics, animation,
proxy, impact marker, interface, and sampled mechanics are unchanged.

The scripted late-silhouette thresholds were generated after the plan and its
mask combines proxy/front/impact pixels, so they were not used as acceptance
authority. Native-pixel inspection independently passes all eight frames. The
semantic snapshot is broad and acceptance-relevant, but should not be
described as proof that every conceivable Blender property is unchanged.

The exact contact sheet SHA-256 is
`dc2437b7151bae77be9ca29ffdbd6ccb5446acf853b4bf2064db1d765706fd84`.
Review was read-only and changed no file or runtime state.
