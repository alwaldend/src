# A99 harness static validation

Validation was deliberately source-only. No Blender or Xvfb process was
started, no `.blend` was copied or opened, no mapped window was queried, no
pixels were captured, and no XTest call or pointer input occurred.

## Bound baseline and exact delta

The four executable A98 baseline digests were rechecked before copying:

```text
62e38212551efd0f516c755cbe73743ca22a5a9c78d6ff6e701c22224bf82767  harness.sh
c5c645125467c7cc163b3af949f7cd8f0d8791520fd5c8552df0bb7055d99683  session.py
4470c880ab60ccab1bc45a8c0e2bd79a7dd503c6e1626541222529ee51c2357e  xtest_drag.c
7e40f5053e581e5290fd90ccb5e2f49a3aeb9755fa8123815f0897a8b4b1ec9b  x11_capture.c
```

An independent `diff -u` of A98 and A99 executable sources produced exactly
these two hunks:

```diff
--- A98/harness.sh
+++ A99/harness.sh
@@ -185,6 +185,7 @@
             bazel_agent run //tools/blender:blender -- \
             --factory-startup \
             --disable-autoexec \
+            "$WORKING_COPY" \
             --python-exit-code 2 \
             --python "$SCRIPT_DIR/session.py" \
--- A98/session.py
+++ A99/session.py
@@ -441,11 +441,6 @@
     if sha256_file(WORKING_COPY) != EXPECTED_INPUT_SHA256:
         raise RuntimeError("working copy is not the exact A98 input")
-    result = bpy.ops.wm.open_mainfile(
-        filepath=str(WORKING_COPY), load_ui=False, use_scripts=False
-    )
-    if "FINISHED" not in result:
-        raise RuntimeError(f"deferred file open failed: {sorted(result)}")
     if Path(bpy.data.filepath).resolve() != WORKING_COPY:
         raise RuntimeError("Blender opened the wrong working copy")
     if sha256_file(WORKING_COPY) != EXPECTED_INPUT_SHA256:
```

`cmp` confirms both C sources are byte-identical to A98. No other executable
source differs.

## Checks

- The exact bound input was read only by `sha256sum` and matched
  `02dd81b24a23a135462044c8b15a7498f743442f71d4de05ae21dae8ba9a1331`.
- `bash -n harness.sh`: pass.
- Python `ast.parse(...)` plus `compile(...)` of `session.py`: pass without
  importing `bpy`, executing the script, or writing bytecode.
- Python AST contains no class definition and no `open_mainfile` call: pass.
- Strict C11 compile/link of `xtest_drag.c` with `-Wall -Wextra -Werror`,
  `libXtst.so.6`, `libX11`, and `libm`: pass.
- Strict C11 compile/link of `x11_capture.c` with `-Wall -Wextra -Werror` and
  `libX11`: pass.
- The task-local strict-build binaries were removed after compilation and were
  never run; runtime binaries can be created only by a future explicit
  `prepare` or `start` action under `harness/run/bin/`.
- Source scan across all four executable sources found no modal-handler add,
  Blender synthetic sculpt-stroke call, XRecord observer, or old sentinel
  class.
- The launcher contains the exact positional `"$WORKING_COPY"` before
  `--python`; `session.py` retains the loaded-filepath assertion and both
  working-copy hash assertions.
- Source attestation retains native Grab `SCENE / 0.050 m / 0.40`.
- `harness.sh` executable mode is `0755`.

## Frozen A99 source digests

```text
da741f58915c6a2fb4e38d5d2150f86250014bf35dbcae27fa5ae89a1b5a5676  harness.sh
c36b836b8068e239a6f9bb52792b97ae03479bdc752f6f6a53df889db05770d6  session.py
4470c880ab60ccab1bc45a8c0e2bd79a7dd503c6e1626541222529ee51c2357e  xtest_drag.c
7e40f5053e581e5290fd90ccb5e2f49a3aeb9755fa8123815f0897a8b4b1ec9b  x11_capture.c
fbe9330aba2f3db3e0b9a71f928a23e900fbc7e575970858a75a7a22b016683f  HANDOFF.md
```

Runtime file loading, foreground window preservation, native sculpt behavior,
external X11 pixel validity, pointer delivery, post-release timer resumption,
geometry effect, fixed render effect, and all A99 gates remain deliberately
unclaimed until the coordinator performs the one permitted launch.
