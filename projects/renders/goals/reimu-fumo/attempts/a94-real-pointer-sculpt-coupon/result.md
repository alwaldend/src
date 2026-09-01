# A94 result — reset on raw-event provenance failure

## Verdict

`RESET`. The first required branch, world `+Y`, failed the hard event-
provenance gate. The `-X` and `-Z` branches were not started. This authoring
modality is not authorized for a whole-plush sculpt by A94.

No Reimu Fumo acceptance criterion was tested or passed. This was an
interface coupon, and the missing raw receipt prevents any geometry or visual
effect from being admissible evidence.

## What passed before input

- Repository-pinned Blender `5.2.1 LTS` ran foreground on a private
  `2560 x 1440` Xvfb through Mesa llvmpipe.
- Wayland was made deliberately unavailable so the task-local X11 probe bound
  mapped window `2097154` exactly.
- The exact A71 source SHA-256 was copied byte-for-byte, deep-copied to a
  single-user target, and the source remained hidden and unchanged.
- Baking the copied target's source transform into only its private mesh
  preserved every world coordinate exactly: both pre/post coordinate digests
  are
  `8893cfd6d8b5d96742fcef2f6ea8fc13750eb1418e40f9e9c6ff4197881a334b`,
  with maximum error `0.0 m`.
- The settled view was `SCULPT`, orthographic, native Essentials `Grab`,
  50 mm scene radius, strength `0.40`, no symmetry or dyntopo.
- The fixed right-view basis projected requested world `+Y` by
  `26.0417486 px` and world `+X` by `0.0 px`; the computed brush margin was
  `260.417486 px`.
- A frozen pre-event Blend, gesture plan, and window-state record were saved.
  The fixed camera renders are nonblank. However, post-stop inspection found
  that the nominal live `before_view3d.png` is entirely black. File existence
  and hashing were checked before input, but its pixels were not inspected;
  this was a missed precondition and independently invalidates the run.
  ImageMagick measured the `2100 x 1389` RGB image with minimum, maximum, and
  mean all exactly `0`; the nonblank `512 x 512` fixed right render had mean
  `0.918054` and maximum `1`.

## Sole gesture and decisive failure

The task-local C injector used the XTest ABI directly and delivered exactly
one planned drag according to its independent receipt:

- mapped window: `2097154`;
- Blender window-local bottom-left path: `(1160, 791)` to `(1186, 791)`;
- X11 screen top-left path: `(1160, 648)` to `(1186, 648)`;
- eight held monotonic motions over `320 ms`;
- one button-1 press and one release;
- injector exit status `0`, wall time `0.297672 s`.

After more than 30 seconds, Blender emitted no
`event_receipt.json`, no post snapshot, and no `DONE` marker. The native sculpt
modal handler therefore prevented the pass-through sentinel from observing or
completing the required raw press/motion/release record. The injector's own
receipt proves what was requested at the X server, but it cannot substitute
for the in-region Blender event receipt required by the plan.

The black live viewport baseline means the pointer event should not have been
sent even before the sentinel failure became visible. It also makes a live
before/after pixel claim impossible. This is recorded as a process failure,
not hidden behind the valid fixed-camera renders.

Per the hard stop, no second event or tuning retry was attempted. The live
process and Xvfb were stopped, and the in-memory effect was discarded without
claiming geometry support, pixel change, direction, or visual coherence.

## Exact surviving evidence

| Artifact | SHA-256 |
| --- | --- |
| `plus_y/READY` | `36a37bf0511c14c042818e4c6ad0de33d1323bb616b2c9d189680bbaaca0caeb` |
| `plus_y/evidence/gesture_plan.json` | `b5c7727de26704deb36a25f4c5781235533d41994f232d1ff1a661f282b8bb82` |
| `plus_y/evidence/window_state_before.json` | `11340d0ad9bc760cb8d795f074c1ee09089531eb52371c54f9cd612827fe184b` |
| `plus_y/evidence/before_view3d.png` (invalid: all black) | `f0e354bbb43ea78d799f8506cccc9ec658b0aed2d85d937dfda8106bd5bf4484` |
| `plus_y/evidence/xtest_injector_receipt.json` | `fcd8e964828380b0383c9157ae55dd08828a0a2c591ad834e0dc2ab3017e328b` |
| `plus_y/snapshots/plus_y_before.blend` | `02dd81b24a23a135462044c8b15a7498f743442f71d4de05ae21dae8ba9a1331` |
| `live_coupon.py` | `8fc4ebd7b1254532a314c02865b81c62600be59ef4c4606b89551f9f40e0d520` |
| `xtest_drag.c` | `072b1f101128026454e79c07e17f05034c3045e4be7a08b857702643b1a937d4` |
| `xtest_drag` | `2115e1f9bfe1922cd63a24aef4360b954a11f9a42d9815fda4b8c2f42c710333` |

`STARTUP_DIAGNOSTICS.md` preserves the pre-event transport corrections. Those
failures occurred before `READY` and do not constitute stroke retries.

## Criteria state

Criteria `criterion-001` through `criterion-008` remain `unverified`. A94
neither changes nor supplies acceptance evidence for the reusable Fumo.

## Resume boundary

Do not rerun this XTest/modal-sentinel family as another parameter variant.
The only useful revisit would require a demonstrably different raw-event
observation interface that can coexist with Blender's native modal sculpt
handler and bind press, held motions, and release before any geometry claim.
