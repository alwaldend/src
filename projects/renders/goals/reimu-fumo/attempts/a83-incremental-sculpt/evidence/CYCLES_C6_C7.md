# A83 cycles C6--C7 evidence

## C6: local proportional lock sculpt

C6 restarted from the retained C1b checkpoint and changed only
`A42 Short right rear lock` through a reversible shape key. It widened and
thinned the middle/lower panel while freezing the root. The result was clearly
visible but failed the named construction cue: three-quarter showed a longer,
narrower dangling tail, and rear showed a pinched strand rather than a broad
stuffed panel with edge thickness and overlap. **UNDO**.

- Candidate blend:
  `sha256:cc9c7e94f27ed5edc5a3336047a5b6b921acb12154df2f561a27b23edcd8a07e`.
- Blind verdict:
  `sha256:89123d31ac2e64f3788a1489d348dfdde0771247af72fcb3cb65b9a75e0102f5`.

## C7: root-preserving drape rotation

C7 again restarted from exact C1b. It preserved the panel's full existing
geometry and rotated only its middle/lower portion ten degrees around the head,
with a zero-displacement root. This corrected C6's extreme tail but did not
pass the baseline comparison: the piece remained a narrow rounded tendril on
the smooth cap, while rear construction was essentially unchanged. **UNDO**.

- Candidate blend:
  `sha256:54680013136465d93956bae5544ed18565c187cf70a0e3a4dabac8b23ce5494f`.
- Blind verdict:
  `sha256:0353fd0c2b70abb9e20e2a21c69e31a4e75d55329bad3d9aefbec7069b123dd5`.

Protected rung 003 and the tracked reusable model retained their expected
hashes throughout.

## Structural conclusion and next method

The separate short rear-lock owner can change silhouette pixels but cannot
create the required broad camera-facing fabric plane while the intact cap
continues to occlude and visually dominate it. Continuing to scale, stretch,
or rotate that lock would repeat the tail/card failure. The next work must
restart from C1b and directly control a bounded cap-side surface in an actual
small sculpt/constructed-panel test. It must still preserve the receiver,
front landmarks, bow seat, and outer bound and must render immediately.
