# Plan

Test the authoring route, not the Reimu candidate.

1. Preserve the failed pinned launch as evidence: the MCP connected, but its
   only Blender window remained `0x0`; screenshot capture failed, and opening
   another main window crashed Blender.
2. Start the already-installed Flatpak Blender 5.1 in its proven graphical
   profile with the loopback MCP bridge. This changes the diagnosed GUI-host
   mechanism while retaining the same `Window.event_simulate` hypothesis.
3. Create and save an untouched disposable sculpt sphere.
4. In a mapped `VIEW_3D`, enqueue one native press/move/release sequence and
   measure vertex-coordinate change. Do not use `brush_stroke` directly.
5. Save the deformed candidate under a new name. Clean-reopen those exact
   bytes with repository-pinned Blender 5.2.1 and produce a verification
   render and geometry receipt.
6. Preserve source bytes, candidate bytes, failure/success metrics, and the
   render digest. Do not open or edit a Reimu model in this attempt.

Stop if the mapped host cannot produce event-driven deformation or if the
saved candidate cannot be reopened and rendered by pinned Blender. Another
corrective trial is allowed only for a newly evidenced failure mechanism; an
equivalent retry is prohibited.
