## ADDED Requirements

### Requirement: Resolve agent ergonomics findings before final preparation

For substantial work or recurring friction, the delivery caller SHALL complete
an agent ergonomics review before final candidate preparation on each supported
forge. The caller SHALL correct actionable in-scope findings and include those
corrections in the candidate. The caller SHALL record each actionable deferred
finding in maintained OpenSpec, with evidence, the responsible component, the
next action, and a verification criterion.

#### Scenario: Review finds an in-scope correction and a deferred improvement

- **WHEN** an ergonomics review identifies an in-scope correction and an improvement outside the current task
- **THEN** the caller includes the correction in the candidate and records the deferred improvement in maintained OpenSpec
- **AND** the caller still performs the regular local review of the final candidate

#### Scenario: Review has no actionable findings

- **WHEN** an ergonomics review finds no actionable issue
- **THEN** delivery proceeds without creating a separate ergonomics artifact
