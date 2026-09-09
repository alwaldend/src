## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: Scope this project to reusable skills

`projects/agents` SHALL own reusable cross-repository skills and their
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
