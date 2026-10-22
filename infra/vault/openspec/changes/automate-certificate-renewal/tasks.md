Deferred; implementation has not started. These unchecked tasks are future
work; this planning request authorizes no scheduler, renewal, or deployment.

## 1. Specify and exercise the renewal lifecycle

- [ ] 1.1 Inspect existing scheduling, reporting, and authentication mechanisms;
      record the selected owner, interval, renewal/alert thresholds, minimum
      credential scope, and independent recovery prerequisites in a reviewed
      design using decision-review.
- [ ] 1.2 Establish an isolated acceptance exercise before code covering normal
      renewal, unavailable issuer, failed installation/reload, and expired TLS;
      verify each failure has an observable expected outcome without live hosts.

## 2. Implement and verify declarative renewal

- [ ] 2.1 Extend the owning role and configuration with scheduling, safe
      activation, and existing reporting integration; verify packaged desired
      state contains the selected policy and source contains only secret references.
- [ ] 2.2 Run the isolated lifecycle and recovery scenarios against the exact
      candidate; save a repeatable sanitized report proving replacement expiry,
      preserved usable material on failure, emitted expiry warnings, and normal
      TLS verification after recovery, then update the operator procedure.
