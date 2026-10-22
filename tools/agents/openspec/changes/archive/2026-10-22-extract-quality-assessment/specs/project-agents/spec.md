## ADDED Requirements

### Requirement: Reusable quality assessment

The repository SHALL provide a discoverable `quality-assessment` skill for
assessing claims, answers, artifacts, and work against their intended purpose.
It SHALL test accuracy and reasoning with relevant evidence, seek strong
counterevidence, reduce agreement and confirmation bias, and distinguish
defects, uncertainty, and preferences. Research effort SHALL be proportional
to stakes and uncertainty. Assessment alone SHALL NOT authorize mutation.

#### Scenario: Ask whether a claim is accurate

- **WHEN** the user asks "is this accurate?" about an identifiable claim
- **THEN** the assessment checks the claim against available evidence, verifies
  material gaps where possible, and gives a calibrated answer with its limits
- **AND** it does not endorse the claim because the user or agent proposed it

#### Scenario: Ask whether an artifact is good

- **WHEN** the user asks "is this good?" about an artifact with a stated purpose
- **THEN** the assessment uses criteria relevant to that purpose, identifies
  strengths and material weaknesses, and suggests useful improvements
- **AND** it does not invent objections, impose tastes as facts, or require a
  decision verdict for a standalone assessment

### Requirement: Decision review composes quality assessment

`decision-review` SHALL use `quality-assessment` as the owner of general
assessment procedure. It SHALL retain decision-specific goal and method
separation, comparison of credible alternatives, material trade-offs, explicit
proceed/revise/ask/refuse verdicts, and reconsideration when evidence changes.
Repository change reviews SHALL retain their existing `review-agent` routing.

#### Scenario: Assess a consequential proposal

- **WHEN** a material proposal needs review before commitment
- **THEN** decision review applies quality assessment and compares credible
  alternatives before issuing its decision verdict
- **AND** the primary agent retains ownership of that verdict and existing
  user authority is preserved

#### Scenario: Review a bounded repository change

- **WHEN** the request is to review a pull request or other bounded change set
- **THEN** the repository's required change-review procedure still applies
- **AND** general quality assessment does not replace that procedure

### Requirement: Remove the standalone question-answering skill

The repository SHALL remove `answer-question` from canonical skill packages,
discovery, and routing. Maintained examples and assessment fixtures SHALL
reference an available skill. Repository policy SHALL continue to preserve the
distinction between a question and authority to mutate state.

#### Scenario: Discover skills after removal

- **WHEN** repository skill discovery is regenerated
- **THEN** it contains quality-assessment and decision-review without
  answer-question or a stale link to its former package
- **AND** the quality-assessment routing no longer overlaps a separate
  generic question-answering skill
