# Verification

The reusable module is `projects/tf_modules/vastai_instance_imagegen`. Its
Terraform CLI acceptance fixture uses the real source and provider against
loopback HTTP only. No Vast.ai credentials, live plans, template edits, or
GPU rentals were used. The selected template remains a caller-owned input.

## Evidence

- Targeted image-generation and legacy-schema fixtures passed, including
  template lookup, quoted environment settings, unsafe Docker options, missing
  or ambiguous metadata, datacenter search filters, empty and over-budget
  searches, offer-result changes without replacement, template replacement,
  pending SSH metadata, and shell-quoted loopback forwarding after refresh.
- `imagegen-lifecycle.json` in the acceptance test's undeclared outputs records
  the forwarding command and fixture lifecycle operations. Reproduce with
  `bazel test //projects/vastai_terraform_provider/test/acceptance:acceptance_test`.
- The source README documents template preparation, defaults, caller packaging,
  SSH access, storage charges, and the remaining live acceptance boundary.

GPU boot, model downloads, live host ports, and actual image generation remain
unverified pending explicit rental approval. Application provisioning must be
reviewed by the template owner; Docker-option validation is not a code sandbox.

## Repository checks

The scoped checks passed: module formatting, the full provider acceptance suite,
the existing infrastructure catalog consumer, and all three owning OpenSpec
validators. The provider suite completed in 64.9 seconds. The prepared candidate
also passed `//:repo_quality_test` (33 tests including the scoped checks in that
invocation) and semantic lint over the module, provider sources and acceptance
package, SDK client, and owner specification exports. `git diff --check` passed.

Local review covers the complete merge-base-to-candidate diff. Exact commit
identities and publication state remain in ignored delivery/review receipts.
The remaining live acceptance boundary is described above and in the module
README; offline success is not evidence of a running GPU application.

## Session ergonomics

Stable task: `add-vastai-instance-imagegen`. Supporting changes are
`support-imagegen-templates` in the provider and `read-instance-templates` in
the SDK. Git and delivery receipts own publication state.

- `template-routing-assumption` (configured): initial image-pinning discussion
  preceded checking provider template support. Corrected to a caller-selected
  template; no new image dependency was declared. For similar tasks, inspect
  provider template and environment support before choosing an image.
- `fixture-provider-resolution` (fixture-tested): initialization attempted
  Registry access for the development provider. The existing acceptance owner
  already demonstrates `terraform get`; the fixture now follows that contract.
- `fixture-source-placement` (fixture-tested): a test function was initially
  inserted into embedded HCL. Corrected before acceptance. Review embedded
  configuration boundaries before running the full suite.
- `publication-keyring-context` (routed): sandboxed GitHub status reported an
  error, while the approved execution context reported keyring authentication
  success. Delivery uses that context without changing credentials or transport.

No shared policy or host configuration was changed.
