# Pinned final sculpt-capability failure evidence

Attempt 018 was the final autonomous sculpt-capability discriminator. It used
only the disposable attempt-017 fixture and did not open or edit any Reimu
model file.

## Exact identities

- Frozen baseline SHA-256:
  `6f49bd4e0a8af6b45870d9d4224a520c1398e52e6e9c42f8fb5bee7b8c17118e`.
- Three-stroke partial input SHA-256:
  `2428de0a0b65e572de9437a8d3ef35f1ee21c18bd9dbf27ff01de1816418c0bd`.
- Terminal 27-stroke file SHA-256:
  `b9ad59c15901b0fe22cd96208e60792e2dc35e8bc0f3e73d0e1a8181697dd6fc`.
- Baseline coordinate digest:
  `41ee23670f67335ac070d95bd782436f53405034f1e24efdebdad709f7d47df2`.
- Partial coordinate digest:
  `f2fdcbdaea90335b9de47861e8ce59b64896ec81853d19ee2ea675725d9fb16e`.
- Terminal coordinate digest:
  `fc3d57571800ecade593e64d2750687f883c7da9ef6177fa9eec622a8346c2e0`.
- Writer receipt SHA-256:
  `a299a04ddf641882c4a1502c5a32ca375eb6effd7e2fa114d402a96de15804cb`.

The writer varied only cumulative identical Flatten dose. Fixture, mask,
front view, brush, strength, radius, trajectory, timing, metric, threshold,
and renderer remained frozen. Eight added blocks delivered 264 events and 24
new strokes, reaching cumulative checkpoints 6 through 27.

## Terminal measurement

Baseline plane variance was `0.014215173048885853`. Variance declined
monotonically at every checkpoint and ended at `0.009603326433540808`, a
`32.44312678772848%` reduction. The unchanged gate was 35 percent. The final
block was the first below the frozen low-response floor, and the hard
27-stroke ceiling terminated the attempt.

At the terminal state:

- 975 vertices changed from the partial input, all inside `REGION_plane`;
- non-plane and `CONTROL_plane` displacement were exactly `0.0`;
- maximum cumulative baseline displacement was
  `0.11849009312857924`, below the `0.20` safety ceiling; and
- root and tip operations were correctly withheld because the plane gate did
  not pass.

Twenty-four native undo operations restored the exact partial coordinate
digest, with zero changed vertices and zero maximum residual displacement.

## Independent pinned verification

Repository-pinned Blender 5.2.1 invocation
`71604ebf-4fb2-420f-875f-23d9bbf844a1` clean-opened the exact baseline,
partial, and terminal bytes. `pinned_audit_018.json`, SHA-256
`ed9e21597d44d601de6532eb472e0adc737c1ce735890113ad59e380330900e7`,
independently reproduced the three variances, all eight ordered blocks,
termination condition, mask counts, locality, displacement ceiling, withheld
root/tip stages, topology digest, and unchanged input hashes.

Pinned invocation `43128522-b4ff-42d2-b510-6148e2c459a5` wrote a fresh render
packet. Its manifest SHA-256 is
`73065da607856f0c6ffc92d1ab2191aa249b80db99370f26f85791b9b4f12267`,
and `READY` has SHA-256
`f6b2aad41c0c6cc88b0d7185d64d8db2c6f89091efde3e07c199b06805a70bf0`.

Image SHA-256 values are:

- baseline front:
  `4affae609bc8794a150bfb1ddc0d4cd44d48d8c935215cce3e549034da21e52f`;
- baseline three-quarter:
  `25f91ee995a03cbbd33b871abc31041ececdf39126c249500855fb925f0b1ec8`;
- partial front:
  `3874b6319fab19155fc7c1627cb81d320737ba77c0c8e500ffd2ec7ee95fe0cb`;
- partial three-quarter:
  `99ab886109f9e32168266c00b3de4e8716fd0dc4324dc20b46e363796fd1ada4`;
- terminal front:
  `5229b2391056e9006f48939c0f066d72f7fcecc3c2ea52fbd65b34059c09e64b`;
  and
- terminal three-quarter:
  `fc97d20848fd7cfacbdef96c35bc662d96e8542257136798fda6b932801dda62`.

The terminal pixels show a flatter center but retain a substantial pillow
bulge and add three horizontal terrace-like ridges. Root, tip, and outer
silhouette remain effectively unchanged. This is a shape-control failure, not
an event-delivery, isolation, save, reopen, render, or undo failure.

## Goal state and resume condition

No model or goal criterion advances. The autonomous broad-transform and
localized fixed-dose routes have both reached their recorded stop conditions.
The blocker is the absence of a sponsor-approved skilled Blender artist or a
genuinely different, independently proven organic authoring capability. Resume
only when one of those is available. Generic desktop control and additional
Flatten dose or parameter tuning do not meet that condition.

Failed binary and image artifacts remain ignored task scratch in this
worktree. The compact writer, pinned, and ergonomics receipts are the durable
evidence.
