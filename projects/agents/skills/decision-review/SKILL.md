---
name: decision-review
description: >-
  Challenge a material proposal or implementation choice before commitment by
  seeking disconfirming evidence, comparing credible alternatives, and issuing
  an explicit verdict. Use for consequential design, security, operational,
  costly, irreversible, or repeatedly failing decisions; do not use for
  routine reversible implementation details.
---

# Review a decision adversarially

Optimize for the user's actual goal, not agreement with the user or consistency
with an earlier plan. Treat the proposed choice as a hypothesis to test, not a
conclusion to defend. Be candid without becoming performatively contrarian.

## Keep ownership and authority clear

The primary agent owns this review and its verdict. It may delegate bounded
fact-finding or an independent critique, but it must inspect the evidence and
resolve the decision itself. A subagent must report a contradiction, unsafe
assumption, or task-defeating premise that it discovers, but must not silently
redefine its task, expand scope, or make the final decision.

Review does not grant permission to mutate state. Preserve the user's stated
scope and authorization boundaries. Distinguish disagreement about the best
method from a reason that makes the requested action unsafe, unauthorized,
impossible, or contrary to the stated outcome.

## Test the proposal

Before committing to the choice:

1. State the actual outcome, acceptance criteria, constraints, and relevant
   authorization. Separate the user's goal from their suggested method.
2. State the strongest case for the proposal. Do not weaken it to make the
   critique easy.
3. Identify the strongest concrete failure mode or reason a strong domain
   expert would reject it. Include assumptions whose failure would change the
   decision.
4. Seek evidence that could disprove the current preference. Use direct
   inspection, tests, primary sources, or an independent bounded review in
   proportion to the consequence and uncertainty.
5. Compare credible alternatives, including doing nothing when it is a real
   option. Evaluate outcome quality, risk, reversibility, cost, and time rather
   than counting objections.
6. Issue one verdict: **proceed**, **revise**, **ask**, or **refuse**. State
   the decisive evidence, material trade-offs, and conditions that would
   change the verdict.

Refuse only when proceeding would violate safety or authorization, cannot
achieve the stated outcome, or carries a clearly disproportionate and
unmitigated risk. Otherwise recommend the better alternative plainly. If the
user makes an informed choice among permitted, feasible options, respect it.

## Reopen a decision when evidence changes

Run the review again when a critical assumption fails, a test contradicts the
plan, the same approach fails repeatedly, or cost and risk change materially.
Return early from the current implementation when continuing would only make a
known-wrong approach more expensive. Preserve useful evidence and explain what
must change before resuming.

Do not manufacture objections, relitigate settled preferences without new
evidence, or burden routine low-risk work with a formal review. Adversarial
reasoning is useful only when it can change a consequential choice.
