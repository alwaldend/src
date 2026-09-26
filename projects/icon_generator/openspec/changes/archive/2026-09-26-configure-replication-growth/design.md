## Context

Replication currently samples two vacant neighbors uniformly and queues each
newborn for a single reproduction decision. This change follows the four
placement tweaks selected by the user; population rules remain unchanged.

## Goals / Non-Goals

Make direction inheritance, forward preference, soft crowding avoidance, and
sideways branching independently configurable. Preserve defaults and generate
comparable seeded previews. Do not introduce one-child births, erase parents,
promise a particular morphology at high density, or add dependencies.

## Decisions

All four options accept finite 0..1 values and default to zero. They are
validated for every algorithm but affect replication only.

Each active cell carries one of eight headings. Roots get a random heading;
a child uses its parent-to-child direction with the inheritance probability,
otherwise a random heading. Forward bias interpolates from uniform weights
to a profile favoring straight and gentle turns over sideways and backward
steps. Inheritance needs forward bias to create persistent paths.

Crowding divides destination weights by a positive penalty based on occupied
neighbors excluding the parent. The first child is considered when scoring the
second. It never makes a valid vacancy impossible or changes reproduction
probability. Branching is the probability that the second child favors a
heading 90 degrees to either side of the first child's direction, retaining
positive weights for all vacancies when that heading is blocked.

Keep the original uniform-selection fast path when forward bias, crowding, and
branching are zero. Direction-only randomness uses a separate seeded stream so
inheritance alone cannot perturb legacy output. Store the heading in the upper
bits of the existing uint32 queue entry: the maximum 4096-square cell grid
uses 24 index bits, leaving room for three heading bits without growing the queue.

## Risks / Trade-offs

- Two-child births can remain bushy, and high reproduction probability can
  still fill the image. Preview multiple settings without promising thin strands.
- Direction weights favor grid directions. Random headings and gentle turns
  reduce regularity but cannot remove the eight-neighbor grid geometry.
- Placement changes can alter later vacancy availability and final population;
  use a common initial sample, seed, and probability for visual comparisons.
- Verify compatibility, parent retention, parity, tiny boundaries, replay,
  sparse crowding morphology, and peak-grid runtime. Write E2E cases first.
