## 1. Generated CLI

- [x] 1.1 Verify an independent master-based worktree, existing dependencies, and
      ownership; record behavioral failures in design.md.
- [x] 1.2 Write built-command E2E before implementation and verify the native test passes.
- [x] 1.3 Generate upstream commands and implement shared JSON/connection lifetime;
      verify native edits, import/export, and all three transports through E2E.
- [x] 1.4 Document usage and review the complete diff; inspect generated registrations
      and native artifacts, resolving in-scope findings.
- [x] 1.5 Format and verify source; select mandatory exact-candidate quality, lint,
      native E2E, and site gates for the separate master-based PR.

The generator produces 269 method commands across 30 nonempty services. Evidence,
artifacts, and delivery receipts belong in out/anki-rpc-cli/; receipts own final
validation and publication state.
