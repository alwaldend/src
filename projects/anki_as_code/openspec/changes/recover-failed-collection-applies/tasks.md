Deferred; implementation has not started. These unchecked tasks are future
work; this planning request authorizes no collection or code changes.

## 1. Define recoverable failure behavior

- [ ] 1.1 Enumerate failures before/after native writes and during binding or
      recovery-information persistence; review output versus journaling through
      decision-review and record the chosen recovery contract and its limits.
- [ ] 1.2 Extend the existing native failure fixture through user-visible
      recovery and retry before implementation; demonstrate the current gap
      with a disposable collection and verifiable identity/state artifacts.

## 2. Implement and validate recovery

- [ ] 2.1 Implement the selected recovery information and supported identity
      reconciliation path; verify handled failures identify completed native
      work and stale or ambiguous recovery refuses unsafe bindings.
- [ ] 2.2 Run failure, recovery, and successful retry against the exact candidate;
      retain a repeatable E2E artifact proving no duplicate deck, stable native
      identities and retained history/media, plus unchanged source configuration.
- [ ] 2.3 Document the supported recovery procedure and limitations; execute
      the documented steps on the disposable fixture and link the observed result.
