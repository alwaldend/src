# A98 runtime failure receipt

## Bound run

- Exact snapshot and task-local working copy SHA-256:
  `02dd81b24a23a135462044c8b15a7498f743442f71d4de05ae21dae8ba9a1331`
- Pinned Blender: `5.2.1 LTS`, foreground
- Private Xvfb display: `:2`
- Task-owned launcher PID: `1978883` (stopped)
- Task-owned Xvfb PID: `1978876` (stopped)

## Failure

The first and only launch reached the deferred exact-file open and rendered
both fixed baseline views. It then failed before `READY` while assigning the
Sculpting workspace:

```text
AttributeError: 'NoneType' object has no attribute 'workspace'
session.py:497: bpy.context.window.workspace =
    bpy.data.workspaces["Sculpting"]
```

`bpy.context.window` was `None` in the timer after the deferred
`bpy.ops.wm.open_mainfile(load_ui=False)` path. The exact failure is in
`harness/run/FATAL.json`, SHA-256
`5f9f737a896559cfe3d80ca70c26e6ba379737e0f1949d8bc6885eca396b613a`.

## Input boundary

No pointer input occurred. All of the following are absent:

- `READY.json`;
- the before-capture latch and external mapped-window capture;
- `state/input_committed`;
- `xtest_injector_receipt.json`;
- `INJECTOR_COMPLETE.json`;
- `DONE`; and
- all post-input artifacts.

The working copy remains byte-identical to the exact input. The immutable
input, rung003 high-water, and tracked reusable Blend rehashed after the stop
as `02dd81...`, `c538a9...`, and `489213...` respectively.

## Stop consequence

A98 prohibited harness repair or relaunch after its first launch. The run
therefore closes as `RESET`. This does not prove native sculpt incapable of a
good deformation; it proves the current autonomous no-sentinel foreground
authoring path failed its final pre-input discriminator. Under A97's declared
consequence, another generator or pointer harness retry is not authorized.

