# Water A01 independent acceptance review

## Verdict

`ACCEPT` candidate
`9e3ae22c93b376d9e7d371388c1ac2f67ffe871dc840d8eacbbf790b285f157c`.
All seven required criteria at revision 2 pass. The generated packet's
`FAIL_RESET` is a verifier overclaim, not a candidate failure.

The reported 14.3716 px warning height violates no frozen criterion or plan.
Criterion 004 requires legibility and 32 px horizontal plus 18 px vertical
silhouette margins. The immutable plan requires readability and specifies no
minimum text height. The 16 px threshold originated only in the generated
contract/verifier and cannot strengthen acceptance after execution. Actual
640 by 360 pixels show the complete high-contrast warning clearly in frames 1,
24, and 36; its measured margins are 143.95 px horizontal and 29.89 px
vertical.

The other three red verifier checks are defects:

- mesh-local Z bounds were compared without the object's +0.125 m parent-space
  translation;
- 1e-8 m tolerances reject ordinary float32 storage for dimensions and the
  water guide; and
- the same tolerance rejects the correctly sampled ballistic path.

## Criterion verdicts

- 001 pass: exact interface names, axes, units, bottom-center support,
  attachment ownership, and sole scene motion owner are present.
- 002 pass: envelope, tank, water guide, clearances, depth, domain bounds, and
  0.100 m headroom pass.
- 003 pass: 48 fps hold/fall/contact path, 0.0512-frame gravity error, no
  rotation, and no post-contact keys pass; the 2.89e-8 m residual is storage
  noise.
- 004 pass: all subjects and the exact legible warning are framed with margins
  above the frozen gates.
- 005 pass: neutral geometry is isolated outside FUMO with no Reimu cues and a
  byte-preserving replacement boundary.
- 006 pass: clean-open inventory has no liquid/particle modifiers, cache,
  liquid outputs, libraries, or post-contact motion.
- 007 pass: corrected bottom-center preflight, 22 exact artifact hashes, exact
  protected hashes, and narrow claim limits pass.

Acceptance proves a reusable neutral scene interface, scale, framing, and
ballistic path through first contact. It does not claim liquid response, a
final animation, or approved-plush integration.
