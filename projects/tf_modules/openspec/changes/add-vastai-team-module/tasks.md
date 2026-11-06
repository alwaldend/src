## 1. Module and caller

- [x] 1.1 Define the reusable team/key module and wire source packaging and documentation.
- [x] 1.2 Exercise two teams through Terraform CLI fixtures, including unchanged plans and cleanup.
- [x] 1.3 Validate formatting, owner OpenSpec, caller packaging, and delivery inputs.

## 2. Existing-team reuse

- [x] 2.1 Write actual-module E2E failure cases for existing-team use, failed activation, and legacy state upgrade before implementation.
- [x] 2.2 Implement existing-team reuse and verify no team or key replacement during upgrade.
- [x] 2.3 Route the credential caller through the module and verify packaged inputs, formatting, and owner specifications.

## 3. Import-based ownership

- [x] 3.1 Before implementation, update E2E cases for importing the module's team, invalid import IDs, scoped-key rejection, activation failure, key-only cleanup, and unchanged legacy addresses.
- [x] 3.2 Remove the existing-team selector, counted resource, and migration block; derive key binding from imported team state and use provider authentication.
- [x] 3.3 Validate the candidate and document import prerequisites and protection of imported teams.
