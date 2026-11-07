# Visual review prompts

Use these prompts when they help assess a reference-matched Blender asset.
Scale the review to the request and available evidence; they are not an
additional approval gate.

- Is the subject recognizable, and are its defining proportions, silhouette,
  materials, and pose consistent with the available reference?
- Are attachments, occlusions, and deformations visibly plausible? Inspect
  another view or an animated pose when it can reveal a hidden problem.
- Which discrepancies matter most to the user's intended use? Distinguish
  observed defects from uncertainty about an unavailable or obscured view.
- For animation, do representative cuts, expressions, motion, and transitions
  follow the source? Inspect the actual output, not just the scene setup.

Independent, implementation-blind image review can reveal missed problems,
particularly after repeated iterations. Give the reviewer the reference and
representative renders without a suggested verdict. Treat its findings as
input to the primary agent's judgment, not as a substitute for the user's
acceptance criteria.

When delivering an approximation, state the material remaining differences.
Do not describe it as an exact match or as satisfying a quality threshold it
has not met. A user may direct delivery of a reviewed imperfect result; do not
add a universal score or export prohibition.
