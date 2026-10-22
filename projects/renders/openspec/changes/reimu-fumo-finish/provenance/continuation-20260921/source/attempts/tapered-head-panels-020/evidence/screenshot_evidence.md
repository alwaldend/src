# Window capture correction

Observed 2026-09-04, same task-isolated Flatpak 5.1.1 writer, loopback MCP
9876, mapped 1280x900 window. No source or host configuration changed.

Desktop Use get_app_state for mutter-x11-frames exposes only Blender's window
frame and title, not internal widgets or screenshot pixels. This is not proof
that Blender itself cannot be captured. Earlier add-on captures displayed
unusable periodic pixels; the exact corruption cause was not established.

The user's startup-popup hypothesis prompted a different capture test:
`bpy.ops.screen.screenshot(filepath=<task-local PNG>)` returned FINISHED.
The coordinator inspected a clear full window, with no popup visible. A
second capture after view-only framing showed the new fumo clearly. The
saved candidate bytes were not changed by framing the viewport.

- First native PNG: `native_window_019.png`, 403961 bytes, SHA-256
  `e49646a62b6772285e1cdfd6406e3194c3fc5b78cba14a63baf0080d36102e1f`.
- Framed native PNG: `native_fumo_019.png`, SHA-256
  `fdf016a29baae92d8780ba830ce7c0c35c6efcd6eece11bde8473a94019b0986`.
- Both remain under ignored `out/reimu_fumo_finish/desktop_astra/`.
- Subsequent add-on window screenshot with a 4000000-byte limit also returned
  readable pixels. This discriminates current capture availability from the
  earlier state; it does not prove which earlier dependency failed.

Read-only implementation inspection found that the add-on calls the same
native screen.screenshot operator, then optionally resizes/re-encodes a PNG
that exceeds its byte limit and returns base64. Thus blaming a separate
framebuffer implementation was unsupported. Timing/redraw state and optional
image processing remain hypotheses, not findings. No restart or screenshot
pipeline modification was needed to obtain the successful current captures.

A startup popup could obstruct ordinary input, but none is visible now and
no evidence establishes it as the cause of earlier corruption. Generic
Desktop Use remains semantically limited for Blender; direct Blender control
and native window screenshots are currently usable. Pinned saved-file renders
remain the acceptance source, not an unsaved viewport.

## Bounded session review

Stable issue `fumo-capture-failure-overgeneralization`, evidence tier live:
one failed capture was generalized into an unusable-window claim before the
native file-capture alternative was tested. Two native captures and one
subsequent add-on capture corrected it. The user received the correction.
Task-local remedy: try one named independent capture test before declaring
all screenshots unavailable, and preserve the exact current positive result.
No shared skill, host setting, or runtime change is proposed.

Discovery also repeated already-known skill/CLI context after compaction and
guessed two incorrect goal-tool paths. The canonical CLI target is
`//projects/goal/cmd/goal`; the delivery target is
`//tools/repo_delivery/cmd/repo_delivery:go`. Keep these exact identities in
the local CURRENT projection. Those reads were friction, not modeling work.
