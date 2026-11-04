## ADDED Requirements

### Requirement: Delivery consumes Codex reviews without initiating them

The delivery procedure SHALL NOT request a Codex review or post a Codex mention
to initiate a review after publication, including after a review-driven rewrite.
It SHALL continue inspecting remote feedback, handling findings with guarded
operations, and reporting review outcomes for the published candidate. Codex
review configuration SHALL remain outside the delivery skill.

#### Scenario: Initial publication

- **WHEN** delivery creates a pull request
- **THEN** it inspects available remote feedback without requesting a Codex
  review, posting a Codex mention to initiate review, or configuring Codex reviews

#### Scenario: Review-driven update

- **WHEN** delivery publishes a correction in response to review feedback
- **THEN** it inspects feedback for the new head without requesting another
  Codex review or posting a Codex mention to initiate one
- **AND** it preserves guarded replies, thread resolution, and honest reporting
  when review completion cannot be observed
