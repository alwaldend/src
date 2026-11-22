## Deployment

Use the `48gb` hardware profile and the `noobai` Ansible asset bundle. Check
team credit and review the Terraform plan before renting. Install the
application bundle on the managed rental and run its saved text-to-image
workflow through the authenticated ComfyUI endpoint.

## Deployment on 2026-11-22

The reviewed Terraform apply created Vast.ai instance `61937807` with zero
changes or destructions. Vast.ai reported one RTX PRO 5000 with 48,935 MiB of
VRAM and an hourly rate of USD 0.9077. The rental has 100 GB of disk. A
refresh-only apply recorded the connection outputs. The direct and proxy SSH
endpoints presented the same Ed25519 host key, and authentication succeeded
over direct SSH with the pinned key.

Ansible installed four model files, the selected nodes and interpolation
engine, and 18 workflows and Apps. The play recap reported 28 successful
tasks, 11 changed tasks, and zero failures or unreachable hosts. The selected
node interfaces became available after restart. All 18 installed JSON files
matched the packaged source by SHA-256. The authenticated endpoint reported
ComfyUI 0.44.0 and an NVIDIA RTX PRO 5000 Blackwell.

The packaged `noobai-text-to-image` workflow completed through the
authenticated API. The downloaded result decoded as a 1024 × 1024 PNG.
Verification of the other workflows covered installation and node availability.

## App follow-up on 2026-11-22

The initial deployment omitted the NoobAI text-to-image App and the Danbooru
Tag Search bundle. Ansible added both Apps and the Tag Search graph to the same
rental. All three new files matched their packaged checksums. ComfyUI listed
both Apps after restart. A test of the NoobAI App produced a 1024 × 1024 PNG,
and the Tag Search App returned 20 results for `blue hair`, with `blue_hair`
first. The existing 18 workflow files matched their prior checksums before the
update.

The local deployment fixture truncates long failure output from the beginning
of the playbook, which hid the failing NoobAI task during the first run. In a
future fixture update, retain the bounded tail of the output so the failing
task and play recap are visible. Verify with an intentionally failing local
installation after a long validation phase.

## Rental release on 2026-11-22

At the operator's request, a reviewed instance-only Terraform apply destroyed
contract `61937807`. A refresh-only apply cleared the instance ID and the
HTTPS and SSH outputs. The stored web login was retained. The Vast.ai team
inventory no longer listed the contract after the destroy.
