## Boundary

Keep infra/local_inference/tf_setup and all Terraform resource addresses intact.
Move its existing backend and Yandex login calls into tf_setup/al.lua with a
stage-owned al_config target. Retain the component AppRole and the
/tf_backend/tf_setup state path. The host AL configuration keeps only SSH access.

The master-based PR consumes the already merged canonical DNS declaration and
shared cloud-init input. It includes no Ansible deployment or runtime artifacts.
The host PR keeps its existing runtime configuration and acceptance records.

## Acceptance and failure cases

- Compare all extracted .tf bytes to the published feature commit.
- Verify provider selection is Yandex-only and backend/auth references unchanged.
- Ensure host packaging has no dependency on tf_setup or Yandex credentials.
- Reject dangling source/module/data references through backend-disabled init
  and Terraform validate, formatting, BUILD checks, owner checks and lint.
- Publish both task-owned branches with exact remote leases. No live plan,
  apply, account provisioning, state migration or host command was included
  in the initial extraction phase.

## Authorized deployment continuation

The operator approved the scoped folder/account prerequisite, one initial VM
boot followed by an immediate stop, reuse of the existing Fedora image, and
public SSH. Use immutable Fedora 43 image fd8iku26nnkveh8s4di1 and the shared
cloud-init document without the earlier Ubuntu group override. Keep image and
administrator CIDRs explicit inputs; permit the authorized IPv4 /0.

The provider cannot manage stopped power state. The authorized initial stop
is a one-time control-plane exception, scoped to the new Terraform-owned VM.
Prepare the stop before applying, verify the instance identity and STOPPED
status, then reconcile with a fresh Terraform plan. All persistent resources
remain owned by this root. Do not run the host playbook: the current host PR
requires Ubuntu and a preinstalled NVIDIA driver, which this Fedora image
does not establish.

## Availability-zone replacement

After the ru-central1-a quota failure, the operator requested a different
availability zone and explicitly authorized replacement of unused component
resources. Proceed with ru-central1-b, the other supported Broadwell V100 zone.
Reusing the current zone would repeat the known quota failure; moving the
unused zonal disk and subnet is the smallest change that can use another zone.
Quota headroom in the new zone remains unverified because the component cannot
read cloud-level quotas, so a successful Terraform plan alone does not establish
that VM creation will succeed.

Temporarily remove the disk's prevent_destroy guard only to execute the
authorized replacement through this Terraform root. The disk has never been
attached to a VM. Keep the KMS key and instance guards, preserve non-zonal
resources when the reviewed plan permits, and restore the disk guard after
the replacement attempt. The initial-boot-and-immediate-stop limit still applies.

The ru-central1-b attempt was rejected at disk allocation. Reopen the decision:
the published incident listing reports that zone unavailable, while the GPU
documentation supports this V100 platform only in a/b. Leave the requested
zone staged with the guard restored and do not retry unchanged restrictions.
Returning to a with resolved quota or selecting another GPU requires a new
operator choice; do not silently change hardware or return to the failed zone.

## Smaller V100 retry in zone a

The operator selected a return to ru-central1-a with a smaller VM. Use the
documented gpu-standard-v2 configuration: one 32 GB V100, eight vCPUs and 48 GiB
RAM. This halves the memory request without changing GPU capacity. Preserve
the 200 GiB encrypted disk, Fedora image, public SSH input and initial stop.
The disk is absent, so its restored guard need not be lifted. Only the unused
zone-b subnet requires replacement. Plan and review this changed input before
another apply; no quota increase is assumed or requested.

## Original size after the operator's quota increase

The operator reported increased quota and retried 48 GiB; allocation then
failed with Not enough resources. Catalog review confirmed that V100 is the
32 GB option and its two platforms are not supported in zones d/e. The live
zone API reports b UP, while this cloud's current b CPU/RAM/SSD quotas remain
zero. Larger GPU alternatives require larger minimum VMs.

The operator then explicitly selected the original configuration after
increasing RAM quota: gpu-standard-v1, eight vCPUs, 96 GiB RAM and one 32 GB
V100 in ru-central1-a. Preserve the existing encrypted disk, subnet, image,
SSH settings and all destruction guards. Review a new plan and stop the VM
immediately after its authorized initial boot. Do not assume allocation
capacity from a successful plan or quota increase.

## Operator-selected preemptible retry

Both regular V100 configurations failed resource allocation. After the
preemptible interruption behavior, maximum 24-hour run and public running
price were explained, the operator explicitly requested trying it. Keep
gpu-standard-v1, one 32 GB V100, eight vCPUs and 96 GiB RAM in zone a, changing
only scheduling_policy.preemptible to true. Preserve the existing six
supporting resources and all destruction guards. Generate a fresh plan,
validate the candidate, and stop the new VM after its initial boot.

The original-size preemptible request also failed allocation, with no VM
created. Retry the previously approved 48 GiB Cascade Lake V100 shape under
the same preemptible policy. This changes the CPU platform and host RAM,
preserving the 32 GB V100, eight CPUs, zone, disk, SSH inputs and stop limit.
It is the remaining smaller V100 variant within the requested GPU capacity.

## Larger GPU trials and current candidate

The operator requested other zones and larger VRAM alternatives, then reported
40 zone-a vCPUs. Test A100 gpu-standard-v3 with 28 vCPUs, 119 GiB RAM and
80 GB VRAM in zone a. Test Gen2 gpu-standard-v3i with 18 vCPUs, 144 GiB RAM
and 80 GB VRAM in zone d under the unused-resource replacement authorization.
The Gen2 attempt was rejected by quota before image compatibility could be
established; the existing Fedora image is legacy BIOS and Gen2 requires UEFI.

The A100 initially inherited the approved preemptible policy. After its
allocation failure, changing preemptible from true to false tested the same
A100 shape as a regular VM. Both failed allocation. Explicit legacy PCI V1
was then tested with both policies, replacing only the unused disk through
Terraform while its guard was temporarily removed. Both failed allocation
again. Restore all destruction guards and retain the last attempted A100
preemptible configuration in zone a with the actual legacy PCI V1 disk.
Restrict final zone validation to a/b, supported by this hardcoded platform.
Do not infer physical GPU stock from the generic allocation error.

Per the operator's explicit instruction, these allocation experiments used
patch then owning apply, with a prepared initial-stop helper. Publication
formatting, tests and semantic lint run once on the final candidate. No
host deployment, geographic-region provisioning or quota writes are included.
