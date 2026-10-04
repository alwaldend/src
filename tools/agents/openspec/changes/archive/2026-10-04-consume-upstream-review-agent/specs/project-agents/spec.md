## MODIFIED Requirements

### Requirement: Evidence-based change review

The repository SHALL require the pinned upstream review-agent skill for bounded
change reviews, including bugs and applicable repository standards. Its upstream
instructions and invocation metadata SHALL remain unchanged. The source archive
SHALL have an immutable commit and integrity pin, with license and provenance
retained. Discovery SHALL be generated from the archive-backed skill target. It SHALL identify the
candidate and comparison baseline, evaluate affected behavior and consumers,
and report actionable findings with evidence and consequence. It SHALL
preserve the distinction between reviewing and authorizing modifications or
external review messages. Each finding SHALL be discrete, actionable,
introduced by the change, demonstrable from code or an applicable contract,
and meaningful to correctness, security, performance, or maintainability.
The reviewer SHALL inspect the complete diff, exclude speculation and harmless
style preferences, and confirm findings through relevant tests or call sites.

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

#### Scenario: Review against a base branch with an ahead upstream

- **WHEN** the requested comparison branch has a configured upstream ahead of
  its local ref
- **THEN** the reviewer selects that upstream and computes the candidate's
  merge base before inspecting the actual changes
- **AND** unrelated base-branch changes are not treated as candidate defects

#### Scenario: Report a qualifying regression

- **WHEN** a concrete introduced defect passes the finding eligibility gate
- **THEN** the result provides a severity-prefixed imperative title, a minimal
  changed file location, and a short explanation of the trigger and consequence
- **AND** it continues reviewing the remaining diff and reports every qualifying
  finding before the overall assessment and validation limits
