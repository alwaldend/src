## Decisions

Proceed with a separate evals/ subproject containing configuration, Bazel composition and deployment instructions. All runtime code and cases belong to projects/comfyui_tools. The existing component owns its credentials; a reusable executable does not introduce a second infrastructure identity or its own Terraform state.

Use the public registry as the model/pin owner. Configure independent model/case workloads, fixed seeds, maximum parallelism, offer constraints, boot/job/download deadlines and a hard lifetime. Require verified Secure Cloud metadata as the existing rental does.

Ephemeral rental lifecycles use the existing Vast SDK, not Terraform resource state. Reuse the existing hardening template and authenticated CLI. The operator selects a standalone cleanup target plus an always-run final eval step. No separately supervised reaper is required. Document that forced process termination skips local finalization and leaves recovery to that cleanup target.

Source-only authority: the earlier no-rental validation instruction remains effective. No live provisioning, secret mutation or persistent service deployment is authorized in this change.

Matrix groups declare compatible hardware profiles, model profiles, workflow purposes and cases with shared seeds/settings. The evaluator expands the axes and pools all compatible workloads by hardware/model pair, so multiple workflows and cases reuse the same installation and model downloads. Each rental is destroyed after its own queued work is retrieved, without waiting for another rental.

The operator subsequently selects the standard >=32 GiB hardware profile and cheapest eligible offers, retaining mandatory datacenter/verification and image compatibility. The eval-only minimum_reliability override removes the inherited 99-percent floor. Personal rental defaults remain with their owning registry. The cancelled first allocation was destroyed; paginated v1 cleanup then confirmed the whole eval namespace empty before the rerun.

The operator explicitly supersedes isolated validation: use the current retained rental and restart or interrupt ComfyUI as needed. This authorizes the one-time in-place first-party node update and ComfyUI restart on instance 60240657; the checked-in startup source is the reconciled desired state. Inspect only the agent's test IDs and public case outputs, not user jobs or uploads. No additional rental is created and this rental remains running.

The operator specifically requests a packaged SSH target for this update. Reuse the host account SSH identity, require strict host-key checking, and pin the provider-reported instance endpoint in task-owned storage before use. No shared host authentication or key files are changed.

After SSH rejected the existing identity, the operator explicitly approves one-time attachment of that existing public key to instance 60240657. Use the instance-scoped Vast endpoint, not the account-wide key registry. Private key material remains on the operator host; the SSH target continues strict host-key verification. This temporary rental access is an approved imperative exception, recorded here; the startup/node code remains the reconciled source for future deployments.

The approved public key was attached only to the retained instance after the host SSH agent proved possession. The agent socket is inaccessible inside the sandbox, so the packaged SSH target uses the existing agent through the authorized host execution boundary; no key passphrase, new identity, account-wide registration or host configuration change is needed. The operator subsequently confirmed no jobs and explicitly requested restart. The immediately preceding aggregate queue count was zero. An atomic update installed first-party source SHA-256 `52766b648c4f02d89b2c63935533a6e8e965c5987993e0b2467c8833e5ac4ef5`; Supervisor restarted only ComfyUI successfully and the new Qwen node appeared in its catalog. Both saved GIF workflows were migrated and read back, preserving existing widget values while selecting previous-frame feedback. No user job content was inspected.

Both corrected Qwen tests completed and their GIFs were downloaded. The operator then explicitly requested destruction, superseding retention. The owning cleanup target succeeded, and independent account-wide v1 enumeration at 2026-11-11 19:48:06 UTC confirmed zero instances, including no rentals outside the evaluation namespace. The operator requests the existing NoobAI XL 1.1 profile on the existing nominal-48-GB RTX PRO 5000 profile for the next personal rental; only the main configuration and its documentation change. No new model, dependency, apply or rental is authorized or performed. Billing evidence and the operator's cost breakdown remain local.

The operator later explicitly requests renting the selected NoobAI configuration and receiving its link, superseding the preceding source-only limit. The owning tf.rent interactive plan added one instance with 100 GB disk, with no changes to the existing generated web login or Vault resources. After reviewing that exact plan, the apply created personal instance 60254657 from source candidate `104d1137a7d2df3fdb3efcf41cc1a521c3a64ea4`. A separately reviewed refresh-only apply recorded its connection outputs, with no remote resource changes. The provider reported RTX PRO 5000, nominal 48 GB, verification=verified and hosting_type=1 within existing offer caps. The instance uses the personal `src_infra_image_gen` label, independent of the cleared eval namespace.

Authenticated native CLI checks at `https://80.188.223.200:39551` succeeded and found all seven expected NoobAI workflows. This verifies service readiness and staging, not image quality or GPU inference compatibility; no image job was submitted. The existing Vault web entry supplies the unchanged login. The requested rental remains running for the operator. No shared host configuration, account-wide SSH registration or additional dependency was changed.

The operator subsequently requested the missing ordinary img2img prompt sequence
and its immediate upload. Add its generated workflow to the existing NoobAI
profile for future provisioning. The authenticated owning CLI saved only
`noobai-image-to-image-prompt-sequence.json` and its readback matched the generated
JSON exactly. Each prompt uses the same original reference; outputs are separate
PNGs. No service restart, job inspection, generation or rental replacement was
performed. Existing saved workflows and runtime settings were preserved.

## HiDream personal rental

The operator requested destruction of NoobAI instance 60254657; the reviewed
resource-targeted release destroyed only that instance, and a reviewed
refresh-only apply cleared its connection outputs while retaining credentials.
The subsequent explicit request selects the existing HiDream-I1 Full model and
48 GB RTX PRO 5000 hardware profile and authorizes renting it.

Include all HiDream-compatible native, style, sequence and GIF workflows.
Compose the missing ordinary img2img sequence with the existing helper: each
prompt starts from the same reference latent and keeps shared negative text,
progress and grouped filenames. Reject accidental previous-frame feedback,
fresh-noise substitution, missing encoder/model links or changed existing
workflow contents before renting. The SDXL IP-Adapter multi-reference graphs
are incompatible with HiDream and are not presented as working HiDream presets.
Keep existing model pins, runtime packages, credentials and host filters.

During the research hold, the eight-workflow inline payload exceeded Vast's
16384-byte encoded startup limit. Research into USO, UMO and joint-sequence
methods added no dependency or model asset. The operator subsequently
confirmed HiDream and explicitly directed deployment
followed by API workflow upload, deferring shared packaging improvements.
Keep all eight workflow names in the profile's canonical catalog and select an
empty startup_workflows list for this profile. Other profiles retain inline
staging by default. The authenticated owning CLI saves each catalog file.
This is an explicitly approved deployment procedure for the current rental;
the JSON definitions and catalog remain the IaC source for future uploads.
The configured-profile fixture verifies the actual inline selection and keeps
the API size guard unchanged. Shared delimiter/diagnostic experiments are
deferred as a task-local patch and are not part of this deployment.

All 34 offline checks and semantic lint passed. The configured startup is
12370 bytes; all seven pre-existing HiDream workflow documents remain unchanged,
and the ordinary img2img sequence adds the eighth file. The operator-approved
interactive rental plan added only instance 60332580, with no changes or
destructions. A reviewed refresh-only apply recorded its connection details;
the URL is https://80.188.223.200:39848. Existing authentication, mandatory
verified/datacenter host constraints and transfer-price limits remain in effect.
All eight authenticated workflow-save commands acknowledged success.

The operator then explicitly requested the URL without verification. No service
readiness, workflow readback, node inventory, model readiness or inference
checks were performed. Successful deployment and save acknowledgements do not
establish image quality or model compatibility. The rental remains running for
the operator.

## Optional native reference models after release

The operator subsequently requested destruction and source-only additions of
the other discussed models. The reviewed targeted plan destroyed HiDream
instance 60332580 only (0 added, 0 changed, 1 destroyed). A reviewed refresh-only
apply cleared its URL, instance ID and SSH output without changing resources.
Component credentials remain deployed. No rental apply or GPU acceptance is
authorized for these additions.

Add USO and UMO-UNO through their published native ComfyUI paths. Both share the
Comfy-Org FLUX.1-dev FP8 all-in-one checkpoint. USO also needs its LoRA,
projector and SigLIP vision encoder; UMO-UNO needs its own LoRA. Pin the public
publisher files in the existing registry. Preserve HiDream as the selected
model and keep the current hardware and mandatory verified/datacenter filters.
No new image, runtime package or custom-node pack is needed. The author's
StoryDiffusion implementation is a separate sampler integration rather than
a checkpoint; the operator subsequently explicitly includes it as well.

Before implementation, extend the existing configured-profile E2E's destination
contract for LoRA and model-patch files. Reject missing catalog outputs,
incorrect checkpoint/LoRA/projector destinations, loss of existing graph
contents and startup payloads above the unchanged 16 KiB limit. ComfyUI owns
model-patch loading; the reusable runtime maps that model kind to its native
directory. Offline packaging does not establish reference fidelity or GPU
compatibility.

The operator questioned the stale SeaArtLab dependency proposal; it has been
withdrawn without declaring the pack. Its latest commit is from July 2024,
and static inspection confirms an attention-forward signature mismatch with
current ComfyUI. The original StoryDiffusion repository last committed in
September 2024. The broader smthemex integration last committed in October
2025 and adds Diffusers/PEFT/Spaces/OmegaConf. Neither has been accepted as a
maintained path for the current runtime. The operator subsequently explicitly
directs that StoryDiffusion be omitted because the integrations are old.
No StoryDiffusion pack, profile, first-party adapter or runtime package has been
added. The source-only model additions are USO and UMO-UNO.

USO and UMO-UNO passed all 34 affected offline/repository checks. Their full
inline workflow catalogs contain six and four files respectively, with
startup payloads below the unchanged 16384-byte limit. The configured-profile
fixtures verify exact model destinations and
record zero live operations. Generated graph inspection confirms fresh output
latents, ordered subject references and separate USO style conditioning;
the seven earlier NoobAI and seven earlier HiDream baseline outputs remain
byte-identical. The current selected model, hardware, image and authentication
are unchanged. No model weights were downloaded and no instance was rented.

## Corpse evaluation and permanent transfer ceiling

The operator explicitly requests running both thorn groups with stronger art
direction and a shocking flesh-covered corpse, followed by complete cleanup.
After restoring the host login, the first controller finds no offers under the
shared USD 0.025/GB transfer ceiling and confirms the namespace empty. The
operator approves USD 0.04/GB for this run, then requests it permanently. The
shared module remains the sole policy owner; current documentation follows it.
Do not change the personal rental, hourly ceiling or required host predicates.

The initial approved-cap run retrieved twelve stills, with a NoobAI readiness
timeout and a stale Chroma offer. A retry limited to those two models retrieved
the remaining eight. All six created rentals across both rendering runs were
destroyed; final namespace enumeration succeeded after each. The retry exited
successfully at 2026-11-12 16:34:42 UTC. Detailed run identities and qualitative
findings belong to the ComfyUI tools change; actual media and receipts remain
in ignored local output. The personal rental was not accessed or modified.

The operator subsequently requests several controlled artistic runs, workflow
diagnosis, disabled cleanup and reuse of existing eval rentals. The reusable
command owns the new explicit retained-run validation and independent result
provenance. Use the current approved profiles and shared host/price limits;
install selected profiles once, then reuse those exact journaled installations
without fallback allocation, personal rental access or automatic cleanup.
Retained instances remain billable between runs and after handoff, as requested.

The operator subsequently restores credit and explicitly authorizes starting
only retained eval instances 60383020, 60383021 and 60383022. Record that
instance-scoped operational exception here; existing startup declarations and
the evaluator remain the desired-state owners. No personal rental, secret,
host configuration or replacement allocation is included in this operation.

## Final retained-experiment disposition

The operator subsequently requested destruction of everything and presentation
of saved results. This supersedes the retention instruction and cancels the
approved but unissued NoobAI/HiDream stop/restart repair. After interrupting the
active controller, the owning cleanup target confirmed the eval namespace
empty. An independent current v1 account listing at 2026-11-12 17:54:37 UTC
returned zero instances and no next page. No personal rental was created or
modified by these experiments.

Twenty additional Chroma stills were retrieved over five batches using the
same installation; the final native-mask batch was interrupted after two
retrievals. NoobAI and HiDream comparisons remained blocked by missing mapped
ports despite provider-running state. The project-owned change records prompt,
workflow and quality findings. No successful NoobAI/HiDream comparison or overall
artistic acceptance is claimed. Images, effective inputs, cleanup evidence and
account-specific cost data remain in ignored task output; no new dependency or
persistent host configuration was introduced.

## Requested NoobAI personal deployment

The operator requests deployment of NoobAI on the existing 48 GB hardware
profile and its connection link. Select the existing NoobAI XL 1.1 bundle,
retaining RTX PRO 5000, 100 GB disk, the USD 1/hour ceiling and mandatory
verified/datacenter and transfer-price constraints. Reuse the component's
existing web login and deploy through its reviewed interactive Terraform apply.
Verify service readiness and the eight saved workflows without submitting image
jobs. Retain this personal rental for the operator. No new dependency, host
configuration or evaluation rental is included.

At 2026-11-13 00:16:51 MSK, the reviewed owning apply created instance 60406482
with 100 GB disk: 1 added, 0 changed, 0 destroyed. Existing Vault login resources
were unchanged. All four offline Terraform/provider checks passed; the NoobAI
packaging artifact reports eight workflows, preserved model destinations and
13,990 startup bytes against the unchanged 16,384-byte limit. The first reviewed
refresh-only apply recorded loading status and SSH connection information,
with no mapped HTTPS port yet and no remote resource changes.

The subsequent reviewed refresh-only apply recorded running status and
`https://93.91.156.86:49649`, changing no remote resources. At 00:20 MSK,
the owning authenticated CLI successfully verified the trusted certificate,
ComfyUI 0.44.0 and NVIDIA RTX PRO 5000 Blackwell, then listed all eight expected
NoobAI workflows with no additional page. No image job was submitted; these
checks establish service readiness and workflow staging, not image quality
or successful model inference. The existing login remains in the component's
web entry. Instance 60406482 remains running for the operator.

## Requested NoobAI release

The operator subsequently requests destruction. At 2026-11-13 00:51:50 UTC,
the owning instance-targeted interactive release destroyed only instance
60406482: 0 added, 0 changed, 1 destroyed. The generated password, Vault web
entry, component identity and dedicated team key remain deployed. No replacement
or evaluation rental was created.

The first confirmation was rejected by automatic approval review because
Terraform's generic destruction warning appeared broader than the requested
scope. That attempt was cancelled without mutation. An independent read-only
text plan and the fresh interactive apply both confirmed the single instance
deletion before the successful confirmation.

At 00:52:33 UTC, a reviewed refresh-only apply cleared the instance ID, HTTPS
URL and SSH forwarding outputs to null, with no remote resource changes.
No rental remains managed by this root. The source retains the selected NoobAI
and 48 GB profiles for a future explicitly authorized deployment.

## Qwen keyframe deployment, 2026-11-15

The operator authorizes one 48 GB Qwen instance and a corrected plush-petting
trial. Select the existing Qwen-Image-Edit-2511 and RTX PRO 5000 profiles, keeping
the USD 1/hour ceiling and mandatory host/transfer constraints. Terraform owns
the personal rental and authenticated access; the existing Ansible role installs
all model files, Qwen workflows, both GIF node directories, the image-sequence
writer and the already approved Video2X/RIFE engine. The initial account SSH lookup failed with a missing-secret error. A scoped Vault apply subsequently confirmed
the bot identity already exists, while the account token lacks read access to
its exact metadata path. The operator explicitly approved deployment of the existing
Vault identity declaration and account SSH registration, restricted to those
resources. After the identity apply found no changes, the operator separately
approved its existing metadata-read policy and account AppRole attachment. The
reviewed Vault apply added that policy and attached it in place, with no deletes.
The broad account plan also proposed unrelated team destruction and
was not applied. Preserve any unrelated instance and inspect only owned
job IDs. Retain this personal rental at handoff under the operator's standing
retention request. Local validation must finish before paid allocation.

### Operator stop and release

The operator stops experimentation and requests destruction. All six owned
Qwen still-image jobs have completed and their original PNGs are retained
locally. No GIF assembly was submitted. The refined press preserves the bare
arm and scene better; both native-mask lift controls still fail to lift the
palm reliably. Do not report successful natural petting.

The reviewed instance-only Terraform destruction completed at 15:21:44 UTC on
2026-11-15: zero additions, zero changes and one destruction, instance60797191.
It retains the web credential and account key. The approved one-time attachment
of the current host public key to this instance ends with that destruction; the
canonical Vault metadata and account registration remain managed by IaC. The
latest workflow refinements remain local when the operator stops the task.

### Rough-pose retry

The operator subsequently requests another rental and trial, with rebase before
allocation. The prior PR is already merged. The task checkout synchronizes to
master 07cc629db7546975d4da619d2a4dfee13c85a8aa and restores its four task-owned
local edits without conflicts. Offline workflow, Ansible and Terraform checks,
repository quality and semantic checks pass before allocation.

The reviewed apply creates one Qwen rental, instance 60828898, and its declared
HTTPS credential resources: three additions, zero changes and zero deletions.
It retains the existing 48 GB RTX PRO 5000, 100 GB disk, price and certified-host
constraints. The GPU reports 48935 MiB. Existing account SSH access works on
this new instance without another key attachment or bootstrap exception.
The shared Ansible role installs all nine Qwen workflows, approved models,
first-party nodes and the existing interpolation tool: ok23, changed8, failed0.

Four owned generation jobs complete and all eight original PNGs are downloaded.
The rough-guide/partial-denoise path improves the lifted hand, but a separate
masked cleanup still leaves background/plush defects. No frame is accepted for
GIF assembly. Preserve the rejected trials and their lineage; do not report a
successful natural-petting animation.

The reviewed instance-only destruction completes at 20:20:38 UTC on 2026-11-15:
zero additions, zero changes and one destruction. A separate lookup of only
instance 60828898 confirms absence at 20:22:20 UTC. Refresh-only apply completes
at 20:22:33 UTC with no resource mutations and clears instance, HTTPS and SSH
connection outputs. The web credential and account key are retained. No other
rental or job inventory is inspected or changed.

## Native structural-reference trial, 2026-11-16

The operator requests one new rental and reference trial. Synchronize the
task-owned branch to master `3da98f0e0cd8708f3578ede4c63c1cf7fda12044` before
allocation. The existing Qwen Edit and 48 GB profiles retain their approved
model pins, 100 GB disk, USD 1/hour ceiling and mandatory host/transfer limits.
Offline repository, workflow, Ansible and Terraform checks and semantic lint
pass; a normal build verifies that the experiment runner actually exists.

The reviewed create plan adds only instance 60886376 at 05:50:13 UTC, retaining
the web credentials and account key. The GPU reports RTX PRO 5000 Blackwell
with 48935 MiB. Existing account SSH access works without another attachment or
bootstrap exception. Direct SSH is verified against keys also returned by the
provider proxy and is used for Ansible. The shared role installs all ten Qwen
workflows, existing approved models, first-party nodes and interpolation tool:
ok23, changed8, unreachable0, failed0, skipped7.

One native Canny preprocessing job and four bounded diffusion jobs complete.
All ten original output PNGs are downloaded before release. The
[workflow owner's trial record](../../../../../projects/comfyui_tools/openspec/changes/add-parallel-evaluations/design.md#native-structural-reference-trial)
owns settings, job lineage and the visual verdict: no accepted lifted frame and
no GIF. No separate ControlNet weight or new dependency is installed. Only
task-owned IDs are inspected; no other rental or job inventory is read.

The reviewed instance-only destruction completes at 06:25:42 UTC: zero additions,
zero changes and one destruction. A reviewed refresh-only apply at 06:26:30 UTC
clears instance, HTTPS and SSH connection outputs, with no resource mutations.
A separate lookup of instance 60886376 confirms absence. Component credentials
and account SSH registration remain deployed. No rental from this trial remains.

### Source-only sketch correction

The operator requests correction of the ambiguous sketch and experiment design.
Select `qwen-sketch-reconstruction` alongside the corrected
`qwen-structure-guide-edit` in the existing Ansible profile. The workflow owner
provides the prepared case input and matched baseline/appearance comparison.
This follow-up changes packaged source only. No rental, live deployment,
credential change or GPU acceptance is performed.

### Retained coherent-sketch trial, 2026-11-16

The operator subsequently requests balance-lookup documentation, another rental
and continued testing, explicitly prohibiting destruction afterward. This
authorizes one new instance with the existing Qwen/48 GB profile and shared
Ansible installation. Retain the instance after testing and provide its link;
the previous disposable-trial cleanup instructions do not apply. Existing
credentials and account SSH registration remain unchanged. A filtered live team
lookup confirms available credit before allocation; amounts remain in private
task evidence. Record live allocation, installation, observed cost and final retention below; do not treat this request
as evidence that a rental already exists.

The reviewed apply creates only instance 60894930 at 07:27:36 UTC: one addition,
zero changes and zero destructions. The provider reports RTX PRO 5000, 48935 MiB,
100 GB disk, verified Secure Cloud hosting and USD 0.905852/hour, with separate
transfer charges. Direct SSH and provider proxy return the same Ed25519 host
key; existing account access needs no bootstrap exception. The reviewed
refresh-only apply at 07:29:37 UTC records connection outputs without changing
remote resources. The endpoint is `https://80.188.223.200:39631`; the generated
login remains at the component's existing Vault web entry.

The shared Ansible role completes with ok=23, changed=8, unreachable=0 and
failed=0. It installs all eleven selected Qwen workflows, existing pinned models
and approved first-party nodes, then passes its native interface checks.
Authenticated ComfyUI reports version 0.44.0 and RTX PRO 5000 Blackwell. Both
case images upload successfully under their canonical filenames. No credential
or host-configuration exception was needed.

The initial owned comparison completes four diffusion jobs and retrieves all six PNGs. Two
matched prompt pairs compare sketch-only against original-photo appearance
conditioning. The default sketch reconstructs the gap; the appearance variant
restores identity but adds raised hair folds. A stronger crown instruction is
rejected in both variants. The workflow owner's record retains detailed visual
acceptance and lineage; no GIF is accepted. The user sees the original raw
comparisons as they finish. Instance 60894930 is retained under the explicit
instruction; no release operation is executed. The exact balance command is
verified with filtered output and links to its account-owned documentation.
The provider lookup at 07:53:13 UTC confirms that this exact instance remains
running at the recorded hourly price. The contemporaneous credit amount remains
in private task evidence and the operator conversation.

The operator then requests a targeted negative prompt on the better default
variant. Two additional owned jobs run on this same retained instance, preserving
the paired sketch-only and original-photo comparison. Both complete and three
more original PNGs are retrieved, bringing this trial to six diffusion jobs and
nine PNGs. The original-photo folds persist, so the negative is not promoted to
the workflow defaults. No replacement rental, dependency or release is involved.

Subsequent operator-approved ear-removal, sketch-animation and photographic
motion tests reuse this same endpoint. Their exact jobs, input lineage and
visual decisions belong to the [workflow owner's motion record](../../../../../projects/comfyui_tools/openspec/changes/add-parallel-evaluations/design.md#approved-sketch-motion).
No new provisioning, model installation, SSH attachment or rental release is
part of those tests. The explicit retention instruction remains in force.

### Retained visible-motion trial, 2026-11-16

The operator subsequently requests destruction of instance 60894930. The
reviewed instance-scoped Terraform apply destroys exactly that resource at
10:36:08 UTC: zero additions, zero changes and one destruction. The exact-instance
provider lookup confirms absence at 10:36:41 UTC. The reviewed refresh-only apply
clears connection outputs at 10:37:24 UTC without changing remote resources.
Component credentials and account SSH registration remain unchanged; downloaded
results are preserved.

After rejecting the small movement in the photographic-guide Qwen sequence,
the operator requests another rental and trial. Prepare the larger guides and
finite endpoint workflows before allocation. The existing Qwen/48 GB profile
creates only instance 60921795 at 12:06:13 UTC: one addition, zero changes and
zero destructions. The provider reports RTX PRO 5000, 48935 MiB and
USD 0.905852/hour, with separate transfer charges. Direct SSH and provider proxy
return the same Ed25519 host key; existing account access requires no bootstrap
exception. The reviewed refresh-only apply records connection outputs at
12:07:58 UTC without changing remote resources. The endpoint is
`https://80.188.223.200:39881`; login remains in the existing component Vault entry.

The checked-in Ansible role completes with ok=23, changed=8, unreachable=0 and
failed=0. Authenticated checks identify ComfyUI 0.44.0, frontend 1.57.6 and the
RTX PRO 5000 Blackwell. Existing pinned Qwen models, workflows and RIFE support
are reused; no dependency, credential or host-configuration exception is added.

Eight owned diffusion jobs and one assembly job finish. Five inspected Qwen
poses form the delivered loop; three defective outputs are excluded. The
[workflow record](../../../../../projects/comfyui_tools/openspec/changes/add-parallel-evaluations/design.md#photographic-pose-guides-and-visible-motion)
owns input preparation, settings, selection and visual limitations. Authenticated
assembly completes at 12:43:50 UTC; all outputs are retrieved. No release is
performed for this rental, and the operator's retention instruction remains in
force. Balance observations stay in private task evidence and the conversation.

During PR completion, a fresh exact-instance lookup fails while obtaining its
Vault credential, before contacting the provider. This does not establish a
change in rental state. The last successful runtime observation above remains
the deployment evidence; no infrastructure or credential mutation is attempted
as part of source delivery.
