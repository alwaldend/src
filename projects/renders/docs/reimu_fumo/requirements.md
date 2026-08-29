# Reimu Fumo requirements and constraints

[Back to goal](README.md)

## Requirement changes

- **2026-08-28:** The user requires a great reference-faithful Reimu Fumo that
  is reusable and animated, created locally in Blender. No external service or
  third-party generated asset may create geometry, textures, materials,
  rigging, or animation. Web research is allowed only for workflow knowledge.
- **2026-08-28:** The user requires the Blender-source migration and the
  persistent goal loop.
- **2026-08-28:** Earlier intermediate approval pauses are superseded by the
  latest instruction to keep the long-running task active and not stop after
  an iteration. Internal quality gates remain mandatory, but they do not pause
  the work.
- **2026-08-28:** Preserving the standalone model and not integrating it into
  the Sisyphus scene remain in scope.
- **2026-08-28:** The 5% landmark tolerance, two-reviewer gate, 8/10 score
  threshold, technical mesh checks, and named animation tests below are
  inferred verification gates, not numeric requirements supplied by the user.
  They operationalize the stated quality, reuse, and animation requirements
  and may be strengthened, but not weakened merely to pass a candidate.
- **2026-08-29:** The user requested a rebase onto the newly merged base and a
  review of updated skill guidance. After `git fetch --prune origin`, both
  direct remote checks still placed `master` at `52667f12`; it was already an
  ancestor of this branch, so `git rebase origin/master` completed as an
  explicit no-op. The complete `goal`, `blender-reference-fidelity`,
  `repo-delivery`, `repo-bazel`, `bazel-nested-module`, `bazel-rules-skill`,
  `full-repo-check`, and `host-bot-diagnostics` guidance was reviewed. The
  model, landmark, and protected-scene bytes remain unchanged.
- **2026-08-29:** Every turn on this long-running goal must generate a concrete
  inspectable artifact. The repository `goal` skill now requires a candidate,
  render, comparison, measurement, test result, script, goal-record update, or
  equivalent evidence on every goal-working turn and requires a resolved direct
  link in the assistant's user-visible session commentary.
- **2026-08-29:** Goal records must live at
  `<subproject-root>/goals/<name>/README.md`, not as a `GOAL.md` beside the
  deliverable. This record therefore moved to
  `projects/renders/goals/reimu_fumo/` and its long-form detail and history are
  split into linked files below the required README entrypoint.
- **2026-08-29:** The rear of the final Fumo must have the reference-faithful
  brown hair silhouette; it may not look bald. A component-only review must use
  either a clearly neutral proxy or a complete hair silhouette, never the
  misleading rejected front-hair/bald-rear hybrid.
- **2026-08-29:** The newly supplied `82f07f2f...png` image is the canonical
  appearance reference, and the plush is `25 cm` tall overall. The actual slow
  180-degree GIF was recovered from the user-supplied Tenor URL as a verified
  30-frame `498 × 498` animation. The canonical front and original GIF must be
  retained in the goal dossier and packed into a hidden review-only collection
  in every surviving working Blender file so later iterations cannot omit
  them. The flattened `da8f04ae...png` upload is only a preview of that GIF.
- **2026-08-30:** Keep the tracked standalone Fumo as one unsplit main Blender
  file. Work on hash-addressed candidate copies, test and review the exact copy,
  and modify or replace the main file only after that candidate is approved.
  Do not split the main file merely for parallel editing.

## Constraints

- Use local Blender through MCP for model creation and editing.
- Do not use an external service or third-party generated asset to create
  geometry, textures, materials, rigging, or animation. Internet research is
  permitted only for workflow knowledge.
- Keep temporary scripts, renders, comparisons, and intermediate `.blend`
  files under the repository-root `out/` directory and do not commit them.
- Do not modify the content of the Sisyphus scene while rebuilding this model.
- Do not integrate the model into the Sisyphus scene as part of this goal.
- Preserve the user-supplied references as the controlling visual target.
- The final binary must be required, non-temporary, and tracked through Git
  LFS.
