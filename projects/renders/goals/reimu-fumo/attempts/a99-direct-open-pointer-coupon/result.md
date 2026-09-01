# A99 result — terminal reset on live brush-state mismatch

## Verdict

`RESET`. The direct-open correction successfully removed A98's foreground
context failure: pinned Blender loaded the exact file, rendered the two fixed
baselines, entered Sculpt mode, emitted READY, and yielded a valid external
mapped-window capture. No pointer input was delivered.

## Decisive pre-input failure

The live capture contradicted the READY receipt. The plan serialized Grab as
world-locked `SCENE / 0.050 m / 0.40`, but the actual Blender UI showed active
Grab at `View / 100 px / 0.40` with a correspondingly smaller cursor. The
complete coupon and right view were visible, so this is not another black-
viewport or wrong-window failure. It is a disagreement between the scripted
brush data-block and the live tool state that would receive XTest input.

Per A99's terminal gate, the coordinator did not call `inject`, did not patch
or relaunch, and stopped both task-owned processes. The exact input, rung003,
and tracked model remain unchanged. No geometry or Fumo criterion was tested.

## Consequence

A99 was the independently authorized final bounded correction. Its failure
establishes that the current autonomous live-pointer/sculpt interface is not
stable enough to author admissible model changes. The main Reimu Fumo goal is
blocked pending a stable external live-sculpt interface or human DCC
authoring. There is no A100 harness, analytical generator retry, or whole-model
rebuild.

Rung003 remains the exact visual high-water. Criteria 001–008 remain
unverified and the tracked reusable asset is not promoted.

