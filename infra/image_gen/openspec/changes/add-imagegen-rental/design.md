## Base checkpoint replacement

The user selected FLUX.2 klein Base 9B and requested redeployment, then explicitly
requested destroying the current rental. Instance 60069855 was destroyed through
its owning targeted interactive apply on 2026-11-10 at 14:05 UTC: 0 added,
0 changed, 1 destroyed. A reviewed refresh-only apply at 14:08 UTC cleared all
connection outputs with zero remote-resource changes. The component AppRole,
team API key, generated web login and model token remain managed.

The flux profile selects the publisher's Base FP8 checkpoint at revision
9ecf2143d71542449960c5584340269c6d401449. Encoder and small-decoder VAE pins are
unchanged. The three native Flux graphs use 50 Euler steps, CFG 4, and an empty
CLIPTextEncode negative prompt instead of zeroing the positive embedding.
Reference latents remain on both branches. These defaults follow the publisher's
Base example; the current upstream Comfy template also confirms the encoded
empty-negative structure and the same model bundle. No runtime package, image
or custom-node dependency changes are required.

The existing Vault token received HTTP 403 from the pinned Base file before any
replacement capacity was created. Account license acceptance or token scope
was resolved by the operator. A subsequent in-memory HEAD returned 302 with the pinned commit and 9,567,278,472-byte checkpoint. The requested replacement
retains the existing USD 1/hour ceiling, Secure Cloud, VRAM floor and exposure.
The operator permits neutral workflow checks on the fresh replacement. Once operator use starts, further inspection or generation requires renewed approval. Earlier distilled workflow results do not validate Base quality.

The user also requested a 96 GB option as source wiring only. The 96gb hardware profile in config.json selects at least 95 GiB of
reported VRAM (covering nominal 96 GB cards), 150 GB disk and a USD 3/hour ceiling.
The normal defaults remain 32 GiB, 100 GB and USD 1/hour. Model selection is
independent. This preset is not selected for the Base replacement and no
high-memory deployment is authorized.

The reviewed Base apply created instance 60080494 on 2026-11-10 at 14:29 UTC:
1 added, 0 changed, 0 destroyed. It selected the same RTX PRO 5000 class,
100 GB disk and existing USD 1/hour ceiling. The operator then requested destroying
the current instance. A reviewed targeted release at 14:31 UTC destroyed it:
0 added, 0 changed, 1 destroyed. A refresh-only apply at 14:32 UTC cleared
connection outputs with zero remote-resource changes. No rental remains running.
No generation job was submitted and no Base runtime or quality acceptance is
claimed. The 96 GB preset was never selected or deployed. Component credentials
remain managed for later use.

## Previous follow-up: profiles and CLI

The user requested retaining Qwen, adding FLUX.2 klein 9B and renting one
instance, then explicitly approved adding the official ComfyUI CLI and a target
that authenticates to the rental. Keep existing identity, credentials, host
requirements and price limits. Flux weights require publisher license acceptance
and authorized Hugging Face download access before any paid rental is created.

The infrastructure target injects its web login and backend context into a short-lived
first-party bridge, which obtains the endpoint through its owning Terraform output
target and removes injected credentials from the upstream CLI child. The tools owner packages the CLI and bridge. Source and tests remain
neutral. No private uploads or prompts are used for acceptance.

## Prompt sequence follow-up

The user requested one node accepting several prompts and producing several
images in one queued workflow. Use a first-party, dependency-free ComfyUI list
node: a multiline prompt list plus shared reference instructions produce paired
prompt and numbered filename-prefix lists. Native downstream list mapping
executes one image per prompt with the same character/style references. This
does not provide previous-panel memory. Stage a sequence workflow for both
model profiles and whitelist only this first-party node directory.

Before implementing, check empty input, blank-line/line splitting, stable input
order, one output per nonempty prompt, paired filename ordering, unchanged shared
references and empty/invalid filename prefixes. Errors must not echo prompts.
The shared module must validate directory names and cannot enable arbitrary
unapproved packs. Verify actual list mapping during neutral live acceptance.

## Current design and deployment status

The current owner is `infra/image_gen`, authenticated as `src_infra_image_gen`.
Its KV namespace contains `vastai`, `web` and `tf_backend/setup`. Vault Terraform
owns the role, operator membership and policies. The existing Vast team root
owns issuance and exact-entry delivery; the component owns generated web login
and rental. The rental receives no Vault or Vast API token. The component
operator can browse its namespace through existing Vault identity membership.
The web and delivered key resources delete all their KV versions on destroy;
their exact metadata-delete permissions do not extend to component state.

The replacement uses fresh state and generated fresh credentials after the
user authorized applying IaC without renting capacity. No migration, import or
reuse of the personal backend was required. The AppRole, policies and credentials
are deployed. The user subsequently authorized one Qwen rental and workflow
acceptance, recorded below. Historical personal-rental receipts describe
retired resources; owning Terraform outputs provide current connection details.

## Legacy teardown receipt, 2026-11-10

The user requested destruction of the instance and then all dedicated resources,
followed by fresh infrastructure IaC with its own AppRole. Owning interactive
applies were reviewed before approval:

| UTC   | Owner                   | Reviewed and completed scope                                                                        |
| ----- | ----------------------- | --------------------------------------------------------------------------------------------------- |
| 09:18 | Former user rental root | Destroy instance 60043006; 0 additions, 0 changes, 1 destroy                                        |
| 09:22 | Former user rental root | Destroy generated password and current web secret; 0 additions, 0 changes, 2 destroys               |
| 09:24 | Vast alwaldend root     | Revoke key 36816419 and delete its current delivery entry; 0 additions, 0 changes, 2 destroys       |
| 09:29 | Vault root              | Remove legacy writer policy from Vast role and delete that policy; 0 additions, 1 change, 1 destroy |

The old rental root now manages no resources. The shared team, its operator seed,
unrelated keys and personal identity remain outside this scope. Legacy KV
resources used soft deletion; protected historical Vault/state versions may
remain. No raw state or secret was printed, saved or moved. The previous
instance URL is inactive. Qwen was not deployed or live-tested before teardown.

## Historical implementation and deployment evidence

The following sections retain earlier decisions and receipts. Personal identity
and credential-preservation statements below describe the retired stack only.

## Context

A personal image generation rental needs generated credentials retained in Vault before paid instance creation.

## Current follow-up

The user destroyed the previous rental and requested a source refactor and PR
publication. Use the ComfyUI wrapper with the full Chroma1-HD-Flash checkpoint at
publisher revision 093c4f63b60507234aa53a073f1f9f565b3f4336, retaining the existing
T5 encoder and VAE. The similarly sized delta file is not a standalone model.
Keep random_password.web and vault_kv_secret_v2.web at their existing addresses;
no new role, token or secret copy is needed. There is no remaining live rental
to migrate. Publish the diagram with the source and update it to this runtime.
This follow-up does not authorize provisioning paid capacity.

## Goals / Non-Goals

Add imagegen/tf_setup with random_password, Vault KV storage, and the reusable exposed-mode module. Reuse the existing user_simeonwarren AppRole, regular user secret namespace, and independent Vault-backed state. Declare a direct published Vast ComfyUI/CUDA 12.9 image tag without a template by default; provide optional declarative public model downloads; allow an optional reviewed environment-only template; retain the USD 1/hour cap. Offline validation SHALL NOT rent capacity; the user separately authorized one live rental after security review.

## Decisions

Reuse the existing provider, pinned Terraform providers, and AL/Vault configuration flow.
Secrets remain sensitive in Terraform and reside in the Vault-backed backend; they are never public outputs.
The caller owns template trust and deployment authority. The 32 GB requirement is interpreted as 32 GiB VRAM.

## Risks / Trade-offs

Offline fixtures verify request and lifecycle behavior, not actual image startup or marketplace availability.
HTTPS uses the Vast-provided instance certificate; clients must verify its trust before submitting credentials.

## Migration Plan

Default SSH behavior is retained; changed startup settings intentionally replace a rental.
Create a dedicated alwaldend API token through the existing alwaldend root and store it in the regular user namespace before renting.

## Open Questions

None for implementation. A live acceptance run must inspect the then-selected image and actual service bindings.

## Default checkpoint

The user explicitly approved RunDiffusion/Juggernaut-XL-v9 as the rental default. Set model_downloads to the publisher's standalone safetensors checkpoint at revision cf419233522daa0b9ea36c3aff98fa2cab1fb0fb. This keeps the reusable template's model list generic and empty. The native provisioner downloads the checkpoint before Forge starts; no model-serving package or authenticated downloader is added. The user may replace the list or supply an empty list.

The publisher describes this as its most popular SDXL checkpoint; this is not an independently verified global ranking of uncensored models. The 7,105,348,188-byte file fits the rental's 100 GB disk and its SDXL architecture fits Forge and the 32 GiB VRAM floor. License terms still apply. Failed public downloads or upstream availability can prevent first boot; a fixed source revision prevents silently selecting another release. Verification covers public URL availability and the existing manifest/lifecycle fixtures, not a live GPU generation.

## Direct image and live deployment

After PR #283 merged, the user requested removal of dynamic image discovery and explicitly requested renting an instance and receiving its URL. The user also required a security/privacy review before deployment. The rental declares the current published Neo/CUDA 12.9 tag directly; Docker Hub lists no moving latest alias. Remove the now-unused external provider, executable, and fixture. Keep all existing GPU, price, hosting-tier, TLS/authentication, model, and Vault requirements. Review the fresh interactive plan before applying. Live endpoint acceptance must not send credentials until TLS is verified.

## Predeployment security verdict

Proceed for personal image generation under the accepted host trust boundary. The user explicitly accepted that Vast and the machine administrator can access instance storage and memory. This is not confidential computing and the stack has no verified no-retention/no-training guarantee. Only account-key SSH and authenticated Forge HTTPS are published; management tools remain local, outbound tunnel/synchronization services are disabled, and the Vast API key is not injected into the container. The generated web credential is retained in Vault/state and necessarily in the instance startup configuration. Images and prompt metadata can remain on instance disk until destruction; host copies cannot be ruled out.

On 2026-11-09, the selected linux/amd64 image manifest was sha256:cb9a046afc1f1c46560c18f68f26856cf0c44f8978a777228c3c6e8199bbabda. The four previously inspected boot, proxy, and Supervisor layers remain unchanged. The changed small Forge layers were downloaded with verified digests; the actual Forge service and launch script are identical to the inspected versions. These checks establish the reviewed exposure contract, not a complete audit of all dependencies. Before sending the real web credential, verify the observed TLS identity; probe unauthenticated and authenticated access and public bindings after boot.

## Authorized live deployment

On 2026-11-09, the owning interactive rental apply displayed 3 additions (random web credential, its regular-user Vault entry, and one constrained GPU instance), 0 changes, and 0 destroys. After reviewing that fresh plan, the authorized apply completed with exactly those counts. Instance 59945579 reports hosting_type=1, verification=verified, one CMP 170HX with 40960 MiB VRAM, 100 GB disk, and API-reported dph_total 0.47703703703703704 USD/hour. The exact configured image is reported by the instance. The username/password entry now exists at imagegen/web; no credential was printed or saved in task scratch.

After initial image loading, the actual port map contains only 22/tcp and 7860/tcp. Live HTTPS validated against the official Vast root CA with hostname verification enabled. Unauthenticated and incorrect-password requests returned 401, the correct Vault login returned 200, and Forge UI configuration contained the configured Juggernaut XL v9 checkpoint. A separate probe encountered a boot-time TLS handshake timeout; repeat verification after boot passed all TLS/authentication/model checks. No TLS verification was bypassed. The owning reviewed refresh-only workflow records connection outputs without changing remote resources. Connection coordinates stay in Vault-backed Terraform outputs and sanitized ignored task evidence rather than maintained source.

## Chroma replacement authorization

The user requested redeployment to Chroma1 and explicitly confirmed that no files
need preservation and the instance may be replaced. The bundle uses Chroma1-HD,
the publisher-recommended T5 XXL FP16 encoder and Flux VAE at fixed public
revisions. Native Forge settings select these components with its flux preset,
40 steps and CFG 3. Keep the existing Vault login, Secure Cloud constraints,
32 GiB floor, 100 GB disk and USD 1/hour ceiling. Changing the first-boot manifest
replaces the current instance; no imperative installation or state edit is needed.

The reviewed interactive apply of candidate f4074515 on 2026-11-09 replaced
instance 59945579 with 59961415: one addition, no changes, one destruction.
The existing random password and Vault web secret were unchanged. The replacement
reports a verified Secure Cloud host (hosting_type=1), one CMP 170HX with
40 GiB VRAM, 100 GB disk and USD 0.477037/hour. Model provisioning and application
readiness are verified separately from contract creation.

At 2026-11-09T21:46:07Z the replacement passed trusted TLS, missing/incorrect
credential rejection (401), correct-login access (200), and authenticated Forge
configuration checks showing the Chroma checkpoint, T5 encoder and VAE. Only SSH
and Forge HTTPS appear in its public port map. The reviewed refresh-only apply
saved its connection outputs with zero resource changes. A subsequent metadata
probe encountered an intermittent TLS handshake timeout; the complete successful
acceptance is retained separately. No certificate check was bypassed.

The user's latest direction was to deploy and try the model directly. Automated
comic-generation and reference-image output assessment are therefore deferred;
UI availability does not establish inference quality or character continuity.

## Release and ComfyUI source verification

The previous Chroma Forge rental was destroyed through the owning targeted
interactive release at the user's request. A reviewed refresh-only apply then
removed connection outputs with zero resource changes. The retained Vault login
and team API key are unchanged. The ComfyUI source update performs no live apply.
Publisher metadata identifies the full Flash checkpoint as 17,800,038,288 bytes
at revision 093c4f63b60507234aa53a073f1f9f565b3f4336. Its encoder and VAE retain
their previously reviewed revisions. CLI fixtures verify native destinations and
credential redaction; they do not verify GPU inference or comic quality.
The maintained dark WebP was rendered and visually inspected: browser request
arrows point to Vault/Caddy, the provisioner is inside the rented instance, and
container titles and edge labels remain readable. The user explicitly requested the usual build-target integration: documentation
consumes its render target directly, with no checked-in raster or copy rule.

## Authorized ComfyUI deployment

The user requested deployment after merging the ComfyUI source, followed by a
rebase. RES4LYF at 6780724de0ddb98e7454d54bb80387836a911ece and its required
Python dependencies were explicitly approved for style-reference workflows.
Vast's native Git/pip provisioner owns installation. The three native workflow
documents are compacted and staged after workspace initialization.

Vast rejected expanded startup arguments above 16,384 bytes. The approved
package versions now use exact CPython 3.12 Linux x86_64 publisher wheels with
hashes, and Terraform checks the rendered size before creation. On 2026-11-10
at 08:19 UTC, the reviewed interactive apply created instance 60043006 from
offer 59544635: one addition, no changes or destruction. Vault login and team
API key were retained. The host reports verified Secure Cloud hosting_type=1,
one CMP 170HX with 40 GiB VRAM, 100 GB disk and USD 0.477037/hour.

Live acceptance passed trusted TLS, missing/wrong-login rejection (401),
authenticated access (200), all three model filenames, both approved RES4LYF
node types and all three staged workflows. Text-to-image, img2img and style-reference
sampling produced visually inspected 1024 x 1024 PNG artifacts. The style test
used the public text-to-image output as reference and generated a different
scene; it establishes execution compatibility, not reliable style isolation.
Vast also maps the image's diagnostic port 3000; its test server is disabled
and a direct probe found no reachable listener. A reviewed refresh-only apply
saved the connection outputs with zero resource changes. The source was rebased
onto master after creation; the delivery receipt owns publication state.

## Ignored style-reference correction

A subsequent user report exposed a gap in the initial execution-only acceptance.
At fixed prompt and seed, the original graph produced identical PNG image data
with guide weight zero and one. The graph omitted the approved pack's
ReChromaPatcher, so native Chroma did not consume its style transformer options.
The corrected document enables that model patch in float32 before sampling.
The active rental receives a separate corrected user document rather than an
infrastructure replacement; preserve current user graphs and instance files.

The corrected path changed output at the same prompt/seed when its guide weight
changed, and when a flat comic reference replaced the watercolor reference.
Both outputs were visually inspected. The live corrected document is named
chroma-flash-style-transfer and preserves the user's recent workflow settings;
it leaves the original document intact and is not overwritten by the old
startup's original filename. Future rentals receive the corrected canonical
style-reference document through the checked-in configuration.

## Requested Qwen replacement

The user selected Qwen-Image-Edit-2511 and requested redeployment, then confirmed
that no files from the current rental require preservation. Keep the existing
Vast image, Vault login, certified hosting constraints and hourly ceiling. Use
the official Comfy-Org mixed-precision checkpoint, Qwen 2.5 VL encoder and Qwen
VAE at fixed publisher revisions, matching the native workflow's default
checkpoint variant. This bundle is about 30.2 GB. Do not install another node
pack, acceleration LoRA or Python package; the default profile no longer needs
RES4LYF. Keep historical neutral Chroma documents as examples only.

Before calling the replacement ready, check these failure cases: mismatched
model/encoder/VAE names; missing or unused reference-image conditioning; omitted
index_timestep_zero reference method; incorrect flow shift or CFG normalization;
image resizing that discards the intended content; out-of-memory inference;
changes to credentials or public management-service exposure. Validate native
single-reference and separate character/style graphs against the published
recipe, then exercise new-action/new-background prompts using public neutral
references. Successful execution alone is insufficient acceptance.

## Replacement source validation

Offline checks passed for the moved rental root, shared modules, team catalog,
Vault formatting and affected OpenSpec workspaces. The catalog fixture applies
the real team key/Vault declarations against mocked providers; the rental
fixture generates a real random password while mocking Vault and the paid
module. Existing real Terraform CLI E2E acceptance covers exposed/SSH lifecycle,
startup compression and credential-safe rejection diagnostics. None contacts
live infrastructure.

Both Qwen graphs have consistent node/slot/link types and model filenames.
The diagram target was rendered and visually inspected with readable nested
containers and the existing dark actor theme. Inventory regeneration was stable
on its second run. New AppRole permissions and Qwen inference still require
separately authorized live acceptance; this delivery makes no runtime claim.

## Credential-only deployment

After merging PR #285 the user explicitly authorized applying IaC without renting
an instance. Inputs matched merged commit 611710800d460b633abc8c27a3d814720c3ffdcf.
Owning interactive applies were reviewed before each approval:

| UTC, 2026-11-10 | Owner                 | Scope and result                                                                                         |
| --------------- | --------------------- | -------------------------------------------------------------------------------------------------------- |
| 09:53           | Vault                 | Component AppRole, identity/group/policies and exact team-writer grant: 10 added, 1 changed, 0 destroyed |
| 09:54           | Vast team             | Key 36936194 and component delivery entry: 2 added, 0 changed, 0 destroyed                               |
| 09:55           | Component credentials | Password created; web entry written but apply failed on metadata read and tainted that entry             |
| 09:57           | Vault                 | Exact web policy corrected to read/delete: 0 added, 1 changed, 0 destroyed                               |
| 09:58           | Component credentials | Replace only the tainted web entry, retaining the generated password: 1 added, 0 changed, 1 destroyed    |

The component backend is initialized at its new setup path. All component-root
applies explicitly targeted random_password.web and vault_kv_secret_v2.web;
no instance module, offer selection or paid capacity was applied. The Vault
owner records the metadata-policy precedence failure and correction. Credential
values were not printed or saved locally. Runtime Qwen verification remains
outside this no-rental deployment.

At 2026-11-10T10:01:20Z, the new AppRole read the dedicated API key and valid
48-character alphanumeric web login. Effective-capability checks passed own
namespace/backend access, exact web metadata read/delete, team-writer login/state
denial and component team-issuer denial. The human operator can browse and read
the namespace. A read-only API request authenticated with the generated key
reported zero team instances and zero image-generation rentals. Credential
values remained in memory and were not printed. This establishes access and
absence of capacity, not ComfyUI readiness or model quality.

## Authorized Qwen rental acceptance

The user explicitly requested renting the instance and checking all supplied
workflows. On 2026-11-10 at 10:28 UTC, the owning interactive apply created
instance 60055445: 1 addition, 0 changes, 0 destructions. Existing credentials
were retained. The host reports verified Secure Cloud hosting_type=1, one CMP
170HX with 40 GiB VRAM, 100 GB disk and USD 0.455037/hour. A reviewed
refresh-only apply saved connection outputs with zero remote-resource changes.

Trusted TLS and HTTP authentication checks passed: missing/wrong login returned
401 and the Vault login returned 200. All three Qwen model filenames, required
native nodes and both staged workflow documents are available. The optional
custom-node pack is disabled. The image also declares diagnostic port 3000;
its mapped port refused a TCP probe. The active workflow library contains only
the two native Qwen documents; historical Chroma examples require their older
models and node pack and are not staged.

Acceptance uses only neutral robot-character and watercolor reference images,
a fixed seed, and a new snowy-park bicycle scene. Both staged native UI graphs
are serialized through the live node schemas and executed by ComfyUI, with
model, reference-latent and encoder connections intact. Output assessment and
its limits are recorded below after execution; availability alone is not
inference or style-fidelity acceptance. No user content is copied into source.

Both active workflows completed with the configured 40 Euler/simple steps,
CFG 4, shift 3.1 and denoise 1. The single-reference run completed in about
117 seconds including initial model load; the two-reference run completed
about 269 seconds after submission, including time queued behind the first.
Both produced visually inspected 1024 x 1024 PNGs. They retained the reference
robot's white body, antenna, black face and cyan eyes while replacing the
book/cafe composition with the requested bicycle and snowy park. The
character/style result uses visibly softer watercolor washes and paper texture,
particularly in the background. This is functional reference-guided generation,
not a guarantee of exact style matching or artifact-free comic continuity.

Neutral outputs and per-job receipts are retained under ignored
out/imagegen/qwen. No raster, user prompt or user upload is committed. Native
workflow execution required no custom nodes, extra packages or runtime repair.

The owning full Terraform plan after startup reported no changes. No image,
model manifest, workflow, access setting or credential replacement was needed
to pass the two active workflow checks.

A separate controlled pair kept seed 184463, the scene prompt, character image,
model and sampler settings fixed while changing only image 2 between a flat
comic and watercolor reference. Both completed at 1024 x 1024. Visual inspection
showed crisp shaded comic rendering for the flat reference and softer pigment
washes and paper texture for the watercolor reference. The robot and bicycle
scene remain recognizable in both. This demonstrates style-image influence
independently of explicit watercolor prompt wording; the character still retains
strong outlines, so the result is approximate style transfer. All four neutral
acceptance jobs completed without runtime errors.

## Subsequent requested release

The user explicitly requested destruction of the current instance. On
2026-11-10 at 11:16 UTC, the owning targeted interactive release destroyed
instance 60055445: 0 added, 0 changed, 1 destroyed. A reviewed refresh-only
apply then cleared instance, HTTPS and SSH connection outputs with zero remote
resource changes. A read-only Vast API lookup confirmed the instance is absent.
The dedicated AppRole, policies, API key, web login and component backend remain
managed. No replacement was rented and no user files were copied from the
instance. Earlier acceptance evidence describes the now-retired rental.

## Flux and CLI source acceptance

The approved official CLI is pinned to 1.25.0 with a separate hash-bearing lock.
Its real executable passed an isolated HTTPS fixture covering workflow submit,
WebSocket completion, status and output download, credential rejection and TLS
hostname rejection. Settings stay under ignored task storage and telemetry is
opted out. During source implementation, the live component AppRole retrieved the existing web login;
the owning Terraform output was null after the preceding release. The later authorized rental and token reconciliation are recorded below.

Publisher metadata pins the Flux checkpoint to 902d9d510b51533e07729f19211414a3648b77d2,
Qwen3 8B encoder to 1138ce884c8e977cdacacbe779784fec4426bb3f and the small-decoder
VAE to a3efc24f613ef42d9428af62fdbd6f5fd8856c4a. The encoder and VAE are public;
the official checkpoint HEAD returned 401 GatedRepo on 2026-11-10. After the private request timed out, the operator stored the token through
the component target. An in-memory HEAD check returned 302 with the pinned
publisher commit, confirming download access. The model-authentication root
reconciled that same value with a reviewed apply: 1 added, 0 changed, 0 destroyed.
Flux GPU workflow acceptance follows the separately reviewed rental below. Source and fixture
acceptance must not be presented as live model or style-fidelity verification.

The first-party Prompt Sequence contract check covers ordered nonempty prompts,
shared instructions, paragraph input and rejected filename prefixes. Native
list mapping is wired into both sequence graphs; the Flux three-panel GPU result
is recorded below. No user uploads or private prompts are
used in repository fixtures or source.

## Authorized Flux rental

After correcting the exact-offer search field to ask_contract_id, the owning
apply with gpu_name=RTX PRO 5000 created instance 60069855 on 2026-11-10 at
12:50 UTC: 1 added, 0 changed, 0 destroyed. Existing web login, dedicated team
key and AppRole were unchanged. The source retains Qwen as a selectable profile.
Startup, trusted HTTPS and actual GPU workflow results are recorded separately;
contract creation alone is not runtime acceptance.

## Flux live acceptance

The host reports verified Secure Cloud hosting_type=1, one RTX PRO 5000 Blackwell
with 48 GB nominal VRAM, 100 GB disk and USD 0.905852/hour. Native downloading
completed for all three fixed-revision model files. Trusted TLS, missing/wrong
login rejection (401), correct login (200), native node availability and all
three staged workflow documents passed. The mapped diagnostic port refused a
TCP connection. A reviewed refresh-only apply saved connection outputs with
zero remote-resource changes.

The new authenticated CLI reports ComfyUI 0.44.0, PyTorch 2.10.0+cu128 and a CUDA
RTX PRO 5000 device. Its live upload, workflow/WebSocket execution and download
commands passed. Single-reference generation produced one 1024-square image in
17.84 seconds including first model load; separate character/style produced one
in 6.52 seconds. Prompt Sequence produced three 1024-square images in input order
in 17.92 seconds, with paired comic_001, comic_002 and comic_003 output prefixes.
These are observed job timings, not a general hardware benchmark.

All five neutral outputs were downloaded and visually inspected. The robot's
white body, antenna, black face and cyan eyes remain recognizable, and cycling,
pond and reading scenes follow the requested actions. Style isolation remains
approximate: strong comic linework dominates the watercolor reference. The pond
image also introduces multiple panels, a duplicate robot and spurious lettering.
Execution and sequence ordering are accepted; exact layout, style and character
continuity are not guaranteed. No private references, prompts or outputs were
committed. Local neutral images, hashes, dimensions and CLI receipts remain in
ignored out/imagegen/flux.

The operator briefly requested destruction and immediately withdrew it. Keep
instance 60069855 running at that stage. The later Base replacement request above supersedes that instruction. Qwen remains an
offline-validated selectable option; its new sequence graph was not GPU-tested
on this Flux-only rental. Historical Qwen acceptance does not establish that new
graph's runtime results.

## Externally supplied model token

The component consumes the huggingface/token Vault entry supplied manually by
the operator or managed by external IaC. It does not own token issuance,
rotation or deletion. The separate tf root only copied an existing token and
introduced redundant Terraform ownership and state.

The operator authorized removing that root and detaching its state ownership
while preserving the Vault secret. Use a temporary declarative removed block
with destroy=false through the owning interactive apply. Check secret metadata
without reading token contents, then remove the transition and confirm the
owning plan is empty before deleting the root and unused main-backend wiring.
Historical protected backend versions are outside this cleanup; do not delete
them or the externally owned token. No rental operation is authorized.

## Token ownership cleanup receipt

On 2026-11-10 at 16:05:18 UTC, the reviewed owning apply relinquished only
vault_kv_secret_v2.huggingface using destroy=false: 0 added, 0 changed,
0 destroyed. The token's Vault current_version remained 2 before and after.
After removing the temporary block, the owning tf.plan reported No changes
at 16:05:41 UTC. The retired main backend's latest state no longer owns the
entry; protected historical versions remain. Token contents were not read.

## Configuration simplification

The component owns one config.json with its image, login name, model and
hardware profile selections and definitions. Terraform reads it directly; no
root input variables, separate tfvars preset or duplicate high-memory targets
are needed. Keep the standard selection and all existing model revisions,
workflow filenames and hardware limits unchanged. Only source validation is
authorized for the new configuration.

## Juggernaut profile restoration

The operator requested restoring Juggernaut and selecting it as the default for
the next rental. Reuse the previously approved RunDiffusion XL v9 checkpoint
at cf419233522daa0b9ea36c3aff98fa2cab1fb0fb, with its original filename. The earlier
move to Chroma followed the operator's request to try another model; it did
not establish that Juggernaut was unsuitable.

Juggernaut is an all-in-one SDXL checkpoint containing its encoders and VAE.
It needs the native models/checkpoints directory and CheckpointLoaderSimple;
Flux/Qwen keep models/diffusion_models. Add a constrained checkpoint-directory
option to the wrapper while preserving all base exposure and host controls.
Native generation uses SDXL nodes plus the existing first-party Prompt Sequence.
Text prompts can describe styles; native img2img also retains input composition
and does not isolate a reference's style or character. The reference-style paths
use only the previously approved pack and requirements described below. No new
external dependency or live operation is authorized.

The follow-up reference-style request restores RES4LYF at its previously approved
6780724de0ddb98e7454d54bb80387836a911ece revision and the byte-identical requirements
from commit 9919bb90a. The Juggernaut profile alone declares them in config.json;
other profiles still install no third-party pack. ComfyUI tools owns the SDXL
style graph and the normalized profile-purpose filenames. No new adapter, vision
encoder or package revision is introduced. GPU style quality remains unverified.

Offline source acceptance: all 14 workflow graphs have reciprocal typed links;
all nine renamed Flux, Qwen and Chroma documents retain their exact graph bytes.
The five Juggernaut graphs use the approved checkpoint filename, and the style
paths include the SDXL patcher, connected guide and fresh latent sampling.
Sequence widget links retain their native widget metadata. The restored lock is
byte-identical to the approved historical lock. Browser duplication and isolated
Terraform CLI acceptance pass. These checks establish source and fixture behavior,
not GPU inference or style quality. No rental or live Vault operation was run.

The Jsonnet authoring refactor retains the same parsed native workflow values
and package-relative filenames. The workflow generator receives only the public
config.json through an explicit file export and owner-local consumer group; no
Vault/backend credential or protected infrastructure artifact is a build input.
Terraform still reads native JSON from declared Bazel data. Empty download arrays
retain the existing workflow filenames and recorded model defaults. Shared
Jsonnet tooling is owned by infra/src/add-jsonnet-tooling; no rental is changed.

## Juggernaut rental and startup transport

On 2026-11-10 the operator explicitly requested a rental and its URL. The first
interactive apply rejected the expanded startup configuration before creating
capacity. The provider now owns opt-in gzip transport and the API byte limit;
the base module supplies plain JSON rather than pre-encoding each file.
The provider CLI fixture passed actual script execution, stable refresh,
external startup drift and oversized-input rejection.

The reviewed apply then created one standard Juggernaut rental with no other
resource changes. The selected certified datacenter offer reported one CMP
170HX, 40 GiB VRAM and approximately USD 0.497/hour including the requested
disk. Reviewed refresh-only applies recorded its connection metadata. The
official CA verified the live HTTPS endpoint, unauthenticated requests returned
401, and the existing Vault login accessed ComfyUI successfully. The native
checkpoint schema included the configured Juggernaut file. No user uploads,
history or outputs were read, and no image-generation jobs were submitted.
GPU generation and style quality remain unverified. Credentials and current
connection coordinates remain in Vault and the owning Terraform outputs.

## Native editor and multi-reference replacement

The operator rejected the growing prompt editor and requested the original
multiline box with blank lines for visual spacing. Stop staging its browser
script; the existing line parser already ignores blank lines. Restore the
upstream float64 SDXL style-patcher setting after the supplied WCT dtype error.
The approved IP-Adapter Plus pack and two public model files add two
multi-reference workflows to the Juggernaut profile. No new pip packages are
required by the general adapter path; the existing image supplies Torch,
Torchvision, NumPy and Pillow.

While the operator was using the instance, all investigation and source checks
remained local or used pinned public source. The operator then explicitly
requested destruction and a new deployment after the changes. A reviewed,
instance-targeted apply destroyed contract 60111260 on 2026-11-10, with one
resource destroyed and no login, key or AppRole changes. Targeted apply left
historical outputs cached; they are not a live connection. The requested
replacement must use the full reviewed plan and publish newly observed outputs.

## Replacement GPU acceptance

The operator explicitly requested tests of all workflows on the rented instance.
The full reviewed plan created replacement 60116333 on 2026-11-10 with the
configured Vast ComfyUI image, one CMP 170HX and 40 GiB VRAM. Reviewed refresh
recorded its new coordinates. Authenticated metadata reported ComfyUI 0.44.0
and Torch 2.10.0; the official root verified HTTPS, unauthenticated requests
returned 401, and the diagnostic port remained closed.

The owning authenticated CLI uploaded three first-party neutral geometric robot
images and ran all seven selected generated Juggernaut workflows. Only reference
filenames and output prefixes changed for fixture isolation; 1024-square resolution,
30-step sampling and guidance defaults were retained. All workflows completed,
producing the expected 13 PNGs. Dimensions, CRC, decompression and nonempty-pixel
checks passed. Both RES4LYF paths completed without the reported dtype failure;
the blank-spaced multi-reference prompt list produced three outputs. Two artifact
downloads timed out during GPU activity and succeeded after the queue finished.

The original WCT guide produced soft images with these simple fixtures. The
three-reference adapter retained the geometric appearance strongly at weight 0.8
but often omitted the requested action. A fixed-seed comparison at weight 0.5
and end_at 0.8 retained the broad reference style while adding the requested
reading scene. These results establish runtime execution and useful tuning,
not exact character identity or reliable scene adherence. The original
reference-heavy defaults remain available; guidance controls are documented.

Repeatable native workflow variants, prompt IDs, per-output hashes and neutral
images remain under ignored out/imagegen/startup/live; acceptance.json records
all seven executions and 13 verified images. No private user content was read.
Credentials remain in Vault and current coordinates in Terraform outputs.

## Requested FLUX.2 dev 32B replacement

After discussing the named FLUX.2 dev 32B model, its native ComfyUI FP8 diffusion
weights and local Mistral text encoder, the operator explicitly requested
destruction and deployment of that larger model. The operator then selected the
prepared larger 96 GB instance for the first run. This directly accepts the
identified model bundle; no new image, custom-node pack or Python dependency is
needed. Pin the Comfy-Org diffusion model, matching Mistral encoder and FLUX.2 VAE
at an observed immutable publisher revision. Preserve all existing model options.

Decision verdict: proceed with FP8 and the existing 96gb hardware profile,
95 GiB VRAM floor, 150 GB disk and USD 3/hour ceiling. The prior 48 GB/FP8 proposal
would be a cheaper trial, but the operator chose more memory for initial use.
Keep text encoding inside the rental. The official ComfyUI example supplies the
native FLUX.2 dev conditioning and sampling contracts; the 9B Klein workflow
must not be relabeled without adapting its different guidance behavior.

A reviewed targeted apply destroyed instance 60116333 on 2026-11-11, with zero
other resources changed. Its former connection outputs are historical until a
full reviewed rental apply and refresh publishes new coordinates. No private
uploads or outputs were inspected during destruction.

The full reviewed apply created instance 60123890 from offer 54022741, adding
one resource and changing no credentials. An initial reviewed refresh recorded
its IP and SSH endpoint before port mappings were available. Subsequent scoped
read-only Vast metadata reported running status, RTX PRO 6000 S, 97,887 MiB
VRAM, approximately 1 TiB host RAM and USD 1.7889/hour including the selected
disk. Its newly observed HTTPS endpoint was supplied to the operator.

The operator then began using the replacement and explicitly prohibited further
access. No new generation jobs, user-content inspection or runtime changes were
performed. Final connection-output refresh, authenticated health/TLS acceptance
and GPU workflow checks were not completed and require renewed authorization.
The source passes offline node, graph, Terraform and provider checks; those do
not establish inference quality for the new 32B deployment.

After operator use began, the operator requested the missing single-reference
workflow. The source bundle now also includes single-reference and matching
sequence JSON. These were compiled and supplied as importable operator artifacts;
they were not installed on the rental. The active initial three-workflow bundle
and the five-workflow desired configuration therefore differ until a later
authorized deployment. This follow-up does not authorize a replacement.

## Full Chroma deployment on RTX PRO 5000

The operator subsequently requested destruction of instance 60127015. Its
reviewed instance-targeted apply destroyed one resource; a reviewed refresh-only
apply cleared its connection outputs while retaining component credentials.
No generation jobs or user-content inspection were performed during teardown.

The operator now explicitly requests the full Chroma deployment, the earlier
fast 48 GB GPU and generation tests. The recorded GPU was RTX PRO 5000
Blackwell, not the slower CMP 170HX. Decision: proceed with a named 48gb hardware
profile, exact GPU filter RTX PRO 5000, 47 GiB reported-memory floor, 100 GB disk
and the existing USD 1/hour ceiling. Do not substitute another GPU silently.

Reuse the previously approved Chroma1-HD revision
0e0c60ece1e82b17cb7f77342d765ba5024c40c0, FP16 T5 and FLUX VAE pins. Keep the
Vast ComfyUI image and existing approved RES4LYF revision and hash-locked
requirements. The full checkpoint has ordinary CFG; use 40 steps and CFG 3 as
the publisher's documented starting point. Native img2img retains composition;
RES4LYF is a style guide, not a guarantee of character identity or prompt fidelity.
Acceptance uses only first-party neutral fixtures and must inspect generated
images rather than equating a completed job with good reference adherence.
The unpublished diagram remains a local preview and is outside this source update.

The reviewed interactive apply created instance 60180986 from offer 44193892 on
2026-11-11: one addition, no changes or destruction. A reviewed refresh-only
apply saved its connection outputs with no remote changes. Authenticated metadata
confirms RTX PRO 5000 Blackwell with 50,752,520,192 bytes of VRAM, ComfyUI 0.44.0
and the exact Prompt Sequence/RES4LYF whitelist. Official-CA TLS verification,
missing-login rejection (401) and an unreachable diagnostic port passed.

All five configured workflows completed at 1024 square, 40 steps and CFG 3 with
neutral first-party fixtures. Text-to-image took 49.16 seconds including initial
loading, img2img 34.17 seconds, style-reference 70.33 seconds, the three-image
native sequence 103.07 seconds and the three-image style sequence 208.62 seconds.
A fixed-seed style-guide-disabled comparison completed in 38.17 seconds. Ten
downloaded PNGs passed dimensions, CRC, decompression and nonempty-pixel checks;
sequence filenames preserve their three input positions. The three archived
Flash graphs also match their parsed pre-refactor baseline exactly.

Visual assessment separates execution from quality. Native text-to-image follows
the neutral robot/book/bench/tree prompt. Img2img retains the original geometric
reference and omits the requested new action. The style guide has a clear visual
effect in the fixed-seed comparison, matching the cream/blue/red palette while
softening detail. Native sequence prompts containing comic wording produce
unwanted page layouts and lettering; style sequences preserve palette but have
unreliable subjects, compositions and actions. These paths are available for
operator tuning, not certified reference fidelity. No private user content was
read, and only neutral prompts and fixtures were used. Reproducible workflow
variants, results, hashes and images remain under ignored
out/imagegen/chroma-hd; acceptance.json owns the detailed local receipt.

### Grouped progress and chained GIF acceptance

The operator prioritizes one media group per sequence over separate queue jobs
and requests a unique prefix for each run. Keep native list mapping, append a
model pass-through to Prompt Sequence, and report the current panel through the
native sampling wrapper. Generate a fresh shared UUID prefix on each execution.
The optional model input preserves saved text-widget and output positions.

Controlled style comparisons favored weight 0.5 through step 15 over weight 1
through all 40 steps: the robot/book/bench composition became clearer while
retaining the reference palette. Those are the HD defaults; archived Flash
settings are unchanged. Source and offline checks were completed before renting
again. No new image, model, custom-node dependency or Python package was added.
The first-party GIF node uses the existing native sampler, VAE, Torch and Pillow.

The reviewed temporary apply created instance 60189648 from offer 53904695 on
2026-11-11, adding one resource and changing no credentials. Authenticated
metadata confirmed ComfyUI 0.44.0 and RTX PRO 5000 Blackwell with
50,752,520,192 bytes of VRAM. All six selected workflows and two controls
completed: native img2img 38.63 seconds, native three-panel sequence 106.56,
style sequence 158.27, style-reference 55.17, text-to-image 38.19, chained
three-frame GIF 135.00, repeated native sequence 106.25 and identical-frame
GIF control 6.46 seconds. These are observed job times, not latency guarantees.

Native WebSocket progress reported Panel 1 of 3 through Panel 3 of 3 for both
sequence variants and reset on the repeated run. Each sequence remained one
job with three output images; repeated runs used different UUID prefixes and
retained 001/002/003 suffixes. Twelve downloaded PNGs passed 1024-square,
CRC and decompression checks. Both GIFs contain three 1024-square frames,
170-millisecond intervals and infinite looping. The zero-denoise control
retained three identical frames rather than merging them. The isolated native
sampler boundary additionally verifies previous-latent feedback and ordered
positive/shared-negative conditioning with real Torch/Pillow media output.

Visual inspection finds clearer requested actions and matching cream/blue/red
colors with the revised guide, with residual softness and character variation
across panels. Chained img2img preserves much of the initial composition at
denoise 0.45; it does not promise smooth motion or precise pose changes.
All fixtures were neutral and first-party. Results and progress receipts remain
under ignored out/imagegen/chroma-hd. After downloading them, the reviewed
instance-only apply destroyed the rental at 09:27:59 UTC; refresh-only apply
cleared its connection outputs. Component credentials remain available for the
next manual rental. No paid instance remains and the diagram is unpublished.
