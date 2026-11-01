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

## Original live acceptance limits

At the original source acceptance, no Terraform plan/apply, Vault write, or
inventory-host connection ran.
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

## DNS execution failure cases and acceptance

Before changing the DNS packaging, record the observed and required cases:

- The DNS runner must start without the absent Yandex account credential.
- DNS must not require an image ID or administrator CIDRs, or substitute dummy
  provisioning values. Full provisioning must retain these required inputs.
- The DNS runfiles must contain only DNS and shared provider/backend Terraform
  sources. Their module address and HTTP backend must remain unchanged.
- A saved DNS plan must contain only the three declared Cloudflare records;
  no GPU, network, disk, KMS, or Yandex resource action is acceptable.
- Apply must consume that exact saved plan; verify public records and a fresh
  no-change scoped plan afterward.

## Authorized DNS execution on 2026-11-03

Vault certificate recovery and the prerequisite AppRole/PKI apply are recorded
in `infra/vault/openspec/changes/add-local-inference-approle/` and PR #242.

The first DNS command stopped on the absent provisioning-only Yandex key.
The corrected DNS target packages only `dns.tf` and `provider.tf`, retains the
same owner HTTP backend and `module.dns` address, and selects no Yandex login.
The full root still packages its Yandex provider and requires the unchanged
image ID and administrator CIDRs. No placeholder provisioning values were used.

The source candidate passed 34 offline tests and affected semantic lint before
the apply. A complete paginated Cloudflare inventory contained 54 records and
none at the requested names. The reviewed saved plan contained exactly three
Cloudflare creates and no other resource actions. Apply completed with
`3 added, 0 changed, 0 destroyed`; a fresh scoped plan exited 0 with no changes.
A complete post-apply inventory contained 57 records and exactly the expected
unproxied CNAME and two NS records. Both authoritative Cloudflare nameservers
returned the Yandex delegation; a public recursive query returned the CNAME.

Only the public delegation and CNAME were deployed. The delegated Yandex zone,
GPU VM, address record, disk, KMS key, and host configuration were not deployed;
this record does not establish a working T3 endpoint. Final publication receipts
own the delivered source revision; private raw plans and state are disposable.

`LOCAL-INFERENCE-DNS-INPUT-COUPLING` is verified live: splitting the runner's
declared inputs removed the unrelated account-key and GPU-variable blockers
without relaxing provisioning inputs. `DELIVERY-INDEX-STAT` was verified by
hashing all 14 reported Go/module files: content and modes matched the index,
and source/index diffs were empty. Refreshing only index stat metadata allowed
preserved-source preparation to continue with an exact remote replacement lease.

## Independent DNS root extraction

The user subsequently selected a separate DNS project and state. The prior
filtered-runfiles acceptance above describes the initial shared-state apply.
`infra/local_inference/dns` now owns the canonical Cloudflare records through
a separate backend; its `separate-dns-terraform-root` change owns the migration
procedure and execution evidence. Provisioning removes Cloudflare configuration
and DNS command variants while retaining the canonical document for its Yandex
zone and VM address. No GPU provisioning or host deployment is authorized by
this extraction.

The independent root migration completed with three imports, zero creates,
updates or destroys, and identical record IDs/values. Legacy state now owns no
DNS records; both plans reported no changes and public DNS was unchanged.
This provisioning cleanup passed 34 offline checks and backend-disabled
Terraform init/validate with only the pinned Yandex provider. The independent
DNS owner's change retains the detailed migration evidence.

## Terraform extraction

The user subsequently requested a separate provisioning PR. The unchanged
Terraform definitions and their existing backend/provider references are now
owned by tf_setup with stage-local AL wiring in that independent change.
PR #233 retains host deployment and signed SSH injection; it declares no cloud
resources or Yandex credential calls. This source split performs no live plan,
apply, state transfer, account provisioning or host deployment.
