## Why

The reference-fidelity skill's fixed scoring gate and mandatory review cycle
can prevent an agent from delivering requested work, even after the user has
chosen to proceed. The guidance should improve visual work without replacing
the user's acceptance and delivery decisions.

## What Changes

- Keep reference comparison, visible-defect checks, and honest reporting.
- Make views, measurements, review depth, and iteration effort proportional
  to the available references and the requested outcome.
- Remove the universal numeric score and export prohibition; honor explicit
  user acceptance criteria and delivery instructions.
- Update the skill's evaluation cases to cover both honest quality reporting
  and user-authorized delivery of an imperfect result.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-agents`: Reference-fidelity guidance supports visual judgment and
  preserves the user's authority over the requested result and delivery.

## Impact

The canonical Blender reference-fidelity skill, its supporting references and
offline evaluation are updated. No new dependency or live configuration change
is required.
