# A99 coordinator static-diff verdict

## Verdict

`GO` for the single permitted launch.

The frozen A99 executable harness differs from A98 by exactly two hunks:

1. `harness.sh` inserts `"$WORKING_COPY"` positionally after
   `--disable-autoexec` and before `--python-exit-code` / `--python`.
2. `session.py` deletes only the five-line deferred `open_mainfile` and
   result-check block. The loaded filepath assertion and both byte-hash
   assertions remain.

Deleting that one new shell line recreates the A98 shell byte-for-byte;
deleting the A98 five-line deferred-open block recreates A99 `session.py`
byte-for-byte. Both C sources are identical to A98. Bash syntax, Python AST,
no-class/no-open-mainfile checks, strict C builds, file modes, and forbidden
source scans pass. `run/` was absent before launch.

Frozen hashes:

- `harness.sh`: `da741f58915c6a2fb4e38d5d2150f86250014bf35dbcae27fa5ae89a1b5a5676`
- `session.py`: `c36b836b8068e239a6f9bb52792b97ae03479bdc752f6f6a53df889db05770d6`
- `xtest_drag.c`: `4470c880ab60ccab1bc45a8c0e2bd79a7dd503c6e1626541222529ee51c2357e`
- `x11_capture.c`: `7e40f5053e581e5290fd90ccb5e2f49a3aeb9755fa8123815f0897a8b4b1ec9b`
- `HANDOFF.md`: `fbe9330aba2f3db3e0b9a71f928a23e900fbc7e575970858a75a7a22b016683f`
- `STATIC_VALIDATION.md`: `d4bb57652e32824d917a43c9f1a73ae8d1a3756654754aaeacd3d7cdcd82dfd6`

The independent verifier and coordinator both reached `GO`. This verdict
authorizes only the one launch; it does not authorize input until READY and a
visually inspected valid external baseline exist.

