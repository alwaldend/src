## 1. Dependency maintenance

- [x] 1.1 Remove the exact Vast.ai requirement, add secure dependency overrides, and regenerate uv.lock through its Bazel target; inspect the resolved versions and hashes.
- [x] 1.2 Update the README to describe lock-owned versions and overrides; verify it matches the manifest and generated lock.

## 2. Verification and delivery

- [x] 2.1 Check every resolved package against PyPI release advisories and verify representative offline CLI behavior with the updated dependencies.
- [x] 2.2 Run OpenSpec validation, repository formatting and quality checks, and affected-target semantic lint; record the observed results.
- [x] 2.3 Prepare task-owned delivery inputs and verify source scope with al-delivery; Git and delivery receipts own final commit and PR state.
