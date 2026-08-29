# Blender execution modes

Choose the mode from the operation's actual context requirements.

## Background is the default

Prefer an existing package-owned `blender_binary` target. For an isolated
one-off script when no such target exists, use this shape from the workspace
root:

```sh
bazel_agent run //tools/blender:blender -- \
  --background \
  --factory-startup \
  --disable-autoexec \
  --python-exit-code 2 \
  --python /absolute/path/to/workspace/out/<task>/script.py \
  -- <script arguments>
```

Use absolute task input/output paths or paths deliberately resolved from the
workspace. `--factory-startup` avoids dependence on user preferences,
`--disable-autoexec` prevents automatic execution embedded in the opened
file, and `--python-exit-code` propagates uncaught script failures. Omit one of
these only when the requested operation demonstrably requires it, and record
the reason.

Batch related views or operations in one Blender process when practical.
Blender startup and scene loading are expensive; several contending render
processes can be slower and less stable than one explicit packet.

## Foreground and MCP are capability-driven

Use foreground Blender only for a UI-context operation. Before launching it:

1. Probe the available Blender MCP connection with a read-only summary.
2. If connected, verify the Blender version, `bpy.app.background`, open file,
   and required area or workspace before reuse.
3. If disconnected and foreground work is required, build or run the pinned
   `//tools/blender:blender` target through `bazel_agent`, isolate
   `BLENDER_USER_CONFIG`, `TMPDIR`, logs, and outputs under `out/<task>/`, and
   start only a loopback-scoped task listener when MCP is actually needed.
   When the user explicitly selected the installed Flatpak live host, start
   that single host under the exception below instead. Reprobe after launch;
   absence of an already-running listener is not a terminal failure.
4. Use the real display when available. Xvfb is a fallback for headless GUI
   context, not a requirement for background work or an already-running
   compatible session.

A long-lived `bazel_agent run` couples Blender to that Bazel client. Plan other
Bazel work around the session and confirm observed lock contention rather than
assuming it. A same-session generated-wrapper launch is a narrow fallback:
first build `//tools/blender:blender` through `bazel_agent`, resolve the
freshly built generic runfiles wrapper for the current checkout and
configuration, execute that wrapper rather than the extracted ELF, and record
its identity.
Never put generated paths in durable commands, reuse them across worktrees, or
apply this exception to package automation that depends on Bazel-provided
runfiles or environment variables.

Do not install an add-on, mutate a normal user profile, or use a host Blender
to make MCP convenient. If the pinned Blender cannot provide the required
interactive capability within the task's authority, preserve the evidence and
report the boundary rather than silently changing toolchains.

### Explicit Flatpak live-host exception

When the user explicitly requests the already-installed Blender Flatpak as a
persistent live MCP host, it may serve as the sole interactive mutation owner
for disposable candidates. It is never the batch verifier, renderer of record,
or deliverable toolchain. Inspect `flatpak info org.blender.Blender` first and
record its version. Start it with online access so the loopback MCP extension
can listen, probe the configured MCP endpoint read-only, and verify the live
version, foreground state, empty or task-owned open file, and workspace before
mutation. If the extension is installed but not enabled, enable it for the
current process only and start `blmcp.server_start`; do not save that preference
or install anything.

If the Flatpak is older than the pinned Blender or the candidate's authoring
version, do not directly open and overwrite the newer source. Start from an
empty live scene or an isolated import/append boundary, preserve the exact
source hash, and treat every saved live snapshot as untrusted until the pinned
background Blender clean-reopens it and proves the required datablocks and
pixels. Keep one live mutation owner, publish immutable snapshots under
`out/<task>/`, and give renderers and reviewers only those snapshots. Stop the
task-owned Flatpak and loopback listener when the live phase ends.

## Settle live UI state across calls

UI transitions can complete after a Blender MCP call returns. For an operator
that needs a specific viewport context:

1. request the workspace or screen transition in one call;
2. inspect a later call and require the settled foreground state and required
   area/region;
3. prepare the exact object, selection, mode, view, tool or brush, and full
   context override in another call;
4. execute the operator only in a later call after preparation has settled;
5. verify the operation's actual effect.

For a `VIEW_3D` operator, reconstruct the current `window`, `screen`, area,
`WINDOW` region, space, region data, scene, view layer, active object, and
object in a `temp_override`. Check `poll()` in that exact override. Do not
reuse stale areas or regions across workspace transitions.

`FINISHED` is not an effect check. Compare operation-appropriate pre/post
state: coordinates and moved-region bounds for deformation, transforms for
object movement, datablock identity for creation, file hashes for saving, or
image hashes and pixels for rendering. Reject zero change and unintended
global change when locality matters.

Stop the task-owned foreground process and listener when the interactive phase
ends. Remove only precisely identified task-created configuration residue;
never clean broad user or temporary directories.
