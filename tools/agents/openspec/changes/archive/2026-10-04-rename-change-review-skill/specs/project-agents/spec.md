## MODIFIED Requirements

### Requirement: Evidence-based change review

The discoverable review-change skill SHALL review a bounded change set for
bugs and violations of applicable repository standards. It SHALL identify the
candidate and comparison baseline, evaluate affected behavior and consumers,
and report actionable findings with evidence and consequence. It SHALL
preserve the distinction between reviewing and authorizing modifications or
external review messages.

#### Scenario: Review a pull request

- **WHEN** the requested change set is a pull request
- **THEN** the review identifies its exact head and base, inspects the applicable
  policy and affected owners, and reports evidenced findings with file locations
- **AND** it does not post comments, fix source, or merge without authorization

#### Scenario: Review working-tree changes

- **WHEN** the requested change set is staged or uncommitted source
- **THEN** the review identifies the baseline and the included index, worktree,
  and untracked content, and bounds claims to that observed candidate
- **AND** it investigates relevant consumers beyond changed lines as needed

#### Scenario: Insufficient evidence or no actionable finding

- **WHEN** inspection cannot prove a suspected defect or finds no actionable issue
- **THEN** the response distinguishes uncertainty and validation limits from
  confirmed bugs without inventing findings or treating passing checks as proof
