## 1. Correct the contract gaps

- [x] 1.1 Write behavioral failure cases before implementation and record the
      expected failures against the previous behavior.
- [x] 1.2 Correct Markdown mutability and document/enforce a defensible budget.
- [x] 1.3 Validate static image format, size, and animation before any upload.
- [x] 1.4 Enforce media readiness without replaying upload or draft POSTs.
- [x] 1.5 Scope persistent media reuse and handle known and unknown lifetimes.
- [x] 1.6 Update documentation and the conversion showcase to match behavior.

## 2. Verify the candidate and attempt once

- [x] 2.1 Pass the E2E suite, relevant consumer checks, lint, and specification
      validation; inspect a freshly converted artifact and independent review.
- [x] 2.2 Run one authorized draft invocation and record its result, with no
      automatic retry and no claim of success without a returned draft ID.

- [x] 2.3 Perform the separately authorized single attempt without proxy variables
      and record the TLS handshake timeout without another invocation.
