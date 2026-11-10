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
