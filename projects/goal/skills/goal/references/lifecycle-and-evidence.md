# Attempt lifecycle and evidence

Use this reference before starting or evaluating an attempt.

## Plan before work

One attempt is a bounded work unit with one reviewable outcome. Record before
implementation:

- the exact goal generation, lifecycle generation, criteria revision, criteria
  digest, and goal-state digest stored as durable attempt bindings;
- the Goal resource version to use only as the next checkpoint's
  compare-and-swap token;
- the uncertainty or stable defect being targeted;
- work type and hypothesis when the work tests one;
- the plan, inputs, parameters, intended result identity, and review packet;
- affected criteria and the fixed regression checks; and
- parent checkpoint or prior attempt when one exists.

Preflight deterministic geometry, numeric gates, controlling references,
paths, APIs, tools, permissions, and whether the proposed artifact can satisfy
its own acceptance gate. Preserve material preflight evidence without turning
discarded prose drafts into fake attempts.

## Bind evidence to what was tested

A criterion verdict is one of `pass`, `fail`, or `unverified`. Evidence must
identify its criterion revision and the exact tested subject:

- content: a digest or immutable version;
- source change: source revision and relevant configuration;
- stateful operation: target identity, environment/config revision, operation
  receipt, and observed postcondition;
- research or decision: source set, retrieval date where relevant, analysis
  artifact, and review method.

A later change that can affect a verdict invalidates it until rerun. Never
combine visual evidence from one candidate with technical evidence from
another. Command success alone is not evidence of the desired postcondition.

## Review and close

Inspect likely failure views and edge cases, not only the best output. For
subjective work, use consistent comparisons and independent review when
available. Judge absolute quality against the target rather than improvement
from a weak predecessor.

Close the attempt with:

1. work actually performed;
2. artifacts and raw verification evidence;
3. a verdict for every affected criterion and regression check;
4. the dominant remaining failure;
5. the decision to accept, refine, or reset; and
6. a process audit covering measurable movement and feedback bottlenecks.

Publish the close only if the goal resource version and lifecycle generation
still match. A stale attempt remains useful evidence, but it cannot silently
become canonical or reopen a closed lifecycle.

## Retry and finish

When a required criterion fails or remains unverified, keep the outcome open,
update the stable failure count, and choose the highest-leverage next attempt.
After the same defect survives twice, materially change strategy.

For final acceptance, freeze one exact result, run every required criterion
and the complete regression set, deliver or export it, and verify that delivery
did not change its identity. If delivery changes the tested subject, rerun the
full evidence plan.
