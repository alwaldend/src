# Attempt 59 — inverse-pattern method reset

[Back to current attempt](../current_attempt.md)

## Scope

Attempt 59 tested one dependency-ordered construction: a deterministic welded
structural cushion, an independently constructed unpressurized hair mantle,
and one lightly filled near-side lock. It intentionally stopped before eyes,
materials, the full rear-leaf system, or scene integration.

The [official Plushie
paper](https://www-ui.is.s.u-tokyo.ac.jp/~takeo/papers/mori_siggraph2007_plushie.pdf)
and [official project
page](https://www.is.ocha.ac.jp/~yuki/plushie/index-e.html) informed the
inflation method. The user-supplied canonical front and verified 30-frame
turn remained the appearance authority.

## Result

**Rejected before candidate geometry. Nothing was rendered or promoted.**

- The exact structural replay fails at cycle 17: maximum dimension
  `1.501540 Wh`, maximum material stretch `6.061705`, and
  `4,788 / 4,826` constraints above `1.02`. A diagnostic continuation keeps
  expanding through cycle 60 rather than approaching a valid equilibrium.
- Equation, sign, weight, winding, topology, and scale audits show no hidden
  arithmetic exception. The failed premise is the unsupported dimensional
  mapping of the paper's `alpha=.02` to `.02 Wh`, amplified by a corner-sliver
  Coons discretization and a finite correction scheme that cannot meet A59's
  strict released-material gate.
- Analytically selected constant, decaying, mesh-relative, stronger-Jacobi,
  and sequential comparisons reduce particular symptoms but still violate
  material or convergence gates. Selecting a transient would hide stretch.
- The corrected mantle/lock capsule independently fails: mantle/lock maximum
  stretch `1.2272/2.4004`, released mantle insertion `.02876 Wh`, lock width
  `.15042 Wh`, landmark error `.08142 Wh`, root gap `.11634 Wh`, and no passing
  exact collision gate. Earlier apparently green evidence measured shaped
  seed geometry rather than construction release and is retired.

The primary runtime decision is
[here](../../../../../out/reimu_fumo_attempt_059_method_reset/SOLVER_RUNTIME_DECISION.md),
with the independent
[divergence audit](../../../../../out/reimu_fumo_attempt_059_method_reset/SOLVER_DIVERGENCE_AUDIT.md),
[scale decision](../../../../../out/reimu_fumo_attempt_059_method_reset/SOLVER_SCALE_DECISION.md),
and [mantle/lock report](../../../../../out/reimu_fumo_attempt_059_method_reset/mantle_lock/REPORT.md).

## Reusable evidence

- The frozen all-reference packet and source ownership survive. Its
  [landmark board](../../../../../out/reimu_fumo_attempt_059_method_reset/reference_inputs/reference_board_owned_landmarks.png)
  includes the canonical front, controlling three-quarter/side/rear turn
  brackets, physical front/side, and older turn/sofa cross-checks.
- The isolated numeric kernel now fails closed and preserves deterministic
  negative evidence. It is not a model generator or likeness proof.
- Separate interfaces exposed invalid contact, collision, pattern, and scale
  assumptions before they contaminated the protected `.blend`.

## Process conclusion

The attempt improved diagnosis but did not move visible asset quality. A new
constraint solver, remesh, and contact representation would be another large
uncertain tool-building branch while the overlay representation has already
failed independently. Result priority therefore ends this method instance.
The next attempt returns to one byte-verified complete Blender copy and changes
one directly visible module at a time under fixed multi-view review.

