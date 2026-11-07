## ADDED Requirements

### Requirement: Proportional Blender reference review

The Blender reference-fidelity skill SHALL guide agents to compare rendered
results with relevant available references, inspect material visible defects,
and report limitations. It SHALL preserve user-specified acceptance criteria
and delivery instructions rather than impose a universal numeric export gate.

#### Scenario: A user requests delivery with known differences

- **WHEN** the user directs delivery of a result that has been visually reviewed
  but still differs from the reference
- **THEN** the skill directs the agent to honor that request and disclose the
  material differences without describing the result as an exact match

#### Scenario: Reference views are incomplete

- **WHEN** the source does not provide a canonical view of part of the subject
- **THEN** the skill directs the agent to use the available views and report
  the uncertainty rather than require an unavailable reference

#### Scenario: Iteration stops improving a model

- **WHEN** repeated edits do not materially improve the same visible defect
- **THEN** the skill directs the agent to reconsider the modeling approach and
  continue useful work toward the user's stated outcome
