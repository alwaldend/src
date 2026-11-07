## Context

See proposal.md. The skill currently requires several canonical views and an
8/10 score in every category before export, even when the source does not
provide those views or the user requests delivery of an imperfect result.

## Goals / Non-Goals

**Goals:** Retain useful visual evidence while scaling review to the actual
references, task, and user-selected acceptance criteria.

**Non-Goals:** Lowering technical asset safety, claiming unsupported likeness,
or overriding a user's request to continue improving a result.

## Decisions

Keep a concise comparison and iteration workflow in the skill. Use the
existing scorecard and plush-construction references only when they help the
task. Replace the universal export gate with a proportional review that
reports material differences and honors the user's delivery instructions.
After a small number of unproductive iterations, reconsider the modeling
approach and continue useful work; do not turn the limit into an automatic
stop. Explicit user thresholds remain controlling.

This retains stronger review when fidelity is essential while avoiding a
fixed score that can displace the user's request. Toolchain and asset-safety
procedures remain in `repo-blender`.

## Risks / Trade-offs

- Less prescriptive scoring could hide weak likeness. Require comparison of
  representative rendered pixels and clear reporting of material defects.
- A delivery request could be mistaken for acceptance of a specific quality
  level. Honor explicit user criteria and do not label an approximate result
  as exact or complete when it does not meet them.
