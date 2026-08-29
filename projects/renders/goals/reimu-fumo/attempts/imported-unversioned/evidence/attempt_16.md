# Reimu Fumo attempt 16

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 16 — correctly owned cushion with short native-sculpt families

**Candidate:** planned
`out/reimu_fumo_attempt_016_owned_cushion/checkpoint_final.blend`, SHA-256
pending, untextured beige under-hair cushion stage, review packet
`attempt-016-owned-cushion-aligned-front-side-and-fixed-views`.

**Failure targeted:** prior attempts judged bare clay against the taller outer
hair envelope, used brush sizes different from their labels, waited until a
large family ended before seeing pixels, and compared unregistered source and
candidate images.

**Hypothesis:** The real plush is built around a compact, broad cushion whose
hair owns the extra crown and rear silhouette. A near-proportion disposable
blank followed by three small, independently checkpointed native-sculpt
families can establish this hidden volume without asking a sphere, box, or
panel generator to be the final form. Exact live-`P` brush sizes and aligned
silhouette evidence make each change observable and falsifiable.

**Plan written before implementation:**

1. Freeze stage ownership before geometry. The beige cushion, not the outer
   hair, must finish at width `0.90–0.95 Wh`, height `0.86–0.91 Wh`, and depth
   `0.68–0.74 Wh`. The later brown hair owns the outer `1.00 Wh` width and
   `1.03 Wh` height. Do not move these bands after seeing Attempt 16.
2. Before candidate creation, generate temporary front and side target guides
   from these fixed world-space silhouette points. The front right half from
   crown to underside is `(0,.445)`, `(.28,.43)`, `(.42,.31)`,
   `(.46,.10)`, `(.455,-.19)`, `(.34,-.38)`, `(0,-.445)`, mirrored on X.
   The side contour in `(y,z)` order is `(0,.445)`, `(-.20,.40)`,
   `(-.31,.27)`, `(-.36,-.05)`, `(-.30,-.30)`, `(-.16,-.42)`,
   `(0,-.445)`, `(.16,-.42)`, `(.29,-.28)`, `(.35,.08)`,
   `(.31,.29)`, and `(.20,.41)`. Record a visible `±0.04 Wh` uncertainty
   band. Create these guides and their hashes before any candidate render.
3. Use a separately declared head-stage camera calibration: square `640 px`
   orthographic cameras centered at `z=0` with scale `1.20 Wh`. This does not
   replace the full-asset `2.36 Wh`, `z=1.055 Wh` contract. Each review packet
   must include raw controlling head crops, the fixed guide, candidate
   silhouette, a `40%` guide/candidate overlay, and an edge-difference image.
   Host-side ImageMagick may assemble temporary evidence after Blender renders;
   Blender must not invoke a command absent from its Flatpak.
4. Start factory-empty and import no rejected geometry. Add one native
   subdivision-four icosphere; scale it once to dimensions
   `0.84 × 0.70 × 0.80 Wh`, apply transforms, and voxel-remesh at `0.025 Wh`.
   This is disposable clay stock only and is never rendered or scored. Require
   one closed, finite, positive-volume surface with one component, coherent
   orientation, no self-overlap, zero modifier, and no other mesh or orphan.
5. Replace discrete brush tiers with exact live projection. For every stroke,
   set `size = round(radius_fraction × P)` and require the executed fraction
   to differ by no more than `0.5/P`. Fix pressure at `1.0` and persist every
   projected endpoint, raycast hit, native sample, exact asset path, setting,
   coordinate hash, topology hash, moved set, maximum displacement, influence
   radius, and signed directional projection. Never label one size as another.
6. Before the candidate, run disposable fresh-stock probes for every intended
   Grab radius (`0.24`, `0.26`, `0.28`, and `0.30 P`) and the intended
   `Scrape/Fill` radius (`0.30 P`). Use short in-silhouette paths. Require
   exact Essentials assets, `FINISHED`, unchanged topology, finite local
   movement, requested size within `0.5/P`, and reflected displacement
   residual below `0.001 Wh` for each enabled X/Y axis. Recompute every
   predicate when loading probe evidence; `pass: true` alone is insufficient.
   Purge all probe data.
7. Family A contains exactly five broad Grab strokes from fixed front view.
   Use spherical smooth falloff, no Dyntopo, front-face restriction,
   automasking, autosmooth, or topology rake. In order: upper shoulder
   `(.22,.20)→(.32,.29)`, `0.30 P`, strength `.72`, X/Y symmetry; lower
   shoulder `(.23,-.19)→(.34,-.28)`, `0.30 P`, `.72`, X/Y; lower-middle
   cheek `(.31,-.02)→(.44,-.05)`, `0.26 P`, `.78`, X/Y; crown
   `(0,.32)→(0,.44)`, `0.28 P`, `.70`, Y; and underside
   `(0,-.31)→(0,-.43)`, `0.28 P`, `.70`, Y. Every start lies at least
   `0.07 Wh` inside the blank projection. Require positive displacement along
   the drag direction, maximum displacement `0.055–0.155 Wh`, and influence
   radius no more than `0.38 Wh`.
8. After every successful stroke, save a new immutable checkpoint and manifest
   whose parent is the preceding manifest hash. On any exception, save a
   separate terminal checkpoint and terminal JSON before re-raising. A stage
   always reopens and hash-checks its parent checkpoint even when Blender
   already reports that filepath. Never mutate a prior metrics or review file.
9. After Family A, require width `0.87–0.97 Wh`, height `0.83–0.94 Wh`, depth
   `0.68–0.73 Wh`, crown and underside notch at most `0.010 Wh`, and maximum
   aligned front-guide error `0.06 Wh`. Render front neutral, silhouette, and
   overlay evidence immediately. An implementation-blind front review must
   find no sphere, egg, balloon, box, corner, cleft, spike, or hard shoulder
   before Family B may run.
10. Family B contains exactly two broad `Scrape/Fill` passes from front view,
    with Y symmetry only: `(-.08,.10)→(.08,.10)` and
    `(-.08,-.12)→(.08,-.12)`, each radius `0.30 P`, strength `.14`, and
    pressure `1.0`. They must reduce excessive center bulge without making a
    plate. Require final front and rear face bow, measured from the center to
    the visible surface at `|x|=.28±.02`, within `0.015–0.085 Wh`; depth must
    remain `0.67–0.74 Wh`. Save and review a new aligned front packet before
    any side stroke.
11. Family C contains exactly six Grab strokes from fixed side view with X
    symmetry only: front/lower `(-.28,-.05)→(-.34,-.08)`, `0.26 P`, `.65`;
    rear/upper `(.28,.05)→(.34,.08)`, `0.26 P`, `.65`; front crown
    `(-.18,.30)→(-.21,.33)`, `0.24 P`, `.58`; rear crown
    `(.18,.30)→(.22,.33)`, `0.24 P`, `.58`; front underside
    `(-.19,-.30)→(-.16,-.33)`, `0.24 P`, `.58`; and rear underside
    `(.19,-.30)→(.16,-.33)`, `0.24 P`, `.58`. Require each maximum
    displacement `0.012–0.075 Wh`, correct signed direction, finite closed
    topology, and influence radius no more than `0.34 Wh`.
12. The final numeric gate is width `0.90–0.95 Wh`, height
    `0.86–0.91 Wh`, depth `0.68–0.74 Wh`, front and side guide maximum error
    at most `0.05 Wh`, face bow `0.015–0.085 Wh`, and axial notch at most
    `0.010 Wh`. Extract the visible front/side contours rather than searching
    all 3D vertices. On the side convex hull, reject any single straight edge
    over `0.10 Wh`, near-horizontal crown roof over `0.08 Wh`, or paired
    near-vertical front/rear run over `0.10 Wh`. Reject non-finite values,
    disconnected or non-manifold topology, non-positive volume, self-overlap,
    dimension scaling, modifiers, or geometry edits outside the three families.
13. Fast review uses the aligned front and side packets before full evidence.
    Full evidence adds neutral, silhouette, and grazing front, side, rear,
    canonical three-quarter, and mirrored three-quarter plus one readable
    perspective presentation render. The implementation-blind record must
    include unlabeled recognition, intended medium, ordered five largest
    discrepancies, per-view and macroform/soft-stuffing/non-primitive/
    presentation scores, and an explicit major-failure boolean.
14. Hard-reject sphere, egg, balloon, capsule, rounded box, mattress, molded
    foam, flat plate, two-lobe split, axial trench, brush dimple, stretched
    spike, pinched pole, roof, parallel walls, hard shoulder, broken oblique
    highlight, or an image packet that hides the silhouette. Require every
    canonical view and macroform/soft-stuffing/non-primitive score at least
    `7/10`, presentation at least `6/10`, no major failure, and no unresolved
    discrepancy above `0.05 Wh`.
15. A PASS accepts only the frozen cushion volume as the target for constructed
    fabric retopology and later hair. The sculpt/remesh topology is disposable,
    not the reusable final asset. A failure ends Attempt 16 without changing a
    radius, stroke, target guide, threshold, camera, crop, or reference.

**Structural difference from Attempt 15:** the target now belongs to the beige
cushion rather than the outer hair; the blank is proportioned disposable clay,
not a claimed final generator; every raycast start is safely interior; every
brush uses and records its exact live-`P` size; pixels arrive after each small
family; checkpoints and manifests are append-only; and references are aligned
by a precomputed guide instead of independently resized in a late montage.

**Result:** Terminal guide-preflight rejection. No Attempt 16 driver or Blender
candidate was created. The front target PNG SHA-256 was
`16d46eac0d4721b2f50fd70cc82b4ac096049e1e15970cbb3a6c4d050f98ad01`;
the side target PNG SHA-256 was
`fcd8b6d8d2483f5a19fb05cd9931b2ea00f97e3a32d30432d3965f5240571d25`;
and the registered-context board SHA-256 was
`a06ddad6c763b84da4b5b418983ca3cb6ce8c05a0e4a07e76cddf7cf84c932e2`.

An implementation-blind contour review scored front `5/10`, side `3/10`, and
reported a major failure. In order, it found excessive egg/lozenge side depth,
no face/rear/gusset transition, constructional conflict between a filleted-box
front and spheroid side, straight and uniform-radius molded-box front cues, and
no seam pinching, stuffing variation, contact compression, or controlled
asymmetry. Static geometry audit independently measured a `0.290 Wh`
near-vertical front-guide wall and `0.280 Wh` crown run. The target itself
therefore required the rounded box, roof, and parallel walls that the plan
hard-rejected.

The audit also found that the side and crown dimensions were invented hidden
cushion boundaries rather than registered source evidence; only stroke starts,
not outward Grab endpoints, can require surface raycast hits; guide error,
face bow, notch, straight-run, and self-overlap metrics lacked deterministic
definitions; and the `0.05 Wh` guide threshold was looser than the default
`3%` critical-landmark fidelity gate. Camera and exact projected brush-size
arithmetic were internally consistent but never reached execution.

**Criterion results:** No modeling, visual, reuse, or animation criterion
passes. The preflight evidence gate passes only in the sense that it rejected
an internally contradictory target before implementation.

**Decision:** Reject Attempt 16 without writing or running its driver. Retain
the guides only as negative evidence. Do not sculpt toward them or adjust their
curves within the same attempt.

### Progress and approach audit after attempt 16

- **Improved:** the target image itself received an implementation-blind gate,
  reducing the feedback loop from a full Blender build to a cheap guide render.
- **Regressed or unchanged:** no model pixels improved; a precise guide merely
  made the wrong inferred construction easier to reproduce.
- **Absolute result:** the proposed contour is a `5/10` rounded box from front
  and a `3/10` egg from side. It is unusable.
- **Invalid assumption:** the hidden beige cushion does not have a measurable
  crown or side boundary in the supplied front photograph. Brown hair occludes
  those regions, and the oblique photograph controls layering and depth order,
  not an orthographic hidden-cushion outline. A precise bare-cushion silhouette
  was false certainty.
- **Evidence against continuing:** the guide contradicts its own vetoes and
  cannot name a visible source for its full boundary. More sculpt accuracy
  would make the wrong target more faithfully wrong.
- **Highest-leverage unresolved problem:** calibrate the coupled beige support
  and brown hair silhouette that is actually visible and measurable, including
  the broad central fringe and side locks, before choosing sculpt operations.
- **Approach decision:** stop requiring an isolated hidden support to resemble
  the finished character. Attempt 17 is a guide-only registered-evidence
  cycle for the coupled head and hair; Blender remains untouched until that
  observable target passes absolutely.
