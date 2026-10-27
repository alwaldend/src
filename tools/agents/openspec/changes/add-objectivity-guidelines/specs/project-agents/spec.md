## ADDED Requirements

### Requirement: Document cognitive objectivity principles

The root `AGENTS.md` SHALL contain the user-supplied Cognitive Objectivity
and Anti-Sycophancy Principles. That guide SHALL remain the canonical policy;
the supporting research link SHALL be retained in agent-system documentation
outside `AGENTS.md`.

#### Scenario: Read the repository operating policy

- **WHEN** an agent reads the root guide
- **THEN** it finds the supplied instructions on leading questions,
  independent premise verification, and direct rebuttals

#### Scenario: Locate supporting research

- **WHEN** a reader opens `tools/agents/README.md`
- **THEN** the supplied Anthropic research link is available as an external
  reference without embedding it in `AGENTS.md`
