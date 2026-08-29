# A73 rung-003 head / hair / bow asset interface

## Verdict

The protected rung has a previously audited, mechanically safe **15-object
head/hair replacement boundary**.  Use that boundary for an A73 visible-layer
coupon.  Do not hide a collection wholesale and do not leave the legacy
fringe, temple turns, cheek locks, or rear leaves visible over a replacement
mantle: those objects encode the rejected receiver and attachment geometry.

The phase distinction matters:

- A73 P0 receiver-only shape veto may hide only
  `Head_Cushion_Manual_Target` and render the new receiver alone.
- The first whole-subject **visible-layer coupon** must hide all 15 exact
  objects in the replacement table below, then build its receiver, thin
  crown/temple mantle, front leaves, cheek locks, free rear leaves, and seam
  witnesses under a new explicit collection/root.
- Preserve the seven face witnesses and the complete bow as frozen visible
  comparison context for that first coupon.  Their existing seating is not a
  valid attachment target: some pieces cross the old shell while others
  float.  After the visible layer passes, reseat the bow as one later rigid
  assembly rather than bridging the new mantle back to contradictory old
  roots.

This follows the exact A66 atomic-boundary audit and the more detailed A69
contact audit.  A70 hid only the old receiver and cap, then retained
incompatible silhouette-owning layers; its successive patches produced the
documented mask/card/helmet failure.  Repeating that boundary would not test
the A73 hypothesis.

## Immutable source and coordinates

- Source: `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/`
  `reimu_fumo_working_rung_003.blend`
- Current/source SHA-256:
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`
- Exact A66 inventory SHA-256:
  `6d60833cfb1388a6203cef987b08a44b5a35ae13539946f25b1e405652aaaeea`
- Exact A69 interface inventory SHA-256:
  `ea4c83d4f2d9128fcebe8552bb79cc6c641991a2b1d4ba07bef86a793834af43`
- Scene/view layer: `Attempt41_Manual_Head_Maquette` / `ViewLayer`.
- Units below are evaluated world-space millimetres.  X is left/right;
  negative Y is front; positive Y is rear; Z is upward.
- All listed objects are parentless and constraint-free.  Absolute world
  transforms are the current assembly mechanism.

No Blender process was needed for this report.  It is a read-only synthesis
of the existing pinned-background A66/A69 inventories and A69/A70 scripts.
The protected source hash was rechecked after inspection.

## Exact 15-object visible-layer replacement boundary

Dimensions are evaluated X x Y x Z.  Bounds are world-space
`X min..max / Y min..max / Z min..max`.  `I` means identity rotation and
scale with zero translation.

| Role | Exact object name | Dimensions (mm) | Evaluated world bounds (mm) | Transform |
| --- | --- | ---: | --- | --- |
| receiver | `Head_Cushion_Manual_Target` | 132.000 x 105.600 x 128.700 | -66.000..66.000 / -54.201..51.399 / 91.500..220.200 | T=(0,0,155.850), R=(0,0,0), S=(1,1,1) |
| crown/cap | `A44 continuous hair cap with smooth opening` | 134.040 x 107.480 x 130.740 | -67.020..67.020 / -55.062..52.419 / 90.480..221.220 | I |
| left forehead leaf | `A44 left temple fringe panel` | 35.582 x 10.031 x 42.516 | -55.992..-20.411 / -55.413..-45.382 / 149.970..192.486 | I |
| left temple turn | `A44 left temple transition panel` | 24.403 x 39.542 x 64.239 | -62.517..-38.113 / -51.108..-11.566 / 101.867..166.106 | I |
| central bang | `A44 off-center main bang panel` | 60.271 x 8.942 x 58.241 | -41.998..18.273 / -57.903..-48.961 / 137.843..196.084 | I |
| right forehead leaf | `A44 right swept fringe panel` | 47.011 x 10.964 x 42.915 | 8.143..55.155 / -56.571..-45.607 / 152.015..194.930 | I |
| right temple turn | `A44 right temple transition panel` | 23.401 x 39.443 x 64.097 | 39.090..62.491 / -51.084..-11.641 / 101.955..166.052 | I |
| left cheek lock | `A45 left tapered flexible cheek lock` | 25.350 x 31.051 x 77.607 | -61.398..-36.048 / -62.116..-31.065 / 82.432..160.039 | I |
| right cheek lock | `A45 right tapered flexible cheek lock` | 29.310 x 31.050 x 77.616 | 33.088..62.398 / -62.115..-31.066 / 82.423..160.039 | I |
| left rear leaf | `A42 Left asymmetric rear lock` | 53.691 x 7.373 x 118.740 | -64.604..-10.912 / 46.371..53.744 / 91.953..210.692 | I |
| central rear leaf | `A42 Off-center main rear lock` | 67.199 x 6.545 x 133.062 | -31.183..36.016 / 50.982..57.527 / 82.886..215.948 | I |
| short right rear leaf | `A42 Short right rear lock` | 54.921 x 6.750 x 109.365 | 6.692..61.613 / 47.514..54.265 / 99.934..209.299 | I |
| left rear seam | `A42 Main lock left seated seam` | 19.333 x 2.349 x 91.099 | -31.175..-11.841 / 55.525..57.875 / 98.927..190.025 | I |
| right rear seam | `A42 Main lock right seated seam` | 11.124 x 2.458 x 77.124 | 24.831..35.955 / 56.121..58.579 / 112.941..190.065 | I |
| crown seam | `Subtle crown center seam` | 0.200 x 29.007 x 42.008 | -0.100..0.100 / -54.689..-25.682 / 177.086..219.094 | I |

Mechanical facts from the audit:

- all 15 objects and their datablocks are single-user;
- none has a parent, child, constraint, driver, shape key, vertex group, or
  object-targeted modifier;
- do not delete or rename them in the coupon; set `hide_viewport` and
  `hide_render` by exact name and preserve them for rollback;
- both A45 cheek-lock meshes also contain their red tie and white-trim faces,
  so replacement locks must recreate those material regions or explicitly
  transfer them before hiding the old meshes;
- use new single-user experimental materials.  The current hair, face, red,
  and trim materials have consumers outside this boundary.

## Frozen face witnesses

Keep these exact objects visible and byte/transforms unchanged for the first
visible-layer coupon.  All have identity transforms.

| Exact object name | Dimensions (mm) | Evaluated world bounds (mm) |
| --- | ---: | --- |
| `A45 left flush composite eye applique` | 28.839 x 5.970 x 23.229 | -39.940..-11.102 / -54.184..-48.214 / 124.769..147.998 |
| `A45 right flush composite eye applique` | 28.838 x 5.977 x 23.229 | 11.102..39.940 / -54.187..-48.211 / 124.769..147.998 |
| `A45 left drooped half-lid stitch` | 17.473 x 3.533 x 1.400 | -36.375..-18.902 / -53.979..-50.445 / 149.048..150.447 |
| `A45 right drooped half-lid stitch` | 17.473 x 3.532 x 1.400 | 18.902..36.375 / -53.977..-50.445 / 149.048..150.447 |
| `A45 left fine upper expression stitch` | 24.571 x 4.947 x 1.220 | -39.572..-15.001 / -54.201..-49.254 / 157.940..159.161 |
| `A45 right fine upper expression stitch` | 24.571 x 4.946 x 1.220 | 15.001..39.572 / -54.201..-49.254 / 157.940..159.161 |
| `A44 tiny neutral embroidered mouth dash` | 3.603 x 0.194 x 0.362 | -1.802..1.802 / -46.894..-46.700 / 110.685..111.047 |

These are comparison witnesses, not final accepted placement: the eye line is
known to be low relative to the controlling reference.  Do not move it during
the A73 macro head-layer test, because that would confound the module verdict.

## Frozen bow / root witnesses

All bow/root objects are parentless.  Every object below is identity transform
except the center tie, whose exact transform is shown.  Dimensions are
evaluated X x Y x Z.

| Exact object name | Dimensions (mm) | Transform / world location |
| --- | ---: | --- |
| `A42 flattened gathered center tie` | 17.500 x 8.000 x 14.000 | T=(0,52.000,220.500), R=(0,0,0), S=(0.8578431,0.83333331,0.77777785) |
| `A42 Left root fold 1` | 45.897 x 11.108 x 8.221 | I; bounds -52.292..-6.395 / 42.279..53.387 / 220.496..228.718 |
| `A42 Left root fold 2` | 47.725 x 21.293 x 3.999 | I; bounds -54.077..-6.352 / 30.346..51.639 / 215.624..219.624 |
| `A42 Right root fold 1` | 46.157 x 9.077 x 8.041 | I; bounds 6.403..52.560 / 45.294..54.371 / 220.516..228.557 |
| `A42 Right root fold 2` | 47.713 x 18.706 x 4.077 | I; bounds 6.358..54.071 / 35.098..53.804 / 215.567..219.644 |
| `A42 Left constructed bow loop` | 123.200 x 53.950 x 61.480 | I; bounds -128.000..-4.800 / 20.884..74.835 / 181.565..243.045 |
| `A42 Right constructed bow loop` | 123.200 x 59.677 x 61.200 | I; bounds 4.800..128.000 / 22.273..81.950 / 187.142..248.342 |
| `A42 Left independent draped bow tail` | 119.767 x 18.471 x 98.127 | I; bounds -130.260..-10.493 / 51.633..70.104 / 123.350..221.477 |
| `A42 Right independent draped bow tail` | 119.566 x 18.471 x 95.957 | I; bounds 10.512..130.078 / 51.633..70.104 / 125.522..221.478 |
| `A42 Left narrow gathered loop ruffle` | 127.802 x 55.178 x 64.759 | I; bounds -131.387..-3.585 / 19.980..75.158 / 179.257..244.016 |
| `A42 Left narrow gathered tail ruffle` | 122.739 x 18.866 x 99.899 | I; bounds -133.180..-10.441 / 51.728..70.593 / 123.329..223.227 |
| `A42 Left white zigzag applique` | 21.055 x 26.251 x 41.727 | I; bounds -111.964..-90.909 / 23.555..49.806 / 194.041..235.768 |
| `A42 Right narrow gathered loop ruffle` | 127.743 x 62.624 x 64.901 | I; bounds 3.643..131.386 / 19.649..82.273 / 184.635..249.536 |
| `A42 Right narrow gathered tail ruffle` | 121.263 x 18.866 x 98.828 | I; bounds 10.432..131.695 / 51.727..70.593 / 125.417..224.245 |
| `A42 Right white zigzag applique` | 20.996 x 29.135 x 41.101 | I; bounds 90.970..111.966 / 26.384..55.520 / 197.735..238.836 |

Do not hide only the red bow meshes: loop ruffles also intersect the old head
and must be evaluated with their receiving crown.  Keep the complete bow
visible for the first regression packet, then either preserve it unchanged or
move the whole bow/root/trim group through one explicit rigid transform.

## Existing attachment and root coordinates

These are evaluated **source-object vertices within 1 mm of the old visible
cap**, not inferred origin points.  They are useful registration witnesses and
diagnostics.  They are not mandatory target coordinates for the new mantle;
forcing a shallower mantle to bridge back to them repeats A69/A70.

| Interface | Samples | Source root/contact band, world mm (X / Y / Z) |
| --- | ---: | --- |
| left rear leaf -> cap | 534 | -42.08..-14.05 / 48.48..52.44 / 127.12..192.83 |
| central rear leaf -> cap | 450 | -26.74..28.82 / 51.73..52.87 / 129.35..184.31 |
| right rear leaf -> cap | 425 | 13.00..39.38 / 49.56..52.69 / 132.82..190.31 |
| left bow loop -> cap | 20 | -51.74..-14.19 / 27.27..44.16 / 199.68..212.01 |
| right bow loop -> cap | 17 | 20.88..51.43 / 29.52..44.71 / 201.79..211.82 |
| left cheek lock -> cap | 7 | -60.84..-58.26 / -40.05..-38.54 / 147.85..148.16 |
| right cheek lock -> cap | 24 | 59.09..62.40 / -40.05..-31.07 / 147.84..160.04 |
| crown center seam -> cap | 312 | about 0 / -54.30..-25.68 / 179.53..219.09 |

Current interface state is contradictory and should be treated as negative
evidence:

- left/right temple turns cross the cap (215/199 evaluated triangle pairs)
  and receiver (112/119 pairs);
- left/right cheek locks cross both cap (50/64 pairs) and receiver (60/73);
- rear leaves cross the cap (401/164/328 pairs); the left and right also cross
  the receiver, while the central leaf floats about 0.913 mm from it;
- bow loops cross cap (156/133 pairs) and receiver (68/59); loop ruffles also
  cross both;
- the nominal center tie floats 5.964 mm from the cap and 6.981 mm from the
  receiver; inner root folds float 3.589/4.696 mm and outer folds
  13.111/13.886 mm from the cap;
- forehead leaf-to-cap spacing alternates between gaps and crossings:
  left fringe 0.955 mm gap, left temple turn crossing, main bang 0.244 mm gap,
  right fringe 1.133 mm gap, right temple turn crossing;
- the cap itself is only about 0.529 mm from the receiver and has a 0.48 mm
  centered Solidify thickness, confirming the rejected nested-shell/helmet
  construction.

## Hidden legacy exclusions

Leave these 41 superseded objects hidden and untouched.  In particular,
`Canonical continuous hair cap` is not the visible A44 cap and must never be
selected by partial name.

`A42 Left restrained catchlight`, `A42 Left rounded dark eye`,
`A42 Left rounded red inset`, `A42 Right restrained catchlight`,
`A42 Right rounded dark eye`, `A42 Right rounded red inset`,
`A44 left cheek band lower white trim`,
`A44 left cheek band upper white trim`, `A44 left fine half-lid stitch`,
`A44 left fine upper expression stitch`, `A44 left red cheek band`,
`A44 left restrained eye catchlight`, `A44 left shallow dark eye applique`,
`A44 left shallow red eye inset`, `A44 left short broad cheek lock`,
`A44 right cheek band lower white trim`,
`A44 right cheek band upper white trim`, `A44 right fine half-lid stitch`,
`A44 right fine upper expression stitch`, `A44 right red cheek band`,
`A44 right restrained eye catchlight`, `A44 right shallow dark eye applique`,
`A44 right shallow red eye inset`, `A44 right short broad cheek lock`,
`Canonical continuous hair cap`, `Left half lid`, `Left lock red tie`,
`Left round eye highlight`, `Left smooth dark eye`, `Left smooth red eye`,
`Left soft cheek lock`, `Left upper expression`, `Right half lid`,
`Right lock red tie`, `Right round eye highlight`, `Right smooth dark eye`,
`Right smooth red eye`, `Right soft cheek lock`, `Right upper expression`,
`Rolled canonical fringe strip`, and `Tiny mouth`.

## A73 coupon construction boundary

1. Append the exact protected scene into a disposable candidate and verify the
   source hash.
2. For receiver-only P0, hide the old receiver/cap and render the new A73 loft
   alone.  Do not claim integration progress from this packet.
3. If P0 passes, start the visible-layer coupon from another exact parent copy.
   Hide all 15 allowlisted objects by exact name, never by collection.
4. Add one new root and replacement collection.  Derive the thin mantle and
   separate leaves from the approved receiver's named surface rails; do not
   use the old world roots as geometric destinations.
5. Keep face and bow witnesses visible and unchanged for the first whole-model
   packet.  A bald gap, exposed receiver, floating/crossing root, or any hidden
   bridge is a reset, not a reason to thicken the mantle.
6. If the head layer is preferred to rung 003 but the bow seating fails, save
   that evidence and open a bounded rigid bow-reseat coupon.  Do not reshape
   the head back into a helmet to preserve the old contacts.
7. Preserve every source object as hidden rollback context and promote only
   explicit approved replacement object/data pairs.

Raw compact extracts used while assembling this report are beside it:
`editable_15.tsv`, `frozen_witness_26.tsv`,
`complete_head_hair_bow_inventory.tsv`, and `object_map.tsv`.
