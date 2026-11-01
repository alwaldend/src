## Source acceptance

This section records the original candidate at `1ef9bbde05cd9885142f4552ba590056ae6d2050`.
The user subsequently selected normal runtime networking and Auto/Full
approval modes. Current behavior and acceptance are maintained in
`../add-t3code-mtls/`; old firewall observations below are historical.

These checks validate prepared source, not deployed infrastructure:

- Backend-disabled initialization and validation passed with the repository's
  Yandex provider 0.203.0 and a scratch copy of the Terraform root.
- The packaged Ansible syntax check passed with a localhost-only fixture.
- The six configuration templates rendered through packaged Ansible into
  ignored task scratch. YAML and rendered TOML parsed successfully.
- Codex 0.166.0 parsed the root-owned managed policy inside a network-isolated
  namespace. Local shell and resumable command tools remained enabled.
- Attempts to enable apps/plugins left both disabled. A supplied remote MCP
  server was reported disabled by `/etc/codex/requirements.toml`.
- Semantic lint and builds passed for the component, its Terraform/Ansible/
  OpenSpec packages, and the affected Vault packages.
- Repository quality, Terraform formatting, and the component/Vault OpenSpec
  checks passed: 30 tests. Delivery reruns the mandatory gates on its exact
  committed candidate; its private receipts own that publication evidence.

Commands and non-secret diagnostics remain under ignored
`out/local_inference/`. The generated OpenSpec workspace inventory was
updated through `//tools/repo_quality:inventory_update`.

## Live acceptance not performed

No Terraform plan/apply, Vault write, or inventory-host connection ran.
GPU quota/capacity, the selected immutable GPU image, runtime archive loading,
model VRAM usage, fresh Codex inference and tool follow-up, actual IPv4/IPv6
filter enforcement, and reboot behavior require separately authorized live
acceptance. The deployment requires a compatible preinstalled proprietary
NVIDIA driver, an image ID, administrator CIDRs, and a verified model artifact.

## Session ergonomics

- `LOCAL-INFERENCE-PROVIDER-SCHEMA`: provider validation rejected instance
  `deletion_protection`; corrected to Terraform `prevent_destroy` before
  acceptance. Evidence tier: configured, provider-validated.
- `LOCAL-INFERENCE-PACKAGE-PROVIDERS`: initial AL and Ansible packaging used
  mismatched attributes/providers; corrected against the owning rule and
  sibling package. Evidence tier: build-tested.
- `LOCAL-INFERENCE-LOCAL-FIXTURE`: scratch rendering needed explicit local
  and remote temporary directories and structured YAML serialization.
  The final localhost-only fixture passed. Evidence tier: fixture-tested.
- Resume from this record, tasks, and delivery receipts. No shared policy or
  skill changes are proposed for these task-local implementation corrections.
