# A94 pre-event startup diagnostics

No XTest pointer event or geometry mutation occurred during these checks.

- Private Xvfb `:94` started successfully at `2560 x 1440 x 24`.
- The repository-pinned wrapper reports Blender `5.2.1 LTS`.
- Default foreground GPU startup stalled before a mapped private-X11 window.
- With Wayland unset and Mesa llvmpipe forced, a factory-startup probe reached
  `background=False` and initialized OpenGL successfully.
- Positional loading of the exact A71 working-copy SHA-256
  `0004a3f0bc4987a250f7028b8697a7f740c70866ad8b60f7181e2b4eafa96400`
  printed the `Read blend` receipt but did not reach a later Python expression
  or mapped X11 window. The bounded probe exited with timeout status `124`
  after 12 seconds.
- The next pre-event transport correction starts at factory settings so the
  X11 window can map, then opens the exact working copy from a deferred timer
  with `load_ui=False`. The one-gesture/no-tuning gate is unchanged.
- The first deferred-open launch did not materialize either a mapped window or
  evidence. The next diagnostic adds `SCRIPT_LOADED` and `TIMER_ENTERED`
  receipts and delays the read-only file open for three seconds. This remains
  before the sentinel/event boundary.
- Removing `--window-geometry` allowed the script and timer to execute. The
  exact file opened, then preparation stopped before Sculpt mode because the
  A71 source carries a non-identity object matrix. The fatal receipt is
  `plus_y/evidence/pre_event_target_identity_failure.json`; immediate Blender
  shutdown also wrote `plus_y/evidence/pre_event_quit_crash.txt`.
- The A76 protocol requires the task-owned target to have an identity world
  matrix and applied scale. The next correction bakes the copied target's
  matrix into only its copied mesh, preserving its world-space surface while
  leaving the A71 source object and data exact. This is target isolation, not
  an event or stroke retry.
- The next pre-event stop exposed a Blender 5.2.1 render-engine identifier
  mismatch before the renderer or Sculpt mode ran: the pinned build exposes
  `BLENDER_EEVEE`, not `BLENDER_EEVEE_NEXT`. Its preserved receipt is
  `plus_y/evidence/pre_event_engine_failure.json`. The setup now selects the
  advertised engine and records exact pre/post world-coordinate digests for
  disposable-target matrix normalization. Failure shutdown is deferred to a
  later timer to avoid quitting from inside `open_mainfile`'s timer callback.
- The fixed before packet then rendered successfully, but the Sculpting
  workspace rejected a context-free `object.select_all` call before Sculpt
  mode. The receipt is
  `plus_y/evidence/pre_event_select_context_failure.json`. Selection setup now
  uses direct per-object selection state before constructing the exact
  `VIEW_3D` override. No sentinel or pointer event had been armed.
- The next preparation reached the configured right view, then a Python
  identity comparison incorrectly rejected Blender's false-valued RNA boolean
  as non-orthographic. The preserved receipt is
  `plus_y/evidence/pre_event_ortho_bool_failure.json`. The check now uses the
  boolean value. Because the exact normalized target and deterministic cameras
  were unchanged, the next setup may reuse the already frozen three-view
  before packet and records that reuse explicitly.
- The corrected boolean exposed the actual state: this workspace retained a
  perspective view after `view_axis(RIGHT)`. The preserved receipt is
  `plus_y/evidence/pre_event_ortho_state_failure.json`. Preparation now sets
  `RegionView3D.view_perspective = ORTHO` explicitly after selecting the right
  axis, as required by the controlling protocol. The event boundary remained
  untouched.
- The first explicit assignment was still overwritten by the asynchronous
  view-axis transition before the later verification. The preserved receipt
  is `plus_y/evidence/pre_event_view_transition_failure.json`. Projection is
  now set in its own later timer after that transition, redrawn, and verified
  only in a subsequent timer, following the repository's settled-state
  protocol. No pointer event was delivered.
- A later projection assignment was again not retained through the configured
  smooth transition. The receipt is
  `plus_y/evidence/pre_event_projection_persistence_failure.json`. The
  task-local factory profile now disables smooth-view animation before
  `view_axis(RIGHT)` and verification binds the explicit
  `view_perspective == ORTHO` enum. This changes only disposable UI transport.
- Orthographic verification then passed. Brush verification stopped only
  because Blender stores the requested decimal radius/strength as finite-
  precision RNA floats while the check required `1e-9` equality. The receipt
  is `plus_y/evidence/pre_event_brush_float_precision_failure.json`. The gate
  now allows `1e-6` numerical representation tolerance and still reports the
  observed values; requested settings remain 50 mm and 0.40.
- All view/brush checks then passed, but projecting the planned `+Y` world
  segment produced a zero-length screen path, proving the view direction had
  not settled to the requested right axis. The receipt is
  `plus_y/evidence/pre_event_axis_projection_failure.json`. The later settled
  view phase now binds Blender's canonical right-view quaternion directly in
  addition to `ORTHO`; the world endpoints and 5 mm distance remain unchanged.
- The initial `view_axis` transition still overwrote that later rotation. The
  receipt is
  `plus_y/evidence/pre_event_direct_axis_persistence_failure.json`. The setup
  now removes the transitional operator entirely, binds the canonical right
  quaternion before `ORTHO` in two separated phases, and verifies both the
  quaternion and its projection: `+Y` must be visible while `+X` must remain
  the depth axis. Endpoint geometry remains unchanged.
- The forward cyclic quaternion `(0.5, 0.5, 0.5, 0.5)` persisted exactly but
  projected `+Y` to zero pixels; in this `RegionView3D` convention it maps
  world `+Y` to view depth. The receipt is
  `plus_y/evidence/pre_event_quaternion_convention_failure.json`. The inverse
  cyclic rotation `(0.5, -0.5, -0.5, -0.5)` maps world `+X` to view depth,
  `+Y` to screen X, and `+Z` to screen Y. The independent projection-axis
  gates remain authoritative.
- The stale-projection receipt was explained by a deeper backend fact: with
  `WAYLAND_DISPLAY` merely unset, the Wayland client defaulted to
  `wayland-0`, so Blender never mapped a window on private Xvfb. A minimal
  probe with `WAYLAND_DISPLAY=a94-invalid` forced fallback to X11 and the
  task-local probe then bound the mapped `2560 x 1440` Blender window.
- The real run now uses that explicit invalid Wayland endpoint, keeps
  `DISPLAY=:94`, and calls `RegionView3D.update()` after both view assignments
  and again before projection. With a live X11 view, the protocol's canonical
  right-view quaternion `(0.5, 0.5, 0.5, 0.5)` is restored; the independent
  `+Y`-visible / `+X`-depth gate will determine correctness before input.
