# A02 execution receipt

## Tool identity

- Repository target: `//tools/blender:blender`
- Blender: `5.2.1 LTS`
- Build hash: `9e2066aef7ef`
- Entry point: `bazel_agent run`
- Mode: pinned background Blender, automatic Blend-file scripts disabled

## Commands

The first launch used task-relative arguments. Blender resolved them from the
Bazel runfiles working directory, could not find the V01 Blend, and exited
before opening or mutating any asset:

```text
bazel_agent run //tools/blender:blender -- --background --disable-autoexec out/fumo_sisyphus_scene_scaffold/fumo_sisyphus_scene_scaffold_v01.blend --python-exit-code 1 --python out/fumo_sisyphus_attempt_002_macro_mask/author_macro_mask.py
```

The corrected launch used absolute task-local paths and exited zero:

```text
bazel_agent run //tools/blender:blender -- --background --disable-autoexec /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/fumo_sisyphus_scene_scaffold/fumo_sisyphus_scene_scaffold_v01.blend --python-exit-code 1 --python /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/fumo_sisyphus_attempt_002_macro_mask/author_macro_mask.py
```

The successful process performed the bounded edit, saved the candidate,
clean-reopened that exact candidate, verified its inventory and protected
hashes, and rendered the first 512-pixel decision frame. It did not save after
the clean-open verification or render.

## Preservation result

The V01 input, controlling reference, tracked Sisyphus Blend, and tracked
Reimu Blend all retained their frozen SHA-256 values. The protected boulder,
placeholder, reference, camera, and lighting object signatures were identical
before authoring and after clean reopen. See `protected_hashes.json` and
`clean_open_inventory.json`.
