## Context

The lesson originates from mlops-course revision
`74dc45062f43f9b4abfbca738ffa6b56db4d5daa`. The latest user request supersedes
the original unchanged-copy and checked-in-sample constraints.

## Goals / Non-Goals

Apply the reviewed fixes, expose the full pinned CSV through Bazel, and verify
training and saved artifacts. Preserve the existing lesson layout and simple
estimator. Do not add dependencies, deploy, or introduce hyperparameter search.

## Decisions

- Use existing locked pandas, NumPy, and scikit-learn; remove obsolete lesson requirements.
- Expose the CSV via the third-party overlay and a lesson data alias.
- Keep reviewed function signatures and add a CLI for explicit data/output/seed/split inputs.
- Filter missing required values and zero coordinate pairs before training the full CSV.
- Produce model.pkl and metrics.json as declared sandboxed Bazel outputs.
- Adapt the proposed smoke test to unittest; exercise the CLI twice on a bounded
  upstream prefix, load the saved estimator, verify finite scores and reproducibility,
  and retain a JSON verification artifact. Also build the full-data model.
- Retain the upstream MIT notice in this record. Do not reply to or resolve review
  threads without the prior required user approval.

## Risks / Trade-offs

The full-data model costs more than the sample test and is a manual build target.
Only the reviewed zero-coordinate rule and required-value filtering are enforced;
this is an exercise, not a comprehensive taxi-data validation system.
