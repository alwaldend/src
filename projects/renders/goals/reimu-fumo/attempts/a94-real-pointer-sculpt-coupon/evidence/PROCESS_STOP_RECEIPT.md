# A94 process-stop receipt

- Stop time: `2026-09-01T15:21:09,757147140+03:00`.
- The `+Y` private Blender session was terminated after the sole XTest drag
  produced no sentinel event receipt or post-processing phase for more than
  30 seconds.
- The Blender PTY required two `SIGQUIT` control characters; the second
  returned exit status `1`.
- The task-owned Xvfb `:94` was then terminated by `SIGQUIT` and returned exit
  status `1`. A subsequent X11 connection attempt produced no display output,
  confirming that the private display was unavailable.
- No second pointer press, motion sequence, release, brush change, camera
  change, or parameter adjustment was sent after the authorized gesture.
- No post-stroke Blend, post render, displacement report, pixel difference, or
  comparison sheet exists. Creating one would have required continuing past
  the hard provenance stop or claiming an unobserved raw event sequence.
- Post-stop visual inspection found that the nominal pre-event
  `before_view3d.png` is entirely black. File existence and hashing had been
  checked before input, but its pixels had not been inspected. It is invalid
  viewport evidence and should itself have stopped the run before the drag.
  The fixed camera renders are nonblank, but they do not repair the missing
  live viewport baseline.

Protected and control inputs remained exact after shutdown:

- A71 source and untouched session working file:
  `0004a3f0bc4987a250f7028b8697a7f740c70866ad8b60f7181e2b4eafa96400`
- Tracked reusable asset:
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`
- Rung003 comparator:
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`

The immutable pre-event snapshot remains available at
`plus_y/snapshots/plus_y_before.blend`, SHA-256
`02dd81b24a23a135462044c8b15a7498f743442f71d4de05ae21dae8ba9a1331`.
