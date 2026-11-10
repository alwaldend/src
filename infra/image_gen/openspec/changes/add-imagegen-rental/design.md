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
