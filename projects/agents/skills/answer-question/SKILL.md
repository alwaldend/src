---
name: answer-question
description: >-
  Answer substantive user questions with evidence while preserving the
  distinction between asking for information and authorizing action. Use for
  mixed question-and-action requests too; quoted transformation payloads do
  not trigger the skill.
---

# Answer a question

## Authority

A question alone authorizes read-only investigation, not changes to files,
settings, pull requests, deployments, or external state. Preserve authority
already granted in the conversation. For mixed requests, answer the question
and carry out the explicitly authorized action. For conditional actions,
establish the condition first; if it remains materially uncertain, explain
what is missing before proceeding. Quoted questions inside material supplied
for transformation are payload unless the user asks to answer them.

## Repository evidence

Inspect the owning artifact at the relevant version. Distinguish working-tree
content, committed history, generated metadata, test observations, built
artifacts, and deployed behavior. When inferring runtime behavior from source,
trace material build selection, generation, packaging, feature flags, and
runtime configuration. Presence at one layer does not establish enablement at
the next. A declared test or catalog entry does not establish a successful
run; cite the actual observation and candidate when making that claim.

Keep negative conclusions within the search or measurement's coverage: “not
found” does not establish “absent.” A cause needs evidence that distinguishes
it from live alternatives; timing alone establishes a hypothesis. Sources
repeating one origin are one evidence chain. State supplied facts at their
given strength, correct false premises, and avoid inventing hidden behavior
or uncertainty after decisive evidence resolves the question.

For quantitative claims, check units, denominators, and assumptions. Do not
infer measured throughput from latency alone; omit unnecessary conversions
or label their relationship assumptions explicitly.

## Answer and judgment

Lead with the answer and decisive observation. Separate verified facts,
inferences, recommendations, and material uncertainty; make references support
the adjacent claim and the version asked about. Match investigation effort to
the stakes and uncertainty; simple stable questions need no research ritual.

For material recommendations, state the criteria, strongest concrete
objection, and trade-offs. In a critique, prioritize by
consequence and confidence, distinguish defects from risks or preferences,
and suggest a concrete improvement for each material issue. Keep the answer
self-contained; disclose that no action was taken only when the context could
otherwise imply a change.
