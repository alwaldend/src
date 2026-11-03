## Source extraction acceptance

During the initial extraction, the three Terraform resource/provider/variable files were compared byte for
byte against published commit `34a6ba90d5cf51c438e3c72e58a2fb592e479d40`.
All matched. Resource addresses, required image/network inputs, shared provider
pin and `/tf_backend/tf_setup` path are preserved.

The new stage owns its AL configuration. Generated normal plan/apply wrappers
built successfully and selected only the pinned Yandex provider. Their merged
configuration retains the existing component AppRole, Yandex credential path
and backend path. Backend-disabled init and Terraform validate succeeded with
no AL plugins or live credentials.

The independent candidate passed 34 repository/owner checks and semantic lint.
The host-only source passed 33 checks and semantic lint; its Ansible package
built successfully. Host Ansible files were unchanged, while its AL wiring now
contains only component authentication and signed SSH injection. Source records
and README boundaries were updated for the split.

No live Terraform plan/apply, DNS/state operation, account provisioning, host
connection or deployment ran. Final prepared candidates receive the mandatory
checks before publication; Git and delivery receipts own exact published heads.

## Deployment continuation

The user then requested applying Terraform after extraction. After being told
Yandex starts newly created VMs, the user explicitly selected one initial boot
followed by stopping the instance. The pinned provider implementation at
`e058dfe2513062914b2360b1c39036fba8a2d14b` calls Instance.Create and exposes
no create-stopped/desired-status input. No VM has been created by this task.

Provisioning preflight stopped on the component Yandex credential fetch.
A read using the component AppRole confirmed HTTP 404 at its declared
account_iam_key path. A personal-token metadata read returned 403 and was not
used as evidence of absence. Neither required image nor administrator-network
input is present in the agent environment; both were requested from the user.

A targeted plan from the existing infra/yandex_cloud/org1/tf owner was saved
and reviewed for only module.approle_folders["src_infra_local_inference"].
It proposes nine additions: a folder, service account, folder-scoped IAM access,
two account keys and four Vault entries. It changes/destroys nothing and
contains no VM actions. Approval for this distinct prerequisite scope was
requested before any mutation. The operator subsequently approved this exact
scope. Applying the saved plan succeeded on 2026-11-03 with nine additions,
zero changes and zero destructions. A fresh plan for the same module reported
no changes. The component AppRole can now read all three required credential
fields; no secret values were displayed. No VM actions occurred in this apply.

The operator requested the existing Fedora image and public SSH. Repository
owners infra/download/yandex and infra/threexui/tf_setup already use immutable
Fedora 43 image fd8iku26nnkveh8s4di1. A live read with the component credentials
confirmed name fedora-43-v20260717, family fedora-43, status READY and legacy
PCI topology v2. It does not advertise an NVIDIA driver. The provisioning
source now reuses the shared Fedora cloud-init document and permits explicit
IPv4 /0; the selected SSH input is 0.0.0.0/0. Host PR #233 still requires Ubuntu
and a preinstalled driver, so its playbook must not be deployed to this image.

Host-only PR #233 was prepared and checked separately. Publication hit one
GitHub TLS handshake timeout; a read-only verification and fresh inspection
confirmed no push occurred and the old remote head was intact. The checked
candidate was then published with its original exact remote lease.

## VM apply result: compute quota exceeded

The Fedora/public-SSH candidate passed 34 checks, semantic lint and offline
Terraform validation. Its saved live plan proposed exactly seven additions,
zero changes and zero destructions: network, subnet, security group, child
DNS zone, KMS key, encrypted 200 GiB SSD and one Broadwell V100 VM with eight
vCPUs and 96 GiB RAM. The prepared initial-stop helper authenticated successfully
and verified that the dedicated component folder contained no instances.

The authorized saved-plan apply on 2026-11-03 created the six supporting
resources but failed to create yandex_compute_instance.host. Yandex returned
ResourceExhausted: the maximum memory quota was exceeded in ru-central1-a,
with one additional quota violation that the error did not identify. The
component's cloud-level QuotaLimit.List read returned HTTP 403; this does not
establish the second quota's identity or value. No quota or access policy
changes were attempted.

The helper's post-apply control-plane read confirmed that no VM was created,
so no VM booted and no stop request was required. A fresh plan after the
failure confirms the six managed resources and proposes only the remaining
instance creation: one addition, zero changes and zero destructions. Creation
was not retried against unchanged quotas. The encrypted disk remains allocated;
the partial deployment is retained in its owning Terraform state.

Continuation requires cloud quota headroom for the configured instance,
including 96 GiB RAM in ru-central1-a, and resolution of the additional
violation. Generate and review a new plan after quotas change, then reuse the
authorized initial-boot-and-stop procedure. Host deployment remains excluded.
Task-private raw plans and logs are removed after retaining sanitized results;
future continuation must generate fresh plans rather than reuse this attempt.

## Session ergonomics

The bounded source record now names both delivered PRs, the prerequisite
result, exact immutable image, public SSH input and remaining live blocker.
Two observed discovery issues were corrected locally: the earlier Ubuntu
assumption was replaced after checking the two existing Fedora consumers,
and oversized combined tool output was replaced by bounded field summaries.
These are live/configured observations, not evidence that the GPU runtime
works. Required checks were repeated only after candidate changes; shared
policy, catalogs and runtime configuration were not modified.

## Authorized availability-zone change

The operator requested another availability zone after the quota failure and
explicitly authorized replacement of the unused component resources. The
reviewed saved plan for ru-central1-b replaces only the existing encrypted
disk and subnet, then creates the VM: three additions, zero changes, two
destructions. It preserves the four non-zonal resources, the component folder
and account, Fedora image, public SSH input and initial-boot-and-stop limit.
The component folder still contains no instances. The disk guard is lifted
temporarily for this exact replacement and will be restored after the attempt.

The apply used prepared candidate 78ea614a406cd0297e227444b09b290fcf817f3a.
It deleted the unused disk and subnet in ru-central1-a and created the subnet
in ru-central1-b. The disk create request then failed with ResourceExhausted:
Resource allocation is restricted. The helper's post-apply read confirmed
that no VM exists. The disk guard was restored immediately after this result.
A fresh saved plan with that guard restored confirms five resources in state:
the network, ru-central1-b subnet, security group, KMS key and child DNS zone.
It proposes only the missing disk and VM: two additions, zero changes and
zero destructions. No boot disk or VM remains allocated.

Stable blocker: compute-allocation-restricted-ru-central1-b. The provider
error is live evidence. The [published incident listing](https://status.yandex.cloud/en/timeline?installation=ru&service=Compute+Cloud)
reports ru-central1-b unavailable and restrictions on resource creation;
its latest indexed update is October 10, so it is supporting context rather
than an independent current API availability measurement. The
[GPU documentation](https://yandex.cloud/en/docs/compute/concepts/gpus),
updated September 28, lists this V100 platform only in zones a/b. Direct
incident-page fetching required JavaScript and did not establish a newer
status. A post-failure zone check cannot justify the earlier selection:
current incident information should have been consulted before moving the
resources. Further unchanged allocation attempts were stopped.

The operator was asked whether to resolve quota in a, leave b staged, or
evaluate another GPU in another zone. Until that choice changes the scope,
the requested ru-central1-b configuration remains staged with destruction
guards restored. The initial-boot-and-immediate-stop authorization persists.
Private plans and raw apply logs were removed after retaining these results
and sanitized address/action summaries.

## Operator-selected smaller V100 in zone a

The operator then selected a return to ru-central1-a and a smaller VM.
The documented gpu-standard-v2 configuration retains one 32 GB V100 and eight
vCPUs, while reducing host RAM from 96 GiB to 48 GiB. The reviewed saved plan
replaces only the unused zone-b subnet and creates the missing encrypted
200 GiB disk and smaller VM in zone a: three additions, zero changes, one
destruction. All three destruction guards are present; no guard exception
is required for this attempt because no disk or VM exists. The stop helper
now verifies gpu-standard-v2 and ru-central1-a for the new VM identity.

The smaller attempt used prepared candidate
adb48b007ecb2ed352d53a4019261d29e2e353be after 34 checks and semantic lint passed.
It replaced the unused subnet and created the encrypted 200 GiB disk in
ru-central1-a, but instance creation returned the same ResourceExhausted
memory-quota error and one further unidentified quota violation. The post-apply
control-plane read confirmed no instance was created, so no VM boot or stop
occurred. The refreshed plan retains six supporting resources and proposes
only the missing gpu-standard-v2 VM: one addition, zero changes, zero
destructions. All storage and instance guards remain present.

Stable blocker: compute-memory-quota-ru-central1-a-48gib. The documented
single-V100 configurations use at least 48 GiB RAM, so further reducing this
V100's memory is unsupported. The earlier cloud-quota read returned 403 and
the checked-in folder module grants folder-scoped admin only; no cloud-wide
access was inferred or added. Continue only after quota headroom or the
operator's hardware selection changes, using a fresh reviewed plan and the
already authorized initial stop. Private raw plans and logs are removed after
retaining sanitized results.

## Quota-increased allocation attempts and current alternatives

A subsequent read using the existing src_infra_yandex_cloud_org1 identity
verified that the management and component cloud IDs match the component
folder's actual parent. The published quota snapshot had enough RAM, CPU,
GPU and dynamic-address headroom for the 48 GiB VM. The component identity's
403 did not preclude this management read, and the second earlier quota
violation remains unidentified. The earlier assertion that increasing quota
was required was too categorical; these are later observations, not proof
of the limits at the earlier failures.

The operator then reported increased quota and selected a fresh 48 GiB retry.
The unchanged gpu-standard-v2 source plan contained only one VM addition.
Allocation operation fhmvss42r60r4mrqmcp7, created 2026-11-03T13:26:03Z,
failed with error code 8 and Not enough resources, with no error details.
The helper's control-plane read confirmed that no VM exists.

After catalog research and another operator-reported RAM quota increase,
the operator explicitly restored gpu-standard-v1 with eight vCPUs, 96 GiB
RAM and one 32 GB V100 in ru-central1-a. Candidate
04ed8fe1f90a767e0ad4ac61d447d300f3d2a1ae passed 34 checks and semantic lint.
The saved plan proposed exactly one addition, no changes and no destructions,
preserving all six supporting resources. Allocation operation
fhmk1pqi4vth1k6f4oj6, created 2026-11-03T13:53:12Z, failed with the same
error code 8 and Not enough resources, with no error details. No VM exists
and no stop was needed. Neither result is reported as a successful VM boot
or evidence of driver/runtime readiness.

Stable blocker: compute-v100-allocation-ru-central1-a. The live zone API
reports a/b/d/e UP and c DOWN, but the cloud's published b CPU/RAM/SSD quotas
remain zero. The authoritative raw platform table uses no.svg for both V100
platforms in d/e, confirming that blank rendered table cells were not evidence
of support. Gen2 is supported in d, with 80 GB VRAM and a minimum 18 vCPUs and
144 GiB RAM; its published GPU quota is zero. The existing Fedora image is
legacy BIOS, while Gen2 requires UEFI, so it is not an immediate image reuse.

Public RUB STREET_PRICE SKU rates from the Compute Cloud - GPU billing
service yield about 80.2456 RUB/hour for a preemptible Broadwell V100 with
eight CPUs and 96 GiB RAM, and 815.0824 RUB/hour for a regular single-GPU
Gen2 VM with 18 CPUs and 144 GiB RAM. These are running-compute sums only,
excluding disk, IP and account-specific terms. Preemptible VMs can be stopped
at any time and run for at most 24 hours per start. Availability is unverified
for both alternatives. The operator was asked to choose interruption behavior,
larger hardware preparation, or waiting for regular capacity before another
configuration change. No automatic retries or larger GPU deployment ran.

## Approved preemptible attempts

After the interruption behavior was explained, the operator explicitly
requested trying preemptible capacity. Candidate
09a7855fbdd441c0f59727757ff45be95fcd6087 used the 96 GiB Broadwell shape with
preemptible true; candidate dca39862790007c5035de705e29dd07276a4083b used the
previously approved 48 GiB Cascade Lake shape under that same policy. Each
passed 34 checks and semantic lint and had an exact one-VM-addition plan,
with no changes, destructions or guard removal.

The 96 GiB operation fhmp4c79geskjtismueq and 48 GiB operation
fhmljpjc6622s8uo7jlk both failed with ResourceExhausted: Not enough resources.
Each helper's post-apply control-plane read confirmed no VM exists, so no
stop was required. The GPU driver and host services remain unverified.

The operator asked about another geographic region. The
[Kazakhstan region-specific platform note](https://yandex.cloud/ru-kz/docs/compute/security/)
permits standard-v3 and standard-v3-t4i only; T4i has 24 GB VRAM. The generic
GPU table also appears in that documentation locale but does not override
the region-specific restriction. No Kazakhstan provisioning ran.

The operator requested faster live iteration. Subsequent authorized trials
should patch the relevant Terraform input and use the owning apply workflow
promptly, with required publication gates afterward, rather than repeating
all publication checks before every cloud-allocation experiment. Private raw
plans and apply logs are removed after retaining sanitized outcomes.

## Larger GPU trials on 2026-11-03

The operator requested other zones, other 32 GB configurations and larger
VRAM variants, and instructed patch then owning apply for allocation trials.
Every attempt used the component's injected credentials and prepared stop
helper; the helper checked that the dedicated folder had no instances before
apply and read it again after failure. No attempt created or booted a VM.

| Platform and zone         | CPU / RAM / VRAM     | Policy and disk topology                                   | Result                                                                                    |
| ------------------------- | -------------------- | ---------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| A100, gpu-standard-v3, a  | 28 / 119 GiB / 80 GB | Preemptible, inherited PCI V2                              | CPU quota rejection; request 64807485-b639-4c84-bad8-853124c8f99d                         |
| Gen2, gpu-standard-v3i, d | 18 / 144 GiB / 80 GB | Preemptible, inherited PCI V2                              | Memory quota plus one unspecified violation; request 0c371d7a-2df3-4702-b37c-a7893c5f96ae |
| A100, gpu-standard-v3, a  | 28 / 119 GiB / 80 GB | Preemptible, inherited PCI V2, after reported CPU increase | Not enough resources; fhmfpvc24icv7r6qg97g                                                |
| A100, gpu-standard-v3, a  | 28 / 119 GiB / 80 GB | Regular, inherited PCI V2                                  | Not enough resources; fhm9pfmd7ihnmo2lvf69                                                |
| A100, gpu-standard-v3, a  | 28 / 119 GiB / 80 GB | Regular, explicit legacy PCI V1                            | Not enough resources; fhmjqqsve5h27nqcg4d5                                                |
| A100, gpu-standard-v3, a  | 28 / 119 GiB / 80 GB | Preemptible, explicit legacy PCI V1                        | Not enough resources; fhmo6g5gtf9hoe2sgibt                                                |

The source retains the last row. "Regular A100" meant changing
scheduling_policy.preemptible from true to false on the same A100 shape;
the original V100 attempts had also been regular. The last trial restored
true, the previously approved interruption policy. The PCI experiment is a
separate disk hardware-generation change, not a consequence of scheduling.
It did not resolve allocation failure. No conclusion about physical GPU stock
or NVIDIA driver compatibility follows from these errors.

The unused subnet and encrypted disk moved to d for the Gen2 attempt and
back to a for the A100 retry, using Terraform. The disk was subsequently
recreated with explicit legacy PCI V1. All three prevent_destroy guards are
restored. The final zone input permits only a/b for the selected A100 platform.
The Gen2 attempt failed admission before the legacy Fedora image's UEFI
compatibility could be tested. The observed Gen2 GPU quota was zero; the API
did not identify it as the unspecified violation, so that remains unconfirmed.

Stable blocker: compute-a100-allocation-ru-central1-a. Six supporting
resources remain; tasks requiring successful creation and STOPPED verification
remain incomplete. No host Ansible or model services ran.

### Session ergonomics

Live evidence: four distinct A100 allocation operations failed after admission;
changing scheduling policy and PCI topology did not resolve allocation. The
operator twice requested direct apply iteration after excessive preparatory
checks. Stable issue allocation-trial-publication-overhead: reuse unchanged
validation evidence during explicitly authorized trials and run required
publication gates on the final candidate. Record the exact shape, policy,
zone and outcome at each attempt so a resumed agent need not reconstruct the
sequence from raw logs. No shared workflow or policy change is included here.

Final read-only plan: six prior resources and no instance, with exactly one
A100 VM addition, zero changes and zero destructions. The planned shape is
gpu-standard-v3, 28 vCPUs, 119 GiB RAM, preemptible true in ru-central1-a.
All three destruction guards are present. Private plan, rendered plan JSON
and raw trial logs were removed after retaining these sanitized outcomes.
