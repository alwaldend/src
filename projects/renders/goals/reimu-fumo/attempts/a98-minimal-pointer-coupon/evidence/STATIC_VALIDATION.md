# A98 harness static validation

Validation was deliberately source-only. No Blender or Xvfb process was
started, no `.blend` was opened or copied, no mapped window was queried, no
pixels were captured, and no XTest call or pointer input occurred.

## Post-validation delta reconciliation

An intermediate validation bound `harness.sh` at SHA-256 `4a416a5c...` and
`HANDOFF.md` at `aeb30be7...`. After that check, the shell lifecycle was
hardened with three never-removed latches for the first session launch, the
before capture, and the after capture. The handoff was updated to describe
those stop semantics. That deliberate post-check delta changed only
`harness.sh` and `HANDOFF.md`; the Blender session and both C sources did not
change.

The complete static gate was then rerun against the final frozen 435-line
`harness.sh`. It passed at `62e38212...`; the current handoff passed its digest
check at `08a6c223...`. A further reconciliation rerun produced the same five
full digests recorded below. The earlier `4a416a5c...` value is superseded and
does not identify the executable handoff.

## Checks

- Exact bound input was read only by `sha256sum` and matched the required
  SHA-256
  `02dd81b24a23a135462044c8b15a7498f743442f71d4de05ae21dae8ba9a1331`.
- `bash -n harness.sh`: pass.
- Python `ast.parse(...)` plus `compile(...)` of `session.py`: pass without
  importing `bpy`, executing the script, or writing bytecode.
- Python AST contains no class definition: pass. The Blender side therefore
  cannot register the A94-style modal sentinel class.
- `cc -std=c11 -O2 -Wall -Wextra -Werror` compile/link of
  `xtest_drag.c` against `libXtst.so.6`, `libX11`, and `libm`: pass.
- Same strict compile/link of `x11_capture.c` against `libX11`: pass.
- Source scan across the four executable sources for a modal-handler add, a
  Blender synthetic sculpt-stroke call, XRecord, alternate `-X`/`-Z`
  branches, or either old sentinel class: no match.
- Source attestation confirms both mutation and settled-state verification of
  native Grab `SCENE / 0.050 m / 0.40`, plus serialization of observed
  `locked_size` into the live gesture plan.
- `harness.sh` executable mode: `0755`.
- Strict-build check binaries were removed after compilation; runtime binaries
  are created only by the future explicit `prepare`/`start` action under
  `harness/run/bin/`.

## Validated source digests

```text
62e38212551efd0f516c755cbe73743ca22a5a9c78d6ff6e701c22224bf82767  harness.sh
c5c645125467c7cc163b3af949f7cd8f0d8791520fd5c8552df0bb7055d99683  session.py
4470c880ab60ccab1bc45a8c0e2bd79a7dd503c6e1626541222529ee51c2357e  xtest_drag.c
7e40f5053e581e5290fd90ccb5e2f49a3aeb9755fa8123815f0897a8b4b1ec9b  x11_capture.c
08a6c2237c67f3647503045265a55d2efbd2f1cf4fef6d0143f91f8aa3737244  HANDOFF.md
```

Runtime Blender APIs, factory-to-file UI transition, native modal sculpt
behavior, external X11 pixel validity, pointer delivery, post-release timer
resumption, geometry effect, fixed render effect, and all A98 acceptance gates
remain deliberately unclaimed until the coordinator executes the exact
handoff once.
