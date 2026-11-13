## 1. Establish failure coverage

- [x] 1.1 Extend the existing real-binary upgrade and restore scenario with
      quoted external repositories, legacy storage, independently incomplete
      backup roots and actual removal before full archive restoration.

## 2. Correct backup validation

- [x] 2.1 Normalize supported INI quoting and reject ambiguous syntax; verify
      quoted and malformed path configurations through the role fixture.
- [x] 2.2 Resolve effective legacy storage and validate startup-file coverage;
      verify omitted roots and relocation fail before service downtime.

## 3. Verify the correction

- [x] 3.1 Run integration acceptance, packaging, semantic lint, explicit Ruff,
      repository quality and OpenSpec validation; record outcomes and limits.
- [x] 3.2 Review the complete corrected source diff and prepare delivery inputs.
      Git and delivery receipts own publication and remote review state.
