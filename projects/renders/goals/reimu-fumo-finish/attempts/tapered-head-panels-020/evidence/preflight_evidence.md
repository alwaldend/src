# Explicit mesh editing preflight

The sole writer is `root-astra-reference-cage`: Blender 5.1.1 in the
already-installed Flatpak, X11 window 1280 by 900, task-isolated unsaved
preferences, existing MCP extension, listener `127.0.0.1:9876`. The host is
the current task worktree and the live process reports PID 2 inside its
Flatpak namespace. The verifier is repository-pinned Blender 5.2.1 LTS,
build `9e2066aef7ef`, through `//tools/blender:blender`.

The native control tools function. Desktop Use exposed a frame tree without
pixels; the Blender screenshot had display artifacts. Saved-file rendering
with pinned Blender is the visual feedback source for this route. No host
preferences or configuration were saved. Enabling the installed extension
needed `default_set=True` to create its in-memory preferences before register.

A selection-based setup test stopped before deformation because selected
faces repopulated vertex selection on entering Edit Mode. Its following undo
returned to the startup cube. It saved no candidate and establishes neither
an organic-shaping result nor failure to undo a completed deformation.
The chosen route consumes explicit mesh coordinates and immutable snapshots,
so it does not rely on that selection path or native undo.

## Measured result

A disposable Blender grid created with `x_subdivisions=5`,
`y_subdivisions=5`, `size=0.2` has 36 vertices. A subdivision modifier at
level 2 and solidify thickness 0.004 make its local deformation visible.
The writer changed only vertex 14 by +0.015 m in Z through its mesh data,
updated the mesh, and saved a fresh candidate. The pinned verifier reopened
both exact files and independently compared stored coordinates.

- Baseline: `ae6bee49cbef21c538ff7f347f299eaa2cb9864928374074212442dfb1e9ccbc`.
- Candidate: `c7fbe514ad54da439a9732a55fbd2e846e2bd61e9c2d27be2e4252a8ec4f7d67`.
- Changed vertices: exactly `[14]`; all 35 others byte-for-byte equal as
  coordinate tuples.
- Stored Z delta: `0.014999999664723873` m.
- Baseline file SHA-256 unchanged after verification.
- Meaningful 320 by 320 workbench render inspected: a softly thickened square
  panel with a localized raised area.
- Render SHA-256:
  `fbf831c7d6302d9546291994784c996503d3c5e342da91322308ee56d3fa069d`.

The ignored raw packet is `out/reimu_fumo_finish/desktop_astra/`; its
`verify_preflight.py` and JSON receipt contain the local audit. Regenerate
the stated test by creating the grid and modifiers above, saving baseline,
applying the one explicit coordinate edit, saving a new file, and comparing
both files in pinned Blender. Image hash comparison requires the same graphics
environment. This proves the consumed edit/save/reopen/render contracts only;
it grants no visual or asset acceptance criterion.
