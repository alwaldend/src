# Agents Specification

## Purpose

Describe the reusable repository skills and their packaging and discovery
contract as owned by `tools/agents`. Component behavior, repository policy,
and requested authority remain at their own owners.

Sources: [project README](../../../README.md),
[project packaging](../../../BUILD.bazel), and
[discovery declaration](../../../../../.agents/BUILD.bazel).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08; durable-work routing includes this OpenSpec migration.

## Requirements

### Requirement: Preserve natural fact ownership

This project SHALL keep component behavior at its owning README and
implementation, repository policy at applicable `AGENTS.md` files, and
requested outcome and action authority in the user interaction. Project
documentation SHALL NOT independently authorize actions.

#### Scenario: A document describes an unimplemented interface

- **WHEN** a document describes an interface without an implemented provider
  and supporting evidence
- **THEN** it states that status explicitly rather than implying a supported
  runtime capability

### Requirement: Package canonical skills separately from discovery

Enabled repository skills SHALL have canonical project-owned directories and
`skill_library` targets. Each skill SHALL grant `//.agents:skill_discovery`
read access. `.agents/BUILD.bazel` SHALL declare the complete discovery set
through `skills_write`, and `.agents/skills/` SHALL expose relative symlinks to
the canonical directories while archive-backed skills are materialized as
regular files beside them. The discovery directory SHALL contain only those
managed entries. The removed goal skill SHALL NOT be reintroduced.

#### Scenario: An agent discovers a reusable procedure

- **WHEN** an enabled skill is selected through `.agents/skills/`
- **THEN** its content comes from its canonical owning project without creating
  a second Bazel package for the discovery link

#### Scenario: Add or rename a skill

- **WHEN** a contributor adds a canonical skill
- **THEN** they add its `:skill` label to the `.agents/BUILD.bazel` declaration
  and grant `//.agents:skill_discovery` visibility
- **AND** `//.agents:write_skills` regenerates the links and
  `//.agents:write_skills_test` verifies the exact state

### Requirement: Use OpenSpec for maintained work

Maintained agent-system work SHALL use this project's OpenSpec changes for its
proposal, acceptance requirements, tasks, and continuation evidence. Migrated
goal history SHALL remain available as historical evidence with its original
acceptance and observation limits.

#### Scenario: Work resumes after migration

- **WHEN** an agent continues maintained repository-agent work
- **THEN** it follows the corresponding OpenSpec change and retained migration
  history instead of creating a new legacy goal record.

### Requirement: Bound claims made from evaluations and catalogs

Documentation SHALL distinguish offline eval-configuration validation from
behavioral correctness, configured capability declarations from runtime
availability, and stored observations from current provider health.

#### Scenario: An offline skill evaluation configuration passes

- **WHEN** the configured Promptfoo validation target succeeds
- **THEN** the result establishes validity of the evaluation harness without
  claiming successful skill routing or real task completion.

### Requirement: Typed Python development guidance

The discoverable repo-python skill SHALL require type hints throughout Python
code and identify the owning dependency manifest, shared tool configuration,
generated files, and pinned validation workflows.

#### Scenario: Add Python implementation

- **WHEN** an agent selects repo-python for a Python change
- **THEN** it is instructed to annotate parameters, return types, and variables,
  preserve component support constraints, and use owning Bazel workflows
- **AND** project support metadata and shared tool settings are distinguished

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

### Requirement: Select maintained work before implementation

The repo-workspace and repo-openspec procedures SHALL apply the repository's
AGENTS.md requirement to select or create an owner-local OpenSpec change for
all repository modifications, including trivial changes. They SHALL preserve
existing authority and keep artifacts proportional to the requested work.

#### Scenario: Correct a single documentation typo

- **WHEN** an agent begins a one-line documentation correction
- **THEN** it verifies isolation and selects a matching owner-local OpenSpec change
  before editing, reusing an active matching record where possible
- **AND** it uses skip_specs for work without a behavioral delta rather than
  inventing requirements or requesting redundant approval

### Requirement: Repository OpenSpec integration

The repo-openspec skill SHALL define owner selection, pinned
CLI invocation, and maintained-work conventions, and SHALL route operation
procedures to the discoverable upstream OpenSpec skills. It SHALL NOT independently
maintain generic create, apply, specification-writing, or archive procedures.
Repository policy SHALL remain owned by AGENTS.md.

#### Scenario: Apply work in a component workspace

- **WHEN** an agent implements a maintained component change
- **THEN** the integration skill selects the owning OpenSpec workspace and pinned
  Bazel invocation, and routes implementation procedure to openspec-apply-change
- **AND** the upstream operation skill owns artifact instructions and task steps

#### Scenario: Validate before archiving

- **WHEN** an implementation is ready for acceptance and archive
- **THEN** the integration skill routes to openspec-verify-change and
  openspec-archive-change while preserving repository validation and authority

### Requirement: Scope this project to reusable skills

`tools/agents` SHALL own reusable cross-repository skills and their
packaging contract. It SHALL NOT maintain a system architecture, roadmap, or
catalog program for agent tooling. Component behavior SHALL stay at its owning
README and implementation, repository policy at the applicable `AGENTS.md`
chain, and requested authority in the user interaction. Build files and
documentation SHALL NOT reference the removed `tools/agents` catalogs, context
capsule, or admission libraries as available interfaces.

#### Scenario: A fresh agent orients in an unfamiliar subtree

- **WHEN** an agent needs applicable policy, owner, and build context
- **THEN** it reads the `AGENTS.md` chain, the nearest owner README, and the
  relevant BUILD and MODULE files
- **AND** it loads the specific skill that matches its task
- **AND** it does not rely on a catalog or capsule that this repository no
  longer provides

#### Scenario: A contributor considers documenting an agent-system design

- **WHEN** a contributor wants to describe agent-system composition in this
  project
- **THEN** the description belongs with the component that owns the behavior,
  or in an OpenSpec change while the work is active
- **AND** this project does not accumulate a maintained architecture or roadmap
  for tooling it does not implement

### Requirement: Evaluate skills without a coverage catalog

Skill evaluation SHALL remain with each canonical skill and SHALL be validated
offline by its `eval_config_test` target. Repository documentation SHALL
distinguish this configuration validation from successful routing or task
completion and SHALL NOT claim behavioral results it did not measure.

#### Scenario: An offline skill evaluation configuration passes

- **WHEN** the configured Promptfoo validation target succeeds
- **THEN** the result establishes validity of the evaluation harness without
  claiming successful skill routing or real task completion

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
