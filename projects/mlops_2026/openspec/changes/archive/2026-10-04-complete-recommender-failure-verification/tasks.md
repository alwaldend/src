## 1. Complete verification

- [x] 1.1 Exercise missing CSVs and corrupt model startup through the executable;
      verify nonzero exits and retain JSON results.
- [x] 1.2 Run the actual training wrapper with malformed CSV input; verify
      failure without an output artifact and retain JSON results.
- [x] 1.3 Render and inspect the client HTTP actor in SVG and WebP; retain SVG
      documentation wiring and present the WebP preview.
- [x] 1.4 Review the complete PR and select receipt-bound quality, lint, and E2E
      checks for publication. Final validation and PR-head verification are owned
      by repo-delivery receipts, which must pass before handoff.
