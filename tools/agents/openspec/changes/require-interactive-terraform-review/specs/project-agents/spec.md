## ADDED Requirements

### Requirement: Guide Terraform review without plan artifacts

The repo-infra skill SHALL apply the root AGENTS.md Terraform review policy
through the owning repository targets. It SHALL distinguish inspection from
deployment authority and SHALL report incompatible saved-plan-only targets
without bypassing their contract.

#### Scenario: Inspect proposed changes

- **WHEN** an agent is authorized to review a live Terraform plan
- **THEN** it uses normal text output without saving a plan or requesting JSON
- **AND** it retains only sanitized summaries as durable evidence

#### Scenario: Apply an authorized change

- **WHEN** an agent is authorized to apply changes in an owning Terraform root
- **THEN** it reviews the plan displayed by the interactive apply and answers
  yes only when that plan matches the authorized scope
- **AND** an earlier standalone plan does not substitute for that review

#### Scenario: A target accepts only a saved plan

- **WHEN** the owning apply target requires a saved plan file
- **THEN** the agent reports the incompatibility without consuming a saved
  plan, bypassing the guard, or switching to a different target or raw Terraform

#### Scenario: Evaluate agent guidance offline

- **WHEN** the skill's offline evaluation configuration is validated
- **THEN** its cases cover inspection, interactive approval, output retention,
  incompatible targets, and recovery after partial failure
- **AND** the result does not claim live infrastructure behavior was tested
