## Why

The shared security workflow must consume this owner's existing inventory.
The parent separately verified an active secure runner missing from that
inventory and authorized recording its confirmed hostname here.

## What Changes

Expose `//infra/forgejo_runner/ansible:inventory` to
`//infra/security/ansible:__pkg__`, preserving existing access.
The target contains only the authoritative `ansible/inventory.yaml`.
Add `secure.forgejo-runner.alwaldend.com` alongside the existing host in
the `forgejo_runner` group. Preserve existing host and Terraform definitions.

## Capabilities

### New Capabilities

- `runner-inventory`: Record the confirmed secure runner in the owner inventory
  so owner deployment and shared security consumers use the same host source.

## Impact

The owner inventory and its BUILD export. Consumer packaging and execution
belong to `infra/security`. No deployment is authorized by this record.

## Evidence

On 2026-10-29 the parent reported independent live XO and root-audit evidence
for the secure runner hostname, Fedora 44, and address `192.168.10.101`, plus
successful security-role SSH access. The inventory records only the confirmed
hostname; it does not duplicate DNS addresses or infer a provisioning change.
The earlier source inspection found only `host1.forgejo-runner.alwaldend.com`
and no source evidence identifying that legacy entry with the secure VM.
