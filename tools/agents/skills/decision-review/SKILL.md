---
name: decision-review
description: >-
  Review a material proposal or implementation choice before commitment using
  quality-assessment, credible alternatives, and an explicit decision verdict.
  Use for consequential design, security, operational, costly, irreversible,
  or repeatedly failing decisions; do not use for routine reversible details
  or standalone questions about an artifact's quality or accuracy.
---

# Review a decision

Read and apply [quality-assessment](../quality-assessment/SKILL.md) for the
shared assessment method, evidence standards, bias control, and authority
boundaries. This skill adds the decision-specific work below. The primary
agent owns the final decision verdict.

## Test the proposal

Before committing to the choice:

1. State the actual outcome, acceptance criteria, constraints, and relevant
   authorization. Separate the user's goal from their suggested method.
2. Apply quality-assessment to the proposal and credible alternatives,
   including doing nothing when it is a real option. Compare outcome quality,
   risk, reversibility, cost, and time against the same criteria.
3. Issue one verdict: **proceed**, **revise**, **ask**, or **refuse**. State
   the decisive evidence, material trade-offs, and conditions that would
   change the verdict.

During iterative design, preserve the user's settled semantics and constraints
in the existing task record or a compact working summary. Evaluate a proposed
API against the user's concrete example before adding lifecycle hooks, global
objects, or general frameworks. Identify the demonstrated requirement each
abstraction serves; speculative future use alone does not establish its value.
A request to simplify the API does not reopen agreed behavior without evidence.

Distinguish disagreement about the best method from a reason that makes the
requested action unsafe, unauthorized, impossible, or contrary to the outcome.
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

Do not relitigate settled preferences without new evidence or burden routine
low-risk work with a formal decision review.
