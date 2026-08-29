# Current attempt

[Back to goal](README.md)

## Attempt 2 — readable neutral rigid drop scaffold

### Failure or uncertainty targeted

Attempt 1 proves mechanics but fails as an inspectable visual scaffold: the
floor and proxy clip nearly white, the settled proxy is too small, and the
required warning does not appear in the render.

### Falsifiable hypothesis

Keeping Attempt 1 geometry, physics, timing, and interface fixed while lowering
light energy, tightening framing, and moving a shorter warning inside the
camera frustum will make the same mechanical proof visually inspectable.

### Frozen plan before implementation

- Preserve Attempt 1 units, frame range, release, rigid bodies, dimensions,
  collections, interface metadata, and sampled frames exactly.
- Reduce area-light energy from `520/250/360 W` to `55/28/40 W`, reduce world
  strength from `.35` to `.12`, and use exposure `-0.35`.
- Move the fixed camera closer while retaining the complete start-to-floor
  range; increase diagnostic resolution to `640 × 360`.
- Replace the off-screen warning with two centered camera-space lines:
  `NEUTRAL PLACEHOLDER` and `RIGID DROP — NO PLUSH`.
- Enlarge the 25 cm scale witness and label without changing its measured
  height.
- Render the same seven frames.  Reject if warning pixels are absent, proxy or
  floor clips, the settled proxy is unreadably small, or any Attempt 1
  technical regression fails.

### Planned review packet

- `fumo_concrete_drop_scaffold.blend`
- `contact_sheet.png`
- `build_result.json`
- seven sampled frame PNGs

### Attempt 2 evidence and decision

Candidate SHA-256:
`a9488c220c5076a3202e61c9897cf3710f24b1abe74fb9edfc4750bfaebfdc26`.
The scripted packet passes all eleven checks and preserves the Attempt 1
mechanics: contact frame `22`, minimum sampled bottom Z `-.00068653 m`, and
late-motion span `0 m`. The protected Reimu and Sisyphus hashes are unchanged.

Absolute review confirms that exposure, warning, impact, and settle improved,
but frames `1` and `12` crop most of the falling proxy. This violates the
frozen complete start-to-floor framing requirement. Verdict: technical pass,
composition reject; a clean-reopen audit alone cannot rescue it.

### Attempt 1 evidence and decision

Candidate hash:
`86630e599525e40663ad01e4bd8f4c5c5f12e9cb127740440a7bbe501b77d292`.
Contact occurs at frame `22`; minimum sampled floor penetration is `.000687 m`;
late motion is `0`; all scripted technical checks pass, while clean reopen is
unverified. The visual packet is rejected for clipped exposure, weak
settled-pose readability, and absent warning text.

### Progress, approach, and process audit

Mechanics, interface, scale, and timing measurably advanced from unverified to
passing.  Visual evidence did not pass absolutely.  Rendering seven frames
took about half a minute, while two environment/API recovery loops were
avoidable setup cost.  Attempt 2 keeps the proved mechanics and changes only
the failed presentation layer; another mechanics rebuild would add risk
without addressing the dominant failure.

Attempt 2 reused the proved mechanics without a whole-scene rebuild, but its
camera correction was insufficient. The next attempt changes only framing,
then repeats the render and clean-reopen gates.
